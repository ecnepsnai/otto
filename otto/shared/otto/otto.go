/*
Package otto an automation toolkit for Unix-like computers.

This package contains the common interfaces and methods shared by the Otto agent and Otto server.
*/
package otto

import (
	"encoding/gob"
	"time"

	"github.com/ecnepsnai/logtic"
	"golang.org/x/crypto/ssh"
)

var log = logtic.Log.Connect("libotto")

// ProtocolVersion the version of the otto protocol
const ProtocolVersion uint32 = 5

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
	gob.Register(MessageCancelAction{})
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
	MessageTypeReadyForData
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

// MessageCancelAction describes a request to cancel a specific action
type MessageCancelAction struct {
	Name string `json:"name"`
}

// MessageActionOutput describes output from an action
type MessageActionOutput struct {
	IsStdErr bool   `json:"is_stderr"`
	Data     []byte `json:"data"`
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
	StdoutLen uint32        `json:"stdout_len"`
	StderrLen uint32        `json:"stderr_len"`
	Elapsed   time.Duration `json:"elapsed"`
}

func (sr ScriptResult) String() string {
	return logtic.StringFromParameters(map[string]interface{}{
		"success":    sr.Success,
		"exec_error": sr.ExecError,
		"code":       sr.Code,
		"stdout_len": sr.StdoutLen,
		"stderr_len": sr.StderrLen,
		"elapsed":    sr.Elapsed.String(),
	})
}

// ScriptOutput describes output from a script
type ScriptOutput struct {
	StdoutLen uint32 `json:"stdout_len"`
	StderrLen uint32 `json:"stderr_len"`
	Data      []byte `json:"data"`
}

func (o *ScriptOutput) Stdout() string {
	if o.StdoutLen == 0 {
		return ""
	}

	return string(o.Data[0:o.StdoutLen])
}

func (o *ScriptOutput) Stderr() string {
	if o.StderrLen == 0 {
		return ""
	}

	return string(o.Data[o.StdoutLen:])
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

const sshChannelName = "otto"

var defaultSSHConfig = ssh.Config{
	KeyExchanges: []string{"curve25519-sha256"},
	Ciphers:      []string{"chacha20-poly1305@openssh.com"},
	MACs:         []string{"hmac-sha2-256-etm@openssh.com"},
}
