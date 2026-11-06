APP_NAME := robohash
BIN_DIR := bin
CMD := ./cmd

.PHONY: all build cross-build build-linux build-darwin build-windows test clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) $(CMD)

cross-build: build-linux build-darwin build-windows

build-linux:
	@mkdir -p $(BIN_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 $(CMD)

build-darwin:
	@mkdir -p $(BIN_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64 $(CMD)
	GOOS=darwin GOARCH=arm64 go build -o $(BIN_DIR)/$(APP_NAME)-darwin-arm64 $(CMD)

build-windows:
	@mkdir -p $(BIN_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe $(CMD)

clean:
	rm -rf $(BIN_DIR)
