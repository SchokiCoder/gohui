// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024 - 2025  Andy Frank Schoknecht

package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/csi"

	"fmt"
	"golang.org/x/term"
	"os"
)

type menuPathNode struct {
	Cursor int
	Menu string
}

type menuPath []menuPathNode

func (mp menuPath) curCursor(
) *int {
	return &mp[len(mp) - 1].Cursor
}

func (mp menuPath) curMenu(
) string {
	return mp[len(mp)-1].Menu
}

type appData struct {
	common.ComAppData
	HuiCfg            huiConfig
	MPath             menuPath
}

const HELP = `Usage: hui [OPTIONS]

Customizable terminal based user-interface for common tasks and personal tastes.

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

Environmental variables:

    PAGER
        sets the used pager, in case feedback exceeds one line of length
`

var (
	AppLicense    string
	AppLicenseUrl string
	AppName       string
	AppNameFormal string
	AppRepo       string
	AppVersion    string
)

func drawMenu(
	contentHeight int,
	curMenu menu,
	cursor int,
	huicfg huiConfig,
	termW int,
) {
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
		hover := i == cursor

		if curMenu.Entries[i].Shell != "" {
			if hover {
				prefix = huicfg.Entry.ShellHoverPrefix
				postfix = huicfg.Entry.ShellHoverPostfix
			} else {
				prefix = huicfg.Entry.ShellPrefix
				postfix = huicfg.Entry.ShellPostfix
			}
		} else if curMenu.Entries[i].ShellSession != "" {
			if hover {
				prefix = huicfg.Entry.ShellSessionHoverPrefix
				postfix = huicfg.Entry.ShellSessionHoverPostfix
			} else {
				prefix = huicfg.Entry.ShellSessionPrefix
				postfix = huicfg.Entry.ShellSessionPostfix
			}
		} else if curMenu.Entries[i].Go != "" {
			if hover {
				prefix = huicfg.Entry.GoHoverPrefix
				postfix = huicfg.Entry.GoHoverPostfix
			} else {
				prefix = huicfg.Entry.GoPrefix
				postfix = huicfg.Entry.GoPostfix
			}
		} else {
			if hover {
				prefix = huicfg.Entry.MenuHoverPrefix
				postfix = huicfg.Entry.MenuHoverPostfix
			} else {
				prefix = huicfg.Entry.MenuPrefix
				postfix = huicfg.Entry.MenuPostfix
			}
		}

		if hover {
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

func handleArgs(
	cfgPath *string,
) bool {
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
			return false

		case "-c":
			fallthrough
		case "--config":
			*cfgPath = os.Args[i + 1]
			i++

		case "-h":
			fallthrough
		case "--help":
			fmt.Printf(HELP)
			return false

		case "-v":
			fallthrough
		case "--version":
			common.PrintVersion(AppName, AppVersion)
			return false

		default:
			fmt.Fprintf(os.Stderr,
				"Unknown argument: \"%v\"\n",
				os.Args[i])
			return false
		}
	}

	return true
}

func handleInput(
	contentHeight int,
	cmdMap common.ScriptCmdMap,
	fnMap common.ScriptFnMap,
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
		panic(fmt.Sprintf("Switching to raw mode failed: %v", err))
	}

	rawInputLen, err = os.Stdin.Read(rawInput)
	if err != nil {
		panic(fmt.Sprintf("Reading from stdin failed: %v", err))
	}
	input = string(rawInput[0:rawInputLen])

	term.Restore(int(os.Stdin.Fd()), canonicalState)

	handleKey(string(input), cmdMap, contentHeight, fnMap, ad)
}

func handleKey(
	key string,
	cmdMap common.ScriptCmdMap,
	contentHeight int,
	fnMap common.ScriptFnMap,
	ad *appData,
) {
	var (
		curCursor = ad.MPath.curCursor()
		curMenu = ad.HuiCfg.Menus[ad.MPath.curMenu()]
		curEntry = &curMenu.Entries[*curCursor]
	)

	if ad.CmdLine.Active {
		common.HandleKeyCmdline(key,
			&ad.Active,
			&ad.CmdLine,
			cmdMap,
			&ad.ComCfg,
			len(curMenu.Entries),
			curCursor,
			&ad.Fb)
		return
	}

	switch key {
	case ad.ComCfg.Keys.Quit:
		ad.Active = false

	case csi.CursorLeft:
		fallthrough
	case ad.ComCfg.Keys.Left:
		if len(ad.MPath) > 1 {
			ad.MPath = ad.MPath[:len(ad.MPath)-1]
		}

	case csi.CursorDown:
		fallthrough
	case ad.ComCfg.Keys.Down:
		if *curCursor < len(curMenu.Entries)-1 {
			*curCursor++
		}

	case csi.CursorUp:
		fallthrough
	case ad.ComCfg.Keys.Up:
		if *curCursor > 0 {
			*curCursor--
		}

	case csi.CursorRight:
		fallthrough
	case ad.ComCfg.Keys.Right:
		if curEntry.Menu != "" {
			ad.MPath = append(ad.MPath, menuPathNode{0, curEntry.Menu})
		} else {
			ad.Fb = "Entry has no menu, can't open."
		}

	case ad.HuiCfg.Keys.Execute:
		if curEntry.Shell != "" {
			ad.Fb = common.HandleShell(curEntry.Shell)
		} else if curEntry.ShellSession != "" {
			ad.Fb = common.HandleShellSession(curEntry.ShellSession)
		} else if curEntry.Go != "" {
			fnMap[curEntry.Go]()
		} else {
			ad.Fb = "Entry has no shell or go, can't execute."
		}

	case ad.ComCfg.Keys.Cmdmode:
		ad.CmdLine.Active = true
		fmt.Printf(csi.CursorShow)

	case csi.PgUp:
		if *curCursor - contentHeight < 0 {
			*curCursor = 0
		} else {
			*curCursor -= contentHeight
		}

	case csi.PgDown:
		if *curCursor + contentHeight >= len(curMenu.Entries) {
			*curCursor = len(curMenu.Entries) - 1
		} else {
			*curCursor += contentHeight
		}

	case csi.Home:
		*curCursor = 0

	case csi.End:
		*curCursor = len(curMenu.Entries) - 1

	case csi.SigInt:
		fallthrough
	case csi.SigTstp:
		ad.Active = false
	}
}

func tick(
	cmdMap common.ScriptCmdMap,
	fnMap  common.ScriptFnMap,
	ad     *appData,
) {
	var (
		contentHeight int
		curMenu menu
		err error
		headerLines []string
		lower string
		termH, termW int
		titleLines []string
	)

	fmt.Print(csi.Clear)
	termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Could not get term size: %v", err))
	}
	curMenu = ad.HuiCfg.Menus[ad.MPath.curMenu()]

	headerLines = common.SplitByLines(termW, ad.HuiCfg.Header)
	titleLines = common.SplitByLines(termW, curMenu.Title)
	lower = common.GenerateLower(ad.CmdLine.Content,
		ad.CmdLine.Active,
		ad.ComCfg,
		&ad.Fb,
		ad.HuiCfg.Pager.Title,
		termW)

	common.DrawUpper(ad.ComCfg, headerLines, termW, titleLines)

	contentHeight = termH -
		len(common.SplitByLines(termW, ad.HuiCfg.Header)) -
		1 -
		len(common.SplitByLines(termW, curMenu.Title)) -
		1
	drawMenu(contentHeight, curMenu, *ad.MPath.curCursor(), ad.HuiCfg, termW)

	csi.SetCursor(1, termH)
	fmt.Printf("%v", lower)
	csi.SetCursorAligned(ad.ComCfg.CmdLine.Alignment,
		(len(ad.ComCfg.CmdLine.Prefix) + len(ad.CmdLine.Content)),
		termW,
		(len(ad.ComCfg.CmdLine.Prefix) + ad.CmdLine.Cursor + 1),
		termH)

	handleInput(contentHeight, cmdMap, fnMap, ad)
}

func main(
) {
	var (
		ad      appData
		cfgPath string
		cmdMap  common.ScriptCmdMap
		fnMap   common.ScriptFnMap
	)

	ad.Active = handleArgs(&cfgPath)
	if ad.Active == false {
		return
	}

	cmdMap = getCmdMap(&ad)
	fnMap = getFnMap(&ad)
	ad.ComAppData = common.NewComAppData(cfgPath)
	ad.MPath = make(menuPath, 1, 8)
	ad.HuiCfg = huiConfigFromFile(&ad, cfgPath, fnMap)

	_, mainMenuExists := ad.HuiCfg.Menus["main"]

	if mainMenuExists == false {
		panic("\"main\" menu not found in config.")
	}
	ad.MPath[0] = menuPathNode{0, "main"}

	if ad.HuiCfg.Events.Start != "" {
		fnMap[ad.HuiCfg.Events.Start]()
	}

	fmt.Printf(csi.CursorHide)
	defer fmt.Printf(csi.CursorShow)
	defer fmt.Printf("%v%v\n", csi.FgDefault, csi.BgDefault)

	for ad.Active {
		tick(cmdMap, fnMap, &ad)
	}

	if ad.HuiCfg.Events.Quit != "" {
		fnMap[ad.HuiCfg.Events.Quit]()
		tick(cmdMap, fnMap, &ad)
	}
}
