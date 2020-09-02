#!/bin/sh
set -e

/otto/otto --data-dir /otto_data -b 0.0.0.0:8080 2>/otto_data/logs/stderr.log