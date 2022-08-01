#!/bin/bash

if [ $1 -eq 0 ] ; then 
    systemctl --no-reload disable --now otto-agent.service &>/dev/null || : 
fi