.PHONY: build-plugins clean-plugins

# Build all plugins as .so files
build-plugins:
	@for dir in ./*/; do $(MAKE) -C $$dir build; done
	@echo "All plugins built successfully"

# Clean built plugins
clean-plugins:
	@for dir in ./*/; do $(MAKE) -C $$dir clean; done
	@echo "All plugins cleaned successfully"

