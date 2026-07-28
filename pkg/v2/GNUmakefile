default: tests

tests: golangci-lint unittest

unittest:
	@sh -c "'$(CURDIR)/scripts/gotest.sh'"

golangci-lint:
	@sh -c "'$(CURDIR)/scripts/golangci_lint_check.sh'"

generate-mocks:
	mockery

check-mockery:
	@sh -c "'$(CURDIR)/scripts/check_mockery.sh'"

oapi-codegen:
	@cd swagger && oapi-codegen --config=oapi-codegen-models.yaml managed-kubernetes.swagger.yaml
	@cd swagger && oapi-codegen --config=oapi-codegen-client.yaml managed-kubernetes.swagger.yaml

.PHONY: tests unittest golangci-lint oapi-codegen
