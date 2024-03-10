// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/csi"

	"fmt"
)

type huiConfig struct {
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

func huiConfigFromFile() huiConfig {
	var ret huiConfig

	common.AnyConfigFromFile(&ret, "hui.toml")

	ret.validateAlignments()
	ret.validateMenus()
	if ret.GoStart != "" {
		validateGo(ret.GoStart)
	}
	if ret.GoQuit != "" {
		validateGo(ret.GoQuit)
	}

	return ret
}

func (c huiConfig) validateAlignments() {
	common.ValidateAlignment(c.EntryAlignment)
}

func (c huiConfig) validateMenus() {
	for i, m := range c.Menus {
		m.validate(i)
	}
}

func validateGo(fnName string) {
	_, fnExists := huiFuncs[fnName]
	if fnExists == false {
		panic(fmt.Sprintf(`Hui Go function "%v" could not be found.`,
			          fnName))
	}
}
