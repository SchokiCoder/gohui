# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024 - 2025  Andy Frank Schoknecht

PAGER_NAME_FORMAL:=Courier
APP_NAME_FORMAL  :=House User Interface
LICENSE          :=GPL-2.0-or-later
LICENSE_URL      :=https://www.gnu.org/licenses/gpl-2.0.html
REPO             :=https://github.com/SchokiCoder/gohui
VERSION          :=v1.4

BIN_DESTDIR:="/usr/local/bin"
CFG_DESTDIR:="/etc/hui"

# uncomment below to install for only current user instead
#BIN_DESTDIR:="$(HOME)/.local/bin"
#CFG_DESTDIR:="$(HOME)/.config/hui"

.PHONY: all build clean install vet

all: vet build

build: courier hui

clean:
	rm -f courier hui

vet:
	go vet ./pager
	go vet ./main

courier:
	go build -ldflags "-X 'main.AppLicense=$(LICENSE)' -X 'main.AppLicenseUrl=$(LICENSE_URL)' -X 'main.AppName=$@' -X 'main.AppNameFormal=$(PAGER_NAME_FORMAL)' -X 'main.AppRepo=$(REPO)' -X 'main.AppVersion=$(VERSION)'" ./pager

hui:
	go build -ldflags "-X 'main.AppLicense=$(LICENSE)' -X 'main.AppLicenseUrl=$(LICENSE_URL)' -X 'main.AppName=$@' -X 'main.AppNameFormal=$(APP_NAME_FORMAL)' -X 'main.AppRepo=$(REPO)' -X 'main.AppVersion=$(VERSION)'" ./main

install: build
	mkdir -p $(BIN_DESTDIR)
	cp courier $(BIN_DESTDIR)
	cp hui $(BIN_DESTDIR)
	mkdir -p $(CFG_DESTDIR)
	cp pkg/common.toml $(CFG_DESTDIR)
	cp pkg/courier.toml $(CFG_DESTDIR)
	cp pkg/hui.toml $(CFG_DESTDIR)

uninstall:
	rm -f $(BIN_DESTDIR)/courier $(BIN_DESTDIR)/hui $(CFG_DESTDIR)/*.toml
	rmdir $(CFG_DESTDIR)
