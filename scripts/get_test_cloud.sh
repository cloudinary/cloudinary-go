#!/usr/bin/env bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

GO_VER=$(go version | head -n 1 | cut -d ' ' -f 3| cut -c 3-);
SDK_VER=$(grep -oiP '(?<=Version \= \")([a-zA-Z0-9\-.]+)(?=")' "${DIR}"/../api/api.go)

go run "${DIR}"/allocate_test_cloud.go "GO ${GO_VER} SDK ${SDK_VER}"
