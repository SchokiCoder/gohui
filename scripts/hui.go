// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package scripts

import (
	"github.com/SchokiCoder/gohui/common"
)

var HuiFuncs = map[string]func(*common.HuiRuntime) {
	"Goodbye": huiGoodbye,
	"PutWordsIntoMyMouth": putWordsIntoMyMouth,
	"Quit": quit,
	"Welcome": huiWelcome,
}

func huiGoodbye(runtime *common.HuiRuntime) {
	runtime.AcceptInput = false
	runtime.Feedback = "Come back soon.\nWe have muffins!"
}

func putWordsIntoMyMouth(runtime *common.HuiRuntime) {
	runtime.CmdMode = true
	runtime.CmdLine = "Surprise"
}

/* Do not touch this!
 * Used by demo cfg.
 */
func quit(runtime *common.HuiRuntime) {
	runtime.Active = false
}

func huiWelcome(runtime *common.HuiRuntime) {
	runtime.Feedback = `Welcome, welcome to HUI.
You have chosen or been chosen to use one of our finest actively developed apps.
I have thought so much of HUI that i elected to pin it on Github.`
}
