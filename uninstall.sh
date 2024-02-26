#!/bin/sh

# SPDX-License-Identifier: GPL-2.0-or-later
# Copyright (C) 2024  Andy Frank Schoknecht

. "./cfg_install.sh"

rm "$DESTDIR$PREFIX/bin/hui" "$DESTDIR$PREFIX/bin/courier"
rm /etc/hui/*
rmdir /etc/hui
