GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GO_PATH:=$(shell go env GOPATH)

init: 	 
	go get && go mod tidy && go fmt ./...

run: test

test: init
	go test ./...

