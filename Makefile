MCP_BIN_DIR ?= $(HOME)/.local/bin
APP_NAME := mcp-assistant

.PHONY: build install deploy

build:
	go build -o bin/$(APP_NAME) main.go
install:
	mkdir -p $(MCP_BIN_DIR)
	cp bin/$(APP_NAME) $(MCP_BIN_DIR)/$(APP_NAME)
deploy:	build install