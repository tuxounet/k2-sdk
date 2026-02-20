PARENT_DIR = $(shell dirname $(CURDIR))
COMPONENT_NAME = $(shell basename $(PARENT_DIR))
OUT_DIR = ../../../components/registry/data/dist/

.PHONY: all install build clean dev

all: build

prepare:
	npm install

build: 
	npm run build

release: build
	mkdir -p $(OUT_DIR)/$(COMPONENT_NAME)
	rm -f $(OUT_DIR)/$(COMPONENT_NAME)/*
	cp dist/* $(OUT_DIR)/$(COMPONENT_NAME)/

clean:
	rm -rf dist node_modules

dev:
	npm run dev
