// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"fmt"
)

const SEQ_CLEAR =      "\033[H\033[2J"
const SEQ_FG_DEFAULT = "\033[H\033[39m"
const SEQ_BG_DEFAULT = "\033[H\033[49m"
const SEQ_CRSR_HIDE =  "\033[?25l"
const SEQ_CRSR_SHOW =  "\033[?25h"

const HEADER = "header"

func set_cursor(x, y uint) {
	fmt.Print("\033[", y, ";", x, "H")
}

func main() {
	for {
		fmt.Print(SEQ_CLEAR)
		fmt.Print(SEQ_FG_DEFAULT, SEQ_BG_DEFAULT, HEADER)
		break
	}
}
