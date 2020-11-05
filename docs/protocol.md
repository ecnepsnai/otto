# Otto Protocol

The Otto client and server communicate using the Otto protocol.

## Message Structure

Each message includes a version number, and encrypted data length.

```
|---------------------------------------------|
| Protocol Version                            |
| (4 bytes - network byte order)              |
|---------------------------------------------|
| Encrypted Data Length                       |
| (4 bytes - network byte order)              |
|---------------------------------------------|
| Encrypted Data                              |
| (binary)                                    |
|                                             |
|                                             |
|                                             |
|                                             |
|---------------------------------------------|
```

## Encryption

Encrypted Data is encrypted with AES-256-GCM. Pre-shared keys are hashed with scrypt.
