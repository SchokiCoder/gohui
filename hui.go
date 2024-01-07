// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"bufio"
	"fmt"
	"os"
)

// Config temporarily as constants
const (
	HEADER = "Example config\n"

	MENUS = [Menu{
		name: "main",
		title: `Main Menu\n
		        ---------`,
		[Entry{
			caption: "Hello...",
			shell: "echo world",
		}]Entry,
	}]Menu
)
// Config temporarily as constants

const (
	SEQ_CLEAR      = "\033[H\033[2J"
	SEQ_FG_DEFAULT = "\033[H\033[39m"
	SEQ_BG_DEFAULT = "\033[H\033[49m"
	SEQ_CRSR_HIDE  = "\033[?25l"
	SEQ_CRSR_SHOW  = "\033[?25h"
)

func set_cursor(x, y uint) {
	fmt.Print("\033[", y, ";", x, "H")
}

func main() {
	var active = true
	var scanner_in = bufio.NewScanner(os.Stdin)
	var writer_err = bufio.NewWriter(os.Stderr)

	for active {
		fmt.Print(SEQ_CLEAR)

		fmt.Print(HEADER)

		if scanner_in.Scan() == false {
			fmt.Fprint(writer_err, "end of input\n")
			active = false
		}

		switch scanner_in.Text() {
		case "q":
			active = false
			break
		}
	}
}
