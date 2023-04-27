# Otto Protocol

The Otto protocol is the shared protocol used for communication between the Otto hosts and the Otto server. The Otto
protocol uses SSH as the transport layer, providing strong and reliable encryption.

## Message structure

```
|---------------------------------------------|
| Protocol Version                            |
| (4 bytes - network byte order)              |
|---------------------------------------------|
| Message Type                                |
| (4 bytes - network byte order)              |
|---------------------------------------------|
| Message Length                              |
| (4 bytes - network byte order)              |
|---------------------------------------------|
| Message Data                                |
| (binary, optional)                          |
|                                             |
|                                             |
|                                             |
|---------------------------------------------|
| Additional Data                             |
| (undetermined, optional)                    |
|                                             |
|                                             |
|                                             |
|---------------------------------------------|
```

A message is a self-contained set of information sent by either the Otto host or Otto server. Each message must always
contain a protocol version, message type, and message length.

The message type field is used to indicate what action to take with the given message. It must be a
[valid enum value](https://pkg.go.dev/github.com/ecnepsnai/otto#MessageType).

Messages can have two forms of data associated with them: a well-defined data structure and additional undefined data.

The message length field of an Otto message describes the length of a defined message structure. Within that structure
may be a length property to describe how much additional data is included.

For example, when uploading a file to an Otto host, metadata about the file is contained within the message data, but
the actual file contents is included as additonal data following the metadata.

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
