# Otto Protocol

The Otto client and server communicate using the Otto protocol. The Otto protocol is a message based system where a
'frame' encapsulates each 'message'. The message itself is encrypted in the frame.

# Protocol Components

## Frame

The goal of the frame is to transport the encrypted message contents.

### Structure

Each frame includes a version number, length of the encrypted data, and the encrypted message.

```
|---------------------------------------------|
| Protocol Version                            |
| (4 bytes - network byte order)              |
|---------------------------------------------|
| Encrypted Data Length                       |
| (4 bytes - network byte order)              |
|---------------------------------------------|
| Encrypted Message Data                      |
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

### Encryption

Message Data is encrypted with AES-256-GCM. Pre-shared keys are hashed with scrypt.

### Structure

Otto messages includes the message type, and binary data that corresponds with the data structure associated with the
message type.

|Message Type|Direction|Description|
|-|-|-|
|`HEARTBEAT_REQUEST`|Server to Client|A heartbeat request|
|`HEARTBEAT_RESPONSE`|Client to Server|A heartbeat response, includes the client version|
|`TRIGGER_ACTION`|Server to Client|A request to trigger a specific action on the client
|`ACTION_OUTPUT`|Client to Server|A portion of, or the entire output (both stdout and stderr) from the action|
|`ACTION_RESULT`|Client to Server|The result of the action|
|`GENERAL_FAILURE`|Client to Server|A message to indicate a general error|

# Process

Except in host registration (which does not use the Otto protocol), Otto Servers always connect to Otto clients, never
the other way around.

When a message is received the recipient must first determine wether or not it can understand the message by checking
the protocol version. Otto clients will refuse messages using different protocol versions.

If the protocol version matches, the client will then collect the length of the encrypted data, then read all of that
data into memory. The otto protocol requires fully-formed messages and does not currently support streaming.

The client will then attempt to decrypt the encrypted data using the configured PSK. If decryption fails, the client
will log out an error and silently close the connection.

With the original message in hand, the client can determine what data type the message data is by using the message type
value. The binary data of the message is a [gob](https://golang.org/pkg/encoding/gob/) encoded byte-slice.

The same connection is re-used for any responses from the Otto client to the server, however some actions may take place
after the client has closed the connection. For example, clients may restart after disconnecting form the server when
requested.
