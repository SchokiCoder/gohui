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

var Version string

func (mp MenuPath) curMenu() string {
	return mp[len(mp) - 1]
}

func drawMenu(curMenu Menu, cursor uint, huicfg HuiCfg) {
	var prefix, postfix string
	var fg common.FgColor
	var bg common.BgColor
	
	for i := uint(0); i < uint(len(curMenu.Entries)); i++ {
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
		
		fmt.Printf("%v%v%v%v%v\n",
		           fg,
		           bg,
		           prefix,
		           curMenu.Entries[i].Caption,
		           postfix)
	}
}

func handleCommand(active   *bool,
                   cmdline  *string,
                   cursor   *uint,
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
				*cursor = uint(num)
			} else {
				*cursor = uint(len(curMenu.Entries) - 1)
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
                 cursor   *uint,
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
               cursor   *uint,
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
		if *cursor < uint(len(curMenu.Entries) - 1) {
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
			*feedback = handleShellSession(curEntry.ShellSession)
		}
	
	case comcfg.KeyCmdmode:
		*cmdmode = true
		fmt.Printf(common.SEQ_CRSR_SHOW)

	case common.SIGINT: fallthrough
	case common.SIGTSTP:
		*active = false
	}
}

func handleKeyCmdline(key      string,
                      active   *bool,
		      cmdline  *string,
		      cmdmode  *bool,
		      comcfg   common.ComCfg,
                      cursor   *uint,
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

func handleShellSession(shell string) string {
	var cmd *exec.Cmd
	var cmderr io.ReadCloser
	var err error
	var strerr []byte

	cmd = exec.Command("sh", "-c", shell)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	cmderr, err = cmd.StderrPipe()
	if err != nil {
		return fmt.Sprintf("Could not get stderr: %s", err)
	}

	err = cmd.Start()
	if err != nil {
		return fmt.Sprintf("Could not start child process: %s", err)
	}

	strerr, err = io.ReadAll(cmderr)
	if err != nil {
		return fmt.Sprintf("Could not read stderr: %s", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Sprintf("Child error: %s", err)
	}

	if len(strerr) > 0 {
		return string(strerr)
	}
	
	return ""
}

func main() {
	var active = true
	var cmdline string = ""
	var cmdmode bool = false
	var comcfg = common.CfgFromFile()
	var cursor uint = 0
	var err error
	var feedback string = fmt.Sprintf("Welcome to hui %v", Version)
	var huicfg = cfgFromFile()
	var lower string
	var menuPath = make(MenuPath, 1, 8)
	var termH, termW int

	_, mainMenuExists := huicfg.Menus["main"]

	if mainMenuExists {
		menuPath[0] = "main"
	} else {
		panic("main menu not found in config")
	}
	
	fmt.Printf(common.SEQ_CRSR_HIDE)
	defer fmt.Printf(common.SEQ_CRSR_SHOW)
	defer fmt.Printf("%v%v", common.SEQ_FG_DEFAULT, common.SEQ_BG_DEFAULT)

	for active {
		fmt.Print(common.SEQ_CLEAR)
		termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			panic(fmt.Sprintf("Could not get term size: %v", err))
		}

		lower = common.GenerateLower(cmdline, cmdmode, comcfg, feedback, termW)

		common.DrawUpper(comcfg, huicfg.Header,
		                 huicfg.Menus[menuPath.curMenu()].Title)
		drawMenu(huicfg.Menus[menuPath.curMenu()], cursor, huicfg)
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
