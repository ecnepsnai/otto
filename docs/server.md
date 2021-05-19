# Server

The Otto server is the central location for all your scripts and hosts.

## Running the Server

### As a Container (Recommended!)

An OCI-compatible container image is the preferred way of running the Otto server. The image can be used with Podman or
Docker.

If you were using Podman, you can run the container with:

```bash
podman run -p 8080:8080 -v <data dir>:/otto_data ghcr.io/ecnepsnai/otto:latest
```

*Note:* Substitute `podman` with `docker` if you're using Docker.

Replace `<data dir>` with a directory where you want Otto to store all server data, or omit the volume parameter
entirely if you don't care about persistence

Navigating to `http://localhost:8080` in your web browser and use the default credentials of `admin`:`admin` to log in.
You will need to change your password the first time you log in.

We recommend using a reverse proxy such as NGINX and configuring TLS.

*Note:* An arm64v8 image is available as of version 0.8.0 on the
[https://github.com/ecnepsnai/otto/releases](Releases) page.

### As a Service

You may also wish to run the Otto server as an executable on any Linux, FreeBSD, NetBSD, or OpenBSD.

### System Requirements

**Hardware:**
- CPU: Any semi-recent amd64/x86_64
- RAM: At least 500MiB of available system memory
- Disk: Varies by use and log retention. 1GiB will be plenty.

**Operating System:**
- Linux kernel version 2.6.23 or later for amd64/x86_64 systems
- OpenBSD stable release
- FreeBSD 10 or later for amd64/x86_64 systems
- NetBSD 8 or later. NetBSD 7 may work but is unsupported due to known and unresolved issues

### Usage

Download the binary for your operating system and start run the `otto` executable.

Command line options are:

```
-d --data-dir <path>        Specify the absolute path to the data directory
-b --bind-addr <socket>     Specify the listen address for the web server
-v --verbose                Set the log level to debug
--no-scheduler              Disable all automatic tasks
```

For example:

```bash
otto -d /usr/share/otto -b 0.0.0.0:8080
```

## Users & Authentication

Otto currently only supports local user accounts. When the server starts up and there are no user accounts it will
create the default account of `admin` with the password `admin`.

You can add users in the Options tab of the web interface. There needs to be at least one user for Otto to function,
but you can delete the `admin` user if you create a new user.

### Resetting a forgotten password

**Reset the Password for Somebody Else**

Any user can change the password for other users simply by editing their user in the options page of the web interface.
All active sessions for that user will be ended if their password is changed by somebody else.

**Reset/Restore the Default Account**

If you have no way to access the Otto service then you will need to reset the user store.

*Note: All users will be deleted and will need to be manually recreated.*

1. Stop the Otto server
2. Navigate to the data directory for the otto server
3. Delete `user.db` and `shadow.db`
4. Start the Otto server

The default account will be recreated and you can log in using `admin`:`admin`.
