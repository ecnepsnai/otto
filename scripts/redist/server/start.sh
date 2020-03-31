#!/bin/bash

rm -f /opt/otto/data/logs/stderr.log
/opt/otto/current/otto -d /opt/otto/data 2>/opt/otto/data/logs/stderr.log