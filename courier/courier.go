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

type couRuntime struct {
	AcceptInput   bool
	Active        bool
	CmdLine       string
	CmdLineCursor int
	CmdLineInsert bool
	CmdLineRowIdx int
	CmdLineRows   [common.CmdlineMaxRows]string
	CmdMode       bool
	ComCfg        common.ComConfig
	Content       string
	CouCfg        couConfig
	Scroll        int
	Feedback      string
	Title         string
}

func newCouRuntime(fnMap common.ScriptFnMap) couRuntime {
	return couRuntime{
		AcceptInput:   true,
		Active:        true,
		CmdLine:       "",
		CmdLineCursor: 0,
		CmdLineInsert: false,
		CmdLineRowIdx: -1,
		CmdMode:       false,
		ComCfg:        common.ComConfigFromFile(),
		Content:       "",
		CouCfg:        couConfigFromFile(fnMap),
		Scroll:        0,
		Feedback:      "",
		Title:         "",
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
	rt couRuntime,
	termW int) {
	var (
		drawRange int = rt.Scroll + contentHeight
	)

	if drawRange > len(contentLines) {
		drawRange = len(contentLines)
	}

	for _, v := range contentLines[rt.Scroll:drawRange] {
		common.Cprinta(rt.CouCfg.Content.Alignment,
			rt.CouCfg.Content.Fg,
			rt.CouCfg.Content.Bg,
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
	rt *couRuntime) {
	var (
		canonicalState *term.State
		err error
		input string
		rawInput = make([]byte, 4)
		rawInputLen int
	)

	if rt.AcceptInput == false {
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

	handleKey(string(input), cmdMap, contentHeight, contentLineCount, rt)
}

func handleKey(key string,
	cmdMap common.ScriptCmdMap,
	contentHeight, contentLineCount int,
	rt *couRuntime) {

	if rt.CmdMode {
		handleKeyCmdline(key, cmdMap, contentLineCount, rt)
		return
	}

	switch key {
	case csi.CursorUp:
		fallthrough
	case rt.ComCfg.Keys.Up:
		if rt.Scroll > 0 {
			rt.Scroll--
		}

	case csi.CursorDown:
		fallthrough
	case rt.ComCfg.Keys.Down:
		if rt.Scroll < contentLineCount {
			rt.Scroll++
		}

	case rt.ComCfg.Keys.Cmdmode:
		rt.CmdMode = true
		fmt.Printf(csi.CursorShow)

	case csi.PgUp:
		if rt.Scroll-contentHeight < 0 {
			rt.Scroll = 0
		} else {
			rt.Scroll -= contentHeight
		}

	case csi.PgDown:
		if rt.Scroll+contentHeight >= contentLineCount {
			rt.Scroll = contentLineCount - 1
		} else {
			rt.Scroll += contentHeight
		}

	case csi.Home:
		rt.Scroll = 0

	case csi.End:
		rt.Scroll = contentLineCount - 1

	case rt.ComCfg.Keys.Quit:
		fallthrough
	case csi.CursorLeft:
		fallthrough
	case rt.ComCfg.Keys.Left:
		fallthrough
	case csi.SigInt:
		fallthrough
	case csi.SigTstp:
		rt.Active = false
	}
}

func handleKeyCmdline(key string,
	cmdMap common.ScriptCmdMap,
	contentLineCount int,
	rt *couRuntime) {

	switch key {
	case rt.ComCfg.Keys.Cmdenter:
		rt.Feedback = common.HandleCommand(&rt.Active,
			rt.CmdLine,
			rt.CmdLineRows[:],
			contentLineCount,
			&rt.Scroll,
			cmdMap)
		fallthrough
	case csi.SigInt:
		fallthrough
	case csi.SigTstp:
		rt.CmdLine = ""
		rt.CmdLineCursor = 0
		rt.CmdLineInsert = false
		rt.CmdLineRowIdx = -1
		rt.CmdMode = false
		fmt.Printf(csi.CursorHide)

	case csi.Backspace:
		if rt.CmdLineCursor > 0 {
			rt.CmdLine = rt.CmdLine[:rt.CmdLineCursor-1] +
				rt.CmdLine[rt.CmdLineCursor:]
			rt.CmdLineCursor--
		}

	case csi.CursorRight:
		if rt.CmdLineCursor < len(rt.CmdLine) {
			rt.CmdLineCursor++
		}

	case csi.CursorUp:
		if rt.CmdLineRowIdx < len(rt.CmdLineRows)-1 {
			rt.CmdLineRowIdx++
			rt.CmdLine = rt.CmdLineRows[rt.CmdLineRowIdx]
			rt.CmdLineCursor = len(rt.CmdLine)
		}

	case csi.CursorLeft:
		if rt.CmdLineCursor > 0 {
			rt.CmdLineCursor--
		}

	case csi.CursorDown:
		if rt.CmdLineRowIdx >= 0 {
			rt.CmdLineRowIdx--
		}
		if rt.CmdLineRowIdx >= 0 {
			rt.CmdLine = rt.CmdLineRows[rt.CmdLineRowIdx]
		} else {
			rt.CmdLine = ""
		}
		rt.CmdLineCursor = len(rt.CmdLine)

	case csi.Home:
		rt.CmdLineCursor = 0

	case csi.Insert:
		rt.CmdLineInsert = !rt.CmdLineInsert

	case csi.Delete:
		if rt.CmdLineCursor < len(rt.CmdLine) {
			rt.CmdLine = rt.CmdLine[:rt.CmdLineCursor] +
				rt.CmdLine[rt.CmdLineCursor+1:]
		}

	case csi.End:
		rt.CmdLineCursor = len(rt.CmdLine)

	default:
		if len(key) == 1 {
			var insertReplace = 0

			if rt.CmdLineInsert == true &&
				rt.CmdLineCursor < len(rt.CmdLine) {
				insertReplace = 1
			}

			rt.CmdLine = rt.CmdLine[:rt.CmdLineCursor] +
				key +
				rt.CmdLine[rt.CmdLineCursor+
					insertReplace:]
			rt.CmdLineCursor++
		}
	}
}

func tick(cmdMap common.ScriptCmdMap, rt *couRuntime) {
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

	headerLines = common.SplitByLines(termW, rt.CouCfg.Header)
	titleLines = common.SplitByLines(termW, rt.Title)
	contentLines = common.SplitByLines(termW, rt.Content)
	lower = common.GenerateLower(rt.CmdLine,
		rt.CmdMode,
		rt.ComCfg,
		&rt.Feedback,
		rt.CouCfg.Pager.Title,
		termW)

	common.DrawUpper(rt.ComCfg, headerLines, termW, titleLines)

	contentHeight = termH -
		len(common.SplitByLines(termW, rt.CouCfg.Header)) -
		1 -
		len(common.SplitByLines(termW, rt.Title)) -
		1

	drawContent(contentLines, contentHeight, *rt, termW)

	csi.SetCursor(1, termH)
	fmt.Printf("%v", lower)
	csi.SetCursor((len(rt.ComCfg.CmdLine.Prefix) + rt.CmdLineCursor + 1),
		termH)

	handleInput(cmdMap, contentHeight, len(contentLines), rt)
}

func main() {
	var cmdMap common.ScriptCmdMap
	var fnMap common.ScriptFnMap
	var rt couRuntime

	cmdMap = getCmdMap(&rt)
	fnMap = getFnMap(&rt)
	rt = newCouRuntime(fnMap)

	rt.Content, rt.Active = handleArgs(&rt.Title)
	if rt.Active == false {
		return
	}

	fmt.Printf(csi.CursorHide)
	defer fmt.Printf(csi.CursorShow)
	defer fmt.Printf("%v%v\n", csi.FgDefault, csi.BgDefault)

	if rt.CouCfg.Events.Start != "" {
		fnMap[rt.CouCfg.Events.Start]()
	}

	for rt.Active {
		tick(cmdMap, &rt)
	}

	if rt.CouCfg.Events.Quit != "" {
		fnMap[rt.CouCfg.Events.Quit]()
		tick(cmdMap, &rt)
	}
}
