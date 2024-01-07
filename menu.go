// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

type EntryContentType uint

const (
	menu EntryContentType = 0
	shell EntryContentType = 1
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
	name string
	title string
	entries []Entry
}
