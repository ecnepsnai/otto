# Scripts

A script is configured on the otto server and is executed on the otto client. Scripts can be in any executable format as long
as the executable itself requires only the first and only parameter being the path to the script.

For example: `bash <script>`, `python <python file>`

## Environment Variables

Environment variables allow you to customize the results of the script on a per host basis.

Envrionment variables can be configured at multiple levels:

1. **Global**. Configured in the options page on the Otto server. They're included in all scripts.
2. **Script**. Configured in the script. These overwrite global variables.
3. **Group**. Configured in the group. These overwrite script variables.
4. **Host**. Configured in the host. These overwrite host variabled.

Lastly, there are a number of implicit variables that are automatically included:

|Key|Value|
|-|-|
|`OTTO_URL`|The absolute URL to the otto server, configured in the options page|
|`OTTO_VERSION`|The version of the otto software|
|`OTTO_HOST_ADDRESS`|The configured address of the host this script is executing on|
|`OTTO_HOST_PORT`|The configured port of the host this script is executing on|
|`OTTO_HOST_PSK`|The configured PSK of the host this script is executing on|