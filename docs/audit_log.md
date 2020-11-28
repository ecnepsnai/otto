# Audit Log

The Otto server maintains a log of various events, both user triggered and automatic, for later reporting in the audit log.

The most recent 20 events can be viewed in the Otto server web interface.

## Events

The following events are recorded by the Otto server

### UserLoggedIn

Event for when a user successfully logs in to the Otto web interface.

|Parameter|Description|
|-|-|
|`username`|The username of the user who logged in|
|`remoteAddr`|The remote IP address of the user|

### UserIncorrectPassword

Event for when an incorrect password is provided while attempting to log in to the Otto web interface.

|Parameter|Description|
|-|-|
|`username`|The username of the user who attempted to log in|
|`remoteAddr`|The remote IP address of the user|

### UserLoggedOut

Event for when a user logs out.

|Parameter|Description|
|-|-|
|`username`|The username of the user who logged out|

### UserAdded

Event for whe a new user is added.

|Parameter|Description|
|-|-|
|`username`|The username of the user|
|`email`|The email address of the user|
|`added_by`|The username of the user who added this new user|

### UserModified

Event for when an existing user is modified.

|Parameter|Description|
|-|-|
|`username`|The username of the user|
|`modified_by`|The username of the user who modified this user|

### UserDeleted

Event for when a user is deleted.

|Parameter|Description|
|-|-|
|`username`|The username of the user|
|`deleted_by`|The username of the user who deleted this user|

### HostAdded

Event for whe a new host is added.

|Parameter|Description|
|-|-|
|`host_id`|The ID of the host|
|`name`|The name of the host|
|`address`|The address of the host|
|`added_by`|The username of the user who added this new host|

### HostModified

Event for when an existing host is modified.

|Parameter|Description|
|-|-|
|`host_id`|The ID of the host|
|`name`|The name of the host|
|`address`|The address of the host|
|`modified_by`|The username of the user who modified this host|

### HostDeleted

Event for when a host is deleted.

|Parameter|Description|
|-|-|
|`host_id`|The ID of the host|
|`name`|The name of the host|
|`address`|The address of the host|
|`deleted_by`|The username of the user who deleted this host|

### HostRegisterSuccess

Event for when a new host successfully registers itself and is added.

|Parameter|Description|
|-|-|
|`host_id`|The ID of the host|
|`name`|The name of the host|
|`address`|The address of the host|
|`uname`|The uname provided by the host|
|`group_id`|The ID of the group that the host was added to|
|`matched_rule_property`|If the host matched a registration rule, the property of the matched rule|
|`matched_rule_pattern`|If the host matched a registration rule, the pattern of the matched rule|
|`matched_rule_group_id`|If the host matched a registration rule, the group ID of the matched rule|

### HostRegisterIncorrectPSK

Event for when host registration fails with an incorrect pre-shared key.

|Parameter|Description|
|-|-|
|`address`|The address of the host|
|`uname`|The uname provided by the host|
|`hostname`|The hostname provided by the host|

### GroupAdded

Event for when a new group is added.

|Parameter|Description|
|-|-|
|`group_id`|The ID of the group|
|`name`|The name of the group|
|`added_by`|The username of the user who added this new group|

### GroupModified

Event for when an existing group is modified.

|Parameter|Description|
|-|-|
|`group_id`|The ID of the group|
|`name`|The name of the group|
|`modified_by`|The username of the user who modified this group|

### GroupDeleted

Event for when a group is deleted.

|Parameter|Description|
|-|-|
|`group_id`|The ID of the group|
|`name`|The name of the group|
|`deleted_by`|The username of the user who deleted this group|

### ScheduleAdded

Event for when a new schedule is added.

|Parameter|Description|
|-|-|
|`schedule_id`|The ID of the schedule|
|`name`|The name of the schedule|
|`script_id`|The ID of the script|
|`pattern`|The frequency pattern of the schedule|
|`added_by`|The username of the user who added this new schedule|

### ScheduleModified

Event for when an existing schedule is modified.

|Parameter|Description|
|-|-|
|`schedule_id`|The ID of the schedule|
|`name`|The name of the schedule|
|`script_id`|The ID of the script|
|`pattern`|The frequency pattern of the schedule|
|`modified_by`|The username of the user who modified this schedule|

### ScheduleDeleted

Event for when a schedule is deleted.

|Parameter|Description|
|-|-|
|`schedule_id`|The ID of the schedule|
|`name`|The name of the schedule|
|`script_id`|The ID of the script|
|`pattern`|The frequency pattern of the schedule|
|`deleted_by`|The username of the user who deleted this schedule|

### AttachmentAdded

Event for when a new attachment is added.

|Parameter|Description|
|-|-|
|`attachment_id`|The ID of the attachment|
|`name`|The name of the attachment|
|`file_path`|The file path of the attachment|
|`mimetype`|The mimetype of the attachment|
|`added_by`|The username of the user who added this new attachment|

### AttachmentModified

Event for when an existing attachment is modified.

|Parameter|Description|
|-|-|
|`attachment_id`|The ID of the attachment|
|`name`|The name of the attachment|
|`file_path`|The file path of the attachment|
|`mimetype`|The mimetype of the attachment|
|`modified_by`|The username of the user who modified this attachment|

### AttachmentDeleted

Event for when a attachment is deleted.

|Parameter|Description|
|-|-|
|`attachment_id`|The ID of the attachment|
|`deleted_by`|The username of the user who deleted this attachment|

### ScriptAdded

Event for when a new script is added.

|Parameter|Description|
|-|-|
|`script_id`|The ID of the script|
|`name`|The name of the script|
|`added_by`|The username of the user who added this new script|

### ScriptModified

Event for when an existing script is modified.

|Parameter|Description|
|-|-|
|`script_id`|The ID of the script|
|`name`|The name of the script|
|`modified_by`|The username of the user who modified this script|

### ScriptDeleted

Event for when a script is deleted.

|Parameter|Description|
|-|-|
|`script_id`|The ID of the script|
|`name`|The name of the script|
|`deleted_by`|The username of the user who deleted this script|

### ScriptRun

Event for when a script is run.

|Parameter|Description|
|-|-|
|`script_id`|The ID of the script|
|`host_id`|The ID of the host this script run on|
|`exit_code`|The return or exit code of the script|
|`schedule_id`|If this script was triggered by a schedule, the ID of that schedule|
|`triggered_by`|If this script was triggered by a user, the username of that user|

### ServerStarted

Event for when the Otto server is started.

|Parameter|Description|
|-|-|
|`args`|The command line arguments passed to the server|

### ServerOptionsModified

Event for when the Otto server options are modified.

|Parameter|Description|
|-|-|
|`config_hash`|The SHA-256 hash of the config file|
|`modified_by`|The username of the user who modified the server options|
