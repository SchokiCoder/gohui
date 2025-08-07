# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024 - 2025  Andy Frank Schoknecht

# dest dir, uncomment below to install for only current user instead
BIN_DESTDIR:="/usr/local/bin"
CFG_DESTDIR:="/etc/hui"
#BIN_DESTDIR:="$(HOME)/.local/bin"
#CFG_DESTDIR:="$(HOME)/.config/hui"

# binaries to be installed
INSTALL_BINARIES:=courier hui

# here starts Makefile content (no touchy, end-user)
PAGER_NAME_FORMAL:=Courier
APP_NAME_FORMAL  :=House User Interface
LICENSE          :=GPL-2.0-or-later
LICENSE_URL      :=https://www.gnu.org/licenses/gpl-2.0.html
REPO             :=https://github.com/SchokiCoder/gohui
VERSION          :=v1.4

.PHONY: all build clean install purge remove vet

all: vet build

build: courier hui

clean:
	rm -f courier hui

install: $(INSTALL_BINARIES)
	mkdir -p $(BIN_DESTDIR)
	cp -t $(BIN_DESTDIR) $(INSTALL_BINARIES)
	mkdir -p $(CFG_DESTDIR)
	cp pkg/common.json $(CFG_DESTDIR)
	cp pkg/courier.json $(CFG_DESTDIR)
	cp pkg/hui.json $(CFG_DESTDIR)

purge: remove
	rm -f $(CFG_DESTDIR)/*.json
	rmdir $(CFG_DESTDIR)

remove:
	rm -f $(BIN_DESTDIR)/courier $(BIN_DESTDIR)/hui

vet:
	go vet ./pager
	go vet ./main

courier:
	go build -ldflags "-X 'main.AppLicense=$(LICENSE)' -X 'main.AppLicenseUrl=$(LICENSE_URL)' -X 'main.AppName=$@' -X 'main.AppNameFormal=$(PAGER_NAME_FORMAL)' -X 'main.AppRepo=$(REPO)' -X 'main.AppVersion=$(VERSION)'" ./pager

hui:
	go build -ldflags "-X 'main.AppLicense=$(LICENSE)' -X 'main.AppLicenseUrl=$(LICENSE_URL)' -X 'main.AppName=$@' -X 'main.AppNameFormal=$(APP_NAME_FORMAL)' -X 'main.AppRepo=$(REPO)' -X 'main.AppVersion=$(VERSION)'" ./main
