// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"fmt"
)

type Entry struct {
	Caption      string
	Menu         string
	Shell        string
	ShellSession string
	Go           string
}

func (e Entry) validate() {
	var numContent = 0

	if e.Shell != "" {
		numContent++
	}

	if e.ShellSession != "" {
		numContent++
	}
	
	if e.Menu != "" {
		numContent++
	}

	if e.Go != "" {
		validateGo(e.Go)
		numContent++
	}

	if numContent < 1 {
		panic(fmt.Sprintf(
`Entry "%v" has no content.
Add a "Shell" value, "ShellSession" value or a "Menu" value.`,
		                  e.Caption))
	} else if numContent > 1 {
		panic(fmt.Sprintf(
`Entry "%v" has too much content.
Use only a "Shell" or a "ShellSession" value or a "Menu" value.`,
		                  e.Caption))
	}
}

type Menu struct {
	Title   string
	Entries []Entry
}

func (m Menu) validate(menuIndex string) {
	if len(m.Entries) <= 0 {
		panic(fmt.Sprintf(`Menu "%v" has no entries.`, menuIndex))
	}

	for _, e := range m.Entries {
		e.validate()
	}
}
