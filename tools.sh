#!/bin/bash

gotest() {
    go test -coverprofile=cover.out \
        -coverpkg=github.com/ArtProcessors/cmdy,github.com/ArtProcessors/cmdy/usage,github.com/ArtProcessors/cmdy/flags,github.com/ArtProcessors/cmdy/arg \
        github.com/ArtProcessors/cmdy/...
}

"$1" "${@:2}"

