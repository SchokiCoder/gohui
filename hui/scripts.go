// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
)

func getCmdMap(rt *huiRuntime) common.ScriptCmdMap {
	sh := func(cmd string) string {
		return common.HandleShell(cmd)
	}

	shs := func(cmd string) string {
		return common.HandleShellSession(cmd)
	}

	return common.ScriptCmdMap{
		"sh":  sh,
		"shs": shs,
	}
}

func getFnMap(rt *huiRuntime) common.ScriptFnMap {
	goodbye := func() {
		rt.AcceptInput = false
		rt.Feedback = "HUI End Feedback Msg"
	}

	putWordsIntoMyMouth := func() {
		rt.CmdMode = true
		rt.CmdLineRowIdx = -1
		rt.CmdLine = "Surprise"
	}

	quit := func() {
		rt.Active = false
	}

	welcome := func() {
		rt.Feedback = `Welcome, welcome to HUI.
You have chosen or been chosen to use one of our finest actively developed apps.
I have thought so much of HUI that i elected to pin it on Github.`
	}

	return common.ScriptFnMap{
		"Goodbye":             goodbye,
		"PutWordsIntoMyMouth": putWordsIntoMyMouth,
		"Quit":                quit,
		"Welcome":             welcome,
	}
}
