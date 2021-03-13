package otto_test

import (
	"bytes"
	"encoding/binary"
	"os"
	"testing"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/secutil"
	fuzz "github.com/google/gofuzz"
)

const psk = "example_psk_please_dont_use_this"

func TestMain(m *testing.M) {
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			logtic.Log.Level = logtic.LevelDebug
			logtic.Open()
		}
	}

	retCode := m.Run()
	os.Exit(retCode)
}

// Perform an end-to-end heartbeat request and reply
func TestHeartbeat(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}

	heartbeatRequest := otto.MessageHeartbeatRequest{ServerVersion: "foo"}
	if err := otto.WriteMessage(otto.MessageTypeHeartbeatRequest, heartbeatRequest, buf, psk); err != nil {
		t.Fatalf("Error writing message: %s", err.Error())
	}

	rMessageType, rHeartbeatRequest, err := otto.ReadMessage(buf, psk)
	if err != nil {
		t.Fatalf("Error reading heartbeat message: %s", err.Error())
	}

	if rMessageType != otto.MessageTypeHeartbeatRequest {
		t.Fatalf("Incorrect message type")
	}
	if _, valid := rHeartbeatRequest.(otto.MessageHeartbeatRequest); !valid {
		t.Fatalf("Incorrect message data type")
	}
}

// Test that otto behaves in an expected mannor when it receives a message that has a valid protocol version number but the rest
// is fuzzed data
func TestUnsupportedProtocolVersion(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	versionBuf := make([]byte, 4)
	binary.PutVarint(versionBuf, -1)
	buf.Write(versionBuf)
	buf.Write(secutil.RandomBytes(128))
	_, _, err := otto.ReadMessage(buf, psk)
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when a unknown message type is seen in in a message
func TestMalformedMessageType(t *testing.T) {
	t.Parallel()

	f := fuzz.New()

	buf := &bytes.Buffer{}
	var messageType uint32
	f.Fuzz(&messageType)

	heartbeatRequest := otto.MessageHeartbeatRequest{ServerVersion: "foo"}
	if err := otto.WriteMessage(messageType, heartbeatRequest, buf, psk); err != nil {
		t.Fatalf("Error writing message: %s", err.Error())
	}

	_, _, err := otto.ReadMessage(buf, psk)
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message type")
	}
}

// Test that otto behaves in an expected mannor when fuzzed data in seen for a specific message type
func TestMalformedMessageData(t *testing.T) {
	t.Parallel()

	f := fuzz.New()

	buf := &bytes.Buffer{}
	var messageData map[string]string
	f.Fuzz(&messageData)

	if err := otto.WriteMessage(otto.MessageTypeHeartbeatRequest, messageData, buf, psk); err != nil {
		t.Fatalf("Error writing message: %s", err.Error())
	}

	_, _, err := otto.ReadMessage(buf, psk)
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when fuzzed data is passed to the reader
func TestMalformedMessageEntireMessage(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	buf.Write(secutil.RandomBytes(128))
	_, _, err := otto.ReadMessage(buf, psk)
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when it receives a message that has a valid protocol version number but the rest
// is fuzzed data
func TestMalformedMessageWithProtocolVersion(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, otto.ProtocolVersion)
	buf.Write(versionBuf)
	buf.Write(secutil.RandomBytes(128))
	_, _, err := otto.ReadMessage(buf, psk)
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when it receives a message that has a valid protocol version number,
// and a valid data length but the rest is fuzzed data
func TestMalformedMessageWithProtocolAndLength(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, otto.ProtocolVersion)
	buf.Write(versionBuf)
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, 128)
	buf.Write(lengthBuf)
	buf.Write(secutil.RandomBytes(128))
	_, _, err := otto.ReadMessage(buf, psk)
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when it receives a message that has a valid protocol version number,
// valid data length, and properly encrypted data but that data is fuzzed.
func TestMalrforedMessageWithEncryptedData(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, otto.ProtocolVersion)
	buf.Write(versionBuf)
	lengthBuf := make([]byte, 4)
	encryptedData, err := secutil.Encrypt(secutil.RandomBytes(128), psk)
	if err != nil {
		panic(err)
	}
	binary.BigEndian.PutUint32(lengthBuf, uint32(len(encryptedData)))
	buf.Write(lengthBuf)
	buf.Write(encryptedData)
	_, _, err = otto.ReadMessage(buf, psk)
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when it receives a single-byte message
func TestSingleByteMessage(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	buf.Write([]byte{'1'})
	_, _, err := otto.ReadMessage(buf, psk)
	if err == nil {
		t.Fatalf("No error seen when one expected for single byte message")
	}
}

// Test that otto behaves in an expected mannor when the reported encrypted data length does not match the actual
// length of the encrypted data both underflow and overflow.
func TestIncorrectDataLength(t *testing.T) {
	t.Parallel()

	// Reported length is shorter than actual
	buf := &bytes.Buffer{}
	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, otto.ProtocolVersion)
	buf.Write(versionBuf)
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, 32)
	buf.Write(lengthBuf)
	buf.Write(secutil.RandomBytes(128))
	if _, _, err := otto.ReadMessage(buf, psk); err == nil {
		t.Fatalf("No error seen when one expected for false message length")
	}

	// Reported length is larger than actual
	buf = &bytes.Buffer{}
	versionBuf = make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, otto.ProtocolVersion)
	buf.Write(versionBuf)
	lengthBuf = make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, 128)
	buf.Write(lengthBuf)
	buf.Write(secutil.RandomBytes(6))
	if _, _, err := otto.ReadMessage(buf, psk); err == nil {
		t.Fatalf("No error seen when one expected for false message length")
	}
}
