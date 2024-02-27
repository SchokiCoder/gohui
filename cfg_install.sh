#!/bin/sh

# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

# Which programs to install
export BINARIES="hui courier"

# Where to install them
export BIN_DESTDIR="/usr/local/bin"
export CFG_DESTDIR="/etc/hui"
# uncomment below to install for current user only
#export BIN_DESTDIR="$HOME/.local/bin"
#export CFG_DESTDIR="$HOME/.config/hui"
