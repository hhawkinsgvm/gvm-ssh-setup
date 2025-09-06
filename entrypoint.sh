#!/usr/bin/env bash
set -e
# REAL_HOME should be set to the host's $HOME when running via Docker
if [ -n "$REAL_HOME" ]; then
  export HOME="$REAL_HOME"
fi
exec /usr/local/bin/gvm-ssh "$@"