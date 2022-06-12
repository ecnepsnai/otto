package otto_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/otto"
	"github.com/ecnepsnai/secutil"
	fuzz "github.com/google/gofuzz"
)

func MockOttoConnection(buf *bytes.Buffer) *otto.Connection {
	return otto.MockConnection(mockConnectionType{buf})
}

type mockConnectionType struct {
	buf *bytes.Buffer
}

func (m mockConnectionType) Read(p []byte) (n int, err error) {
	return m.buf.Read(p)
}

func (m mockConnectionType) Write(p []byte) (n int, err error) {
	return m.buf.Write(p)
}

func (m mockConnectionType) Close() error {
	return nil
}

func TestMain(m *testing.M) {
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			logtic.Log.Level = logtic.LevelDebug
			logtic.Log.Open()
		}
	}

	retCode := m.Run()
	os.Exit(retCode)
}

// Perform an end-to-end heartbeat request and reply
func TestHeartbeat(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	conn := MockOttoConnection(buf)

	heartbeatRequest := otto.MessageHeartbeatRequest{Version: "foo"}
	if err := conn.WriteMessage(otto.MessageTypeHeartbeatRequest, heartbeatRequest); err != nil {
		t.Fatalf("Error writing message: %s", err.Error())
	}

	rMessageType, rHeartbeatRequest, err := conn.ReadMessage()
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
	conn := MockOttoConnection(buf)

	versionBuf := make([]byte, 4)
	binary.PutVarint(versionBuf, -1)
	buf.Write(versionBuf)
	buf.Write(secutil.RandomBytes(128))
	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when a unknown message type is seen in in a message
func TestMalformedMessageType(t *testing.T) {
	t.Parallel()

	f := fuzz.New()

	buf := &bytes.Buffer{}
	conn := MockOttoConnection(buf)

	var messageType uint32
	f.Fuzz(&messageType)

	heartbeatRequest := otto.MessageHeartbeatRequest{Version: "foo"}
	if err := conn.WriteMessage(messageType, heartbeatRequest); err != nil {
		t.Fatalf("Error writing message: %s", err.Error())
	}

	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message type")
	}
}

// Test that otto behaves in an expected mannor when fuzzed data in seen for a specific message type
func TestMalformedMessageData(t *testing.T) {
	t.Parallel()

	f := fuzz.New()

	buf := &bytes.Buffer{}
	conn := MockOttoConnection(buf)

	var messageData map[string]string
	f.Fuzz(&messageData)

	if err := conn.WriteMessage(otto.MessageTypeHeartbeatRequest, messageData); err != nil {
		t.Fatalf("Error writing message: %s", err.Error())
	}

	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when fuzzed data is passed to the reader
func TestMalformedMessageEntireMessage(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	conn := MockOttoConnection(buf)

	buf.Write(secutil.RandomBytes(128))
	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when it receives a message that has a valid protocol version number but the rest
// is fuzzed data
func TestMalformedMessageWithProtocolVersion(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	conn := MockOttoConnection(buf)

	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, otto.ProtocolVersion)
	buf.Write(versionBuf)
	buf.Write(secutil.RandomBytes(128))
	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when it receives a message that has a valid protocol version number,
// and a valid data length but the rest is fuzzed data
func TestMalformedMessageWithProtocolAndLength(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	conn := MockOttoConnection(buf)

	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, otto.ProtocolVersion)
	buf.Write(versionBuf)
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, 128)
	buf.Write(lengthBuf)
	buf.Write(secutil.RandomBytes(128))
	_, _, err := conn.ReadMessage()
	if err == nil {
		t.Fatalf("No error seen when one expected for fuzzed message data")
	}
}

// Test that otto behaves in an expected mannor when it receives a single-byte message
func TestSingleByteMessage(t *testing.T) {
	t.Parallel()

	buf := &bytes.Buffer{}
	conn := MockOttoConnection(buf)

	buf.Write([]byte{'1'})
	_, _, err := conn.ReadMessage()
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
	conn := MockOttoConnection(buf)

	versionBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, otto.ProtocolVersion)
	buf.Write(versionBuf)
	lengthBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, 32)
	buf.Write(lengthBuf)
	buf.Write(secutil.RandomBytes(128))
	if _, _, err := conn.ReadMessage(); err == nil {
		t.Fatalf("No error seen when one expected for false message length")
	}

	// Reported length is larger than actual
	buf = &bytes.Buffer{}
	conn = MockOttoConnection(buf)

	versionBuf = make([]byte, 4)
	binary.BigEndian.PutUint32(versionBuf, otto.ProtocolVersion)
	buf.Write(versionBuf)
	lengthBuf = make([]byte, 4)
	binary.BigEndian.PutUint32(lengthBuf, 128)
	buf.Write(lengthBuf)
	buf.Write(secutil.RandomBytes(6))
	if _, _, err := conn.ReadMessage(); err == nil {
		t.Fatalf("No error seen when one expected for false message length")
	}
}

func TestConnection(t *testing.T) {
	listenerIdentity, err := otto.NewIdentity()
	if err != nil {
		panic(err)
	}
	dialerIdentity, err := otto.NewIdentity()
	if err != nil {
		panic(err)
	}

	l, err := otto.SetupListener(otto.ListenOptions{
		Address:           "127.0.0.1:0",
		Identity:          listenerIdentity.Signer(),
		TrustedPublicKeys: []string{dialerIdentity.PublicKeyString()},
	}, func(c *otto.Connection) {
		messageType, _, err := c.ReadMessage()
		if err != nil {
			t.Errorf("Error reading message: %s", err.Error())
		}
		if messageType != otto.MessageTypeHeartbeatRequest {
			t.Errorf("Unexpected message type")
		}
		c.WriteMessage(otto.MessageTypeHeartbeatResponse, otto.MessageHeartbeatResponse{})
	})
	if err != nil {
		panic(err)
	}
	port := l.Port()
	go l.Accept()
	time.Sleep(5 * time.Millisecond)

	c, err := otto.Dial(otto.DialOptions{
		Network:          "tcp",
		Address:          fmt.Sprintf("127.0.0.1:%d", port),
		Identity:         dialerIdentity.Signer(),
		TrustedPublicKey: listenerIdentity.PublicKeyString(),
	})
	if err != nil {
		t.Fatalf("Error dialing: %s", err.Error())
	}

	c.WriteMessage(otto.MessageTypeHeartbeatRequest, otto.MessageHeartbeatRequest{})
	messageType, _, err := c.ReadMessage()
	if err != nil {
		t.Errorf("Error reading message: %s", err.Error())
	}
	if messageType != otto.MessageTypeHeartbeatResponse {
		t.Errorf("Unexpected message type")
	}
	c.Close()
}
