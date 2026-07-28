#!/usr/bin/env bash

echo "==> Running golangci-lint..."
go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.11.4 run ./...
if [ $? -eq 1 ]; then
    echo ""
    echo "Golangci-lint found suspicious constructs. Please check the reported"; \
    echo "constructs and fix them if necessary before submitting the code for review."; \
    exit 1
fi

exit 0
