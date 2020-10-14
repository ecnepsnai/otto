# Introduction to Otto

Welcome to Otto! Where things are "Ottomatic Beyond Belief!"

Otto is a way for you to run scripts on remote hosts from a single server. It supports most Unix like systems, including: Linux, BSD (Open, Free, and Net), macOS, and Solaris.

There are two components to Otto: The Server and Client.

**The Otto Service**

The otto service is the central location where hosts, groups, and scripts are configured. The otto server connects to clients
to run scripts. All configuration is stored in this central location.

**The Otto Client**

The otto client is a small piece of software that runs on your hosts and accepts requests from the otto server.

# Getting Started

Let's get started using Otto!

## Starting the Otto Service

You can run the client with no configuration options and it will listen to `localhost:8080`. See the server documentation for more information on starting the server software.

Once the server is running, access the URL and log in using the default credentials of `admin` and `admin`. A warning will appear at the top as long as you're logged in using these exact credentials, so we recommend changing that right away from the menu in the top right.

**Warning:** Do not expose the Otto web interface to the internet.

## Starting the Otto Client

Before you can add an host to the server it must first be configured. See the client documentation for instructions on how to configure a client.

Once you've configured the client ensure that it's listening on the port 12444 (unless you changed it) and that your firewall is configured to allow incoming connections to that port.

## Add a Host

On the Otto service web interface navigate to the Hosts list and click "Create New".

Input the friendly-name and the address of the server, the PSK, and check the "Otto Clients" group.

## Add a Group (Optional)

Groups are the primary component for both Scripts and Hosts. Hosts belong to groups, and scripts are assigned to groups.

```
---------           ----------           -----------
| HOSTS |  ------>  | GROUPS |  <------  | SCRIPTS |
---------           ----------           -----------
```

A default group will be created the first time you run Otto, called "Otto Clients". You can rename or delete this group later if you wish.

## Add a Script

Scripts are the actual executed code that is run on the clients. For more detailed information into scripts, see the script documentation.

To create a script: navigate to the Scripts list and click "Create New".

Give the script a name, and add commands to the script body. Assign the script to the "Otto Clients" group.

## Execute a Script

Now that we have a script assigned to a group, and a host that belongs to that group, you will now see that the script you created can now be executed on your host. Wherever you see a green "Play" button you can execute the script and see the results in the web interface.