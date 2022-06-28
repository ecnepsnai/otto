# API

The Otto server is controlled via a REST API. This document details how you may use that API.

## Using the API

All API endpoints, excluding the login endpoint, require that you provide a username and API key. API keys are
associated with a user and must be generated before it can be used. Generate an API key by clicking the "Reset API Key"
button when editing a user.

Each HTTP request must provide the following headers:

- `X-OTTO-USERNAME`: The username of the user
- `X-OTTO-API-KEY`: The API key for the user

Only URLs that start with `/api` can be accessed using these headers.

# Endpoints

## Authentication

**POST /api/login**

**Note:** This should only be used if you cannot use the API key authentication method.

This is the only request that can be performed without an session cookie, or with an expired cookie. Upon successful
authentication, a valid session cookie is required.

Expected Body:
```json
{
    "Username": "",
    "Password": ""
}
```

**POST /api/logout**

Used to terminate an existing session. No body is required. The response should be ignored.

Otto sessions are automatically cleaned up, so logging out is not mandatory.

## Hosts

**GET /api/hosts**



**PUT /api/hosts/host**



**GET /api/hosts/host/:id**



**GET /api/hosts/host/:id/scripts**



**GET /api/hosts/host/:id/groups**



**GET /api/hosts/host/:id/schedules**



**GET /api/hosts/host/:id/id**

Get the server public key for this host.

Example response:
```json
{
    "error": {},
    "code": 200,
    "data": "AAAAC3NzaC1lZDI1NTE5AAAAIOlJYnMAeyFwKSLy3zCy4QCAcnYahCJd12sFAQETOUK3"
}
```

**POST /api/hosts/host/:id/heartbeat**

Trigger a heartbeat on a host

Expected body: none

Example response:
```json
{
    "error": {},
    "code": 200,
    "data": {
        "IsReachable": true,
        "LastAttempt": "2021-08-12T20:01:20.335900268-07:00",
        "Version": "0.10.2",
        "LastReply": "2021-08-12T20:01:20.335900134-07:00",
        "Address": "192.168.0.1",
        "Properties": {
            "hostname": "example.host.foo.bar",
            "kernel_version": "4.18.0-305.10.2.el8_4.x86_64",
            "kernel_name": "Linux",
            "distribution_version": "8.4 (Green Obsidian)",
            "distribution_name": "Rocky Linux"
        }
    }
}
```

**POST /api/hosts/host/:id/id/trust**

Modify the trust for this host.

Expected body:
```json
{
	"Action": "permit",
	"PublicKey": "AAAAC3NzaC1lZDI1NTE5AAAAIOlJYnMAeyFwKSLy3zCy4QCAcnYahCJd12sFAQETOUK3"
}
```

`Action` must be either `permit` or `deny`. The `PublicKey` property is only used when the action is `permit`. If it is
omitted or an empty string, the pending public key is trusted.

Example response:
```json
{
    "code": 200,
    "error": {},
    "data": {
        "ID": "rea_UKwyyQBX",
        "Name": "example",
        "Address": "127.0.0.1",
        "Port": 12444,
        "Enabled": true,
        "Trust": {
            "TrustedIdentity": "AAAAC3NzaC1lZDI1NTE5AAAAIOlJYnMAeyFwKSLy3zCy4QCAcnYahCJd12sFAQETOUK3",
            "UntrustedIdentity": "",
            "LastTrustUpdate": "2022-02-04T19:03:25.322461409-08:00"
        },
        "GroupIDs": [
            "y910Mb38cmud"
        ],
        "Environment": null
    }
}
```

**POST /api/hosts/host/:id/id/rotate**

Rotate the identities for this host.

Expected body: none.

Example response:
```json
{
    "code": 200,
    "error": {},
    "data": true
}
```

**POST /api/hosts/host/:id**



**DELETE /api/hosts/host/:id**



**PUT /api/register**



**GET /api/heartbeat**




## Groups

**GET /api/groups**



**GET /api/groups/membership**



**PUT /api/groups/group**



**GET /api/groups/group/:id**



**GET /api/groups/group/:id/scripts**



**GET /api/groups/group/:id/hosts**



**GET /api/groups/group/:id/schedules**



**POST /api/groups/group/:id/hosts**



**POST /api/groups/group/:id**



**DELETE /api/groups/group/:id**




## Schedules

**GET /api/schedules**



**PUT /api/schedules/schedule**



**GET /api/schedules/schedule/:id**



**GET /api/schedules/schedule/:id/reports**



**GET /api/schedules/schedule/:id/hosts**



**GET /api/schedules/schedule/:id/groups**



**GET /api/schedules/schedule/:id/script**



**POST /api/schedules/schedule/:id**



**DELETE /api/schedules/schedule/:id**




## Scripts

**GET /api/scripts**



**PUT /api/scripts/script**



**GET /api/scripts/script/:id**



**GET /api/scripts/script/:id/hosts**



**GET /api/scripts/script/:id/groups**



**GET /api/scripts/script/:id/schedules**



**GET /api/scripts/script/:id/attachments**



**POST /api/scripts/script/:id/groups**



**POST /api/scripts/script/:id**



**DELETE /api/scripts/script/:id**




## Attachments

**GET /api/attachments**



**PUT /api/attachments**



**GET /api/attachments/attachment/:id**



**GET /api/attachments/attachment/:id/download**

Expected body: None.

Download the attachment file. May return a binary stream, make sure to check the `Content-Type`.

**POST /api/attachments/attachment/:id**



**DELETE /api/attachments/attachment/:id**




## Script Execution

**PUT /api/action/sync**

Executes a script on a single host. Will return a result when the script has exited.

Expected body:
```json
{
    "HostID": "",
    "Action": "",
    "ScriptID": ""
}
```

**WS /api/action/async**

A websocket that can be used to execute a script on a single host and receive live output from the running script.

Upon connecting to the socket, the agent must send a JSON message to start the script:
```json
{
    "HostID": "",
    "Action": "",
    "ScriptID": ""
}
```

The server will respond with messages of the following structure:
```json
{
    "Code": 0,
    "Error": "",
    "Stdout": "",
    "Stderr": "",
    "Result": {},
}
```

`Code` will be:
- 100 for output from the script
- 200 for completion of the script
- 400 for an error

`Stdout` and `Stderr` will be the current, entire text of the standard output and error from the script.

`Result` will only be present on completion of the script and will contain the scripts result

## Users

**GET /api/users**



**PUT /api/users/user**



**GET /api/users/user/:username**



**POST /api/users/user/:username**

Modify the existing user `:username`. Can be used to modify other users.

Expected body:
```json
{
    "Email": "",
    "Enabled": true,
}
```

To change a users password, include the variable `Password` with a string value containing the new password in the
request. The password will not be changed if the `Password` variable is not present, or is an empty string.

**POST /api/users/user/:username/apikey**



**DELETE /api/users/user/:username**

Delete the existing user `:username`.

You cannot delete yourself. If the deleted user has any active sessions, they will be terminated immediately.

## System

**GET /api/state**

Returns the current state of the system. This endpoint is used by the web interface and contains:
- Otto Options
- The Server Version
- The Current User

**GET /api/stats**

Returns statistics about the Otto server.

Example response:
```json
{
    "data": {
        "NumberGroups": 1,
        "NumberHosts": 0,
        "NumberSchedules": 0,
        "NumberScripts": 0,
        "NumberUsers": 0,
        "ReachableHosts": 0,
        "TrustedHosts": 0,
        "UnreachableHosts": 0,
        "UntrustedHosts": 0
    },
    "code": 200
}
```

**GET /api/options**



**POST /api/options**


## Registration Rules

**GET /api/register/rules**



**PUT /api/register/rules/rule**



**GET /api/register/rules/rule/:id**



**POST /api/register/rules/rule/:id**



**DELETE /api/register/rules/rule/:id**



## Events

**GET /api/events**

