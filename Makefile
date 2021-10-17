.PHONY: build
build:
		go build -v ./cmd/chefbook

.DEFAULT_GOAL := build