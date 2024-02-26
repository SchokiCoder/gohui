// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./geninfo.go
package main

import (
	"github.com/SchokiCoder/gohui/common"

	"errors"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strconv"
)

var AppLicense    string
var AppLicenseUrl string
var AppName       string
var AppNameFormal string
var AppRepo       string
var AppVersion    string

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

func drawContent(contentLines  []string,
                 contentHeight int,
                 coucfg        CouCfg,
                 scroll        int) {
	var drawRange int = scroll + contentHeight

	if drawRange > len(contentLines) {
		drawRange = len(contentLines)
	}

	for _, v := range contentLines[scroll:drawRange] {
		common.Cprintf(coucfg.ContentFg, coucfg.ContentBg, "%v\n", v)
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

func handleCommand(active           *bool,
                   cmdline          *string,
                   contentLineCount int,
                   scroll           *int)    string {
	var err error
	var ret string = ""
	var num uint64
	
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
			if int(num) < contentLineCount {
				*scroll = int(num)
			} else {
				*scroll = contentLineCount
			}
		}
	}
	
	*cmdline = ""
	return ret
}

func handleInput(active           *bool,
                 cmdline          *string,
                 cmdmode          *bool,
                 comcfg           common.ComCfg,
                 contentLineCount int,
                 feedback         *string,
                 scroll           *int) {
	var canonicalState *term.State
	var err error
	var input = make([]byte, 1)

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
		handleKey(string(input),
		          active,
		          cmdline,
		          cmdmode,
		          comcfg,
		          contentLineCount,
		          feedback,
		          scroll)
	}
}

func handleKey(key              string,
               active           *bool,
               cmdline          *string,
               cmdmode          *bool,
               comcfg           common.ComCfg,
               contentLineCount int,
               feedback         *string,
               scroll           *int) {
	if *cmdmode {
		handleKeyCmdline(key,
		                 active,
		                 cmdline,
		                 cmdmode,
		                 comcfg,
		                 contentLineCount,
		                 feedback,
		                 scroll)
		return
	}
	
	switch key {
	case comcfg.KeyUp:
		if *scroll > 0 {
			*scroll--
		}

	case comcfg.KeyDown:
		if *scroll < contentLineCount {
			*scroll++
		}


	case comcfg.KeyCmdmode:
		*cmdmode = true
		fmt.Printf(common.SEQ_CRSR_SHOW)

	case comcfg.KeyQuit:
		fallthrough
	case comcfg.KeyLeft:
		fallthrough
	case common.SIGINT:
		fallthrough
	case common.SIGTSTP:
		*active = false
	}
}

func handleKeyCmdline(key              string,
                      active           *bool,
                      cmdline          *string,
                      cmdmode          *bool,
                      comcfg           common.ComCfg,
                      contentLineCount int,
                      feedback         *string,
                      scroll           *int) {
	switch key {
	case comcfg.KeyCmdenter:
		*feedback = handleCommand(active,
		                          cmdline,
		                          contentLineCount,
		                          scroll)
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
	var content string
	var contentLines []string
	var contentHeight int
	var coucfg = cfgFromFile()
	var err error
	var feedback string = fmt.Sprintf("Welcome to %v %v", AppName, AppVersion)
	var lower string
	var scroll int = 0
	var termH, termW int
	var title string

	termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Could not get term size: %v", err))
	}

	content, active = handleArgs(&title)
	contentLines = common.SplitByLines(termW, content)

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
		                             &feedback,
		                             termW)

		common.DrawUpper(comcfg, coucfg.Header, title)

		contentHeight = termH -
		                len(common.SplitByLines(termW, coucfg.Header)) -
		                1 -
		                len(common.SplitByLines(termW, title)) -
		                1

		drawContent(contentLines, contentHeight, coucfg, scroll)

		common.SetCursor(1, termH)
		fmt.Printf("%v", lower)

		handleInput(&active,
		            &cmdline,
		            &cmdmode,
		            comcfg,
		            len(contentLines),
		            &feedback,
		            &scroll)
	}
}
