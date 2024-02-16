// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./genversion.go
package main

import (
	"io"
	"fmt"
	"strconv"
	"strings"
	"os"
	"os/exec"
	
	"golang.org/x/term"
)

type MenuPath []string

const (
	SIGINT  = "\003"
	SIGTSTP = "\004"

	SEQ_CLEAR      = "\033[H\033[2J"
	SEQ_FG_DEFAULT = "\033[39m"
	SEQ_BG_DEFAULT = "\033[49m"
	SEQ_CRSR_HIDE  = "\033[?25l"
	SEQ_CRSR_SHOW  = "\033[?25h"
)

var Version string

func (mp MenuPath) curMenu() string {
	return mp[len(mp) - 1]
}

func drawMenu(cfg Config, curMenu Menu, cursor uint) {
	var prefix, postfix string
	var fg FgColor
	var bg BgColor
	
	for i := uint(0); i < uint(len(curMenu.Entries)); i++ {
		if curMenu.Entries[i].Shell == "" {
			prefix = cfg.EntryMenuPrefix
			postfix = cfg.EntryMenuPostfix
		} else {
			prefix = cfg.EntryShellPrefix
			postfix = cfg.EntryShellPostfix
		}
		
		if i == cursor {
			fg = cfg.EntryHoverFg
			bg = cfg.EntryHoverBg
		} else {
			fg = cfg.EntryFg
			bg = cfg.EntryBg
		}
		
		fmt.Printf("%v%v%v%v%v\n",
		           fg,
		           bg,
		           prefix,
		           curMenu.Entries[i].Caption,
		           postfix)
	}
}

func drawUpper(cfg Config, curMenuName string) {
	fmt.Printf("%v%v%v\n", cfg.HeaderFg, cfg.HeaderBg, cfg.Header)
	fmt.Printf("%v%v%v\n",
	           cfg.TitleFg,
	           cfg.TitleBg,
	           cfg.Menus[curMenuName].Title)
}

func generateLower(cfg      Config,
                   cmdline  string,
                   cmdmode  bool,
                   feedback string,
                   termW    int)    string {
	var ret string
	
	if cmdmode == true {
		ret = fmt.Sprintf("%v%v%v%v",
			          cfg.CmdlineFg,
			          cfg.CmdlineBg,
			          cfg.CmdlinePrefix,
			          cmdline)
	} else {
		feedback = strings.TrimSpace(feedback)
		ret = fmt.Sprintf("%v%v", cfg.FeedbackPrefix, feedback)
		if strNeededLines(ret, termW) > 1 {
			// TODO will become a call to courier later
			ret = cfg.FeedbackPrefix
		}
		
		ret = fmt.Sprintf("%v%v%v",
		                  cfg.FeedbackFg,
		                  cfg.FeedbackBg,
		                  ret)
	}
	
	return ret
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
                 cfg      Config,
                 cmdline  *string,
                 cmdmode  *bool,
                 cursor   *uint,
                 feedback *string,
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
		handle_key(string(input),
		           active,
		           cfg,
		           cmdline,
		           cmdmode,
		           cursor,
		           feedback,
		           menuPath)
	}
}

func handle_key(key      string,
                active   *bool,
                cfg      Config,
		cmdline  *string,
		cmdmode  *bool,
                cursor   *uint,
                feedback *string,
                menuPath *MenuPath) {
	var curMenu = cfg.Menus[menuPath.curMenu()]
	var curEntry = &curMenu.Entries[*cursor]

	if *cmdmode {
		handleKeyCmdline(key,
		                 active,
		                 cfg,
		                 cmdline,
		                 cmdmode,
		                 cursor,
		                 curMenu,
		                 feedback)
		return
	}
	
	switch key {
	case cfg.KeyQuit:
		*active = false

	case cfg.KeyLeft:
		if len(*menuPath) > 1 {
			*menuPath = (*menuPath)[:len(*menuPath) - 1]
			*cursor = 0
		}

	case cfg.KeyDown:
		if *cursor < uint(len(curMenu.Entries) - 1) {
			*cursor++
		}

	case cfg.KeyUp:
		if *cursor > 0 {
			*cursor--
		}

	case cfg.KeyRight:
		if curEntry.Menu != "" {
			*menuPath = append(*menuPath, curEntry.Menu)
			*cursor = 0
		}

	case cfg.KeyExecute:
		if curEntry.Shell != "" {
			*feedback = handleShell(curEntry.Shell)
		}
	
	case cfg.KeyCmdmode:
		*cmdmode = true
		fmt.Printf(SEQ_CRSR_SHOW)

	case SIGINT: fallthrough
	case SIGTSTP:
		*active = false
	}
}

func handleKeyCmdline(key      string,
                      active   *bool,
		      cfg      Config,
		      cmdline  *string,
		      cmdmode  *bool,
                      cursor   *uint,
                      curMenu  Menu,
                      feedback *string) {
	switch key {
	case cfg.KeyCmdenter:
		*feedback = handleCommand(active, cmdline, cursor, curMenu)
		fallthrough
	case SIGINT:
		fallthrough
	case SIGTSTP:
		*cmdmode = false
		*cmdline = ""
		fmt.Printf(SEQ_CRSR_HIDE)

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

func setCursor(x, y int) {
	fmt.Printf("\033[%v;%vH", y, x);
}

func strNeededLines(s string, termW int) uint {
	var ret uint = 0
	var line int = 0

	for i := 0; i < len(s); i++ {
		line++

		if line >= termW {
			line = 0
			ret++
		}
	}

	if line > 0 {
		ret++
	}

	return ret
}

func main() {
	var active = true
	var cfg = cfgFromFile()
	var cmdline string = ""
	var cmdmode bool = false
	var cursor uint = 0
	var err error
	var feedback string = fmt.Sprintf("Welcome to hui %v", Version)
	var lower string
	var menuPath = make(MenuPath, 1, 8)
	var termH, termW int

	_, mainMenuExists := cfg.Menus["main"]

	if mainMenuExists {
		menuPath[0] = "main"
	} else {
		panic("main menu not found in config")
	}
	
	fmt.Printf(SEQ_CRSR_HIDE)
	defer fmt.Printf(SEQ_CRSR_SHOW)
	defer fmt.Printf("%v%v", SEQ_FG_DEFAULT, SEQ_BG_DEFAULT)

	for active {
		fmt.Print(SEQ_CLEAR)
		termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			panic(fmt.Sprintf("Could not get term size: %v", err))
		}

		lower = generateLower(cfg, cmdline, cmdmode, feedback, termW)

		drawUpper(cfg, menuPath.curMenu())
		drawMenu(cfg, cfg.Menus[menuPath.curMenu()], cursor)
		setCursor(1, termH)
		fmt.Printf("%v", lower)

		handleInput(&active,
		            cfg,
		            &cmdline,
		            &cmdmode,
		            &cursor,
		            &feedback,
		            &menuPath)
	}
}
