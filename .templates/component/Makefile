
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
	go install github.com/tuxounet/k2@${K2_LATEST_VERSION}

update-sdk:
	go get -v github.com/tuxounet/k2-sdk
	
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
	
run: prepare
	go run ./main.go

bump-patch: build 
	echo "Bumping version patch"
	@bash ./tools/bump-patch.sh


build: prepare
	if [ -f "./component/Makefile" ]; then \
		$(MAKE) -C ./component build; \
	fi
	go build -o ./.out/${APP_NAME}-${VERSION} ./main.go