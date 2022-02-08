/*
Package otto an automation toolkit for Unix-like computers.

This package contains the common interfaces and methods shared by the Otto client and Otto server.
*/
package otto

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/ecnepsnai/logtic"
	"golang.org/x/crypto/ssh"
)

var log = logtic.Log.Connect("libotto")

// ProtocolVersion the version of the otto protocol
const ProtocolVersion = uint32(2)

func init() {
	gob.Register(MessageHeartbeatRequest{})
	gob.Register(MessageHeartbeatResponse{})
	gob.Register(MessageTriggerAction{})
	gob.Register(MessageCancelAction{})
	gob.Register(MessageActionOutput{})
	gob.Register(MessageActionResult{})
	gob.Register(MessageGeneralFailure{})

	gob.Register(Script{})
	gob.Register(ScriptResult{})
	gob.Register(File{})
}

// Message types
const (
	MessageTypeKeepalive uint32 = iota + 1
	MessageTypeHeartbeatRequest
	MessageTypeHeartbeatResponse
	MessageTypeTriggerAction
	MessageTypeCancelAction
	MessageTypeActionOutput
	MessageTypeActionResult
	MessageTypeGeneralFailure
)

// MessageHeartbeatRequest describes a heartbeat request
type MessageHeartbeatRequest struct {
	ServerVersion string `json:"server_version"`
}

// MessageHeartbeatResponse describes a heartbeat response
type MessageHeartbeatResponse struct {
	ClientVersion string            `json:"client_version"`
	Properties    map[string]string `json:"properties"`
}

// MessageTriggerAction describes an action trigger
type MessageTriggerAction struct {
	Action uint32 `json:"action"`
	Script Script `json:"script"`
	File   File   `json:"file"`
}

// MessageCancelAction describes a request to cancel an action
type MessageCancelAction struct{}

// MessageActionOutput describes output from an action
type MessageActionOutput struct {
	Stdout []byte `json:"stdout"`
	Stderr []byte `json:"stderr"`
}

// MessageActionResult describes the result of a triggered action
type MessageActionResult struct {
	ScriptResult  ScriptResult `json:"script_result"`
	Error         error        `json:"error"`
	File          File         `json:"file"`
	ClientVersion string       `json:"client_version"`
}

// MessageGeneralFailure describes a general failure
type MessageGeneralFailure struct {
	Error error `json:"error"`
}

// Actions
const (
	ActionRunScript uint32 = iota + 1
	ActionReloadConfig
	ActionUploadFile
	ActionUploadFileAndExitClient
	ActionExitClient
	ActionReboot
	ActionShutdown
	ActionUpdateIdentity
)

// Script describes a script
type Script struct {
	Name             string            `json:"name"`
	RunAs            RunAs             `json:"run_as"`
	Environment      map[string]string `json:"environment"`
	WorkingDirectory string            `json:"working_directory"`
	Executable       string            `json:"executable"`
	Files            []File            `json:"files"`
	Data             []byte            `json:"data"`
}

// RunAs describes the user to run a script as
type RunAs struct {
	Inherit bool   `json:"inherit"`
	UID     uint32 `json:"uid"`
	GID     uint32 `json:"gid"`
}

// ScriptResult describes the result of the script
type ScriptResult struct {
	Success   bool          `json:"success"`
	ExecError string        `json:"exec_error"`
	Code      int           `json:"code"`
	Stdout    string        `json:"stdout"`
	Stderr    string        `json:"stderr"`
	Elapsed   time.Duration `json:"elapsed"`
}

func (sr ScriptResult) String() string {
	return logtic.StringFromParameters(map[string]interface{}{
		"success":    sr.Success,
		"exec_error": sr.ExecError,
		"code":       sr.Code,
		"stdout":     strings.ReplaceAll(sr.Stdout, "\n", "\\n"),
		"stderr":     strings.ReplaceAll(sr.Stderr, "\n", "\\n"),
		"elapsed":    sr.Elapsed.String(),
	})
}

// File Describes a file
type File struct {
	Path  string `json:"path"`
	Owner RunAs  `json:"owner"`
	Mode  uint32 `json:"mode"`
	Data  []byte `json:"data"`
}

// RegisterRequest describes a register request
type RegisterRequest struct {
	ClientIdentity string                    `json:"identity"`
	Port           uint32                    `json:"port"`
	Properties     RegisterRequestProperties `json:"properties"`
}

// RegisterRequestProperties describes properties for a register request
type RegisterRequestProperties struct {
	Hostname            string `json:"hostname"`
	KernelName          string `json:"kernel_name"`
	KernelVersion       string `json:"kernel_version"`
	DistributionName    string `json:"distribution_name"`
	DistributionVersion string `json:"distribution_version"`
}

// RegisterResponse describes the response to a register request
type RegisterResponse struct {
	ServerIdentity string   `json:"identity"`
	Scripts        []Script `json:"scripts,omitempty"`
}

func readFrame(r io.Reader) ([]byte, error) {
	versionBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, versionBuf); err != nil {
		if err == io.EOF {
			// Client closed - nothing to read
			return nil, nil
		}

		log.Error("Error reading version: %s", err.Error())
		return nil, err
	}
	version := binary.BigEndian.Uint32(versionBuf)
	if version != ProtocolVersion {
		log.Warn("Unsupported protocol version: %d, wanted: %d", version, ProtocolVersion)
	}

	dataLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, dataLengthBuf); err != nil {
		log.Error("Error reading data length: %s", err.Error())
		return nil, err
	}
	dataLength := binary.BigEndian.Uint32(dataLengthBuf)
	log.PDebug("Read frame", map[string]interface{}{
		"version":     version,
		"data_length": dataLength,
	})
	if dataLength == 0 {
		return []byte{}, nil
	}

	data := make([]byte, dataLength)
	readLength, err := io.ReadFull(r, data)
	if err != nil {
		if err == io.ErrUnexpectedEOF {
			log.PError("Incorrect data length", map[string]interface{}{
				"reported": dataLength,
				"actual":   readLength,
			})
			return nil, fmt.Errorf("bad request length")
		}
		log.Error("Error reading data: %s", err.Error())
		return nil, err
	}

	return data, nil
}

// ReadMessage try to read a message from the given reader. Returns the message type, the message data, or an error
func (c *Connection) ReadMessage() (uint32, interface{}, error) {
	data, err := readFrame(c.w)
	if err != nil {
		log.Error("Error reading data: %s", err.Error())
		return 0, nil, err
	}
	if data == nil {
		return 0, nil, nil
	}

	messageType := binary.BigEndian.Uint32(data[:4])
	log.PDebug("Read message", map[string]interface{}{
		"message_type": messageType,
		"data_length":  len(data) - 4,
	})
	message, err := DecodeMessage(messageType, data[4:])
	if err != nil {
		log.Error("Error decoding message data: %s", err.Error())
		return 0, nil, err
	}

	return messageType, message, nil
}

// WriteMessage try to write a message to the given writer.
func (c *Connection) WriteMessage(messageType uint32, message interface{}) error {
	messageData, err := EncodeMessage(messageType, message)
	if err != nil {
		log.Error("Error encoding message: %s", err.Error())
		return err
	}

	messageTypeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(messageTypeBuf, messageType)

	messageDataLength := len(messageTypeBuf)
	dataLength := messageDataLength + len(messageData)
	data := make([]byte, dataLength)
	i := 0
	for _, b := range messageTypeBuf {
		data[i] = b
		i++
	}
	for _, b := range messageData {
		data[i] = b
		i++
	}

	log.PDebug("Preparing message", map[string]interface{}{
		"message_type":        messageType,
		"message_data_length": messageDataLength,
		"total_length":        dataLength,
	})
	return writeFrame(data, c.w)
}

func writeFrame(data []byte, w io.Writer) error {
	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, ProtocolVersion)

	dataLength := len(data)
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(dataLength))

	replyLength := len(versionBuf) + len(lenBuf) + dataLength
	replyBuf := make([]byte, replyLength)
	i := 0
	for _, b := range versionBuf {
		replyBuf[i] = b
		i++
	}
	for _, b := range lenBuf {
		replyBuf[i] = b
		i++
	}
	for _, b := range data {
		replyBuf[i] = b
		i++
	}

	wrote, err := w.Write(replyBuf)
	log.PDebug("Wrote frame", map[string]interface{}{
		"data_length": dataLength,
		"version":     ProtocolVersion,
		"total":       wrote,
	})
	if wrote != replyLength {
		log.Error("Unable to write all of reply: wrote=%d total=%d", wrote, replyLength)
		return fmt.Errorf("out of space")
	}
	if err != nil {
		log.Error("Error writing data: %s", err.Error())
		return err
	}
	return nil
}

// DecodeMessage try to decode the given message. The returned object should match the message struct for the message
// type.
func DecodeMessage(messageType uint32, data []byte) (interface{}, error) {
	switch messageType {
	case MessageTypeKeepalive:
		return nil, nil
	case MessageTypeHeartbeatRequest:
		message := MessageHeartbeatRequest{}
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
			return nil, err
		}
		return message, nil
	case MessageTypeHeartbeatResponse:
		message := MessageHeartbeatResponse{}
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
			return nil, err
		}
		return message, nil
	case MessageTypeTriggerAction:
		message := MessageTriggerAction{}
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
			return nil, err
		}
		return message, nil
	case MessageTypeCancelAction:
		message := MessageCancelAction{}
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
			return nil, err
		}
		return message, nil
	case MessageTypeActionOutput:
		message := MessageActionOutput{}
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
			return nil, err
		}
		return message, nil
	case MessageTypeActionResult:
		message := MessageActionResult{}
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
			return nil, err
		}
		return message, nil
	case MessageTypeGeneralFailure:
		message := MessageGeneralFailure{}
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
			return nil, err
		}
		return message, nil
	}

	return nil, fmt.Errorf("unknown message type %d", messageType)
}

// EncodeMessage try to encode the given message
func EncodeMessage(messageType uint32, message interface{}) ([]byte, error) {
	if message == nil {
		return []byte{}, nil
	}

	buf := &bytes.Buffer{}

	if err := gob.NewEncoder(buf).Encode(message); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *Connection) Close() error {
	log.PDebug("Connection closed", map[string]interface{}{
		"local_addr":  c.localAddr.String(),
		"remote_addr": c.remoteAddr.String(),
	})
	return c.w.Close()
}

// Identity is DER encoded private key
type Identity []byte

// NewIdentity will generate a new ed25519 identity
func NewIdentity() (Identity, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return privateKeyBytes, nil
}

// ParseIdentity will parse the data as an identity
func ParseIdentity(data []byte) (Identity, error) {
	pkey, err := x509.ParsePKCS8PrivateKey(data)
	if err != nil {
		return nil, err
	}
	if _, err := ssh.NewSignerFromKey(pkey); err != nil {
		return nil, err
	}
	return data, nil
}

// Signer return the SSH signer for the identity
func (i Identity) Signer() ssh.Signer {
	privateKey, err := x509.ParsePKCS8PrivateKey(i)
	if err != nil {
		panic(err)
	}

	signer, err := ssh.NewSignerFromKey(privateKey)
	if err != nil {
		panic(err)
	}

	return signer
}

// PublicKey will return a DER-encoded representation of the public key for this identity
func (i Identity) PublicKey() ssh.PublicKey {
	return i.Signer().PublicKey()
}

// String will return a base64-encoded representation of the identity
func (i Identity) String() string {
	return base64.StdEncoding.EncodeToString(i)
}

// PublicKey will return a base64-encoded representation of the public key for this identity
func (i Identity) PublicKeyString() string {
	return base64.StdEncoding.EncodeToString(i.PublicKey().Marshal())
}

const sshChannelName = "otto"

// Connection describes a connection between the Otto Server and Otto Host
type Connection struct {
	w          io.ReadWriteCloser
	remoteAddr net.Addr
	localAddr  net.Addr
}

func MockConnection(w io.ReadWriteCloser) *Connection {
	return &Connection{
		w: w,
	}
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *Connection) LocalAddr() net.Addr {
	return c.localAddr
}

// ListenOptions describes options for listening
type ListenOptions struct {
	Address          string
	AllowFrom        []net.IPNet
	Identity         ssh.Signer
	TrustedPublicKey string
}

// Listener describes an active listening Otto server
type Listener struct {
	options   ListenOptions
	sshConfig *ssh.ServerConfig
	handle    func(conn *Connection)
	l         net.Listener
}

// SetupListener will prepare a listening socket for incoming connections. No connections are accepted until you call
// Accept().
func SetupListener(options ListenOptions, handle func(conn *Connection)) (*Listener, error) {
	sshConfig := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			log.PDebug("Handshake", map[string]interface{}{
				"public_key": base64.StdEncoding.EncodeToString(pubKey.Marshal()),
			})
			if options.TrustedPublicKey == base64.StdEncoding.EncodeToString(pubKey.Marshal()) {
				log.Debug("Recognized public key")
				return &ssh.Permissions{
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(pubKey),
					},
				}, nil
			}
			log.PWarn("Rejecting connection from untrusted public key", map[string]interface{}{
				"public_key": base64.StdEncoding.EncodeToString(pubKey.Marshal()),
			})
			return nil, fmt.Errorf("unknown public key %x", pubKey.Marshal())
		},
		ServerVersion: fmt.Sprintf("SSH-2.0-OTTO-%d", ProtocolVersion),
	}
	sshConfig.AddHostKey(options.Identity)

	l, err := net.Listen("tcp", options.Address)
	if err != nil {
		return nil, err
	}
	log.Info("Otto client listening on %s", options.Address)
	return &Listener{
		options:   options,
		sshConfig: sshConfig,
		handle:    handle,
		l:         l,
	}, nil
}

// Accept will accpet incoming connections. Blocking.
func (l *Listener) Accept() {
	for {
		c, err := l.l.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				return
			}
			log.PDebug("Error accepting connection", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}
		log.PDebug("Incoming connection", map[string]interface{}{
			"remote_addr": c.RemoteAddr().String(),
		})
		go l.accept(c)
	}
}

// Close will stop the listener.
func (l *Listener) Close() {
	l.l.Close()
}

func (l *Listener) accept(c net.Conn) {
	if len(l.options.AllowFrom) > 0 {
		allow := false
		for _, allowNet := range l.options.AllowFrom {
			if allowNet.Contains(c.RemoteAddr().(*net.TCPAddr).IP) {
				log.PDebug("Connection allowed by rule", map[string]interface{}{
					"remote_addr":     c.RemoteAddr().String(),
					"allowed_network": allowNet.String(),
				})
				allow = true
				break
			}
		}
		if !allow {
			log.PWarn("Rejecting connection from server outside of allowed network", map[string]interface{}{
				"remote_addr":  c.RemoteAddr().String(),
				"allowed_addr": l.options.AllowFrom,
			})
			c.Close()
			return
		}
	}

	_, chans, reqs, err := ssh.NewServerConn(c, l.sshConfig)
	if err != nil {
		log.PError("SSH handshake error", map[string]interface{}{
			"remote_addr": c.RemoteAddr().String(),
			"error":       err.Error(),
		})
		c.Close()
		return
	}

	go ssh.DiscardRequests(reqs)

	newChannel := <-chans

	log.Debug("ssh channel opened")
	if newChannel.ChannelType() != sshChannelName {
		log.PError("Unknown SSH channel", map[string]interface{}{
			"channel_type": newChannel.ChannelType(),
			"remote_addr":  c.RemoteAddr().String(),
		})
		newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
		return
	}
	channel, _, err := newChannel.Accept()
	if err != nil {
		log.PError("SSH channel error", map[string]interface{}{
			"remote_addr": c.RemoteAddr().String(),
			"error":       err.Error(),
		})
		return
	}
	log.PDebug("SSH handshake success", map[string]interface{}{
		"remote_addr": c.RemoteAddr().String(),
	})
	l.handle(&Connection{
		w:          channel,
		remoteAddr: c.RemoteAddr(),
		localAddr:  c.LocalAddr(),
	})
	channel.Close()
}

// DialOptions describes options for dialing to a host
type DialOptions struct {
	Network          string
	Address          string
	Identity         ssh.Signer
	TrustedPublicKey string
	Timeout          time.Duration
}

// Dial will dial the host specified by the options and perform a SSH handshake with it.
func Dial(options DialOptions) (*Connection, error) {
	clientConfig := &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(options.Identity),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			log.PDebug("Handshake", map[string]interface{}{
				"public_key": base64.StdEncoding.EncodeToString(key.Marshal()),
			})
			if options.TrustedPublicKey == base64.StdEncoding.EncodeToString(key.Marshal()) {
				log.Debug("Recognized public key")
				return nil
			}
			log.PWarn("Rejecting connection from untrusted public key", map[string]interface{}{
				"public_key": base64.StdEncoding.EncodeToString(key.Marshal()),
			})
			return fmt.Errorf("unknown public key: %x", key.Marshal())
		},
		HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
		ClientVersion:     fmt.Sprintf("SSH-2.0-OTTO-%d", ProtocolVersion),
		Timeout:           options.Timeout,
	}

	log.PDebug("Dialing", map[string]interface{}{
		"network": options.Network,
		"address": options.Address,
		"timeout": options.Timeout.String(),
	})
	client, err := ssh.Dial(options.Network, options.Address, clientConfig)
	if err != nil {
		log.PError("Error connecting to host", map[string]interface{}{
			"address": options.Address,
			"error":   err.Error(),
		})
		return nil, err
	}

	log.PDebug("Opening channel", map[string]interface{}{
		"address":      options.Address,
		"channel_name": sshChannelName,
	})
	channel, _, err := client.OpenChannel(sshChannelName, nil)
	if err != nil {
		log.PError("Error connecting to host", map[string]interface{}{
			"address": options.Address,
			"error":   err.Error(),
		})
		return nil, err
	}
	log.PDebug("Connected to host", map[string]interface{}{
		"address": options.Address,
	})

	return &Connection{
		w:          channel,
		remoteAddr: client.RemoteAddr(),
		localAddr:  client.LocalAddr(),
	}, nil
}
