#!/bin/sh

# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

. "./cfg_build.sh"

BIN_NAME="courier"

go build -o "$BIN_DIR/$BIN_NAME" -ldflags "-X 'main.Version=$VERSION'" ./courier
