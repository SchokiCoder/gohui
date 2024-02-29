// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"

	"fmt"
)

type HuiCfg struct {
	KeyExecute               string
	EntryMenuPrefix          string
	EntryMenuPostfix         string
	EntryShellPrefix         string
	EntryShellPostfix        string
	EntryShellSessionPrefix  string
	EntryShellSessionPostfix string
	EntryAlignment           string
	Header                   string
	EntryFg                  common.FgColor
	EntryBg                  common.BgColor
	EntryHoverFg             common.FgColor
	EntryHoverBg             common.BgColor
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

	for i, m := range c.Menus {
		if len(m.Entries) <= 0 {
			panic(fmt.Sprintf(`Menu "%v" has no entries.`, i))
		}

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
