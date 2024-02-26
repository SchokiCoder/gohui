// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./genversion.go
package main

import (
	"github.com/SchokiCoder/gohui/common"

	"io"
	"fmt"
	"golang.org/x/term"
	"strconv"
	"os"
	"os/exec"
)

type MenuPath []string

var AppLicense    string
var AppLicenseUrl string
var AppName       string
var AppNameFormal string
var AppRepo       string
var AppVersion    string

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

func (mp MenuPath) curMenu() string {
	return mp[len(mp) - 1]
}

func drawMenu(contentHeight int, curMenu Menu, cursor int, huicfg HuiCfg) {
	var drawBegin int
	var drawEnd int
	var prefix, postfix string
	var fg common.FgColor
	var bg common.BgColor

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
		
		common.Cprintf(fg,
		               bg,
		               "%v%v%v\n",
		               prefix, curMenu.Entries[i].Caption, postfix)
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

func handleCommand(active   *bool,
                   cmdline  *string,
                   cursor   *int,
                   curMenu  Menu)    string {	
	var err error
	var num uint64
	var ret string = ""
	
	switch *cmdline {
	case "q":
		fallthrough
	case "quit":
		fallthrough
	case "exit":
		*active = false

	default:
		num, err = strconv.ParseUint(*cmdline, 10, 32)
		
		if err != nil {
			ret = fmt.Sprintf("Command \"%v\" not recognised",
			                  *cmdline)
		} else {		
			if int(num) < len(curMenu.Entries) - 1 {
				*cursor = int(num)
			} else {
				*cursor = int(len(curMenu.Entries) - 1)
			}
		}
	}
	
	*cmdline = ""
	return ret
}

func handleInput(active   *bool,
                 cmdline  *string,
                 cmdmode  *bool,
                 comcfg   common.ComCfg,
                 cursor   *int,
                 feedback *string,
                 huicfg   HuiCfg,
                 menuPath *MenuPath) {
	var canonicalState *term.State
	var err error
	var input = make([]byte, 1)

	canonicalState, err = term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}

	_, err = os.Stdin.Read(input)
	if err != nil {
		panic(err)
	}

	term.Restore(int(os.Stdin.Fd()), canonicalState)

	for i := 0; i < len(input); i++ {
		handleKey(string(input),
		          active,
		          cmdline,
		          cmdmode,
		          comcfg,
		          cursor,
		          feedback,
		          huicfg,
		          menuPath)
	}
}

func handleKey(key      string,
               active   *bool,
               cmdline  *string,
               cmdmode  *bool,
               comcfg   common.ComCfg,
               cursor   *int,
               feedback *string,
               huicfg   HuiCfg,
               menuPath *MenuPath) {
	var curMenu = huicfg.Menus[menuPath.curMenu()]
	var curEntry = &curMenu.Entries[*cursor]

	if *cmdmode {
		handleKeyCmdline(key,
		                 active,
		                 cmdline,
		                 cmdmode,
		                 comcfg,
		                 cursor,
		                 curMenu,
		                 feedback)
		return
	}
	
	switch key {
	case comcfg.KeyQuit:
		*active = false

	case comcfg.KeyLeft:
		if len(*menuPath) > 1 {
			*menuPath = (*menuPath)[:len(*menuPath) - 1]
			*cursor = 0
		}

	case comcfg.KeyDown:
		if *cursor < len(curMenu.Entries) - 1 {
			*cursor++
		}

	case comcfg.KeyUp:
		if *cursor > 0 {
			*cursor--
		}

	case comcfg.KeyRight:
		if curEntry.Menu != "" {
			*menuPath = append(*menuPath, curEntry.Menu)
			*cursor = 0
		}

	case huicfg.KeyExecute:
		if curEntry.Shell != "" {
			*feedback = handleShell(curEntry.Shell)
		} else if curEntry.ShellSession != "" {
			*feedback = common.HandleShellSession(curEntry.ShellSession)
		}
	
	case comcfg.KeyCmdmode:
		*cmdmode = true
		fmt.Printf(common.SEQ_CRSR_SHOW)

	case common.SIGINT:
		fallthrough
	case common.SIGTSTP:
		*active = false
	}
}

func handleKeyCmdline(key      string,
                      active   *bool,
		      cmdline  *string,
		      cmdmode  *bool,
		      comcfg   common.ComCfg,
                      cursor   *int,
                      curMenu  Menu,
                      feedback *string) {
	switch key {
	case comcfg.KeyCmdenter:
		*feedback = handleCommand(active, cmdline, cursor, curMenu)
		fallthrough
	case common.SIGINT:
		fallthrough
	case common.SIGTSTP:
		*cmdmode = false
		*cmdline = ""
		fmt.Printf(common.SEQ_CRSR_HIDE)

	default:
		*cmdline = fmt.Sprintf("%v%v", *cmdline, string(key))
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

func main() {
	var active = true
	var cmdline string = ""
	var cmdmode bool = false
	var comcfg = common.CfgFromFile()
	var contentHeight int
	var cursor int = 0
	var curMenu Menu
	var err error
	var feedback string = fmt.Sprintf("Welcome to %v %v", AppName, AppVersion)
	var huicfg = cfgFromFile()
	var lower string
	var menuPath = make(MenuPath, 1, 8)
	var termH, termW int

	_, mainMenuExists := huicfg.Menus["main"]

	if mainMenuExists {
		menuPath[0] = "main"
	} else {
		panic("\"main\" menu not found in config.")
	}

	active = handleArgs()
	
	fmt.Printf(common.SEQ_CRSR_HIDE)
	defer fmt.Printf(common.SEQ_CRSR_SHOW)
	defer fmt.Printf("%v%v", common.SEQ_FG_DEFAULT, common.SEQ_BG_DEFAULT)

	for active {
		fmt.Print(common.SEQ_CLEAR)
		termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			panic(fmt.Sprintf("Could not get term size: %v", err))
		}
		curMenu = huicfg.Menus[menuPath.curMenu()]

		lower = common.GenerateLower(cmdline,
		                             cmdmode,
		                             comcfg,
		                             &feedback,
		                             termW)

		common.DrawUpper(comcfg, huicfg.Header, curMenu.Title)

		contentHeight = termH -
		                len(common.SplitByLines(termW, huicfg.Header)) -
		                1 -
		                len(common.SplitByLines(termW, curMenu.Title)) -
		                1
		drawMenu(contentHeight, curMenu, cursor, huicfg)

		common.SetCursor(1, termH)
		fmt.Printf("%v", lower)

		handleInput(&active,
		            &cmdline,
		            &cmdmode,
		            comcfg,
		            &cursor,
		            &feedback,
		            huicfg,
		            &menuPath)
	}
}
