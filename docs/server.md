# Server

The Otto server is the central location for all your scripts and hosts.

## Running the Server

### As a Container (Recommended!)

A container image compatible with Docker and Podman is the preferred way of running the Otto server.

If you were using docker, you can run the container with:

```bash
docker run -p 8080:8080 -v <data dir>:/otto_data -e OTTO_UID=$(id -u) -e OTTO_GID=$(id -g) --name otto otto
```

Navigating to the web interface, use the default credentials of `admin`:`admin` to log in. You'll want to
change that password right away, so head to the Options tab.

### As a Service

You may also wish to run the Otto server as an executable on any Linux, macOS, FreeBSD, or NetBSD host.

Download the approiate binary for your operating system and start run the `otto` executable.

## Accessing the Service

A HTTP web UI will be available on port `8080` (unless changed). The default username and password is *admin* and *admin*.

Once logged in, navigate to Options, scroll down to Users, select "admin" and change the password of the default admin user.
Alternativly you way wish to create a new user and delete the default admin user.