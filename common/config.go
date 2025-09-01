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
	"slices"
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
	EnvVars string
	Pager   string
	Flags   string
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
	Pagers   []pagerConfig
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
				"Config file \"%v\" could not be opened:\n%v\n",
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
		panic(fmt.Sprintf("Config file \"%v\" could not be read:\n%v",
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
	ret.validatePagers()

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

func (c *ComConfig) validatePagers(
) {
	var (
		pagerFound bool
		path = os.Getenv("PATH")
		paths = strings.Split(path, ":")
	)

	paths = append(paths, "") // for local dir

	c.Pagers = append(
		[]pagerConfig{
			pagerConfig{
				EnvVars: "",
				Pager: os.Getenv("PAGER"),
				Flags: ""},
		},
		c.Pagers...)

	for i := 0; i < len(c.Pagers); i++ {
		pagerFound = false

		for j := 0; j < len(paths); j++ {
			_, err := os.Stat(
				filepath.Join(
					paths[j],
					c.Pagers[i].Pager))

			if errors.Is(err, os.ErrNotExist) == false {
				pagerFound = true
				break
			}
		}

		if !pagerFound || c.Pagers[i].Pager == "" {
			c.Pagers = slices.Delete(c.Pagers, i, i + 1)
			i--
		}
	}

	if len(c.Pagers) <= 0 {
		panic(`No pager could be found`)
	}
}
