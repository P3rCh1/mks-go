#!/usr/bin/env bash

echo "==> Checking mocks are up to date..."

go run github.com/vektra/mockery/v2@v2.53.6

if [ -n "$(git diff --stat -- 'mksclient/mocks/')" ]; then
    echo "Mocks are outdated!"
    exit 1
fi

exit 0
