# Client

An Otto client is a individual host that is running the Otto client daemon. Scripts are run on clients.

## Installing the Client

Client binaries are provided by the Otto server at `/clients`. Otto servers only provide the same version of clients as
the server itself.

### System Requirements

**Hardware:**
- CPU: Any semi-recent amd64/x86_64 or arm64 CPU. Generally, if it can run any of the operating systems listed below,
it'll work for Otto. 32-bit CPUs are not supported.
- RAM: At least 250MiB of available system memory
- Disk: Varies by log retention, clients are generally less than 10MiB.

**Network:**
- At least one network interface with a valid IPv4 or IPv6 address
- Must accept incoming connections on the Otto client port (default: 12444)

**Operating System:**
- Linux kernel version 2.6.23 or later for amd64/x86_64 systems, 2.33 or later for arm64 systems
- OpenBSD stable release
- FreeBSD 10 or later for amd64/x86_64 systems, 12 or later for arm64 systems
- NetBSD 8 or later. NetBSD 7 may work but is unsupported due to known and unresolved issues

## Running the Client

The Otto client is a static executable file that supports Linux and most BSD systems.

It works best if you run it as root, but will run as a non-root user.

### Automatic Registration

If enabled on the server, clients can configure themselves by automatically registering with the Otto server.

Please see the [Automatic Registration](automatic_register.md) documentation for further information.

### Manual Configuration

You may also manually configure the Otto client with a configuration file. The configuration file is a JSON file with a
single, top-level object. The `otto_client.conf` configuration file must be in the same directory as the Otto client
binary.

The Otto client offers an interactive setup which can be accessed by running the client with the `-s` argument.

If you choose to manually configure the Otto client, you must also add the host to the Otto server through the web UI.

|Property|Required|Description|
|-|-|-|
|`listen_addr`|No|The address to listen to. Defaults to `0.0.0.0:12444`.|
|`psk`|Yes|The PSK configured for this host on the server.|
|`log_path`|No|The path to a file where the Otto client should log. Defaults to the directory where the otto binary is.|
|`default_uid`|No|The default UID if not specified by the script. Defaults to `0`.|
|`default_gid`|No|The default GID if not specified by the script. Defaults to `0`.|
|`path`|No|The value of $PATH when scripts are executed|
|`allow_from`|No|A CIDR address where connections from Otto Servers will be allowed from. Defaults to `0.0.0.0/0`.|

**Example:**

```json
{
    "psk": "36C1CD5993F64EF0394C0DE9DE12567D",
    "log_path": ".",
    "path":"/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/root/bin",
    "default_uid": 0,
    "default_gid": 0,
    "allow_from": "10.0.0.0/8"
}
```

Please note that the Otto client may update this config file, such as to rotate the PSK.
