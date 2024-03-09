// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
)

type CouCfg struct {
	Header           string
	PagerTitle       string
	ContentAlignment string
	GoStart          string
	GoQuit           string
	ContentFg        common.FgColor
	ContentBg        common.BgColor
}

func cfgFromFile() CouCfg {
	var ret CouCfg

	common.AnyCfgFromFile(&ret, "courier.toml")
	ret.validateAlignments()

	return ret
}

func (c CouCfg) validateAlignments() {
	common.ValidateAlignment(c.ContentAlignment)
}
