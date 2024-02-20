// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./genversion.go
package main

import (
	"github.com/SchokiCoder/gohui/common"

	"fmt"
	"golang.org/x/term"
	"os"
)

var Version string

func drawContent(contentLines []string, coucfg CouCfg) {
	for _, v := range contentLines {
		fmt.Printf("%v%v%v\n", coucfg.ContentFg, coucfg.ContentBg, v)
	}
}

func handleCommand(active *bool, cmdline *string) string {	
	var ret string = ""
	
	switch *cmdline {
	case "q":
		fallthrough
	case "quit":
		fallthrough
	case "exit":
		*active = false

	default:
		ret = fmt.Sprintf("Command \"%v\" not recognised", *cmdline)
	}
	
	*cmdline = ""
	return ret
}

func handleInput(active   *bool,
                 cmdline  *string,
                 cmdmode  *bool,
                 comcfg   common.ComCfg,
                 feedback *string) {
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
		          feedback)
	}
}

func handleKey(key      string,
               active   *bool,
               cmdline  *string,
               cmdmode  *bool,
               comcfg   common.ComCfg,
               feedback *string) {
	if *cmdmode {
		handleKeyCmdline(key,
		                 active,
		                 cmdline,
		                 cmdmode,
		                 comcfg,
		                 feedback)
		return
	}
	
	switch key {
	case comcfg.KeyCmdmode:
		*cmdmode = true
		fmt.Printf(common.SEQ_CRSR_SHOW)

	case comcfg.KeyQuit: fallthrough
	case comcfg.KeyLeft: fallthrough
	case common.SIGINT:  fallthrough
	case common.SIGTSTP:
		*active = false
	}
}

func handleKeyCmdline(key      string,
                      active   *bool,
		      cmdline  *string,
		      cmdmode  *bool,
		      comcfg   common.ComCfg,
                      feedback *string) {
	switch key {
	case comcfg.KeyCmdenter:
		*feedback = handleCommand(active, cmdline)
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

func main() {
	var active = true
	var cmdline string = ""
	var cmdmode bool = false
	var comcfg = common.CfgFromFile()
	var contentLines []string
	var coucfg = cfgFromFile()
	var err error
	var feedback string = fmt.Sprintf("Welcome to courier %v", Version)
	var lower string
	var termH, termW int

	fmt.Printf(common.SEQ_CRSR_HIDE)
	defer fmt.Printf(common.SEQ_CRSR_SHOW)
	defer fmt.Printf("%v%v", common.SEQ_FG_DEFAULT, common.SEQ_BG_DEFAULT)

	for active {
		fmt.Print(common.SEQ_CLEAR)
		termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
		if err != nil {
			panic(fmt.Sprintf("Could not get term size: %v", err))
		}

		lower = common.GenerateLower(cmdline,
		                             cmdmode,
		                             comcfg,
		                             feedback,
		                             termW)

		// TODO add title var
		common.DrawUpper(comcfg, coucfg.Header, "magic title")
		drawContent(contentLines, coucfg)
		common.SetCursor(1, termH)
		fmt.Printf("%v", lower)

		handleInput(&active,
		            &cmdline,
		            &cmdmode,
		            comcfg,
		            &feedback)
	}
}
