#!/bin/sh

# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

. "./cfg_build.sh"
. "./cfg_install.sh"

mkdir "$CFG_DESTDIR"

for BINARY in $BINARIES; do
	./"_build_$BINARY.sh"
	cp "$PKG_DIR/$BINARY" "$BIN_DESTDIR/$BINARY"
	chmod 755 "$BIN_DESTDIR/$BINARY"

	cp "$PKG_DIR/$BINARY.toml" "$CFG_DESTDIR/$BINARY.toml"
done

cp "$PKG_DIR/common.toml" "$CFG_DESTDIR"
