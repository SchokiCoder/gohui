// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./geninfo.go
package main

import (
	"github.com/SchokiCoder/gohui/common"
	"github.com/SchokiCoder/gohui/csi"

	"errors"
	"fmt"
	"golang.org/x/term"
	"io"
	"os"
	"strconv"
	"strings"
)

type couRuntime struct {
	AcceptInput bool
	Active bool
	CmdLine string
	CmdLineCursor int
	CmdLineInsert bool
	CmdMode bool
	Comcfg common.ComConfig
	Content string
	Coucfg couConfig
	Scroll int
	Feedback string
	Title string
}

func newCouRuntime() couRuntime {
	return couRuntime {
		AcceptInput: true,
		Active: true,
		CmdLine: "",
		CmdLineCursor: 0,
		CmdLineInsert: false,
		CmdMode: false,
		Comcfg: common.ComConfigFromFile(),
		Content: "",
		Coucfg: couConfigFromFile(),
		Scroll: 0,
		Feedback: "",
		Title: "",
	}
}

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
                 runtime       couRuntime,
                 termW         int) {
	var drawRange int = runtime.Scroll + contentHeight

	if drawRange > len(contentLines) {
		drawRange = len(contentLines)
	}

	for _, v := range contentLines[runtime.Scroll:drawRange] {
		common.Cprinta(runtime.Coucfg.Content.Alignment,
		               runtime.Coucfg.Content.Fg,
		               runtime.Coucfg.Content.Bg,
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

func handleCommand(contentLineCount int, runtime *couRuntime) string {
	var err error
	var ret string = ""
	var num uint64

	cmdLineParts := strings.SplitN(runtime.CmdLine, " ", 2)
	fn := couCommands[cmdLineParts[0]]
	if fn != nil {
		return fn(cmdLineParts[1], runtime)
	}

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
	runtime.CmdLineCursor = 0
	runtime.CmdLineInsert = false
	return ret
}

func handleInput(contentLineCount int, runtime *couRuntime) {
	var canonicalState *term.State
	var err error
	var input string
	var rawInput = make([]byte, 4)
	var rawInputLen int

	if runtime.AcceptInput == false {
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

	handleKey(string(input), contentLineCount, runtime)
}

func handleKey(key string, contentLineCount int, runtime *couRuntime) {
	if runtime.CmdMode {
		handleKeyCmdline(key, contentLineCount, runtime)
		return
	}
	
	switch key {
	case csi.CURSOR_UP:
		fallthrough
	case runtime.Comcfg.Keys.Up:
		if runtime.Scroll > 0 {
			runtime.Scroll--
		}

	case csi.CURSOR_DOWN:
		fallthrough
	case runtime.Comcfg.Keys.Down:
		if runtime.Scroll < contentLineCount {
			runtime.Scroll++
		}

	case runtime.Comcfg.Keys.Cmdmode:
		runtime.CmdMode = true
		fmt.Printf(csi.CURSOR_SHOW)

	case csi.HOME:
		runtime.Scroll = 0

	case csi.END:
		runtime.Scroll = contentLineCount - 1

	case runtime.Comcfg.Keys.Quit:
		fallthrough
	case csi.CURSOR_LEFT:
		fallthrough
	case runtime.Comcfg.Keys.Left:
		fallthrough
	case csi.SIGINT:
		fallthrough
	case csi.SIGTSTP:
		runtime.Active = false
	}
}

func handleKeyCmdline(key string, contentLineCount int, runtime *couRuntime) {
	switch key {
	case runtime.Comcfg.Keys.Cmdenter:
		runtime.Feedback = handleCommand(contentLineCount, runtime)
		fallthrough
	case csi.SIGINT:
		fallthrough
	case csi.SIGTSTP:
		runtime.CmdMode = false
		runtime.CmdLine = ""
		runtime.CmdLineCursor = 0
		runtime.CmdLineInsert = false
		fmt.Printf(csi.CURSOR_HIDE)

	case csi.BACKSPACE:
		if runtime.CmdLineCursor > 0 {
			runtime.CmdLine = runtime.CmdLine[:runtime.CmdLineCursor - 1] +
			                  runtime.CmdLine[runtime.CmdLineCursor:]
			runtime.CmdLineCursor--
		}

	case csi.CURSOR_RIGHT:
		if runtime.CmdLineCursor < len(runtime.CmdLine) {
			runtime.CmdLineCursor++
		}

	case csi.CURSOR_LEFT:
		if runtime.CmdLineCursor > 0 {
			runtime.CmdLineCursor--
		}

	case csi.HOME:
		runtime.CmdLineCursor = 0

	case csi.INSERT:
		runtime.CmdLineInsert = !runtime.CmdLineInsert

	case csi.DELETE:
		if runtime.CmdLineCursor < len(runtime.CmdLine) {
			runtime.CmdLine = runtime.CmdLine[:runtime.CmdLineCursor] +
			                  runtime.CmdLine[runtime.CmdLineCursor + 1:]
		}

	case csi.END:
		runtime.CmdLineCursor = len(runtime.CmdLine)

	default:
		if len(key) == 1 {
			var insertReplace = 0

			if runtime.CmdLineInsert == true &&
			   runtime.CmdLineCursor < len(runtime.CmdLine) {
				insertReplace = 1
			}

			runtime.CmdLine = runtime.CmdLine[:runtime.CmdLineCursor] +
				          key +
				          runtime.CmdLine[runtime.CmdLineCursor +
				                          insertReplace:]
			runtime.CmdLineCursor++
		}
	}
}

func tick(runtime *couRuntime) {
	var contentLines []string
	var contentHeight int
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

	headerLines = common.SplitByLines(termW, runtime.Coucfg.Header)
	titleLines = common.SplitByLines(termW, runtime.Title)
	contentLines = common.SplitByLines(termW, runtime.Content)
	lower = common.GenerateLower(runtime.CmdLine,
	                             runtime.CmdMode,
	                             runtime.Comcfg,
	                             &runtime.Feedback,
	                             runtime.Coucfg.Pager.Title,
	                             termW)

	common.DrawUpper(runtime.Comcfg, headerLines, termW, titleLines)

	contentHeight = termH -
	                len(common.SplitByLines(termW, runtime.Coucfg.Header)) -
	                1 -
	                len(common.SplitByLines(termW, runtime.Title)) -
	                1

	drawContent(contentLines, contentHeight, *runtime, termW)

	csi.SetCursor(1, termH)
	fmt.Printf("%v", lower)
	csi.SetCursor((len(runtime.Comcfg.CmdLine.Prefix) + runtime.CmdLineCursor + 1),
	              termH)

	handleInput(len(contentLines), runtime)
}

func main() {
	var runtime = newCouRuntime()

	runtime.Content, runtime.Active = handleArgs(&runtime.Title)
	if runtime.Active == false {
		return
	}

	fmt.Printf(csi.CURSOR_HIDE)
	defer fmt.Printf(csi.CURSOR_SHOW)
	defer fmt.Printf("%v%v\n", csi.FG_DEFAULT, csi.BG_DEFAULT)

	if runtime.Coucfg.Events.Start != "" {
		couFuncs[runtime.Coucfg.Events.Start](&runtime)
	}

	for runtime.Active {
		tick(&runtime)
	}

	if runtime.Coucfg.Events.Quit != "" {
		couFuncs[runtime.Coucfg.Events.Quit](&runtime)
		tick(&runtime)
	}
}
