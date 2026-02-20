.PHONY: build-plugins clean-plugins

PLUGIN_DIR := plugins
 
# Build all plugins as .so files
build-plugins:
	
	@for dir in $(PLUGIN_DIR)/*/; do \
		$(MAKE) -C $$dir build  		
	done
	@echo "All plugins built successfully"

# Clean built plugins
clean-plugins:
	@for dir in $(PLUGIN_DIR)/*/; do \
		$(MAKE) -C $$dir clean  		
	done
	@echo "All plugins cleaned successfully"

