#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

set +e
GIT_COMMIT=$(git rev-parse HEAD)
set -e

go build -ldflags "-X main.GitCommit ${GIT_COMMIT}" -o dist/craftctl
