APP_NAME=robohash
BUILD_DIR=build
OUTPUT_DIR=$(BUILD_DIR)/$(APP_NAME)
TAR_FILE=$(APP_NAME).tar.gz
SRC_DIR=./img

# Default target
all: build tar

# Build the Go binary
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(OUTPUT_DIR)
	@go build -o $(OUTPUT_DIR)/$(APP_NAME) main.go

# Archive binary and images
tar: build
	@echo "Creating tar archive..."
	@mkdir -p $(OUTPUT_DIR)/img
	@cp -r $(SRC_DIR)/* $(OUTPUT_DIR)/img/
	@tar -czf $(BUILD_DIR)/$(TAR_FILE) -C $(BUILD_DIR) $(APP_NAME)
	@rm -rf $(OUTPUT_DIR)

# Clean up build files
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)
