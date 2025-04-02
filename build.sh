#!/bin/sh

# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024 - 2025  Andy Frank Schoknecht

. "./cfg_install.sh"

for BINARY in $BINARIES; do
	./"_build_$BINARY.sh"
done
