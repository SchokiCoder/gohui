// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package csi

import (
	"fmt"
)

const (
	SIGINT  = "\003"
	SIGTSTP = "\004"

	BACKSPACE = "\x7f"

	CLEAR = "\033[H\033[2J"
	CURSOR_HIDE  = "\033[?25l"
	CURSOR_SHOW  = "\033[?25h"
	CURSOR_UP = "\033[A"
	CURSOR_DOWN = "\033[B"
	CURSOR_RIGHT = "\033[C"
	CURSOR_LEFT = "\033[D"
	HOME = "\033[1~"
	DELETE = "\033[3~"
	END = "\033[4~"
	FG_DEFAULT = "\033[39m"
	BG_DEFAULT = "\033[49m"
)

func SetCursor(x, y int) {
	fmt.Printf("\033[%v;%vH", y, x);
}
