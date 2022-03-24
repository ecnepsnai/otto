package otto

import (
	"fmt"
	"io"
)

// SendHeartbeat will send a heartbeat request to the host, returning a reply or an error
func (conn *Connection) SendHeartbeat(request MessageHeartbeatRequest) (*MessageHeartbeatResponse, error) {
	if err := conn.WriteMessage(MessageTypeHeartbeatRequest, request); err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	messageType, data, err := conn.ReadMessage()
	if err != nil {
		log.PError("Error reading message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	if messageType != MessageTypeHeartbeatResponse {
		err = fmt.Errorf("incorrect message type %d", messageType)
		log.PError("Error reading message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	response, ok := data.(MessageHeartbeatResponse)
	if !ok {
		err = fmt.Errorf("incorrect response data type")
		log.PError("Error reading message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return &response, nil
}

// RotateIdentity will send a request to rotate the identity on the host, returning a reply or an error
func (conn *Connection) RotateIdentity(request MessageRotateIdentityRequest) (*MessageRotateIdentityResponse, error) {
	if err := conn.WriteMessage(MessageTypeRotateIdentityRequest, request); err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	messageType, data, err := conn.ReadMessage()
	if err != nil {
		log.PError("Error reading message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}
	if messageType != MessageTypeRotateIdentityResponse {
		err = fmt.Errorf("incorrect message type %d", messageType)
		log.PError("Error reading message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	response, ok := data.(MessageRotateIdentityResponse)
	if !ok {
		err = fmt.Errorf("incorrect response data type")
		log.PError("Error reading message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	return &response, nil
}

// TriggerAction will trigger an action on the host. Returns a result or an error.
// actionOutput is called whenever there is output on the action and may be called multiple times.
// Send a message to cancel to abort the action.
func (conn *Connection) TriggerAction(action MessageTriggerAction, actionOutput func(stdout, stderr []byte), cancel chan bool) (*MessageActionResult, error) {
	if err := conn.WriteMessage(MessageTypeTriggerAction, action); err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	go func() {
		for {
			<-cancel
			conn.WriteMessage(MessageTypeCancelAction, MessageCancelAction{})
		}
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err == io.EOF || messageType == 0 {
			return nil, nil
		} else if err != nil {
			log.PError("Error reading message", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}

		switch messageType {
		case MessageTypeActionOutput:
			output := message.(MessageActionOutput)
			if actionOutput != nil {
				actionOutput(output.Stdout, output.Stderr)
			}
		case MessageTypeActionResult:
			result := message.(MessageActionResult)
			return &result, nil
		case MessageTypeGeneralFailure:
			result := message.(MessageGeneralFailure)
			log.PError("General error triggering action on host", map[string]interface{}{
				"error": result.Error,
			})
			return nil, fmt.Errorf("%s", result.Error)
		}
	}
}
