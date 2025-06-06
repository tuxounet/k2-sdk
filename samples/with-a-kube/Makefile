
APP_NAME := WithAKubeSample
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(if $(CI_COMMIT_TAG),$(CI_COMMIT_TAG),v${GIT_BRANCH})
GO_PATH:=$(shell go env GOPATH)
VERSION_FILE := ./app/version.go

init: 
	go get -v

write-version: 
	echo "package app" > ${VERSION_FILE}
	echo "" >> ${VERSION_FILE}
	echo "const (" >> ${VERSION_FILE}
	echo "    AppName    = \"${APP_NAME}\"" >> ${VERSION_FILE}
	echo "    AppVersion = \"${VERSION}\"" >> ${VERSION_FILE}
	echo ")" >> ${VERSION_FILE}

prepare: init write-version	
	go install github.com/swaggo/swag/cmd/swag@v1.16.4
	${GO_PATH}/bin/swag init --instanceName ${APP_NAME}
	${GO_PATH}/bin/swag fmt
	
run: prepare
	go run \
		-ldflags="-X 'app.AppVersion=${VERSION}' -X 'app.AppName=${APP_NAME}'" \
		./main.go
		
ROOT_DIR := $(shell pwd)
enter: 
	cd ${ROOT_DIR}/.data/home && \
		HOME=${ROOT_DIR}/.data/home \
		bash