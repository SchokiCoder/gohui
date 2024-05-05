// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/csi"

	"fmt"
)

type entry struct {
	Caption      string
	Menu         string
	Shell        string
	ShellSession string
	Go           string
}

func (e entry) validate(fnMap common.ScriptFnMap) {
	var numContent = 0

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
		validateGo(e.Go, fnMap)
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

type menu struct {
	Title   string
	Entries []entry
}

func (m menu) validate(fnMap common.ScriptFnMap, menuIndex string) {
	if len(m.Entries) <= 0 {
		panic(fmt.Sprintf(`Menu "%v" has no entries.`, menuIndex))
	}

	for _, e := range m.Entries {
		e.validate(fnMap)
	}
}

type eventsConfig struct {
	Start string
	Quit  string
}

type entryConfig struct {
	Alignment           string
	MenuPrefix          string
	MenuPostfix         string
	ShellPrefix         string
	ShellPostfix        string
	ShellSessionPrefix  string
	ShellSessionPostfix string
	GoPrefix            string
	GoPostfix           string
	Fg                  csi.FgColor
	Bg                  csi.BgColor
	HoverFg             csi.FgColor
	HoverBg             csi.BgColor
}

type keysConfig struct {
	Execute string
}

type pagerConfig struct {
	Title string
}

type huiConfig struct {
	Header string
	Pager  pagerConfig
	Keys   keysConfig
	Entry  entryConfig
	Events eventsConfig
	Menus  map[string]menu
}

func huiConfigFromFile(fnMap common.ScriptFnMap, rt *huiRuntime) huiConfig {
	var ret huiConfig

	common.AnyConfigFromFile(&ret, "hui.toml")

	ret.validateAlignments()
	ret.validateMenus(fnMap)
	if ret.Events.Start != "" {
		validateGo(ret.Events.Start, fnMap)
	}
	if ret.Events.Quit != "" {
		validateGo(ret.Events.Quit, fnMap)
	}

	return ret
}

func (c huiConfig) validateAlignments() {
	common.ValidateAlignment(c.Entry.Alignment)
}

func (c huiConfig) validateMenus(fnMap common.ScriptFnMap) {
	for i, m := range c.Menus {
		m.validate(fnMap, i)
	}
}

func validateGo(fnName string, fnMap common.ScriptFnMap) {
	_, fnExists := fnMap[fnName]
	if fnExists == false {
		panic(fmt.Sprintf(`Hui Go function "%v" could not be found.`,
			fnName))
	}
}
