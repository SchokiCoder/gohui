// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
)

var huiCommands = map[string]func(string, *huiRuntime) string{
	"sh":  shCommand,
	"shs": shsCommand,
}

func shCommand(cmd string, _ *huiRuntime) string {
	return common.HandleShell(cmd)
}

func shsCommand(cmd string, _ *huiRuntime) string {
	return common.HandleShellSession(cmd)
}

var huiFuncs = map[string]func(*huiRuntime){
	"Goodbye":             huiGoodbye,
	"PutWordsIntoMyMouth": putWordsIntoMyMouth,
	"Quit":                quit,
	"Welcome":             huiWelcome,
}

func huiGoodbye(runtime *huiRuntime) {
	runtime.AcceptInput = false
	runtime.Feedback = "HUI End Feedback Msg"
}

func putWordsIntoMyMouth(runtime *huiRuntime) {
	runtime.CmdMode = true
	runtime.CmdLineRowIdx = -1
	runtime.CmdLine = "Surprise"
}

func quit(runtime *huiRuntime) {
	runtime.Active = false
}

func huiWelcome(runtime *huiRuntime) {
	runtime.Feedback = `Welcome, welcome to HUI.
You have chosen or been chosen to use one of our finest actively developed apps.
I have thought so much of HUI that i elected to pin it on Github.`
}
