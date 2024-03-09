// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package scripts

import (
	"github.com/SchokiCoder/gohui/common"
)

var CouFuncs = map[string]func(*common.CouRuntime) {
	"Goodbye": couGoodbye,
	"Welcome": couWelcome,
}

/* Warning: Setting courier's Feedback in the scripts can lead to recursion.
 * If the feedback is longer than one line, courier will be called for it.
 * If you for example have that happen on GoQuit, you will create a fork of
 * courier every time you try to close it.
 * At this point only `pkill courier` may help you.
 */

func couGoodbye(runtime *common.CouRuntime) {
	runtime.CmdMode = true
	runtime.CmdLine = "Come back soon. We have muffins too!"
}

func couWelcome(runtime *common.CouRuntime) {
	runtime.CmdLine = "Eesterexs"
}
