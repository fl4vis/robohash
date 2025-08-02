APP_NAME=robohash
BUILD_DIR=build
TAR_FILE=$(APP_NAME).tar.gz
SRC_DIR=./img

# Default target
all: build tar

# Build the Go binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) main.go

# Archive binary
tar: build
	@echo "Creating tar archive..."
	@tar -czf $(BUILD_DIR)/$(TAR_FILE) -C $(BUILD_DIR) $(APP_NAME)

# Clean up build files
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
