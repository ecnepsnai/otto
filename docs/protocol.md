# Otto Protocol

The Otto server communicates with agents using the Otto protocol. The protocol is a message based system where a
'frame' encapsulates each 'message'. The Otto protocol uses SSH as the transport layer, providing strong and reliable
encryption.

# Protocol Components

## Frame

The goal of the frame is to transport the message contents.

### Structure

Each frame includes a version number, length of the data, and the message.

```
|---------------------------------------------|
| Protocol Version                            |
| (4 bytes - network byte order)              |
|---------------------------------------------|
| Data Length                                 |
| (4 bytes - network byte order)              |
|---------------------------------------------|
| Message Data                                |
| (binary)                                    |
|                                             |
|  |---------------------------------------|  |
|  | Message Type                          |  |
|  | ( 4 bytes - network byte order)       |  |
|  |---------------------------------------|  |
|  | Message Data                          |  |
|  | (binary)                              |  |
|  |                                       |  |
|  |                                       |  |
|  |                                       |  |
|  |---------------------------------------|  |
|                                             |
|---------------------------------------------|
```

## Message

The message contains instructions or results. Messages are defined by a message type, and each type has a corresponding
direction, indicating who is the sender and the recipient.

### Structure

Otto messages includes the message type, and binary data that corresponds with the data structure associated with the
message type.

|Message Type|Direction|Description|
|-|-|-|
|`HEARTBEAT_REQUEST`|Server to Agent|A heartbeat request|
|`HEARTBEAT_RESPONSE`|Agent to Server|A heartbeat response, includes the agent version|
|`TRIGGER_ACTION`|Server to Agent|A request to trigger a specific action on the agent|
|`CANCEL_ACTION`|Server to Agent|Cancel any in-progress action|
|`ACTION_OUTPUT`|Agent to Server|A portion of, or the entire output (both stdout and stderr) from the action|
|`ACTION_RESULT`|Agent to Server|The result of the action|
|`GENERAL_FAILURE`|Agent to Server|A message to indicate a general error|

# Encryption

To security transport messages between agents and hosts, the Otto protocol is designed to use the SSH transport
security protocol.

When the Otto agent starts for the first time, it generates a Ed25519 key, and when that agent is registered, either
manually or automatically, the Otto server also generates a unique Ed25519 key for that agent. The agent then must
be configured to trust the public key from the Otto server.

When the Otto server connects to the Otto agent, the server must use the specific key for that agent. The agent will
verify the public key matches the one it was configured to trust, and if it matches it will allow the connection.

**Note:** While the SSH protocol is used, the Otto agent does not actually use OpenSSH or any SSH identites or
configuration files on the system.

# Process

Except in host registration (which does not use the Otto protocol), Otto Servers always connect to Otto agents, however
the Otto protocol is not request then reply, such as HTTP.

When a message is received the recipient must first determine wether or not it can understand the message by checking
the protocol version. Otto agents will refuse messages using different protocol versions.

If the protocol version matches, the agent will then collect the length of the encrypted data, then read all of that
data into memory. The otto protocol requires fully-formed messages and does not currently support streaming.

With the message in hand, the agent can determine what data type the message data is by using the message type
value. The binary data of the message is a [gob](https://golang.org/pkg/encoding/gob/) encoded byte-slice.

The connection is then used for further messages from either the Otto agent to server. For example, the server may tell
the agent to cancel a running script after it has started.

Some actions will happen after the agent has closed the connection, for example agents may exit after closing the
connection when requested by the server.
