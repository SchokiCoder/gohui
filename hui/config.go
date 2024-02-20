// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"

	"fmt"
)

type HuiCfg struct {
	KeyExecute               string
	EntryFg                  common.FgColor
	EntryBg                  common.BgColor
	EntryHoverFg             common.FgColor
	EntryHoverBg             common.BgColor
	EntryMenuPrefix          string
	EntryMenuPostfix         string
	EntryShellPrefix         string
	EntryShellPostfix        string
	EntryShellSessionPrefix  string
	EntryShellSessionPostfix string
	Header                   string
	Menus                    map[string]Menu
}

func cfgFromFile() HuiCfg {
	var ret HuiCfg

	common.AnyCfgFromFile(&ret, "hui.toml")

	ret.validate()

	return ret
}

func (c HuiCfg) validate() {
	var numContent uint

	for _, m := range c.Menus {		
		for _, e := range m.Entries {
			numContent = 0

			if e.Shell != "" {
				numContent++
			}
			
			if e.ShellSession != "" {
				numContent++
			}
			
			if e.Menu != "" {
				numContent++
			}
		
			if numContent < 1 {
				panic(fmt.Sprintf(
`Entry "%v" has no content.
Add a "Shell" value, "ShellSession" value or a "Menu" value.`,
				                  e.Caption))
			} else if numContent > 1 {
				panic(fmt.Sprintf(
`Entry "%v" has too much content.
Use only a "Shell" or a "ShellSession" value or a "Menu" value.`,
				                  e.Caption))
			} 
		}
	}
}
