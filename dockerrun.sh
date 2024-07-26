#!/usr/bin/env bash
set -eEuo pipefail

docker run \
    -p 9000:9000 \
    --rm \
    -it \
    -v "${HOME}/go/pkg/mod:/go/pkg/mod" \
    -v "${PWD}:/var/tmp/work" \
    --workdir /var/tmp/work \
    golang:1.22 \
    go run .
