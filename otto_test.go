package otto_test

import (
	"bytes"
	"testing"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/otto"
)

func TestHeartbeat(t *testing.T) {
	logtic.Log.Level = logtic.LevelDebug
	logtic.Open()

	psk := "farts"
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
