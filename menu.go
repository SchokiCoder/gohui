// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

type EntryContent uint

const (
	menu EntryContent = 0
	shell EntryContent = 1
)

struct Entry {
	caption: string
	content: EntryContent
}

struct Menu {
	name: string
	title: string
	entries: []Entry
}
