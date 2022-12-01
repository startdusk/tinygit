# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #
# Create the new confirm target.
.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

# Testing. (go install github.com/rakyll/gotest)
.PHONY: test
test: clean
	@gotest -v ./...

.PHONY: clean
	@go mod tidy
	@go fmt ./...
	@go vet ./...

# ==================================================================================== #
# BUILD
# ==================================================================================== #
## build/tinygit: build the cmd/tinygit application
.PHONY: build/tinygit
build/tinygit:
	@echo 'Building cmd/tinygit...'
	go build -ldflags='-s -X main.version=${VERSION}' -o=./bin/tinygit ./cmd/tinygit
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -X main.version=${VERSION}' -o=./bin/linux_amd64/tinygit ./cmd/tinygit

