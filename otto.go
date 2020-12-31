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
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/security"
)

var log = logtic.Connect("libotto")

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
	MessageTypeHeartbeatRequest  uint32 = 1
	MessageTypeHeartbeatResponse uint32 = 2
	MessageTypeTriggerAction     uint32 = 3
	MessageTypeCancelAction      uint32 = 4
	MessageTypeActionOutput      uint32 = 5
	MessageTypeActionResult      uint32 = 6
	MessageTypeGeneralFailure    uint32 = 7
)

// MessageHeartbeatRequest describes a heartbeat request
type MessageHeartbeatRequest struct {
	ServerVersion string
}

// MessageHeartbeatResponse describes a heartbeat response
type MessageHeartbeatResponse struct {
	ClientVersion string
}

// MessageTriggerAction describes an action trigger
type MessageTriggerAction struct {
	Action uint32
	Script Script
	File   File
}

// MessageCancelAction describes a request to cancel an action
type MessageCancelAction struct{}

// MessageActionOutput describes output from an action
type MessageActionOutput struct {
	Stdout []byte
	Stderr []byte
}

// MessageActionResult describes the result of a triggered action
type MessageActionResult struct {
	ScriptResult  ScriptResult
	Error         error
	File          File
	ClientVersion string
}

// MessageGeneralFailure describes a general failure
type MessageGeneralFailure struct {
	Error error
}

// Actions
const (
	ActionRunScript         uint32 = 1
	ActionReloadConfig      uint32 = 2
	ActionUploadFile        uint32 = 3
	ActionUploadFileAndExit uint32 = 4
	ActionExit              uint32 = 5
	ActionReboot            uint32 = 6
	ActionShutdown          uint32 = 7
)

// Script describes a script
type Script struct {
	Name             string
	UID              uint32
	GID              uint32
	Environment      map[string]string
	WorkingDirectory string
	Executable       string
	Files            []File
	Data             []byte
}

// ScriptResult describes the result of the script
type ScriptResult struct {
	Success   bool
	ExecError string
	Code      int
	Stdout    string
	Stderr    string
	Elapsed   time.Duration
}

// File Describes a file
type File struct {
	Path string
	UID  int
	GID  int
	Mode uint32
	Data []byte
}

// RegisterRequest describes a register request
type RegisterRequest struct {
	Address  string `json:"address"`
	PSK      string `json:"psk"`
	Uname    string `json:"uname"`
	Hostname string `json:"hostname"`
	Port     uint32 `json:"port"`
}

// RegisterResponse describes the response to a register request
type RegisterResponse struct {
	PSK string `json:"psk"`
}

func readEncryptedFrame(r io.Reader, psk string) ([]byte, error) {
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

	encryptedData := make([]byte, dataLength)
	readLength, err := io.ReadFull(r, encryptedData)
	if err != nil {
		log.Error("Error reading encrypted data: %s", err.Error())
		return nil, err
	}
	if dataLength != uint32(readLength) {
		log.Error("Incorrect data length. Reported: %d, actual: %d", dataLength, readLength)
		return nil, fmt.Errorf("bad request length")
	}
	log.Debug("Read frame: encryptedLength=%d version=%d total=%d", dataLength, ProtocolVersion, readLength)

	data, err := security.Decrypt(encryptedData, psk)
	if err != nil {
		log.Error("Error decrypting data: %s", err.Error())
		return nil, err
	}

	return data, nil
}

// ReadMessage try to read a message from the given reader. Returns the message type, the message data, or an error
func ReadMessage(r io.Reader, psk string) (uint32, interface{}, error) {
	data, err := readEncryptedFrame(r, psk)
	if err != nil {
		log.Error("Error reading encrypted data: %s", err.Error())
		return 0, nil, err
	}
	if data == nil {
		return 0, nil, nil
	}

	messageType := binary.BigEndian.Uint32(data[:4])
	log.Debug("Read message: messageType=%d dataLength=%d", messageType, len(data)-4)
	message, err := DecodeMessage(messageType, data[4:])
	if err != nil {
		log.Error("Error decoding message data: %s", err.Error())
		return 0, nil, err
	}

	return messageType, message, nil
}

// WriteMessage try to write a message to the given writer.
func WriteMessage(messageType uint32, message interface{}, w io.Writer, psk string) error {
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

	log.Debug("Preparing message: messageType=%d dataLength=%d messageLength=%d", messageType, messageDataLength, dataLength)
	return writeEncryptedFrame(data, psk, w)
}

func writeEncryptedFrame(data []byte, psk string, w io.Writer) error {
	encryptedData, err := security.Encrypt(data, psk)
	if err != nil {
		log.Error("Error encrypting data: %s", err.Error())
		return nil
	}

	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, ProtocolVersion)

	dataLength := len(encryptedData)
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
	for _, b := range encryptedData {
		replyBuf[i] = b
		i++
	}

	wrote, err := w.Write(replyBuf)
	log.Debug("Wrote frame: encryptedLength=%d version=%d total=%d", dataLength, ProtocolVersion, wrote)
	if wrote != replyLength {
		log.Error("Unable to write all of reply: wrote=%d total=%d", wrote, replyLength)
		return fmt.Errorf("out of space")
	}
	if err != nil {
		log.Error("Error writing encrypted data: %s", err.Error())
		return err
	}
	return nil
}

// DecodeMessage try to decode the given message. The returned object should match the message struct for the message
// type.
func DecodeMessage(messageType uint32, data []byte) (interface{}, error) {
	switch messageType {
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

	return nil, fmt.Errorf("unknown message type")
}

// EncodeMessage try to encode the given message
func EncodeMessage(messageType uint32, message interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := gob.NewEncoder(buf).Encode(message); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
