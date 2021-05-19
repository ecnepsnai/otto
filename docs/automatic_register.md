# Automatic Client Registration

If enabled, Otto clients can register themselves with the Otto server and be assigned to groups based on information
about the client.

## Otto Server

To enable client registration, check "Allow Hosts to Register Themselves" in the options page of the Otto web UI.

A register PSK must be specified. This PSK must be specified when telling the client to register itself with the Otto
server.

Registration rules can be added to automatically assign hosts to specific groups based off of information about the
operating system of the host. Each rule must have at least one clause, which is a simple regex test against a predefined
property of the system

Possible properties are:
- **Hostname.** The hostname of the host.
- **Kernel Name.** The name of the kernel running on the host, as determined by running `uname`.
- **Kernel Version.** The version of the kernel running on the host, as determined by running `uname -r`.
- **Distribution Name.** The name of the distribution or variant of the host. The value varies by system.
- **Distribution Version.** The version of the distribution or variant of the host. The value varies by distribution.

Each clause must match for the host to be added to the group specified by the rule. Multiple rules may be applied to
incoming hosts.

For example, you may wish to have a rule that assign hosts to a group for CentOS Linux and another for Ubuntu Linux,
or you may further segregate hosts into specific versions such as CentOS Linux 7 or Ubuntu Linux 20.04.

There is an implicit 'any' rule at the end that will assign the host to a default group, much must be specified.

## Otto Client

Run the Otto client executable with the following environment variables **only once** to register the host:

|Variable|Description|
|-|-|
|`REGISTER_HOST`|The base URL of the Otto server, including the protocol and port (if needed). Must not contain any trailing slash.|
|`REGISTER_PSK`|The register PSK|
|`REGISTER_NO_TLS_VERIFY`|Optional. If `1` then no TLS verification is done when connecting to the server.|
|`OTTO_CLIENT_PORT`|Optional. Specify the port that the Otto client will listen on.|

To aid with registration, running the client with the `-v` argument will print out the property values that are passed
to the Otto server during registration.

**Example:**

```bash
REGISTER_HOST='https://otto.mydomain' REGISTER_PSK='super_secret' ./otto
```

The client will then configure itself and exit with a status code of `0` and will now be ready for normal use.