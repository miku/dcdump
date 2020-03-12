#!/bin/bash
#
# Example script to harvest datacite. Since some results sets get large
# (see: https://gist.github.com/miku/176edd1222fc42ae3b23234bc9d3cd87) using
# minute intervals for all the data is a more expensive, but a bit more robust
# approach.
#
# Note, that manual intervention might still be required, because of an
# unexpected, non-recoverable HTTP 500 or HTTP 403.
#
# $ dcdump.sh [DIR]
#
# If DIR is not given, we create a default `dcdump-YYYY-MM-DD` in the current
# directory.

set -eu -o pipefail

HARVEST_DIR=${1:-dcdump-$(date +"%Y-%m-%d")}
if [ ! -d "$HARVEST_DIR" ]; then
	mkdir -p "$HARVEST_DIR"
fi

# Harvest with one slice per minute (-e), into harvest dir (-d) with a common
# prefix (-p).
PATH="$PATH:$(pwd)" dcdump -i e -d "$HARVEST_DIR" -p "dcdump-"
