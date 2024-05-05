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
	Comcfg        common.ComConfig
	Content       string
	Coucfg        couConfig
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
		Comcfg:        common.ComConfigFromFile(),
		Content:       "",
		Coucfg:        couConfigFromFile(fnMap),
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
		common.Cprinta(rt.Coucfg.Content.Alignment,
			rt.Coucfg.Content.Fg,
			rt.Coucfg.Content.Bg,
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
	case csi.CURSOR_UP:
		fallthrough
	case rt.Comcfg.Keys.Up:
		if rt.Scroll > 0 {
			rt.Scroll--
		}

	case csi.CURSOR_DOWN:
		fallthrough
	case rt.Comcfg.Keys.Down:
		if rt.Scroll < contentLineCount {
			rt.Scroll++
		}

	case rt.Comcfg.Keys.Cmdmode:
		rt.CmdMode = true
		fmt.Printf(csi.CURSOR_SHOW)

	case csi.PGUP:
		if rt.Scroll-contentHeight < 0 {
			rt.Scroll = 0
		} else {
			rt.Scroll -= contentHeight
		}

	case csi.PGDOWN:
		if rt.Scroll+contentHeight >= contentLineCount {
			rt.Scroll = contentLineCount - 1
		} else {
			rt.Scroll += contentHeight
		}

	case csi.HOME:
		rt.Scroll = 0

	case csi.END:
		rt.Scroll = contentLineCount - 1

	case rt.Comcfg.Keys.Quit:
		fallthrough
	case csi.CURSOR_LEFT:
		fallthrough
	case rt.Comcfg.Keys.Left:
		fallthrough
	case csi.SIGINT:
		fallthrough
	case csi.SIGTSTP:
		rt.Active = false
	}
}

func handleKeyCmdline(key string,
	cmdMap common.ScriptCmdMap,
	contentLineCount int,
	rt *couRuntime) {

	switch key {
	case rt.Comcfg.Keys.Cmdenter:
		rt.Feedback = common.HandleCommand(&rt.Active,
			rt.CmdLine,
			rt.CmdLineRows[:],
			contentLineCount,
			&rt.Scroll,
			cmdMap)
		fallthrough
	case csi.SIGINT:
		fallthrough
	case csi.SIGTSTP:
		rt.CmdLine = ""
		rt.CmdLineCursor = 0
		rt.CmdLineInsert = false
		rt.CmdLineRowIdx = -1
		rt.CmdMode = false
		fmt.Printf(csi.CURSOR_HIDE)

	case csi.BACKSPACE:
		if rt.CmdLineCursor > 0 {
			rt.CmdLine = rt.CmdLine[:rt.CmdLineCursor-1] +
				rt.CmdLine[rt.CmdLineCursor:]
			rt.CmdLineCursor--
		}

	case csi.CURSOR_RIGHT:
		if rt.CmdLineCursor < len(rt.CmdLine) {
			rt.CmdLineCursor++
		}

	case csi.CURSOR_UP:
		if rt.CmdLineRowIdx < len(rt.CmdLineRows)-1 {
			rt.CmdLineRowIdx++
			rt.CmdLine = rt.CmdLineRows[rt.CmdLineRowIdx]
			rt.CmdLineCursor = len(rt.CmdLine)
		}

	case csi.CURSOR_LEFT:
		if rt.CmdLineCursor > 0 {
			rt.CmdLineCursor--
		}

	case csi.CURSOR_DOWN:
		if rt.CmdLineRowIdx >= 0 {
			rt.CmdLineRowIdx--
		}
		if rt.CmdLineRowIdx >= 0 {
			rt.CmdLine = rt.CmdLineRows[rt.CmdLineRowIdx]
		} else {
			rt.CmdLine = ""
		}
		rt.CmdLineCursor = len(rt.CmdLine)

	case csi.HOME:
		rt.CmdLineCursor = 0

	case csi.INSERT:
		rt.CmdLineInsert = !rt.CmdLineInsert

	case csi.DELETE:
		if rt.CmdLineCursor < len(rt.CmdLine) {
			rt.CmdLine = rt.CmdLine[:rt.CmdLineCursor] +
				rt.CmdLine[rt.CmdLineCursor+1:]
		}

	case csi.END:
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

	fmt.Print(csi.CLEAR)
	termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Could not get term size: %v", err))
	}

	headerLines = common.SplitByLines(termW, rt.Coucfg.Header)
	titleLines = common.SplitByLines(termW, rt.Title)
	contentLines = common.SplitByLines(termW, rt.Content)
	lower = common.GenerateLower(rt.CmdLine,
		rt.CmdMode,
		rt.Comcfg,
		&rt.Feedback,
		rt.Coucfg.Pager.Title,
		termW)

	common.DrawUpper(rt.Comcfg, headerLines, termW, titleLines)

	contentHeight = termH -
		len(common.SplitByLines(termW, rt.Coucfg.Header)) -
		1 -
		len(common.SplitByLines(termW, rt.Title)) -
		1

	drawContent(contentLines, contentHeight, *rt, termW)

	csi.SetCursor(1, termH)
	fmt.Printf("%v", lower)
	csi.SetCursor((len(rt.Comcfg.CmdLine.Prefix) + rt.CmdLineCursor + 1),
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

	fmt.Printf(csi.CURSOR_HIDE)
	defer fmt.Printf(csi.CURSOR_SHOW)
	defer fmt.Printf("%v%v\n", csi.FG_DEFAULT, csi.BG_DEFAULT)

	if rt.Coucfg.Events.Start != "" {
		fnMap[rt.Coucfg.Events.Start]()
	}

	for rt.Active {
		tick(cmdMap, &rt)
	}

	if rt.Coucfg.Events.Quit != "" {
		fnMap[rt.Coucfg.Events.Quit]()
		tick(cmdMap, &rt)
	}
}
