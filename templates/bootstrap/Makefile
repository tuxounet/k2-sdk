
APP_NAME :="{{ .name }}"
GOPATH := $(shell go env GOPATH)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(if $(CI_COMMIT_TAG),$(CI_COMMIT_TAG),v${GIT_BRANCH})
GO_PATH:=$(shell go env GOPATH)
VERSION_FILE := ./app/version.go
K2_LATEST_VERSION := $(shell curl --silent "https://api.github.com/repos/tuxounet/k2/tags" | jq -r '.[0].name')

init:
	go get -v


update-k2:	
	@go install github.com/tuxounet/k2@${K2_LATEST_VERSION}
	${GOPATH}/bin/k2 apply

plan:
	${GOPATH}/bin/k2 plan  --inventory ./k2.inventory.yaml

apply:
	${GOPATH}/bin/k2  apply --inventory ./k2.inventory.yaml

destroy:
	${GOPATH}/bin/k2  destroy --inventory ./k2.inventory.yaml


write-version: 
	echo "package app" > ${VERSION_FILE}
	echo "" >> ${VERSION_FILE}
	echo "const (" >> ${VERSION_FILE}
	echo "    AppName    = \"${APP_NAME}\"" >> ${VERSION_FILE}
	echo "    AppVersion = \"${VERSION}\"" >> ${VERSION_FILE}
	echo ")" >> ${VERSION_FILE}

prepare: init write-version	
	go get -v
	go mod tidy
	go install github.com/swaggo/swag/cmd/swag@v1.16.4
	${GO_PATH}/bin/swag init --instanceName ${APP_NAME} --parseInternal --parseDepth 1 --exclude "refs/"
	${GO_PATH}/bin/swag fmt --exclude "refs/"
	
run: prepare
	go run ./main.go

bump-patch: build 
	echo "Bumping version patch"
	@bash ./tools/bump-patch.sh