// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/csi"

	"fmt"
)

type couConfig struct {
	Header           string
	PagerTitle       string
	ContentAlignment string
	GoStart          string
	GoQuit           string
	ContentFg        csi.FgColor
	ContentBg        csi.BgColor
}

func couConfigFromFile() couConfig {
	var ret couConfig

	common.AnyConfigFromFile(&ret, "courier.toml")
	ret.validateAlignments()
	if ret.GoStart != "" {
		validateGo(ret.GoStart)
	}
	if ret.GoQuit != "" {
		validateGo(ret.GoQuit)
	}

	return ret
}

func (c couConfig) validateAlignments() {
	common.ValidateAlignment(c.ContentAlignment)
}

func validateGo(fnName string) {
	_, fnExists := couFuncs[fnName]
	if fnExists == false {
		panic(fmt.Sprintf(`Courier Go function "%v" could not be found.`,
			          fnName))
	}
}
