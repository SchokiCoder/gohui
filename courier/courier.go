// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

//go:generate go ./genversion.go
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

var Version string

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

func handleArgs(title *string) string {
	var err error
	var f *os.File
	var nextIsTitle = false
	var path string

	if len(os.Args) < 2 {
		panic("No filepath argument given.\n")
	}

	for _, v := range os.Args[1:] {
		switch v {
		case "-t":      fallthrough
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
		panic(fmt.Sprintf("File could not be found: \"%v\", \"%v\"\n",
		                  path,
		                  err))
	} else if err != nil {
		panic(fmt.Sprintf("File could not be opened: \"%v\", \"%v\"\n",
		                  path,
		                  err))
	}

	ret, err := io.ReadAll(f)

	if err != nil {
		panic(fmt.Sprintf("File could not be read: \"%v\", \"%v\"\n",
		                  path,
		                  err))
	}

	return string(ret)
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

	case comcfg.KeyQuit: fallthrough
	case comcfg.KeyLeft: fallthrough
	case common.SIGINT:  fallthrough
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
	var contentLines []string
	var contentHeight int
	var coucfg = cfgFromFile()
	var err error
	var feedback string = fmt.Sprintf("Welcome to courier %v", Version)
	var lower string
	var scroll int = 0
	var termH, termW int
	var title string

	termW, termH, err = term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic(fmt.Sprintf("Could not get term size: %v", err))
	}

	contentLines = common.SplitByLines(termW, handleArgs(&title))

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
