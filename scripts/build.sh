#!/bin/bash -e

export GO_BUILD_DELIVERABLE="${GO_BUILD_DELIVERABLE:-p2iwd}"

curl https://gist.githubusercontent.com/030/620a95e7a699c4db3e76b2b8b0309909/raw/e218edf140117c1b368a3fcdb7d6110a073dd0e7/go-build.sh -o go-build.sh
chmod +x go-build.sh
./go-build.sh p2iwd "${GO_BUILD_DELIVERABLE}"