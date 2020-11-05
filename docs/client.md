# Client

An Otto client is a individual host that is running the Otto client daemon. Scripts are run on clients.

## Installing the Client

Client binaries are provided on any Otto server server at `/clients`.

## Running the Client

The Otto client is a static executable file that supports any *nix like system (Linux, BSD, macOS, Solaris).

It works best if you run it as root, but will run as a non-root user.

### Automatic Registration

If enabled on the server, clients can configure themselves by automatically registering with the Otto server.

For server configuration information, see the server documentation.

Run the Otto client executable with the following environment variables **only once** to register the host:

|Variable|Description|
|-|-|
|`REGISTER_HOST`|The base URL of the Otto server, including the protocol and port (if needed). Must not contain any trailing slash.|
|`REGISTER_PSK`|The register PSK|
|`REGISTER_NO_TLS_VERIFY`|Optional. If `1` then no TLS verification is done when connecting to the server.|
|`OTTO_CLIENT_PORT`|Optional. Specify the port that the Otto client will listen on.|

**Example:**

```bash
REGISTER_HOST='https://otto.mydomain' REGISTER_PSK='super_secret' ./otto
```

The client will then configure itself and exit with a status code of `0` and will now be ready for normal use.

### Manual Configuration

You may also manually configure the Otto client with a configuration file. The configuration file is a JSON file with a
single, top-level object. The `otto_client.conf` configuration file must be in the same directory as the Otto client binary.

|Property|Required|Description|
|-|-|-|
|`listen_addr`|No|The address to listen to. Defaults to `0.0.0.0:12444`.|
|`psk`|Yes|The PSK configured for this host on the server|
|`log_path`|No|The path to a file where the Otto client should log. Defaults to the directory where the otto binary is.|
|`default_uid`|No|The default UID if not specified by the script. Defaults to `0`.|
|`default_gid`|No|The default GID if not specified by the script. Defaults to `0`.|
|`path`|No|The value of $PATH when scripts are executed|

**Example:**

```json
{
    "psk": "36C1CD5993F64EF0394C0DE9DE12567D",
    "log_path": ".",
    "path":"/usr/local/sbin:/usr/local/bin:/sbin:/bin:/usr/sbin:/usr/bin:/root/bin",
    "default_uid": 0,
    "default_gid": 0
}
```