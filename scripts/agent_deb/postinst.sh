#!/bin/bash

if [[ $1 -eq 1 ]]; then
    systemctl --no-reload preset otto-agent.service &>/dev/null || true
fi
