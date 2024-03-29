# Scripts

A script is configured on the Otto server and is executed on the Otto agent. Scripts can be in any executable format as
long as the executable itself requires only the first and only parameter being the path to the script.

For example: `bash <script>`, `python <python file>`.

The Otto agent can run multiple unique scripts in parallel and scripts can be aborted during execution. Aborted scripts
are terminated with SIGTERM.

By default scripts run with their working directory set to a temporary directory, unless the working directory is
specified by the script configuration.

When a script is running, its standard output and error are copied to a file within a temporary directory created by the
agent.

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

When creating an environment variable you can mark the variable as "hidden". This will hide the value of the variable in
the web interface. Take note, however, that hidden environment variables are not obfuscated from script output. Take
care not to print any hidden environment variables to stdout ot stderr.

## Attachments

You can attach files to scripts that will be uploaded and placed on hosts at specified paths. Attachments are uploaded
each time a script runs, and will overwrite any existing files at the same path.

By default the attachment is uploaded before the script runs, however you can optionally specify that the attachment
should be uploaded after the script has completed. In this case, the file will be written only if the script exited with
no error. If the script exited with an error, the file is not uploaded.

If the destination directory for the attachment does not exist the Otto agent will create it. The otto agent will use
the default mode for the directory, and set the owner to the same as the attachment.

Attachments can be owned by specific UID/GID or inherit the UID/GID that the script runs as, and have a specific
permission mode. Attachments have a maximum file size of 100MiB.

Attachment files are compared against a SHA-256 hash calculated when the file is uploaded to the Otto server. This check
occurs both on the server and the client. Do not modify the attachment file inside the Otto server's data directory
directly, as the attachment will be deleted automatically. Always update the attachment through the Otto server web UI.
