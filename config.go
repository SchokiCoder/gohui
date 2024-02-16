// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"errors"
	"io"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type Config struct {
	KeyLeft           string
	KeyDown           string
	KeyUp             string
	KeyRight          string
	KeyExecute        string
	KeyQuit           string
	KeyCmdmode        string
	KeyCmdenter       string
	HeaderFg          FgColor
	HeaderBg          BgColor
	TitleFg           FgColor
	TitleBg           BgColor
	EntryFg           FgColor
	EntryBg           BgColor
	EntryHoverFg      FgColor
	EntryHoverBg      BgColor
	FeedbackFg        FgColor
	FeedbackBg        BgColor
	CmdlineFg         FgColor
	CmdlineBg         BgColor
	CmdlinePrefix     string
	FeedbackPrefix    string
	EntryMenuPrefix   string
	EntryMenuPostfix  string
	EntryShellPrefix  string
	EntryShellPostfix string
	Header            string
	Menus             map[string]Menu
}

func cfgFromFile() Config {
	type path struct {
		EnvVar string
		Core   string
	}

	var curPath string
	var i int
	var err error
	var f *os.File
	var found bool = false
	var paths = []path {
		path {"",                "/etc/hui/hui.toml"},
		path {"XDG_CONFIG_HOME", "/hui/hui.toml"},
		path {"HOME",            "/.config/hui/hui.toml"},
		path {"HOME",            "/.hui/hui.toml"},
		path {"",                "hui.toml"},
	}
	var prefix string
	var ret Config

	for _, v := range paths {
		if v.EnvVar != "" {
			prefix = os.Getenv(v.EnvVar)

			if prefix == "" {
				continue
			}

			curPath = fmt.Sprintf("%v%v", prefix, v.Core)
		} else {
			curPath = v.Core
		}

		f, err = os.Open(curPath)

		if errors.Is(err, os.ErrNotExist) {
			continue
		} else if err != nil {
			fmt.Fprintf(os.Stderr,
			            "Config file could not be opened: \"%v\", \"%v\"\n",
			            paths, err)
		} else {
			found = true
			break
		}
	}

	if found == false {
		panic("No config file could be found\n")
	}

	str, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("Config file could not be read: \"%v\", \"%v\"\n",
		                  paths[i], err))
	}

	err = toml.Unmarshal(str, &ret)
	if err != nil {
		panic(err)
	}

	ret.validate()

	return ret
}

func (c Config) validate() {
	for _, m := range c.Menus {
		for _, e := range m.Entries {
			if e.Shell == "" && e.Menu == "" {
				panic(fmt.Sprintf(
`Entry "%v" has no content.
Add a "Shell" value or a "Menu" value.`,
				                  e.Caption))
			}
			
			if e.Shell != "" && e.Menu != "" {
				panic(fmt.Sprintf(
`Entry "%v" has too much content.
Use only a "Shell" or a "Menu" value.`,
				                  e.Caption))
			} 
		}
	}
}
