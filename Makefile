# Change these variables as necessary.
BINARY_NAME ?= test
MAIN_PACKAGE_PATH := ./
TEST_BASE_PATH := ./...


## audit: run quality control checks
.PHONY: code-check
code-check:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.56.1 run ./...

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test-local
test-local:
	docker compose up -d --build
	./run_tests.sh
