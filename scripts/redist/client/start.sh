#!/bin/bash

rm -f /opt/otto/stderr.log
cp /opt/otto/otto /opt/otto/.otto_running
/opt/otto/.otto_running 2>/opt/otto/stderr.log