// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

type Entry struct {
	Caption string
	Menu  string
	Shell string
}

type Menu struct {
	Title   string
	Entries []Entry
}
