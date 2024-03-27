ifeq (, $(shell which testifylint))
$(error "No 'testifylint' on PATH, consider doing: go install github.com/Antonboom/testifylint@latest")
endif

ifeq (, $(shell which golint))
$(error "No 'golint' on PATH, consider doing: go install golang.org/x/lint/golint@latest")
endif

ifeq (, $(shell which golangci-lint))
$(error "No 'golangci-lint' on PATH, consider following these instructions: https://golangci-lint.run/welcome/install/#local-installation")
endif

.PHONY: lint
lint:
	golint -set_exit_status ./...
	testifylint ./...
	golangci-lint run
