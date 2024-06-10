// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./geninfo.go
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

func newAppData(fnMap common.ScriptFnMap) appData {
	return appData {
		ComAppData: common.NewComAppData(),
		Content:    "",
		CouCfg:     couConfigFromFile(fnMap),
		Scroll:     0,
		Title:      "",
	}
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

    -h --help
        prints this message then exits

    -t --title TITLE
        takes an argument and prints given string as title below the header

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
`

func drawContent(contentLines []string,
	contentHeight int,
	ad appData,
	termW int) {
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

func handleArgs(title *string) (string, bool) {
	var err error
	var f *os.File
	var nextIsTitle = false
	var path string

	if len(os.Args) < 2 {
		panic("Not enough arguments given.")
	}

	for _, v := range os.Args[1:] {
		switch v {
		case "-v":
			fallthrough
		case "--version":
			common.PrintVersion(AppName, AppVersion)
			return "", false

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

		case "-h":
			fallthrough
		case "--help":
			fmt.Printf(HELP)
			return "", false

		case "-t":
			fallthrough
		case "--title":
			nextIsTitle = true

		default:
			if nextIsTitle {
				*title = v
				nextIsTitle = false
			} else {
				path = v
			}
		}
	}

	f, err = os.Open(path)
	defer f.Close()

	if errors.Is(err, os.ErrNotExist) {
		panic(fmt.Sprintf("File \"%v\" could not be found: %v",
			path,
			err))
	} else if err != nil {
		panic(fmt.Sprintf("File \"%v\" could not be opened: %v",
			path,
			err))
	}

	ret, err := io.ReadAll(f)
	if err != nil {
		panic(fmt.Sprintf("File \"%v\" could not be read: %v",
			path,
			err))
	}

	return string(ret), true
}

func handleInput(cmdMap common.ScriptCmdMap,
	contentHeight int,
	contentLineCount int,
	ad *appData) {
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
		panic(fmt.Sprintf("Switching to raw mode failed: %v", err))
	}

	rawInputLen, err = os.Stdin.Read(rawInput)
	if err != nil {
		panic(fmt.Sprintf("Reading from stdin failed: %v", err))
	}
	input = string(rawInput[0:rawInputLen])

	term.Restore(int(os.Stdin.Fd()), canonicalState)

	handleKey(string(input), cmdMap, contentHeight, contentLineCount, ad)
}

func handleKey(key string,
	cmdMap common.ScriptCmdMap,
	contentHeight, contentLineCount int,
	ad *appData) {

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

func tick(cmdMap common.ScriptCmdMap, ad *appData) {
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
		panic(fmt.Sprintf("Could not get term size: %v", err))
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

func main() {
	var cmdMap common.ScriptCmdMap
	var fnMap common.ScriptFnMap
	var ad appData

	cmdMap = getCmdMap(&ad)
	fnMap = getFnMap(&ad)
	ad = newAppData(fnMap)

	ad.Content, ad.Active = handleArgs(&ad.Title)
	if ad.Active == false {
		return
	}

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
