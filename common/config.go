// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package common

import (
	"errors"
	"io"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
)

type ComCfg struct {
	KeyLeft                  string
	KeyDown                  string
	KeyUp                    string
	KeyRight                 string
	KeyQuit                  string
	KeyCmdmode               string
	KeyCmdenter              string
	HeaderFg                 FgColor
	HeaderBg                 BgColor
	TitleFg                  FgColor
	TitleBg                  BgColor
	FeedbackFg               FgColor
	FeedbackBg               BgColor
	CmdlineFg                FgColor
	CmdlineBg                BgColor
	CmdlinePrefix            string
	FeedbackPrefix           string
}

func AnyCfgFromFile(cfg interface{}, cfgFileName string) {
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
		path {"",                "/etc/hui/"},
		path {"XDG_CONFIG_HOME", "/hui/"},
		path {"HOME",            "/.config/hui/"},
		path {"HOME",            "/.hui/"},
		path {"",                ""},
	}
	var prefix string

	for _, v := range paths {
		if v.EnvVar != "" {
			prefix = os.Getenv(v.EnvVar)

			if prefix == "" {
				continue
			}

			curPath = fmt.Sprintf("%v%v%v", prefix, v.Core, cfgFileName)
		} else {
			curPath = fmt.Sprintf("%v%v", v.Core, cfgFileName)
		}

		f, err = os.Open(curPath)
		defer f.Close()

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

	err = toml.Unmarshal(str, cfg)
	if err != nil {
		panic(err)
	}
}

func CfgFromFile() ComCfg {
	var ret ComCfg

	AnyCfgFromFile(&ret, "common.toml")

	return ret
}
