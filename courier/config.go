// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/csi"

	"fmt"
)

type contentConfig struct {
	Alignment string
	Fg        csi.FgColor
	Bg        csi.BgColor
}

type eventsConfig struct {
	Start string
	Quit  string
}

type pagerConfig struct {
	Title string
}

type couConfig struct {
	Header  string
	Pager   pagerConfig
	Content contentConfig
	Events  eventsConfig
}

func couConfigFromFile() couConfig {
	var ret couConfig

	common.AnyConfigFromFile(&ret, "courier.toml")
	ret.validateAlignments()
	if ret.Events.Start != "" {
		validateGo(ret.Events.Start)
	}
	if ret.Events.Quit != "" {
		validateGo(ret.Events.Quit)
	}

	return ret
}

func (c couConfig) validateAlignments() {
	common.ValidateAlignment(c.Content.Alignment)
}

func validateGo(fnName string) {
	_, fnExists := couFuncs[fnName]
	if fnExists == false {
		panic(fmt.Sprintf(`Courier Go function "%v" could not be found.`,
			fnName))
	}
}
