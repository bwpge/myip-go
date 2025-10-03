APP_NAME := myip
VERSION := 0.1.0
BUILD_DIR := tmp
BIN_PATH := $(BUILD_DIR)/$(APP_NAME)
INSTALL_PATH := /usr/local/bin/$(APP_NAME)
SERVICE_NAME := $(APP_NAME).service
SERVICE_PATH := /etc/systemd/system/$(SERVICE_NAME)

DEB_BUILD_DIR := $(BUILD_DIR)/deb
DEB_PACKAGE := $(APP_NAME)
DEB_VERSION := $(VERSION)
DEB_REV ?= 1
DEB_MAINTAINER := bwpge <bwpge.dev@gmail.com>
DEB_DEPENDS :=
DEB_ARCH := amd64
DEB_HOMEPAGE := https://github.com/bwpge/$(APP_NAME)-go
DEB_DESCRIPTION := A simple go API for returning the requestor's IP address


.PHONY: all build build-deb clean

all: build

build:
	go build -o $(BIN_PATH) .

build-deb: build
	cp -r build/deb $(DEB_BUILD_DIR)
	mkdir -p $(DEB_BUILD_DIR)/usr/bin/
	cp $(BIN_PATH) $(DEB_BUILD_DIR)/usr/bin/$(APP_NAME)
	chmod +x $(DEB_BUILD_DIR)/usr/bin/$(APP_NAME)
	@echo "Package: $(DEB_PACKAGE)" > $(DEB_BUILD_DIR)/DEBIAN/control
	@echo "Version: $(DEB_VERSION)" >> $(DEB_BUILD_DIR)/DEBIAN/control
	@echo "Maintainer: $(DEB_MAINTAINER)" >> $(DEB_BUILD_DIR)/DEBIAN/control
	@echo "Depends: $(DEB_DEPENDS)" >> $(DEB_BUILD_DIR)/DEBIAN/control
	@echo "Architecture: $(DEB_ARCH)" >> $(DEB_BUILD_DIR)/DEBIAN/control
	@echo "Homepage: $(DEB_HOMEPAGE)" >> $(DEB_BUILD_DIR)/DEBIAN/control
	@echo "Description: $(DEB_DESCRIPTION)" >> $(DEB_BUILD_DIR)/DEBIAN/control
	chmod 0755 $(DEB_BUILD_DIR)
	dpkg -b $(DEB_BUILD_DIR) $(BUILD_DIR)/$(DEB_PACKAGE)_$(DEB_VERSION)-$(DEB_REV)_$(DEB_ARCH).deb

clean:
	rm -rf $(BUILD_DIR)
