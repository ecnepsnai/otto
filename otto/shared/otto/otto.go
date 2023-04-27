/*
Package otto an automation toolkit for Unix-like computers.

This package contains the common interfaces and methods shared by the Otto agent and Otto server.
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
	"golang.org/x/crypto/ssh"
)

var log = logtic.Log.Connect("libotto")

// ProtocolVersion the version of the otto protocol
const ProtocolVersion = uint32(4)

func init() {
	gob.Register(ScriptInfo{})
	gob.Register(ScriptResult{})
	gob.Register(FileInfo{})
	gob.Register(MessageGeneralFailure{})
	gob.Register(MessageHeartbeatRequest{})
	gob.Register(MessageHeartbeatResponse{})
	gob.Register(MessageRotateIdentityRequest{})
	gob.Register(MessageRotateIdentityResponse{})
	gob.Register(MessageTriggerActionRunScript{})
	gob.Register(MessageTriggerActionUploadFile{})
	gob.Register(MessageActionOutput{})
	gob.Register(MessageActionResult{})
}

type MessageType uint32

// Otto message types
const (
	MessageTypeGeneralFailure = MessageType(0)
	MessageTypeKeepalive      = MessageType(iota)
	MessageTypeHeartbeatRequest
	MessageTypeHeartbeatResponse
	MessageTypeRotateIdentityRequest
	MessageTypeRotateIdentityResponse
	MessageTypeTriggerActionRunScript
	MessageTypeTriggerActionReloadConfig
	MessageTypeTriggerActionUploadFile
	MessageTypeTriggerActionExitAgent
	MessageTypeTriggerActionReboot
	MessageTypeTriggerActionShutdown
	MessageTypeCancelAction
	MessageTypeActionOutput
	MessageTypeActionResult
)

// MessageHeartbeatRequest describes a heartbeat request
type MessageHeartbeatRequest struct {
	Version string `json:"server_version"`
	Nonce   string `json:"nonce"`
}

// MessageHeartbeatResponse describes a heartbeat response
type MessageHeartbeatResponse struct {
	AgentVersion string            `json:"agent_version"`
	Properties   map[string]string `json:"properties"`
	Nonce        string            `json:"nonce"`
}

// MessageTriggerActionRunScript
type MessageTriggerActionRunScript struct {
	ScriptInfo
}

// MessageTriggerActionUploadFile
type MessageTriggerActionUploadFile struct {
	FileInfo
}

// MessageActionOutput describes output from an action
type MessageActionOutput struct {
	Stdout []byte `json:"stdout"`
	Stderr []byte `json:"stderr"`
}

// MessageActionResult describes the result of a triggered action
type MessageActionResult struct {
	ScriptResult ScriptResult `json:"script_result"`
	Error        string       `json:"error"`
	AgentVersion string       `json:"agent_version"`
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

// ScriptInfo describes information about a script
type ScriptInfo struct {
	Name             string            `json:"name"`
	RunAs            RunAs             `json:"run_as"`
	Environment      map[string]string `json:"environment"`
	WorkingDirectory string            `json:"working_directory"`
	Executable       string            `json:"executable"`
	Length           uint64            `json:"length"`
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

// FileInfo describes information about file
type FileInfo struct {
	Path        string `json:"path"`
	Owner       RunAs  `json:"owner"`
	Mode        uint32 `json:"mode"`
	Checksum    string `json:"checksum"`
	AfterScript bool   `json:"after_script"`
	Length      uint64 `json:"length"`
}

// RegisterRequest describes a register request
type RegisterRequest struct {
	AgentIdentity string                    `json:"identity"`
	Port          uint32                    `json:"port"`
	Nonce         string                    `json:"nonce"`
	Properties    RegisterRequestProperties `json:"properties"`
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
	ServerIdentity string `json:"identity"`
	Nonce          string `json:"nonce"`
}

// ReadMessage try to read a message from the given reader. Returns the message type, the message data, or an error.
// Depending on the message type, there may be additional data to read following the message. It is up to the caller to
// continue reading any additional data.
func (c *Connection) ReadMessage() (MessageType, interface{}, error) {
	versionBuf := make([]byte, 4)
	if _, err := io.ReadFull(c.w, versionBuf); err != nil {
		if err == io.EOF {
			// Agent closed - nothing to read
			return 0, nil, err
		}

		log.Error("Error reading version: %s", err.Error())
		return 0, nil, err
	}
	version := binary.BigEndian.Uint32(versionBuf)
	if version > ProtocolVersion {
		log.PError("Unsupported protocol version", map[string]interface{}{
			"frame_version":     version,
			"supported_version": ProtocolVersion,
		})
		return 0, nil, fmt.Errorf("unsupported protocol version %d", version)
	}
	if version < ProtocolVersion {
		log.Warn("Unsupported protocol version: %d, wanted: %d", version, ProtocolVersion)
	}
	messageTypeBuf := make([]byte, 4)
	if _, err := io.ReadFull(c.w, messageTypeBuf); err != nil {
		if err == io.EOF {
			// Agent closed - nothing to read
			return 0, nil, err
		}

		log.Error("Error reading message type: %s", err.Error())
		return 0, nil, err
	}
	messageType := MessageType(binary.BigEndian.Uint32(messageTypeBuf))
	dataLengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(c.w, dataLengthBuf); err != nil {
		if err == io.EOF {
			// Agent closed - nothing to read
			return 0, nil, err
		}

		log.Error("Error reading data length: %s", err.Error())
		return 0, nil, err
	}
	dataLength := binary.BigEndian.Uint32(dataLengthBuf)
	if dataLength == 0 {
		return messageType, nil, nil
	}

	log.PDebug("Read message", map[string]interface{}{
		"version":      version,
		"message_type": messageType,
		"data_length":  dataLength,
	})

	decoder := gob.NewDecoder(c.w)

	switch messageType {
	case MessageTypeGeneralFailure:
		message := MessageGeneralFailure{}
		if err := decoder.Decode(&message); err != nil {
			log.Error("Error decoding MessageTypeGeneralFailure: %s", err.Error())
			return 0, nil, err
		}
		return messageType, message, nil
	case MessageTypeHeartbeatRequest:
		message := MessageHeartbeatRequest{}
		if err := decoder.Decode(&message); err != nil {
			log.Error("Error decoding MessageTypeHeartbeatRequest: %s", err.Error())
			return 0, nil, err
		}
		return messageType, message, nil
	case MessageTypeHeartbeatResponse:
		message := MessageHeartbeatResponse{}
		if err := decoder.Decode(&message); err != nil {
			log.Error("Error decoding MessageTypeHeartbeatResponse: %s", err.Error())
			return 0, nil, err
		}
		return messageType, message, nil
	case MessageTypeRotateIdentityRequest:
		message := MessageRotateIdentityRequest{}
		if err := decoder.Decode(&message); err != nil {
			log.Error("Error decoding MessageTypeRotateIdentityRequest: %s", err.Error())
			return 0, nil, err
		}
		return messageType, message, nil
	case MessageTypeRotateIdentityResponse:
		message := MessageRotateIdentityResponse{}
		if err := decoder.Decode(&message); err != nil {
			log.Error("Error decoding MessageTypeRotateIdentityResponse: %s", err.Error())
			return 0, nil, err
		}
		return messageType, message, nil
	case MessageTypeTriggerActionRunScript:
		message := MessageTriggerActionRunScript{}
		if err := decoder.Decode(&message); err != nil {
			log.Error("Error decoding MessageTypeTriggerActionRunScript: %s", err.Error())
			return 0, nil, err
		}
		return messageType, message, nil
	case MessageTypeTriggerActionUploadFile:
		message := MessageTriggerActionUploadFile{}
		if err := decoder.Decode(&message); err != nil {
			log.Error("Error decoding MessageTypeTriggerActionUploadFile: %s", err.Error())
			return 0, nil, err
		}
		return messageType, message, nil
	case MessageTypeActionOutput:
		message := MessageActionOutput{}
		if err := decoder.Decode(&message); err != nil {
			log.Error("Error decoding MessageTypeActionOutput: %s", err.Error())
			return 0, nil, err
		}
		return messageType, message, nil
	case MessageTypeActionResult:
		message := MessageActionResult{}
		if err := decoder.Decode(&message); err != nil {
			log.Error("Error decoding MessageTypeActionResult: %s", err.Error())
			return 0, nil, err
		}
		return messageType, message, nil
	}
	log.Error("Unknown message type '%d'", messageType)
	return messageType, nil, fmt.Errorf("unknown message type %d", messageType)
}

// ReadData will read len(p) bytes from the connection.
func (c *Connection) ReadData(p []byte) (int, error) {
	return c.w.Read(p)
}

// WriteMessage try to write a message to the given writer.
func (c *Connection) WriteMessage(messageType MessageType, message interface{}) error {
	messageData := []byte{}
	messageLength := uint32(0)
	if message != nil {
		data, err := encodeMessageData(message)
		if err != nil {
			log.Error("Error encoding message: %s", err.Error())
			return err
		}
		messageData = data
		messageLength = uint32(len(data))
	}
	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, ProtocolVersion)
	messageTypeBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(messageTypeBuf, uint32(messageType))
	messageLengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(messageLengthBuf, messageLength)

	log.PDebug("Preparing message", map[string]interface{}{
		"message_type":        messageType,
		"message_data_length": messageLength,
	})

	if _, err := c.w.Write(versionBuf); err != nil {
		log.Error("Error writing version: %s", err.Error())
		return err
	}
	if _, err := c.w.Write(messageTypeBuf); err != nil {
		log.Error("Error writing message type: %s", err.Error())
		return err
	}
	if _, err := c.w.Write(messageLengthBuf); err != nil {
		log.Error("Error writing message length: %s", err.Error())
		return err
	}
	if messageLength > 0 {
		if _, err := c.w.Write(messageData); err != nil {
			log.Error("Error writing message data: %s", err.Error())
			return err
		}
	}
	return nil
}

// WriteData will write raw data to the connection. This data must only be written after a message that is appropriate
// for raw data.
func (c *Connection) WriteData(p []byte) (int, error) {
	return c.w.Write(p)
}

func encodeMessageData(message interface{}) ([]byte, error) {
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

var defaultSSHConfig = ssh.Config{
	KeyExchanges: []string{"curve25519-sha256"},
	Ciphers:      []string{"chacha20-poly1305@openssh.com"},
	MACs:         []string{"hmac-sha2-256-etm@openssh.com"},
}
