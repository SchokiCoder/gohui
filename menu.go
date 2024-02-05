// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

type EntryContentType uint

const (
	ECT_MENU EntryContentType = 0
	ECT_SHELL EntryContentType = 1
)

type EntryContent struct {
	EcType EntryContentType
	Menu   string
	Shell  string
}

type Entry struct {
	Caption string
	Content EntryContent
}

type Menu struct {
	Title   string
	Entries []Entry
}
