// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./geninfo.go
package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/scripts"

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
                 scroll        int,
                 termW         int) {
	var drawRange int = scroll + contentHeight

	if drawRange > len(contentLines) {
		drawRange = len(contentLines)
	}

	for _, v := range contentLines[scroll:drawRange] {
		common.Cprinta(coucfg.ContentAlignment,
		               coucfg.ContentFg,
		               coucfg.ContentBg,
		               termW,
		               v)
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

func handleCommand(contentLineCount int, runtime *common.CouRuntime) string {
	var err error
	var ret string = ""
	var num uint64
	
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
			if int(num) < contentLineCount {
				runtime.Scroll = int(num)
			} else {
				runtime.Scroll = contentLineCount
			}
		}
	}
	
	runtime.CmdLine = ""
	return ret
}

func handleInput(comcfg           common.ComCfg,
                 contentLineCount int,
                 runtime          *common.CouRuntime) {
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
		handleKey(string(input),
		          comcfg,
		          contentLineCount,
		          runtime)
	}
}

func handleKey(key              string,
               comcfg           common.ComCfg,
               contentLineCount int,
               runtime          *common.CouRuntime) {
	if runtime.CmdMode {
		handleKeyCmdline(key, comcfg, contentLineCount, runtime)
		return
	}
	
	switch key {
	case comcfg.KeyUp:
		if runtime.Scroll > 0 {
			runtime.Scroll--
		}

	case comcfg.KeyDown:
		if runtime.Scroll < contentLineCount {
			runtime.Scroll++
		}


	case comcfg.KeyCmdmode:
		runtime.CmdMode = true
		fmt.Printf(common.SEQ_CRSR_SHOW)

	case comcfg.KeyQuit:
		fallthrough
	case comcfg.KeyLeft:
		fallthrough
	case common.SIGINT:
		fallthrough
	case common.SIGTSTP:
		runtime.Active = false
	}
}

func handleKeyCmdline(key              string,
                      comcfg           common.ComCfg,
                      contentLineCount int,
                      runtime          *common.CouRuntime) {
	switch key {
	case comcfg.KeyCmdenter:
		runtime.Feedback = handleCommand(contentLineCount, runtime)
		fallthrough
	case common.SIGINT:
		fallthrough
	case common.SIGTSTP:
		runtime.CmdMode = false
		runtime.CmdLine = ""
		fmt.Printf(common.SEQ_CRSR_HIDE)

	default:
		runtime.CmdLine = fmt.Sprintf("%v%v", runtime.CmdLine, string(key))
	}
}

func tick(comcfg common.ComCfg,
          content string,
          coucfg *CouCfg,
          runtime *common.CouRuntime,
          title string) {
	var contentLines []string
	var contentHeight int
	var err error
	var headerLines []string
	var lower string
	var termH, termW int
	var titleLines []string

	fmt.Print(common.SEQ_CLEAR)
	termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Could not get term size: %v", err))
	}

	headerLines = common.SplitByLines(termW, coucfg.Header)
	titleLines = common.SplitByLines(termW, title)
	contentLines = common.SplitByLines(termW, content)
	lower = common.GenerateLower(runtime.CmdLine,
	                             runtime.CmdMode,
	                             comcfg,
	                             &runtime.Feedback,
	                             coucfg.PagerTitle,
	                             termW)

	common.DrawUpper(comcfg, headerLines, termW, titleLines)

	contentHeight = termH -
	                len(common.SplitByLines(termW, coucfg.Header)) -
	                1 -
	                len(common.SplitByLines(termW, title)) -
	                1

	drawContent(contentLines, contentHeight, *coucfg, runtime.Scroll, termW)

	common.SetCursor(1, termH)
	fmt.Printf("%v", lower)

	handleInput(comcfg, len(contentLines), runtime)
}

func main() {
	var comcfg = common.CfgFromFile()
	var content string
	var coucfg = cfgFromFile()
	var runtime = common.NewCouRuntime()
	var title string

	content, runtime.Active = handleArgs(&title)

	fmt.Printf(common.SEQ_CRSR_HIDE)
	defer fmt.Printf(common.SEQ_CRSR_SHOW)
	defer fmt.Printf("%v%v", common.SEQ_FG_DEFAULT, common.SEQ_BG_DEFAULT)

	if coucfg.GoStart != "" {
		scripts.CouFuncs[coucfg.GoStart](&runtime)
	}

	for runtime.Active {
		tick(comcfg, content, &coucfg, &runtime, title)
	}

	if coucfg.GoQuit != "" {
		scripts.CouFuncs[coucfg.GoQuit](&runtime)
		tick(comcfg, content, &coucfg, &runtime, title)
	}
}
