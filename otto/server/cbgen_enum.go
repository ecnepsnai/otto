package server

// This file is was generated automatically by Codegen v1.12.3
// Do not make changes to this file as they will be lost

const (
	// Ping the host
	AgentActionPing = "ping"
	// Run the script on the host
	AgentActionRunScript = "run_script"
	// Reload the configuration of the agent
	AgentActionReloadConfig = "reload_config"
	// Exit the agent on the host
	AgentActionExitAgent = "exit_agent"
	// Reboot the host
	AgentActionReboot = "reboot"
	// Power off the host
	AgentActionShutdown = "shutdown"
)

// AllAgentAction all AgentAction values
var AllAgentAction = []string{
	AgentActionPing,
	AgentActionRunScript,
	AgentActionReloadConfig,
	AgentActionExitAgent,
	AgentActionReboot,
	AgentActionShutdown,
}

// AgentActionMap map AgentAction keys to values
var AgentActionMap = map[string]string{
	AgentActionPing:         "ping",
	AgentActionRunScript:    "run_script",
	AgentActionReloadConfig: "reload_config",
	AgentActionExitAgent:    "exit_agent",
	AgentActionReboot:       "reboot",
	AgentActionShutdown:     "shutdown",
}

// IsAgentAction is the provided value a valid AgentAction
func IsAgentAction(q string) bool {
	_, k := AgentActionMap[q]
	return k
}

// ForEachAgentAction call m for each AgentAction
func ForEachAgentAction(m func(value string)) {
	for _, v := range AllAgentAction {
		m(v)
	}
}

const (
	// UserLoggedIn event
	EventTypeUserLoggedIn = "UserLoggedIn"
	// UserIncorrectPassword event
	EventTypeUserIncorrectPassword = "UserIncorrectPassword"
	// UserLoggedOut event
	EventTypeUserLoggedOut = "UserLoggedOut"
	// UserAdded event
	EventTypeUserAdded = "UserAdded"
	// UserModified event
	EventTypeUserModified = "UserModified"
	// UserResetPassword event
	EventTypeUserResetPassword = "UserResetPassword"
	// UserResetAPIKey event
	EventTypeUserResetAPIKey = "UserResetAPIKey"
	// UserDeleted event
	EventTypeUserDeleted = "UserDeleted"
	// UserPermissionDenied event
	EventTypeUserPermissionDenied = "UserPermissionDenied"
	// HostAdded event
	EventTypeHostAdded = "HostAdded"
	// HostModified event
	EventTypeHostModified = "HostModified"
	// HostDeleted event
	EventTypeHostDeleted = "HostDeleted"
	// HostRegisterSuccess event
	EventTypeHostRegisterSuccess = "HostRegisterSuccess"
	// HostRegisterIncorrectKey event
	EventTypeHostRegisterIncorrectKey = "HostRegisterIncorrectKey"
	// HostTrustModified event
	EventTypeHostTrustModified = "HostTrustModified"
	// HostIdentityRotated event
	EventTypeHostIdentityRotated = "HostIdentityRotated"
	// HostBecameReachable event
	EventTypeHostBecameReachable = "HostBecameReachable"
	// HostBecameUnreachable event
	EventTypeHostBecameUnreachable = "HostBecameUnreachable"
	// GroupAdded event
	EventTypeGroupAdded = "GroupAdded"
	// GroupModified event
	EventTypeGroupModified = "GroupModified"
	// GroupDeleted event
	EventTypeGroupDeleted = "GroupDeleted"
	// ScheduleAdded event
	EventTypeScheduleAdded = "ScheduleAdded"
	// ScheduleModified event
	EventTypeScheduleModified = "ScheduleModified"
	// ScheduleDeleted event
	EventTypeScheduleDeleted = "ScheduleDeleted"
	// AttachmentAdded event
	EventTypeAttachmentAdded = "AttachmentAdded"
	// AttachmentModified event
	EventTypeAttachmentModified = "AttachmentModified"
	// AttachmentDeleted event
	EventTypeAttachmentDeleted = "AttachmentDeleted"
	// ScriptAdded event
	EventTypeScriptAdded = "ScriptAdded"
	// ScriptModified event
	EventTypeScriptModified = "ScriptModified"
	// ScriptDeleted event
	EventTypeScriptDeleted = "ScriptDeleted"
	// ScriptRun event
	EventTypeScriptRun = "ScriptRun"
	// ServerStarted event
	EventTypeServerStarted = "ServerStarted"
	// ServerOptionsModified event
	EventTypeServerOptionsModified = "ServerOptionsModified"
	// RegisterRuleAdded event
	EventTypeRegisterRuleAdded = "RegisterRuleAdded"
	// RegisterRuleModified event
	EventTypeRegisterRuleModified = "RegisterRuleModified"
	// RegisterRuleDeleted event
	EventTypeRegisterRuleDeleted = "RegisterRuleDeleted"
)

// AllEventType all EventType values
var AllEventType = []string{
	EventTypeUserLoggedIn,
	EventTypeUserIncorrectPassword,
	EventTypeUserLoggedOut,
	EventTypeUserAdded,
	EventTypeUserModified,
	EventTypeUserResetPassword,
	EventTypeUserResetAPIKey,
	EventTypeUserDeleted,
	EventTypeUserPermissionDenied,
	EventTypeHostAdded,
	EventTypeHostModified,
	EventTypeHostDeleted,
	EventTypeHostRegisterSuccess,
	EventTypeHostRegisterIncorrectKey,
	EventTypeHostTrustModified,
	EventTypeHostIdentityRotated,
	EventTypeHostBecameReachable,
	EventTypeHostBecameUnreachable,
	EventTypeGroupAdded,
	EventTypeGroupModified,
	EventTypeGroupDeleted,
	EventTypeScheduleAdded,
	EventTypeScheduleModified,
	EventTypeScheduleDeleted,
	EventTypeAttachmentAdded,
	EventTypeAttachmentModified,
	EventTypeAttachmentDeleted,
	EventTypeScriptAdded,
	EventTypeScriptModified,
	EventTypeScriptDeleted,
	EventTypeScriptRun,
	EventTypeServerStarted,
	EventTypeServerOptionsModified,
	EventTypeRegisterRuleAdded,
	EventTypeRegisterRuleModified,
	EventTypeRegisterRuleDeleted,
}

// EventTypeMap map EventType keys to values
var EventTypeMap = map[string]string{
	EventTypeUserLoggedIn:             "UserLoggedIn",
	EventTypeUserIncorrectPassword:    "UserIncorrectPassword",
	EventTypeUserLoggedOut:            "UserLoggedOut",
	EventTypeUserAdded:                "UserAdded",
	EventTypeUserModified:             "UserModified",
	EventTypeUserResetPassword:        "UserResetPassword",
	EventTypeUserResetAPIKey:          "UserResetAPIKey",
	EventTypeUserDeleted:              "UserDeleted",
	EventTypeUserPermissionDenied:     "UserPermissionDenied",
	EventTypeHostAdded:                "HostAdded",
	EventTypeHostModified:             "HostModified",
	EventTypeHostDeleted:              "HostDeleted",
	EventTypeHostRegisterSuccess:      "HostRegisterSuccess",
	EventTypeHostRegisterIncorrectKey: "HostRegisterIncorrectKey",
	EventTypeHostTrustModified:        "HostTrustModified",
	EventTypeHostIdentityRotated:      "HostIdentityRotated",
	EventTypeHostBecameReachable:      "HostBecameReachable",
	EventTypeHostBecameUnreachable:    "HostBecameUnreachable",
	EventTypeGroupAdded:               "GroupAdded",
	EventTypeGroupModified:            "GroupModified",
	EventTypeGroupDeleted:             "GroupDeleted",
	EventTypeScheduleAdded:            "ScheduleAdded",
	EventTypeScheduleModified:         "ScheduleModified",
	EventTypeScheduleDeleted:          "ScheduleDeleted",
	EventTypeAttachmentAdded:          "AttachmentAdded",
	EventTypeAttachmentModified:       "AttachmentModified",
	EventTypeAttachmentDeleted:        "AttachmentDeleted",
	EventTypeScriptAdded:              "ScriptAdded",
	EventTypeScriptModified:           "ScriptModified",
	EventTypeScriptDeleted:            "ScriptDeleted",
	EventTypeScriptRun:                "ScriptRun",
	EventTypeServerStarted:            "ServerStarted",
	EventTypeServerOptionsModified:    "ServerOptionsModified",
	EventTypeRegisterRuleAdded:        "RegisterRuleAdded",
	EventTypeRegisterRuleModified:     "RegisterRuleModified",
	EventTypeRegisterRuleDeleted:      "RegisterRuleDeleted",
}

// IsEventType is the provided value a valid EventType
func IsEventType(q string) bool {
	_, k := EventTypeMap[q]
	return k
}

// ForEachEventType call m for each EventType
func ForEachEventType(m func(value string)) {
	for _, v := range AllEventType {
		m(v)
	}
}

const (
	// IPv4 or IPv6 as chosen by the system automatically
	IPVersionOptionAuto = "auto"
	// IPv4 only
	IPVersionOptionIPv4 = "ipv4"
	// IPv6 only
	IPVersionOptionIPv6 = "ipv6"
)

// AllIPVersionOption all IPVersionOption values
var AllIPVersionOption = []string{
	IPVersionOptionAuto,
	IPVersionOptionIPv4,
	IPVersionOptionIPv6,
}

// IPVersionOptionMap map IPVersionOption keys to values
var IPVersionOptionMap = map[string]string{
	IPVersionOptionAuto: "auto",
	IPVersionOptionIPv4: "ipv4",
	IPVersionOptionIPv6: "ipv6",
}

// IsIPVersionOption is the provided value a valid IPVersionOption
func IsIPVersionOption(q string) bool {
	_, k := IPVersionOptionMap[q]
	return k
}

// ForEachIPVersionOption call m for each IPVersionOption
func ForEachIPVersionOption(m func(value string)) {
	for _, v := range AllIPVersionOption {
		m(v)
	}
}

const (
	// Hostname
	RegisterRulePropertyHostname = "hostname"
	// Kernel Name
	RegisterRulePropertyKernelName = "kernel_name"
	// Kernel Version
	RegisterRulePropertyKernelVersion = "kernel_version"
	// Distribution Name
	RegisterRulePropertyDistributionName = "distribution_name"
	// Distribution Version
	RegisterRulePropertyDistributionVersion = "distribution_version"
)

// AllRegisterRuleProperty all RegisterRuleProperty values
var AllRegisterRuleProperty = []string{
	RegisterRulePropertyHostname,
	RegisterRulePropertyKernelName,
	RegisterRulePropertyKernelVersion,
	RegisterRulePropertyDistributionName,
	RegisterRulePropertyDistributionVersion,
}

// RegisterRulePropertyMap map RegisterRuleProperty keys to values
var RegisterRulePropertyMap = map[string]string{
	RegisterRulePropertyHostname:            "hostname",
	RegisterRulePropertyKernelName:          "kernel_name",
	RegisterRulePropertyKernelVersion:       "kernel_version",
	RegisterRulePropertyDistributionName:    "distribution_name",
	RegisterRulePropertyDistributionVersion: "distribution_version",
}

// IsRegisterRuleProperty is the provided value a valid RegisterRuleProperty
func IsRegisterRuleProperty(q string) bool {
	_, k := RegisterRulePropertyMap[q]
	return k
}

// ForEachRegisterRuleProperty call m for each RegisterRuleProperty
func ForEachRegisterRuleProperty(m func(value string)) {
	for _, v := range AllRegisterRuleProperty {
		m(v)
	}
}

const (
	RequestResponseCodeOutput    = 100
	RequestResponseCodeKeepalive = 101
	RequestResponseCodeError     = 400
	RequestResponseCodeFinished  = 200
)

// AllRequestResponseCode all RequestResponseCode values
var AllRequestResponseCode = []int{
	RequestResponseCodeOutput,
	RequestResponseCodeKeepalive,
	RequestResponseCodeError,
	RequestResponseCodeFinished,
}

// RequestResponseCodeMap map RequestResponseCode keys to values
var RequestResponseCodeMap = map[int]int{
	RequestResponseCodeOutput:    100,
	RequestResponseCodeKeepalive: 101,
	RequestResponseCodeError:     400,
	RequestResponseCodeFinished:  200,
}

// IsRequestResponseCode is the provided value a valid RequestResponseCode
func IsRequestResponseCode(q int) bool {
	_, k := RequestResponseCodeMap[q]
	return k
}

// ForEachRequestResponseCode call m for each RequestResponseCode
func ForEachRequestResponseCode(m func(value int)) {
	for _, v := range AllRequestResponseCode {
		m(v)
	}
}

const (
	// All hosts executed the script successfully
	ScheduleResultSuccess = 0
	// Some hosts did not execute the script successfully
	ScheduleResultPartialSuccess = 1
	// No hosts executed the script successfully
	ScheduleResultFail = 2
)

// AllScheduleResult all ScheduleResult values
var AllScheduleResult = []int{
	ScheduleResultSuccess,
	ScheduleResultPartialSuccess,
	ScheduleResultFail,
}

// ScheduleResultMap map ScheduleResult keys to values
var ScheduleResultMap = map[int]int{
	ScheduleResultSuccess:        0,
	ScheduleResultPartialSuccess: 1,
	ScheduleResultFail:           2,
}

// IsScheduleResult is the provided value a valid ScheduleResult
func IsScheduleResult(q int) bool {
	_, k := ScheduleResultMap[q]
	return k
}

// ForEachScheduleResult call m for each ScheduleResult
func ForEachScheduleResult(m func(value int)) {
	for _, v := range AllScheduleResult {
		m(v)
	}
}

// Permission level for users to run scripts
const (
	// No scripts can be executed
	ScriptRunLevelNone = 0
	// Only scripts mark as read only can be executed
	ScriptRunLevelReadOnly = 1
	// All scripts can be executed
	ScriptRunLevelReadWrite = 2
)

// AllScriptRunLevel all ScriptRunLevel values
var AllScriptRunLevel = []int{
	ScriptRunLevelNone,
	ScriptRunLevelReadOnly,
	ScriptRunLevelReadWrite,
}

// ScriptRunLevelMap map ScriptRunLevel keys to values
var ScriptRunLevelMap = map[int]int{
	ScriptRunLevelNone:      0,
	ScriptRunLevelReadOnly:  1,
	ScriptRunLevelReadWrite: 2,
}

// IsScriptRunLevel is the provided value a valid ScriptRunLevel
func IsScriptRunLevel(q int) bool {
	_, k := ScriptRunLevelMap[q]
	return k
}

// ForEachScriptRunLevel call m for each ScriptRunLevel
func ForEachScriptRunLevel(m func(value int)) {
	for _, v := range AllScriptRunLevel {
		m(v)
	}
}
