GO ?= go

.PHONY: test
test:
	$(GO) test \
		./internal/... \
		-v -count=1 \
		-coverprofile coverage.out \
		-covermode=atomic

.PHONY: coverage
coverage: test
	$(GO) tool cover -html=coverage.out

.PHONY: lint
lint:
	$(GO) vet ./...
	$(GO) tool golint -set_exit_status ./...
	$(GO) tool testifylint -v -enable-all ./...
	$(GO) tool golangci-lint run

.PHONY: install-dependencies
install-dependencies:
	$(GO) install tool
