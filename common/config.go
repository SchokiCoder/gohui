// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package common

import (
	"errors"
	"io"
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"strings"
)

type ComCfg struct {
	AppPager                 string
	KeyLeft                  string
	KeyDown                  string
	KeyUp                    string
	KeyRight                 string
	KeyQuit                  string
	KeyCmdmode               string
	KeyCmdenter              string
	CmdlinePrefix            string
	FeedbackPrefix           string
	HeaderAlignment          string
	TitleAlignment           string
	HeaderFg                 FgColor
	HeaderBg                 BgColor
	TitleFg                  FgColor
	TitleBg                  BgColor
	FeedbackFg               FgColor
	FeedbackBg               BgColor
	CmdlineFg                FgColor
	CmdlineBg                BgColor
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
	ret.validateAlignments()
	ret.validatePager()

	return ret
}

func ValidateAlignment(alignment string) {
	switch alignment {
	case "left":
	case "center":
	case "centered":
	case "right":

	default:
		panic(fmt.Sprintf(`Unknown alignment "%v" in config.`, alignment))
	}
}

func (c ComCfg) validateAlignments() {
	ValidateAlignment(c.HeaderAlignment)
	ValidateAlignment(c.TitleAlignment)
}

func (c ComCfg) validatePager() {
	var pagerExists = false
	var path = os.Getenv("PATH")

	_, err := os.Stat(c.AppPager)
	if errors.Is(err, os.ErrNotExist) == false {
		return
	}

	for _, v := range strings.Split(path, ":") {
		_, err := os.Stat(fmt.Sprintf("%v/%v", v, c.AppPager))
		if errors.Is(err, os.ErrNotExist) == false {
			pagerExists = true
			break
		} 
	}

	if pagerExists == false {
		panic(fmt.Sprintf(`The configured pager "%v" can not be found.`,
		                  c.AppPager))
	}
}
