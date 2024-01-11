// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"fmt"
	"os"
	
	"golang.org/x/term"
)

// Config temporarily hacked into
type Config struct {
	entry_menu_prefix   string
	entry_menu_postfix  string
	entry_shell_prefix  string
	entry_shell_postfix string
	header              string
	menus               []Menu
}

var g_cfg = Config{
	entry_menu_prefix:   "> [",
	entry_menu_postfix:  "]",
	entry_shell_prefix:  "> ",
	entry_shell_postfix: "",
	
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
			}},
			Entry{
			caption: "My final message...",
			content: EntryContent {
				ectype: ECT_SHELL,
				shell: "echo goodbye",
			}},
		},
	}},
}
// Config temporarily hacked into

type MenuStack []*Menu

func (self MenuStack) CurMenu() *Menu {
	return self[len(self) - 1]
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

func handle_input(active *bool, cursor *uint, mstack *MenuStack) {
	var input = make([]byte, 1)

	_, err := os.Stdin.Read(input)
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(input); i++ {
		handle_key(input[i], active, cursor, mstack)
	}
}

func handle_key(key byte, active *bool, cursor *uint, mstack *MenuStack) {
	switch key {
	case 'q':
		*active = false

	case 'j':
		if *cursor < uint(len(mstack.CurMenu().entries) - 1) {
			*cursor++
		}

	case 'k':
		if *cursor > 0 {
			*cursor--
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
	var mstack MenuStack = make(MenuStack, 1, 5)

	mstack[0] = &cfg.menus[len(cfg.menus) - 1]

	for active {
		fmt.Print(SEQ_CLEAR)

		draw_upper(cfg.header, mstack.CurMenu().title)
		draw_menu(cfg, *mstack.CurMenu(), cursor)

		canonical_state, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}

		handle_input(&active, &cursor, &mstack)

		term.Restore(int(os.Stdin.Fd()), canonical_state)
	}
}
