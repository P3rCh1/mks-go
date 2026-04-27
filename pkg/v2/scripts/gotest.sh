#!/usr/bin/env bash

echo "==> Running go test..."
go test -v -coverprofile=coverage.out $(go list ./... | grep -v mksclient)

if [ $? -eq 1 ]; then
    exit 1
fi

exit 0