# Variables
BINARY_NAME=task-timer
SOURCE_FILES=main.go
INSTALL_PATH=/usr/local/bin

# Build the binary
build:
	go build -o $(BINARY_NAME) $(SOURCE_FILES)
# Clean the build
clean:
	rm -f $(BINARY_NAME)
# Install the binary
install: build
	sudo mv $(BINARY_NAME) $(INSTALL_PATH)
# Uninstall the binary
uninstall:
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
.PHONY: build clean install uninstall