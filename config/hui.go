// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package config

import (
	"github.com/SchokiCoder/gohui/csi"

	"fmt"
)

type HuiCfg struct {
	Header                   string
	PagerTitle               string
	KeyExecute               string
	EntryMenuPrefix          string
	EntryMenuPostfix         string
	EntryShellPrefix         string
	EntryShellPostfix        string
	EntryShellSessionPrefix  string
	EntryShellSessionPostfix string
	EntryGoPrefix            string
	EntryGoPostfix           string
	EntryAlignment           string
	GoStart                  string
	GoQuit                   string
	EntryFg                  csi.FgColor
	EntryBg                  csi.BgColor
	EntryHoverFg             csi.FgColor
	EntryHoverBg             csi.BgColor
	Menus                    map[string]Menu
}

func HuiCfgFromFile() HuiCfg {
	var ret HuiCfg

	anyCfgFromFile(&ret, "hui.toml")

	ret.validateAlignments()
	ret.validateMenus()

	return ret
}

func (c HuiCfg) validateAlignments() {
	validateAlignment(c.EntryAlignment)
}

func (c HuiCfg) validateMenus() {
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

			if e.Go != "" {
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
