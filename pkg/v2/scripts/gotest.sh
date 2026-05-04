#!/usr/bin/env bash

echo "==> Running go test..."
go test -v -coverprofile=coverage.out ./...
test_exit=$?

if [ -f coverage.out ]; then
    grep -v "\.gen\.go:" coverage.out > coverage.out.tmp && mv coverage.out.tmp coverage.out
fi

if [ $test_exit -eq 1 ]; then
    exit 1
fi

exit 0
