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
	ComCfg        common.ComConfig
	Cursor        int
	Feedback      string
	HuiCfg        huiConfig
	Menupath      menuPath
}

func newHuiRuntime(fnMap common.ScriptFnMap) huiRuntime {
	var ret = huiRuntime{
		AcceptInput:   true,
		Active:        true,
		CmdLine:       "",
		CmdLineCursor: 0,
		CmdLineInsert: false,
		CmdLineRowIdx: -1,
		CmdMode:       false,
		ComCfg:        common.ComConfigFromFile(),
		Cursor:        0,
		Feedback:      "",
		Menupath:      make(menuPath, 1, 8),
	}

	ret.HuiCfg = huiConfigFromFile(fnMap, &ret)

	return ret
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

func handleInput(contentHeight int,
	cmdMap common.ScriptCmdMap,
	fnMap common.ScriptFnMap,
	rt *huiRuntime) {
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

	handleKey(string(input), cmdMap, contentHeight, fnMap, rt)
}

func handleKey(key string,
	cmdMap common.ScriptCmdMap,
	contentHeight int,
	fnMap common.ScriptFnMap,
	rt *huiRuntime) {
	var (
		curMenu = rt.HuiCfg.Menus[rt.Menupath.curMenu()]
		curEntry = &curMenu.Entries[rt.Cursor]
	)

	if rt.CmdMode {
		common.HandleKeyCmdline(key,
			&rt.Active,
			&rt.CmdLine,
			&rt.CmdLineCursor,
			&rt.CmdLineInsert,
			&rt.CmdLineRowIdx,
			rt.CmdLineRows[:],
			cmdMap,
			&rt.CmdMode,
			&rt.ComCfg,
			len(curMenu.Entries),
			&rt.Cursor,
			&rt.Feedback)
		return
	}

	switch key {
	case rt.ComCfg.Keys.Quit:
		rt.Active = false

	case csi.CursorLeft:
		fallthrough
	case rt.ComCfg.Keys.Left:
		if len(rt.Menupath) > 1 {
			rt.Menupath = rt.Menupath[:len(rt.Menupath)-1]
			rt.Cursor = 0
		}

	case csi.CursorDown:
		fallthrough
	case rt.ComCfg.Keys.Down:
		if rt.Cursor < len(curMenu.Entries)-1 {
			rt.Cursor++
		}

	case csi.CursorUp:
		fallthrough
	case rt.ComCfg.Keys.Up:
		if rt.Cursor > 0 {
			rt.Cursor--
		}

	case csi.CursorRight:
		fallthrough
	case rt.ComCfg.Keys.Right:
		if curEntry.Menu != "" {
			rt.Menupath = append(rt.Menupath, curEntry.Menu)
			rt.Cursor = 0
		}

	case rt.HuiCfg.Keys.Execute:
		if curEntry.Shell != "" {
			rt.Feedback = common.HandleShell(curEntry.Shell)
		} else if curEntry.ShellSession != "" {
			rt.Feedback = common.HandleShellSession(curEntry.ShellSession)
		} else if curEntry.Go != "" {
			fnMap[curEntry.Go]()
		}

	case rt.ComCfg.Keys.Cmdmode:
		rt.CmdMode = true
		fmt.Printf(csi.CursorShow)

	case csi.PgUp:
		if rt.Cursor-contentHeight < 0 {
			rt.Cursor = 0
		} else {
			rt.Cursor -= contentHeight
		}

	case csi.PgDown:
		if rt.Cursor+contentHeight >= len(curMenu.Entries) {
			rt.Cursor = len(curMenu.Entries) - 1
		} else {
			rt.Cursor += contentHeight
		}

	case csi.Home:
		rt.Cursor = 0

	case csi.End:
		rt.Cursor = len(curMenu.Entries) - 1

	case csi.SigInt:
		fallthrough
	case csi.SigTstp:
		rt.Active = false
	}
}

func tick(cmdMap common.ScriptCmdMap, fnMap common.ScriptFnMap, rt *huiRuntime) {
	var contentHeight int
	var curMenu menu
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
	curMenu = rt.HuiCfg.Menus[rt.Menupath.curMenu()]

	headerLines = common.SplitByLines(termW, rt.HuiCfg.Header)
	titleLines = common.SplitByLines(termW, curMenu.Title)
	lower = common.GenerateLower(rt.CmdLine,
		rt.CmdMode,
		rt.ComCfg,
		&rt.Feedback,
		rt.HuiCfg.Pager.Title,
		termW)

	common.DrawUpper(rt.ComCfg, headerLines, termW, titleLines)

	contentHeight = termH -
		len(common.SplitByLines(termW, rt.HuiCfg.Header)) -
		1 -
		len(common.SplitByLines(termW, curMenu.Title)) -
		1
	drawMenu(contentHeight, curMenu, rt.Cursor, rt.HuiCfg, termW)

	csi.SetCursor(1, termH)
	fmt.Printf("%v", lower)
	csi.SetCursorAligned(rt.ComCfg.CmdLine.Alignment,
		(len(rt.ComCfg.CmdLine.Prefix) + len(rt.CmdLine)),
		termW,
		(len(rt.ComCfg.CmdLine.Prefix) + rt.CmdLineCursor + 1),
		termH)

	handleInput(contentHeight, cmdMap, fnMap, rt)
}

func main() {
	var cmdMap common.ScriptCmdMap
	var fnMap common.ScriptFnMap
	var rt huiRuntime

	cmdMap = getCmdMap(&rt)
	fnMap = getFnMap(&rt)
	rt = newHuiRuntime(fnMap)

	_, mainMenuExists := rt.HuiCfg.Menus["main"]

	if mainMenuExists == false {
		panic("\"main\" menu not found in config.")
	}
	rt.Menupath[0] = "main"

	rt.Active = handleArgs()
	if rt.Active == false {
		return
	}

	if rt.HuiCfg.Events.Start != "" {
		fnMap[rt.HuiCfg.Events.Start]()
	}

	fmt.Printf(csi.CursorHide)
	defer fmt.Printf(csi.CursorShow)
	defer fmt.Printf("%v%v\n", csi.FgDefault, csi.BgDefault)

	for rt.Active {
		tick(cmdMap, fnMap, &rt)
	}

	if rt.HuiCfg.Events.Quit != "" {
		fnMap[rt.HuiCfg.Events.Quit]()
		tick(cmdMap, fnMap, &rt)
	}
}
