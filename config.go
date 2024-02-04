// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

// Config temporarily hacked into
type Config struct {
	key_left            byte
	key_down            byte
	key_up              byte
	key_right           byte
	key_execute         byte
	key_quit            byte
	key_cmdmode         byte
	key_cmdenter        byte
	header_fg           FgColor           
	header_bg           BgColor
	title_fg            FgColor
	title_bg            BgColor
	entry_fg            FgColor
	entry_bg            BgColor
	entry_hover_fg      FgColor
	entry_hover_bg      BgColor
	feedback_fg         FgColor
	feedback_bg         BgColor
	cmdline_fg          FgColor
	cmdline_bg          BgColor
	cmdline_prefix      string
	feedback_prefix     string
	entry_menu_prefix   string
	entry_menu_postfix  string
	entry_shell_prefix  string
	entry_shell_postfix string
	header              string
	menus               map[string]Menu
}

var g_cfg = Config{
	key_left: 'h',
	key_down: 'j',
	key_up: 'k',
	key_right: 'l',
	key_execute: 'L',
	key_quit: 'q',
	key_cmdmode: ':',
	key_cmdenter: '\r',

	header_fg: FgColor {
		active: false,
		r: 0,
		g: 0,
		b: 0,
	},

	header_bg: BgColor {
		active: false,
		r: 0,
		g: 0,
		b: 0,
	},

	title_fg: FgColor {
		active: false,
		r: 0,
		g: 0,
		b: 0,
	},

	title_bg: BgColor {
		active: false,
		r: 0,
		g: 0,
		b: 0,
	},

	entry_fg: FgColor {
		active: false,
		r: 0,
		g: 0,
		b: 0,
	},

	entry_bg: BgColor {
		active: false,
		r: 0,
		g: 0,
		b: 0,
	},

	entry_hover_fg: FgColor {
		active: true,
		r: 0,
		g: 0,
		b: 0,
	},

	entry_hover_bg: BgColor {
		active: true,
		r: 255,
		g: 255,
		b: 255,
	},

	feedback_fg: FgColor {
		active: true,
		r: 175,
		g: 175,
		b: 175,
	},

	feedback_bg: BgColor {
		active: false,
		r: 0,
		g: 0,
		b: 0,
	},

	cmdline_fg: FgColor {
		active: false,
		r: 0,
		g: 0,
		b: 0,
	},

	cmdline_bg: BgColor {
		active: false,
		r: 0,
		g: 0,
		b: 0,
	},

	entry_menu_prefix:   "> [",
	entry_menu_postfix:  "]",
	entry_shell_prefix:  "> ",
	entry_shell_postfix: "",
	cmdline_prefix: ":",
	feedback_prefix: "<",

	header: "Example config\n",

	menus: map[string]Menu {
		"main": Menu {
			title:
`Main Menu
---------`,
			entries: []Entry {
				Entry {
					caption: "echo to temp",
					content: EntryContent {
						ectype: ECT_SHELL,
						shell: "echo gotest >> ~/temp",
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
					caption: "print short",
					content: EntryContent {
						ectype: ECT_SHELL,
						shell: "echo short",
					},
				},
				Entry {
					caption: "print long",
					content: EntryContent {
						ectype: ECT_SHELL,
						shell: "echo loooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong",
					},
				},
				
				Entry {
					caption: "print too long",
					content: EntryContent {
						ectype: ECT_SHELL,
						shell: "echo looooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooong",
					},
				},
			}, 
		},
	},
}
// Config temporarily hacked into
