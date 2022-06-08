# Contributing to Otto

## Things you'll need

Otto can be developed on any modern Linux or macOS system. So far we've tested development on Fedora 33 and macOS Big
Sur. It may be possible to develop on Windows, but we'd recommend using WSL2 instead if you use Windows.

- The latest version of Golang (seriously, the very latest release)
- The latest LTS version of node.js
- Podman or Docker (podman is preferred)

# Project Structure

## Server Web UI (aka Frontend)

The web UI for the Otto server is a React.JS web application written in Typescript and packaged using Webpack.

Directory structure:

```
/frontend
    /css   -> Contains all Sass .scss files
    /html  -> Contains the HTML for the login page and application frame
    /img   -> Contains all images
    /src
        /components   -> Contains all shared React components
        /pages        -> Contains react components that make up the visible "pages" of the web UI.
        /services     -> Contains shared services for utilities that are not react components.
        /types        -> Contains typescript definitions of backend API types.
```

## libotto

libotto is a golang library that defines the common data structures shared between the Otto client and server. This is
all contained within the `/otto` directory.

## Otto Client

The otto client is a small golang application with no runtime requirements.

Source code for the otto client is located in `/otto/cmd/client`

## Otto Server

The otto server is a golang application that powers the otto web UI and is what actually interacts with the otto
clients.

Directory structure:

```
/scripts
    /codegen   -> Contains definitions for the cbgen golang code generator.
/otto
    /server        -> Contains all golang code for the server
        /environ   -> Library that defines an environment variable. Broken off for easier testing.
```

The actual executable for the server itself is located in `/otto/cmd/server`. The executable is only responsible for starting
the server and capturing signals from the system.

# Running a Development Build

To run a development build of the server you first need to get everything prepared. The `install_backend.sh` and
`install_frontend.sh` scripts will prepare the codebase and compile a debug version of the app.

To run the Otto server:

```bash
./run.sh -v
```
*Note: Changes to the backend server software are not automatically reflected in the running instance. You will need to quit and restart the server to show any new changes.*

To automatically compile any changes to the front-end:

```bash
cd frontend
node start_webpack.js --watch 
```

# Releasing the Application

To compile release artifacts use the provided release script.

```bash
cd scripts
./release.sh <version to release>
```
*Note: If you get the error `docker: command not found...` you need to tell the script to use `docker`, as it will default to using `podman`. Specify the `DOCKER` environment variable with the docker executable path.*
