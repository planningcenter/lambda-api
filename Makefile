ROOT_DIR := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

.PHONY: test
test:
	cd ${ROOT_DIR} && go test -v ./...

.PHONY: clean
clean:
	cd ${ROOT_DIR} && go mod tidy
