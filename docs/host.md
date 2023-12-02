# Host

An Otto host is a individual host that is running the Otto agent daemon. Scripts are run on hosts by the agent.

## Installing the Agent

Agent binaries are provided by the Otto server at `/agents/`. Otto servers only provide the same version of agent as
the server itself.

### System Requirements

**Hardware:**
- CPU: AMD64: x86-64-v2 or newer. ARM64: ARMv8 or newer.
- RAM: At least 250MiB of available system memory
- Disk: Varies by log retention, agents are generally less than 10MiB.

**Network:**
- At least one network interface with a valid IPv4 or IPv6 address
- Must accept incoming connections on the Otto agent port (default: 12444)

**Operating System:**
- Linux kernel version 2.6.23 or later for amd64/x86_64 systems, 2.33 or later for arm64 systems
- OpenBSD stable release
- FreeBSD 10 or later for amd64/x86_64 systems, 12 or later for arm64 systems
- NetBSD 8 or later. NetBSD 7 may work but is unsupported due to known and unresolved issues

## Running the Agent

The Otto agent is a static executable file that supports Linux and most BSD systems.

It works best if you run it as root, but will run as a non-root user.

### Automatic Registration

If enabled on the server, agents can configure themselves by automatically registering with the Otto server.

Please see the [Automatic Registration](automatic_register.md) documentation for further information.

### Interactive Setup

The Otto agent offers an interactive setup which can be accessed by running the agent with the `-s` argument.

### Manual Configuration

You may also manually configure the Otto agent with a configuration file. The configuration file is a JSON file with a
single, top-level object. The `otto_agent.conf` configuration file must be in the same directory as the Otto agent
binary.

|Property|Required|Type|Description|Default Value|
|-|-|-|-|-|
|`listen_addr`|No|string|The address to listen to.|`0.0.0.0:12444`|
|`identity_path`|No|string|The full path to where the agent identity will be saved.|`.otto_id.der`|
|`server_identity`|Yes|string|The public key from the Otto server.||
|`log_path`|No|string|The directory where the Otto agent log will be saved. Do not specify a file name.|`.`|
|`default_uid`|No|number|The UID for scripts to run as.|The current UID of the Otto agent process|
|`default_gid`|No|number|The GID for scripts to run as.|The current GID of the Otto agent process|
|`path`|No|string|The value of the $PATH environment variable used when running scripts.|Value of `$PATH`|
|`allow_from`|No|[]string|Array of CIDR addresses where connections from Otto Servers will be allowed from.|`["0.0.0.0/0", "::/0"]`
|`script_timeout`|No|number|Maximum number of seconds a script can run before it is automatically aborted. Passing a negative number disables the timeout.|600 (10 minutes)|
|`reboot_command`|No|string|Path to executable to run when rebooting the host.|`/usr/sbin/reboot`|
|`shutdown_command`|No|string|Path to executable to run when shutting down the host.|`/usr/sbin/halt`|

**Example:**

```json
{
    "listen_addr": "0.0.0.0:12444",
    "identity_path": ".otto_id.der",
    "server_identity": "AAAAC3NzaC1lZDI1NTE5AAAAIOnAN2JvtaL7AHsQlfj0IXmxHJSh6/3gKSP7lYwIDszZ",
    "log_path": ".",
    "default_uid": 1000,
    "default_gid": 1000,
    "path": "/usr/local/bin:/usr/local/sbin:/usr/bin:/usr/sbin",
    "allow_from": [
        "192.168.0.0/16"
    ]
}

```

Please note that the Otto agent may update this config file, such as to rotate the server identity.

## Identity Management

An identity refers to a private and public key used as part of the Otto protocol. The Otto agent maintains an identity
that is used when the Otto server connects to the Otto agent. The Otto server also maintains a unique identity for each
Otto host.

The Otto agent must be configured to trust the public key from the Otto host. You can access the servers unique public
key on the Otto server web interface on the Trust menu of a host.

Should you need to manually update the identity of the server on an Otto agent, you can do so by running the agent
software with the `-s` argument. For example: `./otto -s <server identity>`. Alternatively you can manually edit the
Otto agent config file.

