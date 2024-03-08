// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package scripts

import (
	"github.com/SchokiCoder/gohui/common"
)

var FuncMap = map[string]func(*common.HuiRuntime) {
	"PutWordsIntoMyMouth": PutWordsIntoMyMouth,
	"Quit": Quit,
}

func PutWordsIntoMyMouth(runtime *common.HuiRuntime) {
	runtime.CmdMode = true
	runtime.CmdLine = "Surprise"
}

/* Do not touch this!
 * Used by demo cfg.
 */
func Quit(runtime *common.HuiRuntime) {
	runtime.Active = false
}
