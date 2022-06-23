/*
Package otto an automation toolkit for Unix-like computers.

This package contains the common interfaces and methods shared by the Otto client and Otto server.
*/
package otto

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/ecnepsnai/logtic"
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
	gob.Register(MessageRotateIdentityRequest{})
	gob.Register(MessageRotateIdentityResponse{})
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
	MessageTypeRotateIdentityRequest
	MessageTypeRotateIdentityResponse

	MessageTypeGeneralFailure = uint32(0xFFFFFFFF)
)

// MessageHeartbeatRequest describes a heartbeat request
type MessageHeartbeatRequest struct {
	Version string `json:"server_version"`
	Nonce   string `json:"nonce"`
}

// MessageHeartbeatResponse describes a heartbeat response
type MessageHeartbeatResponse struct {
	ClientVersion string            `json:"client_version"`
	Properties    map[string]string `json:"properties"`
	Nonce         string            `json:"nonce"`
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
	Error         string       `json:"error"`
	File          File         `json:"file"`
	ClientVersion string       `json:"client_version"`
}

// MessageRotateIdentityRequest describes a request to rotate an identity
type MessageRotateIdentityRequest struct {
	PublicKey string `json:"public_key"`
}

// MessageRotateIdentityResponse describes the response for rotating an identity
type MessageRotateIdentityResponse struct {
	PublicKey string `json:"public_key"`
	Error     string `json:"error"`
}

// MessageGeneralFailure describes a general failure
type MessageGeneralFailure struct {
	Error string `json:"error"`
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
	case MessageTypeRotateIdentityRequest:
		message := MessageRotateIdentityRequest{}
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&message); err != nil {
			return nil, err
		}
		return message, nil
	case MessageTypeRotateIdentityResponse:
		message := MessageRotateIdentityResponse{}
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

const sshChannelName = "otto"

// Connection describes a connection between the Otto Server and Otto Host
type Connection struct {
	id             int
	w              io.ReadWriteCloser
	remoteAddr     net.Addr
	localAddr      net.Addr
	remoteIdentity []byte
	localIdentity  []byte
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

func (c *Connection) RemoteIdentity() []byte {
	return c.remoteIdentity
}

func (c *Connection) LocalIdentity() []byte {
	return c.localIdentity
}

func (c *Connection) Close() error {
	log.PDebug("Connection closed", map[string]interface{}{
		"id":          c.id,
		"local_addr":  c.localAddr.String(),
		"remote_addr": c.remoteAddr.String(),
	})
	return c.w.Close()
}
