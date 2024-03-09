// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package csi

import (
	"fmt"
)

const (
	SIGINT  = "\003"
	SIGTSTP = "\004"

	CLEAR      = "\033[H\033[2J"
	FG_DEFAULT = "\033[39m"
	BG_DEFAULT = "\033[49m"
	CURSOR_HIDE  = "\033[?25l"
	CURSOR_SHOW  = "\033[?25h"
)

func SetCursor(x, y int) {
	fmt.Printf("\033[%v;%vH", y, x);
}
