// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	KeyLeft            byte
	KeyDown            byte
	KeyUp              byte
	KeyRight           byte
	KeyExecute         byte
	KeyQuit            byte
	KeyCmdmode         byte
	KeyCmdenter        byte
	HeaderFg           FgColor
	HeaderBg           BgColor
	TitleFg            FgColor
	TitleBg            BgColor
	EntryFg            FgColor
	EntryBg            BgColor
	EntryHoverFg      FgColor
	EntryHoverBg      BgColor
	FeedbackFg         FgColor
	FeedbackBg         BgColor
	CmdlineFg          FgColor
	CmdlineBg          BgColor
	CmdlinePrefix      string
	FeedbackPrefix     string
	EntryMenuPrefix   string
	EntryMenuPostfix  string
	EntryShellPrefix  string
	EntryShellPostfix string
	Header              string
	Menus               map[string]Menu
}

func cfgFromFile() Config {
	var i int
	var err error
	var f *os.File
	var found bool = false
	var paths = []string {
		"/etc/hui/hui.toml",
		"~/.config/hui/hui.toml",
		"~/.hui/hui.toml",
		"./hui.toml",
	}
	var ret Config

	for i = 0; i < len(paths); i++ {
		f, err = os.Open(paths[i])
		
		if err != nil {
			found = true
			break
		} else if err != os.ErrNotExist {
			fmt.Fprintf(os.Stderr,
			            "Config file could not be read: %v\n",
			            paths[i])
		}
	}
	
	if found == false {
		panic("No config file could be found\n")
	}
	
	dec := toml.NewDecoder(f)
	dec.DisallowUnknownFields()
	
	err = dec.Decode(ret)
	if err != nil {
		panic(err)
	}
	
	return ret
}

