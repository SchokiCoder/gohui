// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./geninfo.go
package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/csi"

	"fmt"
	"golang.org/x/term"
	"os"
	"strconv"
	"strings"
)

type menuPath []string

func (mp menuPath) curMenu() string {
	return mp[len(mp)-1]
}

type huiRuntime struct {
	AcceptInput   bool
	Active        bool
	CmdLine       string
	CmdLineCursor int
	CmdLineInsert bool
	CmdLineRowIdx int
	CmdLineRows   [common.CmdlineMaxRows]string
	CmdMode       bool
	Comcfg        common.ComConfig
	Cursor        int
	Feedback      string
	Huicfg        huiConfig
	Menupath      menuPath
}

func newHuiRuntime() huiRuntime {
	return huiRuntime{
		AcceptInput:   true,
		Active:        true,
		CmdLine:       "",
		CmdLineCursor: 0,
		CmdLineInsert: false,
		CmdLineRowIdx: -1,
		CmdMode:       false,
		Comcfg:        common.ComConfigFromFile(),
		Cursor:        0,
		Feedback:      "",
		Huicfg:        huiConfigFromFile(),
		Menupath:      make(menuPath, 1, 8),
	}
}

const HELP = `Usage: hui [OPTIONS]

Customizable terminal based user-interface for common tasks and personal tastes.

Options:

    -a --about
        prints program name, version, license and repository information then exits

    -h --help
        prints this message then exits

    -v --version
        prints version information then exits

Default keybinds:

    q
        quit the program

    h
        go back

    j
        go down

    k
        go up

    l
        go into

    L
        execute

    :
        enter the internal command line

Internal commands:

    q quit exit
        quit the program

    *number*
        when given a positive number, it is used as a line number to scroll to
`

var (
	AppLicense    string
	AppLicenseUrl string
	AppName       string
	AppNameFormal string
	AppRepo       string
	AppVersion    string
)

func drawMenu(contentHeight int,
	curMenu menu,
	cursor int,
	huicfg huiConfig,
	termW int) {
	var (
		drawBegin       int
		drawEnd         int
		prefix, postfix string
		fg              csi.FgColor
		bg              csi.BgColor
	)

	if len(curMenu.Entries) > contentHeight {
		drawBegin = cursor
		drawEnd = cursor + contentHeight
		if drawEnd > len(curMenu.Entries) {
			drawEnd = len(curMenu.Entries)
		}
	} else {
		drawBegin = 0
		drawEnd = len(curMenu.Entries)
	}

	for i := drawBegin; i < len(curMenu.Entries) && i < drawEnd; i++ {
		if curMenu.Entries[i].Shell != "" {
			prefix = huicfg.Entry.ShellPrefix
			postfix = huicfg.Entry.ShellPostfix
		} else if curMenu.Entries[i].ShellSession != "" {
			prefix = huicfg.Entry.ShellSessionPrefix
			postfix = huicfg.Entry.ShellSessionPostfix
		} else if curMenu.Entries[i].Go != "" {
			prefix = huicfg.Entry.GoPrefix
			postfix = huicfg.Entry.GoPostfix
		} else {
			prefix = huicfg.Entry.MenuPrefix
			postfix = huicfg.Entry.MenuPostfix
		}

		if i == cursor {
			fg = huicfg.Entry.HoverFg
			bg = huicfg.Entry.HoverBg
		} else {
			fg = huicfg.Entry.Fg
			bg = huicfg.Entry.Bg
		}

		common.Cprinta(huicfg.Entry.Alignment,
			fg,
			bg,
			termW,
			fmt.Sprintf("%v%v%v",
				prefix,
				curMenu.Entries[i].Caption,
				postfix))
	}
}

func handleArgs() bool {
	for _, v := range os.Args[1:] {
		switch v {
		case "-v":
			fallthrough
		case "--version":
			common.PrintVersion(AppName, AppVersion)
			return false

		case "-a":
			fallthrough
		case "--about":
			common.PrintAbout(AppLicense,
				AppLicenseUrl,
				AppName,
				AppNameFormal,
				AppRepo,
				AppVersion)
			return false

		case "-h":
			fallthrough
		case "--help":
			fmt.Printf(HELP)
			return false

		default:
			fmt.Fprintf(os.Stderr, "Unknown argument: %v", v)
		}
	}

	return true
}

func handleCommand(contentLineCount int, rt *huiRuntime) string {
	var err error
	var num uint64
	var ret string = ""

	cmdLineParts := strings.SplitN(rt.CmdLine, " ", 2)
	fn := huiCommands[cmdLineParts[0]]
	if fn != nil {
		return fn(cmdLineParts[1], rt)
	}

	switch rt.CmdLine {
	case "q":
		fallthrough
	case "quit":
		fallthrough
	case "exit":
		rt.Active = false

	default:
		num, err = strconv.ParseUint(rt.CmdLine, 10, 32)

		if err != nil {
			ret = fmt.Sprintf("Command \"%v\" not recognised",
				rt.CmdLine)
		} else {
			if int(num) < contentLineCount {
				rt.Cursor = int(num)
			} else {
				rt.Cursor = contentLineCount
			}
		}
	}

	for i := 0; i < len(rt.CmdLineRows)-1; i++ {
		rt.CmdLineRows[len(rt.CmdLineRows)-1-i] =
			rt.CmdLineRows[len(rt.CmdLineRows)-1-i-1]
	}
	rt.CmdLineRows[0] = rt.CmdLine
	return ret
}

func handleInput(contentHeight int, rt *huiRuntime) {
	var canonicalState *term.State
	var err error
	var input string
	var rawInput = make([]byte, 4)
	var rawInputLen int

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

	handleKey(string(input), contentHeight, rt)
}

func handleKey(key string, contentHeight int, rt *huiRuntime) {
	var curMenu = rt.Huicfg.Menus[rt.Menupath.curMenu()]
	var curEntry = &curMenu.Entries[rt.Cursor]

	if rt.CmdMode {
		handleKeyCmdline(key, len(curMenu.Entries), rt)
		return
	}

	switch key {
	case rt.Comcfg.Keys.Quit:
		rt.Active = false

	case csi.CURSOR_LEFT:
		fallthrough
	case rt.Comcfg.Keys.Left:
		if len(rt.Menupath) > 1 {
			rt.Menupath = rt.Menupath[:len(rt.Menupath)-1]
			rt.Cursor = 0
		}

	case csi.CURSOR_DOWN:
		fallthrough
	case rt.Comcfg.Keys.Down:
		if rt.Cursor < len(curMenu.Entries)-1 {
			rt.Cursor++
		}

	case csi.CURSOR_UP:
		fallthrough
	case rt.Comcfg.Keys.Up:
		if rt.Cursor > 0 {
			rt.Cursor--
		}

	case csi.CURSOR_RIGHT:
		fallthrough
	case rt.Comcfg.Keys.Right:
		if curEntry.Menu != "" {
			rt.Menupath = append(rt.Menupath, curEntry.Menu)
			rt.Cursor = 0
		}

	case rt.Huicfg.Keys.Execute:
		if curEntry.Shell != "" {
			rt.Feedback = common.HandleShell(curEntry.Shell)
		} else if curEntry.ShellSession != "" {
			rt.Feedback = common.HandleShellSession(curEntry.ShellSession)
		} else if curEntry.Go != "" {
			huiFuncs[curEntry.Go](rt)
		}

	case rt.Comcfg.Keys.Cmdmode:
		rt.CmdMode = true
		fmt.Printf(csi.CURSOR_SHOW)

	case csi.PGUP:
		if rt.Cursor-contentHeight < 0 {
			rt.Cursor = 0
		} else {
			rt.Cursor -= contentHeight
		}

	case csi.PGDOWN:
		if rt.Cursor+contentHeight >= len(curMenu.Entries) {
			rt.Cursor = len(curMenu.Entries) - 1
		} else {
			rt.Cursor += contentHeight
		}

	case csi.HOME:
		rt.Cursor = 0

	case csi.END:
		rt.Cursor = len(curMenu.Entries) - 1

	case csi.SIGINT:
		fallthrough
	case csi.SIGTSTP:
		rt.Active = false
	}
}

func handleKeyCmdline(key string, contentLineCount int, rt *huiRuntime) {
	switch key {
	case rt.Comcfg.Keys.Cmdenter:
		rt.Feedback = handleCommand(contentLineCount, rt)
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

func tick(rt *huiRuntime) {
	var contentHeight int
	var curMenu menu
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
	curMenu = rt.Huicfg.Menus[rt.Menupath.curMenu()]

	headerLines = common.SplitByLines(termW, rt.Huicfg.Header)
	titleLines = common.SplitByLines(termW, curMenu.Title)
	lower = common.GenerateLower(rt.CmdLine,
		rt.CmdMode,
		rt.Comcfg,
		&rt.Feedback,
		rt.Huicfg.Pager.Title,
		termW)

	common.DrawUpper(rt.Comcfg, headerLines, termW, titleLines)

	contentHeight = termH -
		len(common.SplitByLines(termW, rt.Huicfg.Header)) -
		1 -
		len(common.SplitByLines(termW, curMenu.Title)) -
		1
	drawMenu(contentHeight, curMenu, rt.Cursor, rt.Huicfg, termW)

	csi.SetCursor(1, termH)
	fmt.Printf("%v", lower)
	csi.SetCursor((len(rt.Comcfg.CmdLine.Prefix) + rt.CmdLineCursor + 1),
		termH)

	handleInput(contentHeight, rt)
}

func main() {
	var rt = newHuiRuntime()

	_, mainMenuExists := rt.Huicfg.Menus["main"]

	if mainMenuExists == false {
		panic("\"main\" menu not found in config.")
	}
	rt.Menupath[0] = "main"

	rt.Active = handleArgs()
	if rt.Active == false {
		return
	}

	fmt.Printf(csi.CURSOR_HIDE)
	defer fmt.Printf(csi.CURSOR_SHOW)
	defer fmt.Printf("%v%v\n", csi.FG_DEFAULT, csi.BG_DEFAULT)

	if rt.Huicfg.Events.Start != "" {
		huiFuncs[rt.Huicfg.Events.Start](&rt)
	}

	for rt.Active {
		tick(&rt)
	}

	if rt.Huicfg.Events.Quit != "" {
		huiFuncs[rt.Huicfg.Events.Quit](&rt)
		tick(&rt)
	}
}
