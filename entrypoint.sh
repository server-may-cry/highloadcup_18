#!/bin/sh
set -ex

ulimit -m 1950000
# 2097152

exec "$@"
