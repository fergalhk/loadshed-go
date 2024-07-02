#!/usr/bin/env bash

do_request() {
    REQ_NO="${1}"
    URL_PATH="${2}"
    START="$(gdate +%s)"
    curl "localhost:9000${URL_PATH}" >/dev/null 2>&1
    END="$(gdate +%s)"

    echo "Request ${REQ_NO} to "${URL_PATH}" took $(("${END}" - "${START}"))s" >&2
}

CHILDREN=()
for i in $(seq 0 39); do
    do_request "${i}" /do-something-important &
    CHILDREN+=("${!}")
    do_request "${i}" /livez &
    CHILDREN+=("${!}")
done

for PID in "${CHILDREN[@]}"; do
    wait "${PID}"
done
