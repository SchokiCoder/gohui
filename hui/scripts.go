// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
)

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
		ad.AcceptInput = false
		ad.Fb = "HUI End Feedback Msg"
	}

	putWordsIntoMyMouth := func() {
		ad.CmdLine.Active = true
		ad.CmdLine.RowIdx = -1
		ad.CmdLine.Content = "Surprise"
	}

	quit := func() {
		ad.Active = false
	}

	welcome := func() {
		ad.Fb = `Welcome, welcome to HUI.
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
