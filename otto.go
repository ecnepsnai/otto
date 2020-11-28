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
const ProtocolVersion = uint32(1)

func init() {
	gob.Register(Request{})
	gob.Register(Reply{})
	gob.Register(Script{})
	gob.Register(ScriptResult{})
	gob.Register(File{})
}

// Request describes an otto request
type Request struct {
	Action uint32
	Script Script
	File   File
}

// Reply describes the reply to an otto request
type Reply struct {
	Error        error
	ScriptResult ScriptResult
	File         File
	Version      string
}

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

const (
	// ActionPing ping action
	ActionPing uint32 = 1
	// ActionRunScript run script action
	ActionRunScript uint32 = 2
	// ActionReloadConfig reload config action
	ActionReloadConfig uint32 = 3
	// ActionUploadFile save the provided file on the remote host
	ActionUploadFile uint32 = 4
	// ActionUploadFileAndExit save the provided file on the remote host and exit. Used to update otto clients.
	ActionUploadFileAndExit uint32 = 5
	// ActionExit exit the otto client
	ActionExit uint32 = 6
	// ActionReboot reboot the host
	ActionReboot uint32 = 7
	// ActionShutdown power down the host
	ActionShutdown uint32 = 8
)

func readEncryptedMessage(r io.Reader, psk string) ([]byte, error) {
	dataLengthBuf := make([]byte, 4)

	if _, err := io.ReadFull(r, dataLengthBuf); err != nil {
		log.Error("Error reading data length: %s", err.Error())
		return nil, err
	}
	dataLength := binary.BigEndian.Uint32(dataLengthBuf)
	log.Debug("Data length: %d", dataLength)

	encryptedData := make([]byte, dataLength)
	readLength, err := io.ReadFull(r, encryptedData)
	if err != nil {
		log.Error("Error reading encrypted data: %s", err.Error())
		return nil, err
	}
	log.Debug("Read length: %#v", readLength)
	log.Debug("Err: %#v", err)
	if dataLength != uint32(readLength) {
		log.Error("Incorrect data length. Reported: %d, actual: %d", dataLength, readLength)
		return nil, fmt.Errorf("bad request length")
	}

	data, err := security.Decrypt(encryptedData, psk)
	if err != nil {
		log.Error("Error decrypting data: %s", err.Error())
		return nil, err
	}

	return data, nil
}

// ReadRequest try to read a request from the given reader
func ReadRequest(r io.Reader, psk string) (*Request, error) {
	versionBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, versionBuf); err != nil {
		log.Error("Error reading version: %s", err.Error())
		return nil, err
	}
	version := binary.BigEndian.Uint32(versionBuf)
	log.Debug("Protocol version: %d", version)
	if version != ProtocolVersion {
		log.Warn("Unsupported protocol version: %d, wanted: %d", version, ProtocolVersion)
	}

	data, err := readEncryptedMessage(r, psk)
	if err != nil {
		log.Error("Error reading encrypted data: %s", err.Error())
		return nil, err
	}

	request := Request{}
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&request); err != nil {
		log.Error("Error decoding data as request: %s", err.Error())
		return nil, err
	}

	return &request, nil
}

// ReadReply try to read a reply from the given reader
func ReadReply(r io.Reader, psk string) (*Reply, error) {
	versionBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, versionBuf); err != nil {
		log.Error("Error reading version: %s", err.Error())
		return nil, err
	}
	version := binary.BigEndian.Uint32(versionBuf)
	log.Debug("Protocol version: %d", version)
	if version != ProtocolVersion {
		log.Warn("Unsupported protocol version: %d, wanted: %d", version, ProtocolVersion)
	}

	data, err := readEncryptedMessage(r, psk)
	if err != nil {
		log.Error("Error reading encrypted data: %s", err.Error())
		return nil, err
	}

	reply := Reply{}
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&reply); err != nil {
		log.Error("Error decoding data as reply: %s", err.Error())
		return nil, err
	}

	return &reply, nil
}

func writeEncryptedMessage(data []byte, psk string, w io.Writer) error {
	encryptedData, err := security.Encrypt(data, psk)
	if err != nil {
		log.Error("Error encrypting data: %s", err.Error())
		return nil
	}

	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, ProtocolVersion)

	dataLength := len(encryptedData)
	lenBuf := make([]byte, 4)
	log.Debug("Encrypted data length: %d", dataLength)
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
	log.Debug("Wrote %d bytes", wrote)
	if err != nil {
		log.Error("Error writing encrypted data: %s", err.Error())
		return err
	}
	return nil
}

// WriteRequest write the request to the writer
func WriteRequest(r Request, psk string, w io.Writer) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(r); err != nil {
		log.Error("Error encoding request: %s", err.Error())
		return nil
	}

	if err := writeEncryptedMessage(buf.Bytes(), psk, w); err != nil {
		log.Error("Error writing encrypted message: %s", err.Error())
		return err
	}

	return nil
}

// WriteReply write the reply to the writer
func WriteReply(r Reply, psk string, w io.Writer) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(r); err != nil {
		log.Error("Error encoding reply: %s", err.Error())
		return nil
	}

	if err := writeEncryptedMessage(buf.Bytes(), psk, w); err != nil {
		log.Error("Error writing encrypted message: %s", err.Error())
		return err
	}

	return nil
}
