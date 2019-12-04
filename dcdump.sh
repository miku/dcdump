#!/bin/bash

# Example mini script to harvest datacite in three partial runs to mitigate API
# errors due to huge result sets
# (https://gist.github.com/miku/176edd1222fc42ae3b23234bc9d3cd87).

set -eu -o pipefail

HARVEST_DIR=${1:-$(mktemp -d -t dcdump-XXXXXXXXXX)}
echo >&2 "harvest dir: $HARVEST_DIR"

function finish() {
	echo >&2 "harvest dir: $HARVEST_DIR"
}
trap finish EXIT

PATH="$PATH:$(pwd)" dcdump -s 2018-01-01 -e '2019-07-31 23:59:59' -i daily -d tmp -p 'part-01-'
PATH="$PATH:$(pwd)" dcdump -s 2019-08-01 -e '2019-08-03 23:59:59' -i e -d tmp -p 'part-02-'
PATH="$PATH:$(pwd)" dcdump -s 2019-08-04 -i daily -d tmp -p 'part-03-' # implicit -e, today

echo >&2 "done: $HARVEST_DIR"
