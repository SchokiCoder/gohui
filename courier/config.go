// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
)

type CouCfg struct {
	ContentFg common.FgColor
	ContentBg common.BgColor
	Header    string
}

func cfgFromFile() CouCfg {
	var ret CouCfg

	common.AnyCfgFromFile(&ret, "courier.toml")

	return ret
}
