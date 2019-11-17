SHELL := /bin/sh

all: compile

compile:
	go build

test:
	go test ./...

.PHONY: force
