#!/bin/sh

# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

. "./cfg_build.sh"

BIN_NAME="hui"
BIN_NAME_FORMAL="House User Interface"

go build -o "$PKG_DIR/$BIN_NAME" \
	-ldflags "-X 'main.AppLicense=$LICENSE' -X 'main.AppLicenseUrl=$LICENSE_URL' -X 'main.AppName=$BIN_NAME' -X 'main.AppNameFormal=$BIN_NAME_FORMAL' -X 'main.AppRepo=$REPO' -X 'main.AppVersion=$VERSION'" \
	./hui
