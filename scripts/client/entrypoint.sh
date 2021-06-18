#!/bin/sh
set -e

cd /otto
if [ -z "${REGISTER_HOST}" ]; then
cat > otto_client.conf<< EOF
{
    "listen_addr": "0.0.0.0:12444",
    "psk": "${OTTO_PSK}",
    "log_path": ".",
    "default_uid": 1000,
    "default_gid": 1000,
    "path": "/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
    "allow_from": "0.0.0.0/0"
}
EOF
else
    export REGISTER_DONT_EXIT_ON_FINISH=1
fi

./otto