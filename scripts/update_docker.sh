#!/bin/bash
set -e

AMD64_DIGEST=$(curl -Ss 'https://hub.docker.com/v2/repositories/library/alpine/tags?page_size=25&page=1&ordering=&name=latest' | jq -r '.results.[0].images[] | select(.architecture=="amd64") | .digest')
ARM64_DIGEST=$(curl -Ss 'https://hub.docker.com/v2/repositories/library/alpine/tags?page_size=25&page=1&ordering=&name=latest' | jq -r '.results.[0].images[] | select(.architecture=="arm64") | .digest')

echo "ALPINE_HASH_AMD64=$(echo ${AMD64_DIGEST} | cut -d ':' -f2)" > alpine_hash.txt
echo "ALPINE_HASH_ARM64=$(echo ${ARM64_DIGEST} | cut -d ':' -f2)" >> alpine_hash.txt
