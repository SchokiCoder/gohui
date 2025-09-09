// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024 - 2025  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/csi"

	"errors"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
)

type appData struct {
	common.ComAppData
	Content           string
	CouCfg            couConfig
	Scroll            int
	Title             string
}

var (
	AppLicense    string
	AppLicenseUrl string
	AppName       string
	AppNameFormal string
	AppRepo       string
	AppVersion    string
)

const HELP = `Usage: courier [OPTIONS] FILE

Small customizable pager, written for and usually distributed with hui.

Options:

    -a --about
        prints program name, version, license and repository information then exits

    -c --config
        takes an argument as additional path for config dir search

    -h --help
        prints this message then exits

    -v --version
        prints version information then exits

Default keybinds:

    q, h
        quit the program

    j
        go down

    k
        go up

    :
        enter the internal command line

Internal commands:

    q quit exit
        quit the program

    *number*
        when given a positive number, it is used as a line number to scroll to

Environmental variables:

    PAGER
        sets the used pager, in case feedback exceeds one line of length

    PAGERTITLE
        sets the title
`

func drawContent(
	contentLines []string,
	contentHeight int,
	ad appData,
	termW int,
) {
	var (
		drawRange int = ad.Scroll + contentHeight
	)

	if drawRange > len(contentLines) {
		drawRange = len(contentLines)
	}

	for _, v := range contentLines[ad.Scroll:drawRange] {
		common.Cprinta(ad.CouCfg.Content.Alignment,
			ad.CouCfg.Content.Fg,
			ad.CouCfg.Content.Bg,
			termW,
			v)
	}
}

func handleArgs(
	cfgPath *string,
) (string, bool) {
	var filepath string

	if len(os.Args) < 2 {
		panic("Not enough arguments given")
	}

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-a":
			fallthrough
		case "--about":
			common.PrintAbout(AppLicense,
				AppLicenseUrl,
				AppName,
				AppNameFormal,
				AppRepo,
				AppVersion)
			return "", false

		case "-c":
			fallthrough
		case "--config":
			*cfgPath = os.Args[i + 1]
			i++

		case "-h":
			fallthrough
		case "--help":
			fmt.Printf(HELP)
			return "", false

		case "-v":
			fallthrough
		case "--version":
			common.PrintVersion(AppName, AppVersion)
			return "", false

		default:
			if os.Args[i][0] == '-' {
				panic(`Unknown argument "` + os.Args[i] + `"`)
			}

			filepath = os.Args[i]
		}
	}

	if filepath = "" {
		panic("No filepath has been given")
	}

	return filepath, true
}

func handleInput(
	cmdMap common.ScriptCmdMap,
	contentHeight int,
	contentLineCount int,
	ad *appData,
) {
	var (
		canonicalState *term.State
		err error
		input string
		rawInput = make([]byte, 4)
		rawInputLen int
	)

	if ad.AcceptInput == false {
		return
	}

	canonicalState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Switching to raw mode failed:\n%v", err))
	}

	rawInputLen, err = os.Stdin.Read(rawInput)
	if err != nil {
		panic(fmt.Sprintf("Reading from stdin failed:\n%v", err))
	}
	input = string(rawInput[0:rawInputLen])

	term.Restore(int(os.Stdin.Fd()), canonicalState)

	handleKey(string(input), cmdMap, contentHeight, contentLineCount, ad)
}

func handleKey(
	key string,
	cmdMap common.ScriptCmdMap,
	contentHeight, contentLineCount int,
	ad *appData,
) {
	if ad.CmdLine.Active {
		common.HandleKeyCmdline(key,
			&ad.Active,
			&ad.CmdLine,
			cmdMap,
			&ad.ComCfg,
			contentLineCount,
			&ad.Scroll,
			&ad.Fb)
		return
	}

	switch key {
	case csi.CursorUp:
		fallthrough
	case ad.ComCfg.Keys.Up:
		if ad.Scroll > 0 {
			ad.Scroll--
		}

	case csi.CursorDown:
		fallthrough
	case ad.ComCfg.Keys.Down:
		if (ad.Scroll + 1) < contentLineCount {
			ad.Scroll++
		}

	case ad.ComCfg.Keys.Cmdmode:
		ad.CmdLine.Active = true
		fmt.Printf(csi.CursorShow)

	case csi.PgUp:
		if ad.Scroll-contentHeight < 0 {
			ad.Scroll = 0
		} else {
			ad.Scroll -= contentHeight
		}

	case csi.PgDown:
		if ad.Scroll+contentHeight >= contentLineCount {
			ad.Scroll = contentLineCount - 1
		} else {
			ad.Scroll += contentHeight
		}

	case csi.Home:
		ad.Scroll = 0

	case csi.End:
		ad.Scroll = contentLineCount - 1

	case ad.ComCfg.Keys.Quit:
		fallthrough
	case csi.CursorLeft:
		fallthrough
	case ad.ComCfg.Keys.Left:
		fallthrough
	case csi.SigInt:
		fallthrough
	case csi.SigTstp:
		ad.Active = false
	}
}

func readfile(
	filepath string,
) string {
	f, err := os.Open(filepath)
	defer f.Close()

	if errors.Is(err, os.ErrNotExist) {
		panic(fmt.Sprintf("File \"%v\" could not be found:\n%v",
			filepath,
			err))
	} else if err != nil {
		panic(fmt.Sprintf("File \"%v\" could not be opened:\n%v",
			filepath,
			err))
	}

	ret, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("File \"%v\" could not be read:\n%v",
			filepath,
			err))
	}

	return string(ret)
}

func tick(
	cmdMap common.ScriptCmdMap,
	ad *appData,
) {
	var contentLines []string
	var contentHeight int
	var err error
	var headerLines []string
	var lower string
	var termH, termW int
	var titleLines []string

	fmt.Print(csi.Clear)
	termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Could not get term size:\n%v", err))
	}

	headerLines = common.SplitByLines(termW, ad.CouCfg.Header)
	titleLines = common.SplitByLines(termW, ad.Title)
	contentLines = common.SplitByLines(termW, ad.Content)
	lower = common.GenerateLower(ad.CmdLine.Content,
		ad.CmdLine.Active,
		ad.ComCfg,
		&ad.Fb,
		ad.CouCfg.Pager.Title,
		termW)

	common.DrawUpper(ad.ComCfg, headerLines, termW, titleLines)

	contentHeight = termH -
		len(common.SplitByLines(termW, ad.CouCfg.Header)) -
		1 -
		len(common.SplitByLines(termW, ad.Title)) -
		1

	drawContent(contentLines, contentHeight, *ad, termW)

	csi.SetCursor(1, termH)
	fmt.Printf("%v", lower)
	csi.SetCursorAligned(ad.ComCfg.CmdLine.Alignment,
		(len(ad.ComCfg.CmdLine.Prefix) + len(ad.CmdLine.Content)),
		termW,
		(len(ad.ComCfg.CmdLine.Prefix) + ad.CmdLine.Cursor + 1),
		termH)

	handleInput(cmdMap, contentHeight, len(contentLines), ad)
}

func main(
) {
	var (
		ad       appData
		cfgPath  string
		cmdMap   common.ScriptCmdMap
		filepath string
		fnMap    common.ScriptFnMap
	)

	filepath, ad.Active = handleArgs(&cfgPath)
	if ad.Active == false {
		return
	}

	ad.Content = readfile(filepath)

	ad.Title = os.Getenv("PAGERTITLE")

	cmdMap = getCmdMap(&ad)
	fnMap = getFnMap(&ad)
	ad.ComAppData = common.NewComAppData(cfgPath)
	ad.CouCfg = couConfigFromFile(cfgPath, fnMap)

	if ad.CouCfg.Events.Start != "" {
		fnMap[ad.CouCfg.Events.Start]()
	}

	fmt.Printf(csi.CursorHide)
	defer fmt.Printf(csi.CursorShow)
	defer fmt.Printf("%v%v\n", csi.FgDefault, csi.BgDefault)

	for ad.Active {
		tick(cmdMap, &ad)
	}

	if ad.CouCfg.Events.Quit != "" {
		fnMap[ad.CouCfg.Events.Quit]()
		tick(cmdMap, &ad)
	}
}
