#!/bin/sh
set -e

if [ "$1" = 'app' ]; then
    [[ -n "$SESSION_KEY" ]] || (echo "SESSION_KEY property not passed; exiting"; exit 1)
    exec app
fi

exec "$@"
