// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024 - 2025  Andy Frank Schoknecht

package common

import (
	"github.com/SchokiCoder/gohui/csi"

	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type cmdlineConfig struct {
	Alignment string
	Prefix    string
	Fg        csi.FgColor
	Bg        csi.BgColor
}

type feedbackConfig struct {
	Alignment string
	Prefix    string
	Fg        csi.FgColor
	Bg        csi.BgColor
}

type headerConfig struct {
	Alignment string
	Fg        csi.FgColor
	Bg        csi.BgColor
}

type pagerConfig struct {
	Fallback string
	Name     string `json:"IWillBeOverwrittenAnyway"`
}

type keysConfig struct {
	Left     string
	Down     string
	Up       string
	Right    string
	Quit     string
	Cmdmode  string
	Cmdenter string
}

type titleConfig struct {
	Alignment string
	Fg        csi.FgColor
	Bg        csi.BgColor
}

type ComConfig struct {
	Pager    pagerConfig
	Keys     keysConfig
	Header   headerConfig
	Title    titleConfig
	CmdLine  cmdlineConfig
	Feedback feedbackConfig
}

func AnyConfigFromFile(
	cfg interface{},
	cfgFileName string,
	customPath string,
) {
	type path struct {
		EnvVar string
		Core   string
	}

	var curPath string
	var i int
	var err error
	var f *os.File
	var found bool = false
	var paths = []path{
		path{"", customPath},
		path{"", "/etc/hui/"},
		path{"XDG_CONFIG_HOME", "/hui/"},
		path{"HOME", "/.config/hui/"},
		path{"HOME", "/.hui/"},
		path{"", ""},
	}
	var prefix string

	if len(customPath) == 0 {
		paths = paths[1:]
	}

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
				curPath, err)
		} else {
			found = true
			break
		}
	}

	if found == false {
		panic("No config file could be found")
	}

	str, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("Config file could not be read: \"%v\", \"%v\"",
			paths[i], err))
	}

	err = json.Unmarshal(str, cfg)
	if err != nil {
		panic(err)
	}
}

func ComConfigFromFile(
	customPath string,
) ComConfig {
	var ret ComConfig

	AnyConfigFromFile(&ret, "common.json", customPath)
	ret.validateAlignments()
	ret.validatePager()

	return ret
}

func ValidateAlignment(
	alignment string,
) {
	switch alignment {
	case "left":
	case "center":
	case "centered":
	case "right":

	default:
		panic(fmt.Sprintf(`Unknown alignment "%v" in config`, alignment))
	}
}

func (c ComConfig) validateAlignments(
) {
	ValidateAlignment(c.Header.Alignment)
	ValidateAlignment(c.Title.Alignment)
	ValidateAlignment(c.CmdLine.Alignment)
	ValidateAlignment(c.Feedback.Alignment)
}

func (c *ComConfig) validatePager(
) {
	var (
		pagers = [...]string{
			os.Getenv("PAGER"),
			c.Pager.Fallback,
			"less",
			"more",
		}
		path = os.Getenv("PATH")
		paths = strings.Split(path, ":")
	)

	// In case the fallback needs no prepend to be there
	paths = append(paths, "")

	for i := 0; i < len(pagers); i++ {
		if len(pagers[i]) <= 0 {
			continue
		}

		for j := 0; j < len(paths); j++ {
			_, err := os.Stat(filepath.Join(paths[j], pagers[i]))

			if errors.Is(err, os.ErrNotExist) == false {
				c.Pager.Name = pagers[i]
				return
			}
		}
	}

	panic(`No pager could not be found`)
}
