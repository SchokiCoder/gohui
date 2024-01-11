// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"fmt"
)

type FgColor struct {
	active bool
	r, g, b uint
}

func (c FgColor) String() string {
	if c.active {
		return fmt.Sprintf("\033[38;2;%v;%v;%vm", c.r, c.g, c.b)
	} else {
		return SEQ_FG_DEFAULT
	}
}

type BgColor struct {
	active bool
	r, g, b uint
}

func (c BgColor) String() string {
	if c.active {
		return fmt.Sprintf("\033[48;2;%v;%v;%vm", c.r, c.g, c.b)
	} else {
		return SEQ_BG_DEFAULT
	}
}

var BgWhite = BgColor {
	true,
	255,
	255,
	255,
}

var FgWhite = FgColor {
	true,
	255,
	255,
	255,
}

var BgBlack = BgColor {
	true,
	0,
	0,
	0,
}

var FgBlack = FgColor {
	true,
	0,
	0,
	0,
}
