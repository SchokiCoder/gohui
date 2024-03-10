// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./genversion.go
package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/csi"

	"io"
	"fmt"
	"golang.org/x/term"
	"strconv"
	"os"
	"os/exec"
)

type menuPath []string

func (mp menuPath) curMenu() string {
	return mp[len(mp) - 1]
}

type huiRuntime struct {
	AcceptInput bool
	Active bool
	CmdLine string
	CmdMode bool
	Comcfg common.ComConfig
	Cursor int
	Feedback string
	Huicfg huiConfig
	Menupath menuPath
}

func newHuiRuntime() huiRuntime {
	return huiRuntime {
		AcceptInput: true,
		Active: true,
		CmdLine: "",
		CmdMode: false,
		Comcfg: common.ComConfigFromFile(),
		Cursor: 0,
		Feedback: "",
		Huicfg: huiConfigFromFile(),
		Menupath: make(menuPath, 1, 8),
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

var AppLicense    string
var AppLicenseUrl string
var AppName       string
var AppNameFormal string
var AppRepo       string
var AppVersion    string

func drawMenu(contentHeight int,
              curMenu Menu,
              cursor int,
              huicfg huiConfig,
              termW int) {
	var drawBegin int
	var drawEnd int
	var prefix, postfix string
	var fg csi.FgColor
	var bg csi.BgColor

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
			prefix = huicfg.EntryShellPrefix
			postfix = huicfg.EntryShellPostfix
		} else if curMenu.Entries[i].ShellSession != "" {
			prefix = huicfg.EntryShellSessionPrefix
			postfix = huicfg.EntryShellSessionPostfix
		} else if curMenu.Entries[i].Go != "" {
			prefix = huicfg.EntryGoPrefix
			postfix = huicfg.EntryGoPostfix
		} else {
			prefix = huicfg.EntryMenuPrefix
			postfix = huicfg.EntryMenuPostfix
		}
		
		if i == cursor {
			fg = huicfg.EntryHoverFg
			bg = huicfg.EntryHoverBg
		} else {
			fg = huicfg.EntryFg
			bg = huicfg.EntryBg
		}
		
		common.Cprinta(huicfg.EntryAlignment,
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

func handleCommand(curMenu Menu, runtime *huiRuntime) string {
	var err error
	var num uint64
	var ret string = ""
	
	switch runtime.CmdLine {
	case "q":
		fallthrough
	case "quit":
		fallthrough
	case "exit":
		runtime.Active = false

	default:
		num, err = strconv.ParseUint(runtime.CmdLine, 10, 32)
		
		if err != nil {
			ret = fmt.Sprintf("Command \"%v\" not recognised",
			                  runtime.CmdLine)
		} else {		
			if int(num) < len(curMenu.Entries) - 1 {
				runtime.Cursor = int(num)
			} else {
				runtime.Cursor = int(len(curMenu.Entries) - 1)
			}
		}
	}
	
	runtime.CmdLine = ""
	return ret
}

func handleInput(runtime *huiRuntime) {
	var canonicalState *term.State
	var err error
	var input = make([]byte, 1)

	if runtime.AcceptInput == false {
		return
	}

	canonicalState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Switching to raw mode failed: %v", err))
	}

	_, err = os.Stdin.Read(input)
	if err != nil {
		panic(fmt.Sprintf("Reading from stdin failed: %v", err))
	}

	term.Restore(int(os.Stdin.Fd()), canonicalState)

	for i := 0; i < len(input); i++ {
		handleKey(string(input), runtime)
	}
}

func handleKey(key string, runtime *huiRuntime) {
	var curMenu = runtime.Huicfg.Menus[runtime.Menupath.curMenu()]
	var curEntry = &curMenu.Entries[runtime.Cursor]

	if runtime.CmdMode {
		handleKeyCmdline(key, curMenu, runtime)
		return
	}
	
	switch key {
	case runtime.Comcfg.KeyQuit:
		runtime.Active = false

	case runtime.Comcfg.KeyLeft:
		if len(runtime.Menupath) > 1 {
			runtime.Menupath = runtime.Menupath[:len(runtime.Menupath) - 1]
			runtime.Cursor = 0
		}

	case runtime.Comcfg.KeyDown:
		if runtime.Cursor < len(curMenu.Entries) - 1 {
			runtime.Cursor++
		}

	case runtime.Comcfg.KeyUp:
		if runtime.Cursor > 0 {
			runtime.Cursor--
		}

	case runtime.Comcfg.KeyRight:
		if curEntry.Menu != "" {
			runtime.Menupath = append(runtime.Menupath, curEntry.Menu)
			runtime.Cursor = 0
		}

	case runtime.Huicfg.KeyExecute:
		if curEntry.Shell != "" {
			runtime.Feedback = handleShell(curEntry.Shell)
		} else if curEntry.ShellSession != "" {
			runtime.Feedback = common.HandleShellSession(curEntry.ShellSession)
		} else if curEntry.Go != "" {
			huiFuncs[curEntry.Go](runtime)
		}
	
	case runtime.Comcfg.KeyCmdmode:
		runtime.CmdMode = true
		fmt.Printf(csi.CURSOR_SHOW)

	case csi.SIGINT:
		fallthrough
	case csi.SIGTSTP:
		runtime.Active = false
	}
}

func handleKeyCmdline(key     string,
                      curMenu Menu,
                      runtime *huiRuntime) {
	switch key {
	case runtime.Comcfg.KeyCmdenter:
		runtime.Feedback = handleCommand(curMenu, runtime)
		fallthrough
	case csi.SIGINT:
		fallthrough
	case csi.SIGTSTP:
		runtime.CmdMode = false
		runtime.CmdLine = ""
		fmt.Printf(csi.CURSOR_HIDE)

	default:
		runtime.CmdLine = fmt.Sprintf("%v%v",
		                              runtime.CmdLine,
		                              string(key))
	}
}

func handleShell(shell string) string {
	var cmd *exec.Cmd
	var cmderr io.ReadCloser
	var cmdout io.ReadCloser
	var err error
	var strerr []byte
	var strout []byte

	cmd = exec.Command("sh", "-c", shell)

	cmderr, err = cmd.StderrPipe()
	if err != nil {
		return fmt.Sprintf("Could not get stderr: %s", err)
	}

	cmdout, err = cmd.StdoutPipe()
	if err != nil {
		return fmt.Sprintf("Could not get stdout: %s", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Sprintf("Could not start child process: %s", err)
	}

	strerr, err = io.ReadAll(cmderr)
	if err != nil {
		return fmt.Sprintf("Could not read stderr: %s", err)
	}

	strout, err = io.ReadAll(cmdout)
	if err != nil {
		return fmt.Sprintf("Could not read stdout: %s", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Sprintf("Child error: %s", err)
	}

	if len(strerr) > 0 {
		return string(strerr)
	} else {
		return string(strout)
	}
}

func tick(runtime *huiRuntime) {
	var contentHeight int
	var curMenu Menu
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
	curMenu = runtime.Huicfg.Menus[runtime.Menupath.curMenu()]

	headerLines = common.SplitByLines(termW, runtime.Huicfg.Header)
	titleLines = common.SplitByLines(termW, curMenu.Title)
	lower = common.GenerateLower(runtime.CmdLine,
	                             runtime.CmdMode,
	                             runtime.Comcfg,
	                             &runtime.Feedback,
	                             runtime.Huicfg.PagerTitle,
	                             termW)

	common.DrawUpper(runtime.Comcfg, headerLines, termW, titleLines)

	contentHeight = termH -
	                len(common.SplitByLines(termW, runtime.Huicfg.Header)) -
	                1 -
	                len(common.SplitByLines(termW, curMenu.Title)) -
	                1
	drawMenu(contentHeight, curMenu, runtime.Cursor, runtime.Huicfg, termW)

	csi.SetCursor(1, termH)
	fmt.Printf("%v", lower)

	handleInput(runtime)
}

func main() {
	var runtime = newHuiRuntime()

	_, mainMenuExists := runtime.Huicfg.Menus["main"]

	if mainMenuExists == false {
		panic("\"main\" menu not found in config.")
	}
	runtime.Menupath[0] = "main"

	runtime.Active = handleArgs()
	
	fmt.Printf(csi.CURSOR_HIDE)
	defer fmt.Printf(csi.CURSOR_SHOW)
	defer fmt.Printf("%v%v", csi.FG_DEFAULT, csi.BG_DEFAULT)

	if runtime.Huicfg.GoStart != "" {
		huiFuncs[runtime.Huicfg.GoStart](&runtime)
	}

	for runtime.Active {
		tick(&runtime)
	}

	if runtime.Huicfg.GoQuit != "" {
		huiFuncs[runtime.Huicfg.GoQuit](&runtime)
		tick(&runtime)
	}
}
