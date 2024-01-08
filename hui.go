// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"bufio"
	"fmt"
	"os"
)

// Config temporarily hacked into
type Config struct {
	header string
	menus  []Menu
}

var cfg = Config{
	header: "Example config\n",

	menus: []Menu {Menu{
		name: "main",
		title:
`Main Menu
---------`,
		entries: []Entry {Entry{
			caption: "Hello...",
			content: EntryContent {
				ectype: ECT_SHELL,
				shell: "echo world",
			},
		}},
	}},
}
// Config temporarily hacked into

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
	var cfg = cfg
	var cur_menu *Menu
	var scanner_in = bufio.NewScanner(os.Stdin)
	var writer_err = bufio.NewWriter(os.Stderr)
	var menu_path = []*Menu {&cfg.menus[len(cfg.menus) - 1]}

	for active {
		cur_menu = menu_path[len(menu_path) - 1] 

		fmt.Print(SEQ_CLEAR)

		fmt.Print(cfg.header, "\n")
		fmt.Print(cur_menu.title, "\n")

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
