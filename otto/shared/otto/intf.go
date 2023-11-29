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

func (conn *Connection) TriggerActionRunScript(script ScriptInfo, scriptReader io.ReadCloser, actionOutput func(stdout, stderr []byte)) (*MessageActionResult, *ScriptOutput, error) {
	if err := conn.WriteMessage(MessageTypeTriggerActionRunScript, script); err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil, err
	}

	if messageType, _, _ := conn.ReadMessage(); messageType != MessageTypeReadyForData {
		log.Error("Unexpected message from agent when waiting for MessageTypeReadyForData: %d", messageType)
		return nil, nil, fmt.Errorf("error writing script data: unexpected message from agent %d", messageType)
	}

	wrote, err := io.Copy(conn.w, scriptReader)
	if err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, nil, err
	}
	log.PDebug("Wrote script data", map[string]interface{}{
		"script_length": wrote,
	})
	scriptReader.Close()

	for {
		messageType, message, err := conn.ReadMessage()
		if err == io.EOF || messageType == 0 {
			return nil, nil, nil
		} else if err != nil {
			log.PError("Error reading message", map[string]interface{}{
				"error": err.Error(),
			})
			return nil, nil, err
		}

		switch messageType {
		case MessageTypeActionOutput:
			output := message.(MessageActionOutput)
			log.Debug("Recieved %dB of output from script", len(output.Data))
			if actionOutput != nil {
				if output.IsStdErr {
					actionOutput([]byte{}, output.Data)
				} else {
					actionOutput(output.Data, []byte{})
				}
			}
		case MessageTypeActionResult:
			result := message.(MessageActionResult)
			outputLen := result.ScriptResult.StdoutLen + result.ScriptResult.StderrLen
			var output ScriptOutput
			if outputLen > 0 {
				outputData := make([]byte, outputLen)
				if err := conn.WriteMessage(MessageTypeReadyForData, nil); err != nil {
					log.Error("Error sending message: %s", err.Error())
					return nil, nil, err
				}
				if len, err := conn.ReadData(outputData); err != nil {
					log.PError("Error reading output from script", map[string]interface{}{
						"error":      err.Error(),
						"output_len": outputLen,
						"read_len":   len,
					})
					return &result, nil, nil
				}
				output = ScriptOutput{
					StdoutLen: result.ScriptResult.StdoutLen,
					StderrLen: result.ScriptResult.StderrLen,
					Data:      outputData,
				}
			} else {
				output = ScriptOutput{
					StdoutLen: 0,
					StderrLen: 0,
				}
			}
			return &result, &output, nil
		case MessageTypeGeneralFailure:
			result := message.(MessageGeneralFailure)
			log.PError("General error triggering action on host", map[string]interface{}{
				"error": result.Error,
			})
			return nil, nil, fmt.Errorf("%s", result.Error)
		}
	}
}

func (conn *Connection) TriggerActionReloadConfig() error {
	if err := conn.WriteMessage(MessageTypeTriggerActionReloadConfig, nil); err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return nil
}

func (conn *Connection) TriggerActionUploadFile(file FileInfo, fileReader io.Reader) error {
	if err := conn.WriteMessage(MessageTypeTriggerActionUploadFile, file); err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	if _, err := io.Copy(conn.w, fileReader); err != nil {
		log.PError("Error writing file data", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	if err := conn.WriteFinished(); err != nil {
		log.PError("Error writing file data", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return nil
}

func (conn *Connection) TriggerActionExitAgent() error {
	if err := conn.WriteMessage(MessageTypeTriggerActionExitAgent, nil); err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return nil
}

func (conn *Connection) TriggerActionReboot() error {
	if err := conn.WriteMessage(MessageTypeTriggerActionReboot, nil); err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return nil
}

func (conn *Connection) TriggerActionShutdown() error {
	if err := conn.WriteMessage(MessageTypeTriggerActionShutdown, nil); err != nil {
		log.PError("Error writing message", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	return nil
}

func (conn *Connection) ActionOutput() {}

func (conn *Connection) ActionResult() {}
