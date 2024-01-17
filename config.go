// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

// Config temporarily hacked into
type Config struct {
	entry_menu_prefix   string
	entry_menu_postfix  string
	entry_shell_prefix  string
	entry_shell_postfix string
	header              string
	menus               map[string]Menu
}

var g_cfg = Config{
	entry_menu_prefix:   "> [",
	entry_menu_postfix:  "]",
	entry_shell_prefix:  "> ",
	entry_shell_postfix: "",
	
	header: "Example config\n",

	menus: map[string]Menu {
		"main": Menu {
			title:
`Main Menu
---------`,
			entries: []Entry {
				Entry {
					caption: "Hello...",
					content: EntryContent {
						ectype: ECT_SHELL,
						shell: "echo world",
					},
				},

				Entry {
					caption: "Submenu",
					content: EntryContent {
						ectype: ECT_MENU,
						menu: "submenu",
					},
				},
			},
		},
	
		"submenu": Menu {
			title: "Submenu",
			entries: []Entry {
				Entry {
					caption: "Welcome to sub",
					content: EntryContent {
						ectype: ECT_SHELL,
						shell: "echo nothing",
					},
				},
			}, 
		},
	},
}
// Config temporarily hacked into
