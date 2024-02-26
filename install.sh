#!/bin/sh

# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

. "./cfg_build.sh"
. "./cfg_install.sh"

mkdir "/etc/hui"

for BINARY in $BINARIES; do
	./"build_$BINARY.sh"
	cp "$PKG_DIR/$BINARY" "$DESTDIR$PREFIX/bin"
	chmod 755 "$DESTDIR$PREFIX/bin/$BINARY"

	cp "$PKG_DIR/$BINARY.toml" "/etc/hui/"
done

cp "$PKG_DIR/common.toml" "/etc/hui/"
