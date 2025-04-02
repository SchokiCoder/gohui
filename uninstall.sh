#!/bin/sh

# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024 - 2025  Andy Frank Schoknecht

. "./cfg_install.sh"

for BINARY in $BINARIES; do
	rm "$BIN_DESTDIR/$BINARY"
	rm "$CFG_DESTDIR/$BINARY.toml"
done

rm "$CFG_DESTDIR/common.toml"
rmdir "$CFG_DESTDIR"
