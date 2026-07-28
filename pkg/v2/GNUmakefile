default: tests

tests: golangci-lint unittest

unittest:
	@sh -c "'$(CURDIR)/scripts/gotest.sh'"

golangci-lint:
	@sh -c "'$(CURDIR)/scripts/golangci_lint_check.sh'"

generate-mocks:
	go run github.com/vektra/mockery/v2@v2.53.6

oapi-codegen:
	@cd swagger && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.6.0 --config=oapi-codegen-models.yaml managed-kubernetes.swagger.yaml
	@cd swagger && go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.6.0 --config=oapi-codegen-client.yaml managed-kubernetes.swagger.yaml

.PHONY: tests unittest golangci-lint oapi-codegen
