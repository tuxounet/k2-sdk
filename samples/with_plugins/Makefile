.PHONY: build-plugins clean-plugins

PLUGIN_DIR := plugins
PLUGIN_OUTPUT_DIR := plug_host/.data/plugins

# Build all plugins as .so files
build-plugins:
	@mkdir -p $(PLUGIN_OUTPUT_DIR)
	@for dir in $(PLUGIN_DIR)/*/; do \
		plugin_name=$$(basename $$dir); \
		echo "Building plugin: $$plugin_name"; \
		(cd $$dir && go build -tags plugin -buildmode=plugin -o ../../$(PLUGIN_OUTPUT_DIR)/$$plugin_name.so .); \
	done
	@echo "All plugins built successfully"

# Clean built plugins
clean-plugins:
	rm -rf $(PLUGIN_OUTPUT_DIR)/*.so

# Build single plugin (usage: make build-plugin PLUGIN=plugin_one)
build-plugin:
	@mkdir -p $(PLUGIN_OUTPUT_DIR)
	cd $(PLUGIN_DIR)/$(PLUGIN) && go build -tags plugin -buildmode=plugin -o ../../$(PLUGIN_OUTPUT_DIR)/$(PLUGIN).so .
	@echo "Plugin $(PLUGIN) built successfully"

run: build-plugins
	cd plug_host && go run main.go