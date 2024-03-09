// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package config

import (
	"github.com/SchokiCoder/gohui/csi"
)

type CouCfg struct {
	Header           string
	PagerTitle       string
	ContentAlignment string
	GoStart          string
	GoQuit           string
	ContentFg        csi.FgColor
	ContentBg        csi.BgColor
}

func CouCfgFromFile() CouCfg {
	var ret CouCfg

	anyCfgFromFile(&ret, "courier.toml")
	ret.validateAlignments()

	return ret
}

func (c CouCfg) validateAlignments() {
	validateAlignment(c.ContentAlignment)
}
