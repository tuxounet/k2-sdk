GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GO_PATH:=$(shell go env GOPATH)
VERSION := $(if $(CI_COMMIT_TAG),$(CI_COMMIT_TAG),v${GIT_BRANCH})
VERSION_FILE := ./version.txt


write-version: 
	echo ${VERSION} > ${VERSION_FILE}
init: 	 
	go get && go mod tidy && go fmt ./...

run: test

test: init
	go test ./...

build: write-version
	mkdir -p ./.out
	go build  -o ./.out/k2-sdk ./main.go

	
bump-patch: build 
	echo "Bumping version patch"
	@bash ./tools/bump-patch.sh