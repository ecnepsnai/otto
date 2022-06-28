# Scripts

A script is configured on the Otto server and is executed on the Otto agent. Scripts can be in any executable format as
long as the executable itself requires only the first and only parameter being the path to the script.

For example: `bash <script>`, `python <python file>`.

An Otto agent can run multiple scripts in parallel, and scripts can be aborted by the user during execution. Aborted
scripts are killed with SIGTERM.

## Environment Variables

Environment variables are run-time variables that are passed to the script when it is run on the host. Variables can be
configured in a couple different locations, and cascade down overwriting any duplicate keys.

1. **Global**. Configured in the options page on the Otto server. These are included in all scripts.
2. **Script**. Configured in the script. These overwrite global variables.
3. **Group**. Configured in the group. These overwrite script variables.
4. **Host**. Configured in the host. These overwrite host variables.

For example, you may want to have a script that sets a users password. The script will contain a default password but
individual groups could specify a different password that would be used by the script.

Lastly, there are a number of implicit variables that are automatically included and can not be overwritten:

|Key|Value|
|-|-|
|`OTTO_SERVER_URL`|The absolute URL to the Otto server, configured in the options page.|
|`OTTO_SERVER_VERSION`|The version of the Otto software.|
|`OTTO_HOST_ADDRESS`|The configured address of the host this script is executing on.|
|`OTTO_HOST_PORT`|The configured port of the host this script is executing on.|

When creating an environment variable you can mark the variable as "secret". This will hide the value of the variable in
the web interface.

## Attachments

You can attach files to scripts that will be uploaded and placed on hosts at specified paths before the script is run.
Attachments are uploaded each time a script runs, and will overwrite any existing files at the same path.

If the destination directory for the attachment does not exist the Otto agent will create it. The otto agent will use
the default mode for the directory, and set the owner to the same as the attachment.

Attachments can be owned by specific UID/GID or inherit the UID/GID that the script runs as, and have a specific
permission mode. Attachments have a maximum file size of 100MiB.
