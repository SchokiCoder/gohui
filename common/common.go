// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024 - 2025  Andy Frank Schoknecht

package common

import (
	"github.com/SchokiCoder/gohui/csi"

	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type ComAppData struct {
	AcceptInput bool
	Active      bool
	CmdLine     CmdLine
	ComCfg      ComConfig
	Fb          Feedback
}

func NewComAppData(
	customPath string,
) ComAppData {
	return ComAppData {
		AcceptInput:   true,
		Active:        true,
		CmdLine:       NewCmdLine(),
		ComCfg:        ComConfigFromFile(customPath),
		Fb:            "",
	}
}

type CmdLine struct {
	Active  bool
	Content string
	Cursor  int
	Insert  bool
	RowIdx  int
	Rows    [CmdlineMaxRows]string
}

func NewCmdLine(
) CmdLine {
	return CmdLine {
		Active:  false,
		Content: "",
		Cursor:  0,
		Insert:  false,
		RowIdx:  -1,
	}
}

type (
	Feedback     string
	ScriptCmd    func(cmd string) Feedback
	ScriptFn     func()
	ScriptCmdMap map[string]ScriptCmd
	ScriptFnMap  map[string]ScriptFn
)

const (
	CmdlineMaxRows = 10
)

func callPager(
	fb Feedback,
	pager string,
	pagerTitle string,
) Feedback {
	var (
		err error
		flags string
		shCall string
		tempFile *os.File
		tempFileContent string
		tempFilePath string
	)

	tempFile, err = os.CreateTemp("", "huiFeedback")
	if err != nil {
		panic("Could not create a temp file for feedback")
	}
	defer os.Remove(tempFile.Name())
	tempFilePath = tempFile.Name()

	tempFileContent = string(fb)
	if tempFileContent[len(tempFileContent)-1] != '\n' {
		tempFileContent = fmt.Sprintf("%v\n", tempFileContent)
	}

	_, err = io.WriteString(tempFile, tempFileContent)
	if err != nil {
		panic("Could not write feedback to temp file")
	}

	pagerTitle = strings.ReplaceAll(pagerTitle, "\"", "\\\"")
	if pager == "less" {
		flags = "-fr"
	}

	shCall = fmt.Sprintf("PAGERTITLE=\"%v\" %v %v %v", pagerTitle, pager, flags, tempFilePath)

	return HandleShellSession(shCall)
}

func handleCommand(
	active           *bool,
	cmdLine          CmdLine,
	contentLineCount int,
	cursor           *int,
	customCmds       map[string]ScriptCmd,
) Feedback {
	var (
		cmdLineParts []string
		err          error
		fn           ScriptCmd
		num          uint64
		ret          string = ""
	)

	if cmdLine.Content == "" {
		return ""
	}

	cmdLineParts = strings.SplitN(cmdLine.Content, " ", 2)
	fn = customCmds[cmdLineParts[0]]
	if fn != nil {
		return fn(cmdLineParts[1])
	}

	switch cmdLine.Content {
	case "q":
		fallthrough
	case "quit":
		fallthrough
	case "exit":
		*active = false

	default:
		num, err = strconv.ParseUint(cmdLine.Content, 10, 32)

		if err != nil {
			ret = fmt.Sprintf("Command \"%v\" not recognised",
				cmdLine.Content)
		} else {
			if int(num) < contentLineCount {
				*cursor = int(num)
			} else {
				*cursor = contentLineCount - 1
			}
		}
	}

	for i := 0; i < len(cmdLine.Rows)-1; i++ {
		cmdLine.Rows[len(cmdLine.Rows)-1-i] =
			cmdLine.Rows[len(cmdLine.Rows)-1-i-1]
	}
	cmdLine.Rows[0] = cmdLine.Content
	return Feedback(ret)
}

func HandleKeyCmdline(
	key              string,
	active           *bool,
	cmdLine          *CmdLine,
	cmdMap           ScriptCmdMap,
	comCfg           *ComConfig,
	contentLineCount int,
	cursor           *int,
	fb               *Feedback,
) {
	switch key {
	case comCfg.Keys.Cmdenter:
		*fb = handleCommand(active,
			*cmdLine,
			contentLineCount,
			cursor,
			cmdMap)
		fallthrough
	case csi.SigInt:
		fallthrough
	case csi.SigTstp:
		*cmdLine = NewCmdLine()
		fmt.Printf(csi.CursorHide)

	case csi.Backspace:
		if cmdLine.Cursor > 0 {
			cmdLine.Content =
				(cmdLine.Content)[:cmdLine.Cursor-1] +
				(cmdLine.Content)[cmdLine.Cursor:]
			cmdLine.Cursor--
		}

	case csi.CursorRight:
		if cmdLine.Cursor < len(cmdLine.Content) {
			cmdLine.Cursor++
		}

	case csi.CursorUp:
		if cmdLine.RowIdx < len(cmdLine.Rows)-1 {
			if cmdLine.Rows[cmdLine.RowIdx+1] != "" {
				cmdLine.RowIdx++
				cmdLine.Content = cmdLine.Rows[cmdLine.RowIdx]
				cmdLine.Cursor = len(cmdLine.Content)
			}
		}

	case csi.CursorLeft:
		if cmdLine.Cursor > 0 {
			cmdLine.Cursor--
		}

	case csi.CursorDown:
		if cmdLine.RowIdx >= 0 {
			cmdLine.RowIdx--
		}
		if cmdLine.RowIdx >= 0 {
			cmdLine.Content = cmdLine.Rows[cmdLine.RowIdx]
		} else {
			cmdLine.Content = ""
		}
		cmdLine.Cursor = len(cmdLine.Content)

	case csi.Home:
		cmdLine.Cursor = 0

	case csi.Insert:
		cmdLine.Insert = !(cmdLine.Insert)

	case csi.Delete:
		if cmdLine.Cursor < len(cmdLine.Content) {
			cmdLine.Content =
				(cmdLine.Content)[:cmdLine.Cursor] +
				(cmdLine.Content)[cmdLine.Cursor+1:]
		}

	case csi.End:
		cmdLine.Cursor = len(cmdLine.Content)

	default:
		if len(key) == 1 {
			var insertReplace = 0

			if cmdLine.Insert == true &&
				cmdLine.Cursor < len(cmdLine.Content) {
				insertReplace = 1
			}

			cmdLine.Content = (cmdLine.Content)[:cmdLine.Cursor] +
				key +
				(cmdLine.Content)[cmdLine.Cursor+insertReplace:]
			cmdLine.Cursor++
		}
	}
}

func HandleShell(
	shell string,
) Feedback {
	var cmd *exec.Cmd
	var cmderr io.ReadCloser
	var cmdout io.ReadCloser
	var err error
	var strerr []byte
	var strout []byte

	cmd = exec.Command("sh", "-c", shell)

	cmderr, err = cmd.StderrPipe()
	if err != nil {
		return Feedback(fmt.Sprintf("Could not get stderr: %s", err))
	}

	cmdout, err = cmd.StdoutPipe()
	if err != nil {
		return Feedback(fmt.Sprintf("Could not get stdout: %s", err))
	}

	err = cmd.Start()
	if err != nil {
		return Feedback(
			fmt.Sprintf("Could not start child process: %s", err))
	}

	strerr, err = io.ReadAll(cmderr)
	if err != nil {
		return Feedback(fmt.Sprintf("Could not read stderr: %s", err))
	}

	strout, err = io.ReadAll(cmdout)
	if err != nil {
		return Feedback(fmt.Sprintf("Could not read stdout: %s", err))
	}

	err = cmd.Wait()
	if err != nil {
		return Feedback(fmt.Sprintf("Child error: %s", err))
	}

	if len(strerr) > 0 {
		return Feedback(strerr)
	} else {
		return Feedback(strout)
	}
}

func HandleShellSession(
	shell string,
) Feedback {
	var cmd *exec.Cmd
	var cmderr io.ReadCloser
	var err error
	var strerr []byte

	cmd = exec.Command("sh", "-c", shell)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	cmderr, err = cmd.StderrPipe()
	if err != nil {
		return Feedback(fmt.Sprintf("Could not get stderr: %s", err))
	}

	err = cmd.Start()
	if err != nil {
		return Feedback(
			fmt.Sprintf("Could not start child process: %s", err))
	}

	strerr, err = io.ReadAll(cmderr)
	if err != nil {
		return Feedback(fmt.Sprintf("Could not read stderr: %s", err))
	}

	fmt.Printf("%v%v\n", csi.FgDefault, csi.BgDefault)
	fmt.Printf(csi.CursorShow)
	defer fmt.Printf(csi.CursorHide)
	defer fmt.Printf("%v", csi.Clear)

	err = cmd.Wait()
	if err != nil {
		return Feedback(fmt.Sprintf("Child error: %s", err))
	}

	if len(strerr) > 0 {
		return Feedback(strerr)
	}

	return ""
}
