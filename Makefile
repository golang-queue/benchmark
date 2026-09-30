GO ?= go
TOOLS_MOD := -modfile=go.tools.mod

## install-tools: download tool dependencies (golangci-lint) via tools modfile
install-tools:
	$(GO) mod download $(TOOLS_MOD)

## fmt: format go files using golangci-lint
fmt:
	$(GO) tool $(TOOLS_MOD) golangci-lint fmt

## lint: run golangci-lint to check for issues
lint:
	$(GO) tool $(TOOLS_MOD) golangci-lint run

testing:
	go test -v -run=^$ -benchmem -count=2  -bench .

.PHONY: install-tools fmt lint testing
