# Server

The Otto server is the central location for all your scripts and hosts.

## Running the Server

### As a Container (Recommended!)

A container image compatible with Docker and Podman is the preferred way of running the Otto server.

If you were using docker, you can run the container with:

```bash
docker run -p 8080:8080 -v <data dir>:/otto_data otto:latest
```

Replace `<data dir>` with a directory where you want Otto to store all server data, or omit the volume parameter entirely if you don't care about persistence

Navigating to `http://localhost:8080` in your web browser and use the default credentials of `admin`:`admin` to log in. As long as you are using the default credentials a warning will appear at the top of the page, so be sure to change the password right away by clicking "admin" in the top right and selecting "Edit User" to change your password.

### As a Service

You may also wish to run the Otto server as an executable on any Linux, macOS, FreeBSD, or NetBSD host.

Download the binary for your operating system and start run the `otto` executable.

Command line options are:

```
-d --data-dir <path>        Specify the absolute path to the data directory
-b --bind-addr <socket>     Specify the listen address for the web server
-v --verbose                Set the log level to debug
--no-scheduler              Disable all automatic tasks
```

## Users & Authentication

Otto currently only supports local user accounts. When the server starts up, if there are no user accounts it will create the default account of `admin` with the password `admin`.

You can add users in the Options tab of the web interface. There needs to be at least one user for Otto to function, but you can delete the `admin` user if you create a new user.

### Resetting a forgotten password

**Reset the Password for Somebody Else**

Any user can change the password for other users simply by editing their user in the options page of the web interface.

**Reset/Restore the Default Account**

If you have no way to access the Otto service:

1. Stop the Otto server
2. Navigate to the data directory for the otto server
3. Delete the files `user.db` and `session.db`
4. Start the Otto server

The default account will be recreated and you can log in using `admin`:`admin`.