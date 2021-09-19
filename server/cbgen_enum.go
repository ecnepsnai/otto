package server

// This file is was generated automatically by Codegen v1.8.0
// Do not make changes to this file as they will be lost

const (

	// ClientActionPing Ping the host
	ClientActionPing = "ping"

	// ClientActionRunScript Run the script on the host
	ClientActionRunScript = "run_script"

	// ClientActionExitClient Exit the client on the host
	ClientActionExitClient = "exit_client"

	// ClientActionReboot Reboot the host
	ClientActionReboot = "reboot"

	// ClientActionShutdown Power off the host
	ClientActionShutdown = "shutdown"

	// ClientActionUpdatePSK Update the client PSK
	ClientActionUpdatePSK = "update_psk"
)

// AllClientAction all ClientAction values
var AllClientAction = []string{

	ClientActionPing,

	ClientActionRunScript,

	ClientActionExitClient,

	ClientActionReboot,

	ClientActionShutdown,

	ClientActionUpdatePSK,
}

// ClientActionMap map ClientAction keys to values
var ClientActionMap = map[string]string{

	ClientActionPing: "ping",

	ClientActionRunScript: "run_script",

	ClientActionExitClient: "exit_client",

	ClientActionReboot: "reboot",

	ClientActionShutdown: "shutdown",

	ClientActionUpdatePSK: "update_psk",
}

// ClientActionNameMap map ClientAction keys to values
var ClientActionNameMap = map[string]string{

	"Ping": "ping",

	"RunScript": "run_script",

	"ExitClient": "exit_client",

	"Reboot": "reboot",

	"Shutdown": "shutdown",

	"UpdatePSK": "update_psk",
}

// IsClientAction is the provided value a valid ClientAction
func IsClientAction(q string) bool {
	_, k := ClientActionMap[q]
	return k
}

// ClientActionSchema the ClientAction schema.
var ClientActionSchema = []map[string]interface{}{

	{
		"name":        "Ping",
		"description": "Ping the host",
		"value":       "ping",
	},

	{
		"name":        "RunScript",
		"description": "Run the script on the host",
		"value":       "run_script",
	},

	{
		"name":        "ExitClient",
		"description": "Exit the client on the host",
		"value":       "exit_client",
	},

	{
		"name":        "Reboot",
		"description": "Reboot the host",
		"value":       "reboot",
	},

	{
		"name":        "Shutdown",
		"description": "Power off the host",
		"value":       "shutdown",
	},

	{
		"name":        "UpdatePSK",
		"description": "Update the client PSK",
		"value":       "update_psk",
	},
}

const (

	// EventTypeUserLoggedIn UserLoggedIn event
	EventTypeUserLoggedIn = "UserLoggedIn"

	// EventTypeUserIncorrectPassword UserIncorrectPassword event
	EventTypeUserIncorrectPassword = "UserIncorrectPassword"

	// EventTypeUserLoggedOut UserLoggedOut event
	EventTypeUserLoggedOut = "UserLoggedOut"

	// EventTypeUserAdded UserAdded event
	EventTypeUserAdded = "UserAdded"

	// EventTypeUserModified UserModified event
	EventTypeUserModified = "UserModified"

	// EventTypeUserResetPassword UserResetPassword event
	EventTypeUserResetPassword = "UserResetPassword"

	// EventTypeUserResetAPIKey UserResetAPIKey event
	EventTypeUserResetAPIKey = "UserResetAPIKey"

	// EventTypeUserDeleted UserDeleted event
	EventTypeUserDeleted = "UserDeleted"

	// EventTypeHostAdded HostAdded event
	EventTypeHostAdded = "HostAdded"

	// EventTypeHostModified HostModified event
	EventTypeHostModified = "HostModified"

	// EventTypeHostDeleted HostDeleted event
	EventTypeHostDeleted = "HostDeleted"

	// EventTypeHostRegisterSuccess HostRegisterSuccess event
	EventTypeHostRegisterSuccess = "HostRegisterSuccess"

	// EventTypeHostRegisterIncorrectKey HostRegisterIncorrectKey event
	EventTypeHostRegisterIncorrectKey = "HostRegisterIncorrectKey"

	// EventTypeGroupAdded GroupAdded event
	EventTypeGroupAdded = "GroupAdded"

	// EventTypeGroupModified GroupModified event
	EventTypeGroupModified = "GroupModified"

	// EventTypeGroupDeleted GroupDeleted event
	EventTypeGroupDeleted = "GroupDeleted"

	// EventTypeScheduleAdded ScheduleAdded event
	EventTypeScheduleAdded = "ScheduleAdded"

	// EventTypeScheduleModified ScheduleModified event
	EventTypeScheduleModified = "ScheduleModified"

	// EventTypeScheduleDeleted ScheduleDeleted event
	EventTypeScheduleDeleted = "ScheduleDeleted"

	// EventTypeAttachmentAdded AttachmentAdded event
	EventTypeAttachmentAdded = "AttachmentAdded"

	// EventTypeAttachmentModified AttachmentModified event
	EventTypeAttachmentModified = "AttachmentModified"

	// EventTypeAttachmentDeleted AttachmentDeleted event
	EventTypeAttachmentDeleted = "AttachmentDeleted"

	// EventTypeScriptAdded ScriptAdded event
	EventTypeScriptAdded = "ScriptAdded"

	// EventTypeScriptModified ScriptModified event
	EventTypeScriptModified = "ScriptModified"

	// EventTypeScriptDeleted ScriptDeleted event
	EventTypeScriptDeleted = "ScriptDeleted"

	// EventTypeScriptRun ScriptRun event
	EventTypeScriptRun = "ScriptRun"

	// EventTypeServerStarted ServerStarted event
	EventTypeServerStarted = "ServerStarted"

	// EventTypeServerOptionsModified ServerOptionsModified event
	EventTypeServerOptionsModified = "ServerOptionsModified"

	// EventTypeRegisterRuleAdded RegisterRuleAdded event
	EventTypeRegisterRuleAdded = "RegisterRuleAdded"

	// EventTypeRegisterRuleModified RegisterRuleModified event
	EventTypeRegisterRuleModified = "RegisterRuleModified"

	// EventTypeRegisterRuleDeleted RegisterRuleDeleted event
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

	EventTypeHostAdded,

	EventTypeHostModified,

	EventTypeHostDeleted,

	EventTypeHostRegisterSuccess,

	EventTypeHostRegisterIncorrectKey,

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

	EventTypeUserLoggedIn: "UserLoggedIn",

	EventTypeUserIncorrectPassword: "UserIncorrectPassword",

	EventTypeUserLoggedOut: "UserLoggedOut",

	EventTypeUserAdded: "UserAdded",

	EventTypeUserModified: "UserModified",

	EventTypeUserResetPassword: "UserResetPassword",

	EventTypeUserResetAPIKey: "UserResetAPIKey",

	EventTypeUserDeleted: "UserDeleted",

	EventTypeHostAdded: "HostAdded",

	EventTypeHostModified: "HostModified",

	EventTypeHostDeleted: "HostDeleted",

	EventTypeHostRegisterSuccess: "HostRegisterSuccess",

	EventTypeHostRegisterIncorrectKey: "HostRegisterIncorrectKey",

	EventTypeGroupAdded: "GroupAdded",

	EventTypeGroupModified: "GroupModified",

	EventTypeGroupDeleted: "GroupDeleted",

	EventTypeScheduleAdded: "ScheduleAdded",

	EventTypeScheduleModified: "ScheduleModified",

	EventTypeScheduleDeleted: "ScheduleDeleted",

	EventTypeAttachmentAdded: "AttachmentAdded",

	EventTypeAttachmentModified: "AttachmentModified",

	EventTypeAttachmentDeleted: "AttachmentDeleted",

	EventTypeScriptAdded: "ScriptAdded",

	EventTypeScriptModified: "ScriptModified",

	EventTypeScriptDeleted: "ScriptDeleted",

	EventTypeScriptRun: "ScriptRun",

	EventTypeServerStarted: "ServerStarted",

	EventTypeServerOptionsModified: "ServerOptionsModified",

	EventTypeRegisterRuleAdded: "RegisterRuleAdded",

	EventTypeRegisterRuleModified: "RegisterRuleModified",

	EventTypeRegisterRuleDeleted: "RegisterRuleDeleted",
}

// EventTypeNameMap map EventType keys to values
var EventTypeNameMap = map[string]string{

	"UserLoggedIn": "UserLoggedIn",

	"UserIncorrectPassword": "UserIncorrectPassword",

	"UserLoggedOut": "UserLoggedOut",

	"UserAdded": "UserAdded",

	"UserModified": "UserModified",

	"UserResetPassword": "UserResetPassword",

	"UserResetAPIKey": "UserResetAPIKey",

	"UserDeleted": "UserDeleted",

	"HostAdded": "HostAdded",

	"HostModified": "HostModified",

	"HostDeleted": "HostDeleted",

	"HostRegisterSuccess": "HostRegisterSuccess",

	"HostRegisterIncorrectKey": "HostRegisterIncorrectKey",

	"GroupAdded": "GroupAdded",

	"GroupModified": "GroupModified",

	"GroupDeleted": "GroupDeleted",

	"ScheduleAdded": "ScheduleAdded",

	"ScheduleModified": "ScheduleModified",

	"ScheduleDeleted": "ScheduleDeleted",

	"AttachmentAdded": "AttachmentAdded",

	"AttachmentModified": "AttachmentModified",

	"AttachmentDeleted": "AttachmentDeleted",

	"ScriptAdded": "ScriptAdded",

	"ScriptModified": "ScriptModified",

	"ScriptDeleted": "ScriptDeleted",

	"ScriptRun": "ScriptRun",

	"ServerStarted": "ServerStarted",

	"ServerOptionsModified": "ServerOptionsModified",

	"RegisterRuleAdded": "RegisterRuleAdded",

	"RegisterRuleModified": "RegisterRuleModified",

	"RegisterRuleDeleted": "RegisterRuleDeleted",
}

// IsEventType is the provided value a valid EventType
func IsEventType(q string) bool {
	_, k := EventTypeMap[q]
	return k
}

// EventTypeSchema the EventType schema.
var EventTypeSchema = []map[string]interface{}{

	{
		"name":        "UserLoggedIn",
		"description": "UserLoggedIn event",
		"value":       "UserLoggedIn",
	},

	{
		"name":        "UserIncorrectPassword",
		"description": "UserIncorrectPassword event",
		"value":       "UserIncorrectPassword",
	},

	{
		"name":        "UserLoggedOut",
		"description": "UserLoggedOut event",
		"value":       "UserLoggedOut",
	},

	{
		"name":        "UserAdded",
		"description": "UserAdded event",
		"value":       "UserAdded",
	},

	{
		"name":        "UserModified",
		"description": "UserModified event",
		"value":       "UserModified",
	},

	{
		"name":        "UserResetPassword",
		"description": "UserResetPassword event",
		"value":       "UserResetPassword",
	},

	{
		"name":        "UserResetAPIKey",
		"description": "UserResetAPIKey event",
		"value":       "UserResetAPIKey",
	},

	{
		"name":        "UserDeleted",
		"description": "UserDeleted event",
		"value":       "UserDeleted",
	},

	{
		"name":        "HostAdded",
		"description": "HostAdded event",
		"value":       "HostAdded",
	},

	{
		"name":        "HostModified",
		"description": "HostModified event",
		"value":       "HostModified",
	},

	{
		"name":        "HostDeleted",
		"description": "HostDeleted event",
		"value":       "HostDeleted",
	},

	{
		"name":        "HostRegisterSuccess",
		"description": "HostRegisterSuccess event",
		"value":       "HostRegisterSuccess",
	},

	{
		"name":        "HostRegisterIncorrectKey",
		"description": "HostRegisterIncorrectKey event",
		"value":       "HostRegisterIncorrectKey",
	},

	{
		"name":        "GroupAdded",
		"description": "GroupAdded event",
		"value":       "GroupAdded",
	},

	{
		"name":        "GroupModified",
		"description": "GroupModified event",
		"value":       "GroupModified",
	},

	{
		"name":        "GroupDeleted",
		"description": "GroupDeleted event",
		"value":       "GroupDeleted",
	},

	{
		"name":        "ScheduleAdded",
		"description": "ScheduleAdded event",
		"value":       "ScheduleAdded",
	},

	{
		"name":        "ScheduleModified",
		"description": "ScheduleModified event",
		"value":       "ScheduleModified",
	},

	{
		"name":        "ScheduleDeleted",
		"description": "ScheduleDeleted event",
		"value":       "ScheduleDeleted",
	},

	{
		"name":        "AttachmentAdded",
		"description": "AttachmentAdded event",
		"value":       "AttachmentAdded",
	},

	{
		"name":        "AttachmentModified",
		"description": "AttachmentModified event",
		"value":       "AttachmentModified",
	},

	{
		"name":        "AttachmentDeleted",
		"description": "AttachmentDeleted event",
		"value":       "AttachmentDeleted",
	},

	{
		"name":        "ScriptAdded",
		"description": "ScriptAdded event",
		"value":       "ScriptAdded",
	},

	{
		"name":        "ScriptModified",
		"description": "ScriptModified event",
		"value":       "ScriptModified",
	},

	{
		"name":        "ScriptDeleted",
		"description": "ScriptDeleted event",
		"value":       "ScriptDeleted",
	},

	{
		"name":        "ScriptRun",
		"description": "ScriptRun event",
		"value":       "ScriptRun",
	},

	{
		"name":        "ServerStarted",
		"description": "ServerStarted event",
		"value":       "ServerStarted",
	},

	{
		"name":        "ServerOptionsModified",
		"description": "ServerOptionsModified event",
		"value":       "ServerOptionsModified",
	},

	{
		"name":        "RegisterRuleAdded",
		"description": "RegisterRuleAdded event",
		"value":       "RegisterRuleAdded",
	},

	{
		"name":        "RegisterRuleModified",
		"description": "RegisterRuleModified event",
		"value":       "RegisterRuleModified",
	},

	{
		"name":        "RegisterRuleDeleted",
		"description": "RegisterRuleDeleted event",
		"value":       "RegisterRuleDeleted",
	},
}

const (

	// IPVersionOptionAuto IPv4 or IPv6 as chosen by the system automatically
	IPVersionOptionAuto = "auto"

	// IPVersionOptionIPv4 IPv4 only
	IPVersionOptionIPv4 = "ipv4"

	// IPVersionOptionIPv6 IPv6 only
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

// IPVersionOptionNameMap map IPVersionOption keys to values
var IPVersionOptionNameMap = map[string]string{

	"Auto": "auto",

	"IPv4": "ipv4",

	"IPv6": "ipv6",
}

// IsIPVersionOption is the provided value a valid IPVersionOption
func IsIPVersionOption(q string) bool {
	_, k := IPVersionOptionMap[q]
	return k
}

// IPVersionOptionSchema the IPVersionOption schema.
var IPVersionOptionSchema = []map[string]interface{}{

	{
		"name":        "Auto",
		"description": "IPv4 or IPv6 as chosen by the system automatically",
		"value":       "auto",
	},

	{
		"name":        "IPv4",
		"description": "IPv4 only",
		"value":       "ipv4",
	},

	{
		"name":        "IPv6",
		"description": "IPv6 only",
		"value":       "ipv6",
	},
}

const (

	// RegisterRulePropertyHostname Hostname
	RegisterRulePropertyHostname = "hostname"

	// RegisterRulePropertyKernelName Kernel Name
	RegisterRulePropertyKernelName = "kernel_name"

	// RegisterRulePropertyKernelVersion Kernel Version
	RegisterRulePropertyKernelVersion = "kernel_version"

	// RegisterRulePropertyDistributionName Distribution Name
	RegisterRulePropertyDistributionName = "distribution_name"

	// RegisterRulePropertyDistributionVersion Distribution Version
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

	RegisterRulePropertyHostname: "hostname",

	RegisterRulePropertyKernelName: "kernel_name",

	RegisterRulePropertyKernelVersion: "kernel_version",

	RegisterRulePropertyDistributionName: "distribution_name",

	RegisterRulePropertyDistributionVersion: "distribution_version",
}

// RegisterRulePropertyNameMap map RegisterRuleProperty keys to values
var RegisterRulePropertyNameMap = map[string]string{

	"Hostname": "hostname",

	"KernelName": "kernel_name",

	"KernelVersion": "kernel_version",

	"DistributionName": "distribution_name",

	"DistributionVersion": "distribution_version",
}

// IsRegisterRuleProperty is the provided value a valid RegisterRuleProperty
func IsRegisterRuleProperty(q string) bool {
	_, k := RegisterRulePropertyMap[q]
	return k
}

// RegisterRulePropertySchema the RegisterRuleProperty schema.
var RegisterRulePropertySchema = []map[string]interface{}{

	{
		"name":        "Hostname",
		"description": "Hostname",
		"value":       "hostname",
	},

	{
		"name":        "KernelName",
		"description": "Kernel Name",
		"value":       "kernel_name",
	},

	{
		"name":        "KernelVersion",
		"description": "Kernel Version",
		"value":       "kernel_version",
	},

	{
		"name":        "DistributionName",
		"description": "Distribution Name",
		"value":       "distribution_name",
	},

	{
		"name":        "DistributionVersion",
		"description": "Distribution Version",
		"value":       "distribution_version",
	},
}

const (

	// ScheduleResultSuccess All hosts executed the script successfully
	ScheduleResultSuccess = 0

	// ScheduleResultPartialSuccess Some hosts did not execute the script successfully
	ScheduleResultPartialSuccess = 1

	// ScheduleResultFail No hosts executed the script successfully
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

	ScheduleResultSuccess: 0,

	ScheduleResultPartialSuccess: 1,

	ScheduleResultFail: 2,
}

// ScheduleResultNameMap map ScheduleResult keys to values
var ScheduleResultNameMap = map[string]int{

	"Success": 0,

	"PartialSuccess": 1,

	"Fail": 2,
}

// IsScheduleResult is the provided value a valid ScheduleResult
func IsScheduleResult(q int) bool {
	_, k := ScheduleResultMap[q]
	return k
}

// ScheduleResultSchema the ScheduleResult schema.
var ScheduleResultSchema = []map[string]interface{}{

	{
		"name":        "Success",
		"description": "All hosts executed the script successfully",
		"value":       0,
	},

	{
		"name":        "PartialSuccess",
		"description": "Some hosts did not execute the script successfully",
		"value":       1,
	},

	{
		"name":        "Fail",
		"description": "No hosts executed the script successfully",
		"value":       2,
	},
}

// AllEnums map of all enums
var AllEnums = map[string]interface{}{

	"ClientAction": ClientActionSchema,

	"EventType": EventTypeSchema,

	"IPVersionOption": IPVersionOptionSchema,

	"RegisterRuleProperty": RegisterRulePropertySchema,

	"ScheduleResult": ScheduleResultSchema,
}
