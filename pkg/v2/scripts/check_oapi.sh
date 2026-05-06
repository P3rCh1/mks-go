#!/usr/bin/env bash

echo "==> Checking oapi generated code is up to date..."

cd swagger
oapi-codegen --config=oapi-codegen-models.yaml managed-kubernetes.swagger.yaml
oapi-codegen --config=oapi-codegen-client.yaml managed-kubernetes.swagger.yaml
cd ..

if [ -n "$(git diff --stat -- 'mksclient/')" ]; then
    echo "oapi generated code is outdated!"
    exit 1
fi

exit 0
