package otto

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"sync"
)

// ReadMessage try to read a message from the given reader. Returns the message type, the message data, or an error.
// Depending on the message type, there may be additional data to read following the message. It is up to the caller to
// continue reading any additional data.
func (c *Connection) ReadMessage() (MessageType, interface{}, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	headerBuf := make([]byte, 4*3)
	if _, err := io.ReadFull(c.w, headerBuf); err != nil {
		if err == io.EOF {
			// Agent closed - nothing to read
			return 0, nil, err
		}

		log.Error("Error reading header: %s", err.Error())
		return 0, nil, err
	}

	version := binary.BigEndian.Uint32(headerBuf[0:4])
	messageType := MessageType(binary.BigEndian.Uint32(headerBuf[4:8]))
	dataLength := binary.BigEndian.Uint32(headerBuf[8:12])

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

	log.PDebug("Read message", map[string]interface{}{
		"version":      version,
		"message_type": messageType,
		"data_length":  dataLength,
	})

	if dataLength == 0 {
		return messageType, nil, nil
	}

	messageData := make([]byte, dataLength)
	if len, err := io.ReadFull(c.w, messageData); err != nil {
		log.Error("Error reading message data: %s. Need %dB got %dB", err.Error(), dataLength, len)
		return 0, nil, err
	}

	decoder := gob.NewDecoder(bytes.NewBuffer(messageData))

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
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.w.Read(p)
}

// WriteMessage try to write a message to the given writer.
func (c *Connection) WriteMessage(messageType MessageType, message interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	messageData := []byte{}
	messageLength := uint32(0)
	if message != nil {
		data, l, err := encodeMessageData(message)
		if err != nil {
			log.Error("Error encoding message: %s", err.Error())
			return err
		}
		messageData = data
		messageLength = l
	}

	headerBuf := make([]byte, 4*3)
	binary.BigEndian.PutUint32(headerBuf[0:], ProtocolVersion)
	binary.BigEndian.PutUint32(headerBuf[4:], uint32(messageType))
	binary.BigEndian.PutUint32(headerBuf[8:], messageLength)

	log.PDebug("Preparing message", map[string]interface{}{
		"protocol_version":    ProtocolVersion,
		"message_type":        messageType,
		"message_data_length": messageLength,
	})

	if wrote, err := c.w.Write(headerBuf); err != nil {
		log.Error("Error writing header: %s. Wrote: %d", err.Error(), wrote)
		return err
	}
	if messageLength > 0 {
		if wrote, err := c.w.Write(messageData); err != nil {
			log.Error("Error writing message data: %s. Wrote: %d", err.Error(), wrote)
			return err
		}
	}
	return nil
}

// WriteData will write raw data to the connection. This data must only be written after a message that is appropriate
// for raw data. Note, you must call Connection.WriteFinished() when you've written all your data.
func (c *Connection) WriteData(p []byte) (int, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.w.Write(p)
}

// Copy will copy all data fron src reader to the connection
func (c *Connection) Copy(src io.Reader) (int64, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return io.Copy(c.w, src)
}

// WriteFinished will signal that all data has been written. This does not close the connection, but
// the caller's side cannot write any more data.
func (c *Connection) WriteFinished() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.w.CloseWrite()
}

func encodeMessageData(message interface{}) ([]byte, uint32, error) {
	buf := &bytes.Buffer{}

	if err := gob.NewEncoder(buf).Encode(message); err != nil {
		return nil, 0, err
	}

	return buf.Bytes(), uint32(buf.Len()), nil
}

type ReadWriteCloserFinisher interface {
	io.ReadWriteCloser
	CloseWrite() error
}

// Connection describes a connection between the Otto Server and Otto Host
type Connection struct {
	id             int
	w              ReadWriteCloserFinisher
	remoteAddr     net.Addr
	localAddr      net.Addr
	remoteIdentity []byte
	localIdentity  []byte
	mutex          sync.Mutex
}

func MockConnection(w ReadWriteCloserFinisher) *Connection {
	return &Connection{
		w:     w,
		mutex: sync.Mutex{},
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
