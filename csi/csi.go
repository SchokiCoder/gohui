// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package csi

import (
	"fmt"
)

const (
	SigInt  = "\003"
	SigTstp = "\004"

	Backspace = "\x7f"

	Clear = "\033[H\033[2J"
	CursorHide  = "\033[?25l"
	CursorShow  = "\033[?25h"
	CursorUp = "\033[A"
	CursorDown = "\033[B"
	CursorRight = "\033[C"
	CursorLeft = "\033[D"
	Home = "\x1b[H"
	Insert = "\033[2~"
	Delete = "\033[3~"
	PgUp = "\033[5~"
	PgDown = "\033[6~"
	End = "\x1b[F"
	FgDefault = "\033[39m"
	BgDefault = "\033[49m"
)

func SetCursor(x, y int) {
	fmt.Printf("\033[%v;%vH", y, x);
}
