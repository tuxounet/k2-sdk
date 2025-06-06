GOPATH := $(shell go env GOPATH)

prepare:
	if [ -f ./Makefile.inc ]; then \
		make -f Makefile.inc prepare; \
	fi
	go install github.com/swaggo/swag/cmd/swag@v1.16.4
	$(shell go env GOPATH)/bin/swag init --instanceName {{ .name }} --parseInternal --parseDepth 1  --generalInfo component.go
	$(shell go env GOPATH)/bin/swag fmt  

build: prepare
	if [ -f ./Makefile.inc ]; then \
		make -f Makefile.inc build; \
	fi
	if [ -d "./ui" ]; then \
		if [ -f "./ui/package.json" ]; then \
			cd ui && npm install && npm run build; \
		fi; \
	fi
	