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

func (mp MenuPath) CurMenu() string {
	return mp[len(mp) - 1]
}

func draw_menu(cfg Config, cur_menu Menu, cursor uint) {
	var prefix, postfix string
	var fg FgColor
	var bg BgColor
	
	for i := uint(0); i < uint(len(cur_menu.Entries)); i++ {
		switch cur_menu.Entries[i].Content.EcType {
		case ECT_MENU:
			prefix = cfg.EntryMenuPrefix
			postfix = cfg.EntryMenuPostfix
		
		case ECT_SHELL:
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
		           cur_menu.Entries[i].Caption,
		           postfix)
	}
}

func draw_upper(cfg Config, cur_menu_name string) {
	fmt.Printf("%v%v%v\n", cfg.HeaderFg, cfg.HeaderBg, cfg.Header)
	fmt.Printf("%v%v%v\n",
	           cfg.TitleFg,
	           cfg.TitleBg,
	           cfg.Menus[cur_menu_name].Title)
}

func generate_lower(cfg Config,
                    cmdline string,
                    cmdmode bool,
                    feedback string,
                    term_w int) string {
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
		if strNeededLines(ret, term_w) > 1 {
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

func handle_command(active   *bool,
                    cmdline  *string,
                    cursor   *uint,
                    cur_menu Menu) string {	
	var ret string = ""
	
	switch *cmdline {
	case "q":
		fallthrough
	case "quit":
		fallthrough
	case "exit":
		*active = false

	default:
		num, err := strconv.ParseUint(*cmdline, 10, 32)
		
		if err != nil {
			ret = fmt.Sprintf("Command \"%v\" not recognised",
			                  *cmdline)
		} else {		
			if int(num) < len(cur_menu.Entries) - 1 {
				*cursor = uint(num)
			} else {
				*cursor = uint(len(cur_menu.Entries) - 1)
			}
		}
	}
	
	*cmdline = ""
	return ret
}

func handle_input(active    *bool,
                  cfg       Config,
                  cmdline   *string,
                  cmdmode   *bool,
                  cursor    *uint,
                  feedback  *string,
                  menu_path *MenuPath) {
	var input = make([]byte, 1)

	canonical_state, raw_err := term.MakeRaw(int(os.Stdin.Fd()))
	if raw_err != nil {
		panic(raw_err)
	}

	_, read_err := os.Stdin.Read(input)
	if read_err != nil {
		panic(read_err)
	}

	term.Restore(int(os.Stdin.Fd()), canonical_state)

	for i := 0; i < len(input); i++ {
		handle_key(string(input),
		           active,
		           cfg,
		           cmdline,
		           cmdmode,
		           cursor,
		           feedback,
		           menu_path)
	}
}

func handle_key(key       string,
                active    *bool,
                cfg       Config,
		cmdline   *string,
		cmdmode   *bool,
                cursor    *uint,
                feedback  *string,
                menu_path *MenuPath) {
	var cur_menu = cfg.Menus[menu_path.CurMenu()]
	var cur_entry = &cur_menu.Entries[*cursor]

	if *cmdmode {
		handle_key_cmdline(key,
		                   active,
		                   cfg,
		                   cmdline,
		                   cmdmode,
		                   cursor,
		                   cur_menu,
		                   feedback)
		return
	}
	
	switch key {
	case cfg.KeyQuit:
		*active = false

	case cfg.KeyLeft:
		if len(*menu_path) > 1 {
			*menu_path = (*menu_path)[:len(*menu_path) - 1]
			*cursor = 0
		}

	case cfg.KeyDown:
		if *cursor < uint(len(cur_menu.Entries) - 1) {
			*cursor++
		}

	case cfg.KeyUp:
		if *cursor > 0 {
			*cursor--
		}

	case cfg.KeyRight:
		if cur_entry.Content.EcType == ECT_MENU {
			*menu_path = append(*menu_path, cur_entry.Content.Menu)
			*cursor = 0
		}

	case cfg.KeyExecute:
		if cur_entry.Content.EcType == ECT_SHELL {
			*feedback = handle_shell(cur_entry.Content.Shell)
		}
	
	case cfg.KeyCmdmode:
		*cmdmode = true
		fmt.Printf(SEQ_CRSR_SHOW)

	case SIGINT: fallthrough
	case SIGTSTP:
		*active = false
	}
}

func handle_key_cmdline(key       string,
                        active    *bool,
		        cfg       Config,
		        cmdline   *string,
		        cmdmode   *bool,
                        cursor    *uint,
                        cur_menu  Menu,
                        feedback  *string) {
	switch key {
	case cfg.KeyCmdenter:
		*feedback = handle_command(active, cmdline, cursor, cur_menu)
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

func handle_shell(shell string) string {
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

func set_cursor(x, y int) {
	fmt.Printf("\033[%v;%vH", y, x);
}

func strNeededLines(s string, term_w int) uint {
	var ret uint = 0
	var line int = 0

	for i := 0; i < len(s); i++ {
		line++

		if line >= term_w {
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
	var menu_path = make(MenuPath, 1, 8)
	var term_h, term_w int

	_, main_menu_exists := cfg.Menus["main"]

	if main_menu_exists {
		menu_path[0] = "main"
	} else {
		panic("main menu not found in config")
	}
	
	fmt.Printf(SEQ_CRSR_HIDE)
	defer fmt.Printf(SEQ_CRSR_SHOW)
	defer fmt.Printf("%v%v", SEQ_FG_DEFAULT, SEQ_BG_DEFAULT)

	for active {
		fmt.Print(SEQ_CLEAR)
		term_w, term_h, err = term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			panic(fmt.Sprintf("Could not get term size: %v", err))
		}

		lower = generate_lower(cfg, cmdline, cmdmode, feedback, term_w)

		draw_upper(cfg, menu_path.CurMenu())
		draw_menu(cfg, cfg.Menus[menu_path.CurMenu()], cursor)
		set_cursor(1, term_h)
		fmt.Printf("%v", lower)

		handle_input(&active,
		             cfg,
		             &cmdline,
		             &cmdmode,
		             &cursor,
		             &feedback,
		             &menu_path)
	}
}

