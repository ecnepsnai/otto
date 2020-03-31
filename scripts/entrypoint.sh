#!/bin/sh
set -e

if [ -z "$OTTO_UID" ]; then
	OTTO_UID=$(cat /etc/passwd | grep otto | cut -d: -f3)
	echo "OTTO_UID variable not specified, defaulting to otto user id ($OTTO_UID)"
fi

if [ -z "$OTTO_GID" ]; then
	OTTO_GID=$(cat /etc/group | grep otto | cut -d: -f3)
	echo "OTTO_GID variable not specified, defaulting to otto user group id ($OTTO_GID)"
fi

usermod -u $OTTO_UID -g $OTTO_GID --non-unique otto > /dev/null 2>&1
chown -R $OTTO_UID:$OTTO_GID /otto_data

exec su otto -s /bin/sh -c '/otto/otto --data-dir /otto_data -b 0.0.0.0:8080'