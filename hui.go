// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"fmt"
	"os"
	"os/exec"
	
	"golang.org/x/term"
)

type MenuPath []string

func (mp MenuPath) CurMenu() string {
	return mp[len(mp) - 1]
}

const (
	SIGINT  = 3
	SIGTSTP = 4

	SEQ_CLEAR      = "\033[H\033[2J"
	SEQ_FG_DEFAULT = "\033[H\033[39m"
	SEQ_BG_DEFAULT = "\033[H\033[49m"
	SEQ_CRSR_HIDE  = "\033[?25l"
	SEQ_CRSR_SHOW  = "\033[?25h"
)

func draw_menu(cfg Config, cur_menu Menu, cursor uint) {
	var prefix, postfix string
	var fg FgColor
	var bg BgColor
	
	for i := uint(0); i < uint(len(cur_menu.entries)); i++ {
		switch cur_menu.entries[i].content.ectype {
		case ECT_MENU:
			prefix = cfg.entry_menu_prefix
			postfix = cfg.entry_menu_postfix
		
		case ECT_SHELL:
			prefix = cfg.entry_shell_prefix
			postfix = cfg.entry_shell_postfix
		}
		
		if i == cursor {
			fg = FgBlack
			bg = BgWhite
		} else {
			fg = FgWhite
			bg = BgBlack
		}
		
		fmt.Printf("%v%v%v%v%v\n",
		           fg,
		           bg,
		           prefix,
		           cur_menu.entries[i].caption,
		           postfix)
	}
	
	fmt.Printf("%v%v", SEQ_FG_DEFAULT, SEQ_BG_DEFAULT)
}

func draw_upper(header, title string) {
		fmt.Print(header, "\n")
		fmt.Print(title, "\n")
}

func handle_input(active *bool, cfg Config, cursor *uint, menu_path *MenuPath) {
	var input = make([]byte, 1)

	canonical_state, raw_err := term.MakeRaw(int(os.Stdin.Fd()))
	if raw_err != nil {
		panic(raw_err)
	}

	_, read_err := os.Stdin.Read(input)
	if read_err != nil {
		panic(read_err)
	}

	term.Restore(int(os.Stdin.Fd()), canonical_state)

	for i := 0; i < len(input); i++ {
		handle_key(input[i], active, cfg, cursor, menu_path)
	}
}

func handle_key(key       byte,
                active    *bool,
                cfg       Config,
                cursor    *uint,
                menu_path *MenuPath) {
	var cur_menu = cfg.menus[menu_path.CurMenu()]
	var cur_entry = &cur_menu.entries[*cursor]

	switch key {
	case 'q':
		*active = false

	case 'h':
		if len(*menu_path) > 1 {
			*menu_path = (*menu_path)[:len(*menu_path) - 1]
			*cursor = 0
		}

	case 'j':
		if *cursor < uint(len(cur_menu.entries) - 1) {
			*cursor++
		}

	case 'k':
		if *cursor > 0 {
			*cursor--
		}

	case 'l':
		if cur_entry.content.ectype == ECT_MENU {
			*menu_path = append(*menu_path, cur_entry.content.menu)
			*cursor = 0
		}

	case 'L':
		if cur_entry.content.ectype == ECT_SHELL {
			var cmd = exec.Command("sh", "-c", cur_entry.content.shell)
			var starterr = cmd.Start()
			if starterr != nil {
				// TODO no panic, give feedback
				panic(fmt.Sprintf("Could not start child process: %s", starterr))
			}

			var waiterr = cmd.Wait()
			if waiterr != nil {
				// TODO no panic, give feedback
				panic(fmt.Sprintf("Child error: %s", waiterr))
			}
		}

	case SIGINT: fallthrough
	case SIGTSTP:
		*active = false
	}
}

func set_cursor(x, y uint) {
	fmt.Print("\033[", y, ";", x, "H")
}

func main() {
	var active = true
	var cfg = g_cfg
	var cursor uint = 0
	var menu_path = make(MenuPath, 1, 8)

	_, main_menu_exists := cfg.menus["main"]

	if main_menu_exists {
		menu_path[0] = "main"
	} else {
		panic("main menu not found in config")
	}

	for active {
		fmt.Print(SEQ_CLEAR)

		draw_upper(cfg.header, cfg.menus[menu_path.CurMenu()].title)
		draw_menu(cfg, cfg.menus[menu_path.CurMenu()], cursor)

		handle_input(&active, cfg, &cursor, &menu_path)
	}
}
