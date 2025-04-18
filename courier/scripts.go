// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024 - 2025  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
)

/* Warning: Setting courier's Feedback in the scripts can lead to recursion.
 * If the feedback is longer than one line, courier will be called for it.
 * If you for example have that happen on GoQuit, you will create a fork of
 * courier every time you try to close it.
 * On GoStart, you would absolutely FORK BOMB yourself.
 * At this point only `pkill courier` may help you.
 */

func getCmdMap(ad *appData) common.ScriptCmdMap {
	sh := func(cmd string) common.Feedback {
		return common.HandleShell(cmd)
	}

	shs := func(cmd string) common.Feedback {
		return common.HandleShellSession(cmd)
	}

	return common.ScriptCmdMap{
		"sh":  sh,
		"shs": shs,
	}
}

func getFnMap(ad *appData) common.ScriptFnMap {
	goodbye := func() {
		ad.CmdLine.Active = true
		ad.CmdLine.Content = "Courier End CmdLine Msg"
	}

	welcome := func() {
		ad.CmdLine.Content = "Eesterexs"
	}

	return common.ScriptFnMap{
		"Goodbye":             goodbye,
		"Welcome":             welcome,
	}
}
