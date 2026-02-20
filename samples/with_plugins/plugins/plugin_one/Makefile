.PHONY: build-plugins clean-plugins
PLUGIN := $(notdir $(CURDIR))
PLUGIN_DIR := plugins
PLUGINS_OUTPUT_DIR := ../../plug_host/registry/dist
LOCAL_PLUGIN := plugin_one
 
# Build specific plugin (usage: make build-plugin PLUGIN=plugin_name)
build:
	@mkdir -p $(PLUGINS_OUTPUT_DIR)
	go build -tags plugin -buildmode=plugin -o $(PLUGINS_OUTPUT_DIR)/$(PLUGIN).so .
	@echo "Plugin $(PLUGIN) built successfully"

run: build
	cd plug_host && go run main.go


clean: 
	rm -f $(PLUGINS_OUTPUT_DIR)/$(PLUGIN).so
	@echo "Plugin $(PLUGIN) cleaned successfully"