APP_NAME := myip
BUILD_DIR := build
BIN_PATH := $(BUILD_DIR)/$(APP_NAME)

.PHONY: all build install service clean

all: build

build:
	go build -o $(BIN_PATH) .

clean:
	rm -rf $(BUILD_DIR)
