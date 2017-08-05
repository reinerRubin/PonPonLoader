.PHONY: build

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
BIN_DIR := $(ROOT_DIR)/bin

build:
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/pon-pon-loader $(ROOT_DIR)/cmd/loader/main.go
