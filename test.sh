#!/usr/bin/env bash

do_request() {
    REQ_NO="${1}"
    START="$(gdate +%s)"
    curl 'localhost:9000/do-something-important' >/dev/null 2>&1
    END="$(gdate +%s)"

    echo "Request ${REQNO} took $(("${END}" - "${START}"))s" >&2
}

CHILDREN=()
for i in $(seq 0 39); do
    do_request &
    CHILDREN+=("${!}")
done

for PID in "${CHILDREN[@]}"; do
    wait "${PID}"
done
