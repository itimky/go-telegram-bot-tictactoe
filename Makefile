APP_NAME := go-telegram-bot-tictactoe
APP_VERSION ?= $(if $(GIT_TAG),$(GIT_TAG)-$(GIT_HASH)$(GIT_DIRTY),$(GIT_BRANCH)$(GIT_DIRTY))

GO_FILES := $(shell find . -type f -name '*.go')
BIN ?= ./bin
BUILD_ENVPARAMS := CGO_ENABLED=0
BUILDER_PATH := "github.com/itimky/go-telegram-bot-tictactoe/pkg/buildinfo"
LDFLAGS := -ldflags "-X $(BUILDER_PATH).appName=$(APP_NAME) \
	-X $(BUILDER_PATH).appVersion=$(APP_VERSION) \
	-X $(BUILDER_PATH).buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) \
 	-X $(BUILDER_PATH).buildNumber=$(BUILD_NUMBER) \
	-X $(BUILDER_PATH).gitHash=$(GIT_HASH) \
	-X $(BUILDER_PATH).gitBranch=$(GIT_BRANCH)"

.PHONY: run
run:
	go run -race $(LDFLAGS) ./cmd/$(APP_NAME)

.PHONY: build
build:
	@$(BUILD_ENVPARAMS) go build -o $(BIN)/$(APP_NAME) $(LDFLAGS) ./cmd/$(APP_NAME)

.PHONY: format
format:
	@gofmt -s -w $(GO_FILES)
