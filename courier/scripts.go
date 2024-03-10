// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

var couFuncs = map[string]func(*couRuntime) {
	"Goodbye": couGoodbye,
	"Welcome": couWelcome,
}

/* Warning: Setting courier's Feedback in the scripts can lead to recursion.
 * If the feedback is longer than one line, courier will be called for it.
 * If you for example have that happen on GoQuit, you will create a fork of
 * courier every time you try to close it.
 * On GoStart, you would absolutely FORK BOMB yourself.
 * At this point only `pkill courier` may help you.
 */

func couGoodbye(runtime *couRuntime) {
	runtime.CmdMode = true
	runtime.CmdLine = "Courier End CmdLine Msg"
}

func couWelcome(runtime *couRuntime) {
	runtime.CmdLine = "Eesterexs"
}