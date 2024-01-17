// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

type EntryContentType uint

const (
	ECT_MENU EntryContentType = 0
	ECT_SHELL EntryContentType = 1
)

type EntryContent struct {
	ectype EntryContentType
	menu   string
	shell  string
}

type Entry struct {
	caption string
	content EntryContent
}

type Menu struct {
	title   string
	entries []Entry
}
