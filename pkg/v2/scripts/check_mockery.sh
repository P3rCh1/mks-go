#!/usr/bin/env bash

echo "==> Checking mocks are up to date..."

mockery

if [ -n "$(git diff --stat -- 'mksclient/mocks/')" ]; then
    echo "Mocks are outdated!"
    exit 1
fi

exit 0
