GO ?= go
COVERAGE=coverage.out

.PHONY: all check view-coverage

all:

view-coverage: check
	$(GO) tool cover -html=$(COVERAGE)

check:
	$(GO) test -v -cover -coverpkg=./... -coverprofile=$(COVERAGE) ./...


