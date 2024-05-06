// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

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

type ScriptCmd    func(cmd string) string
type ScriptFn     func()
type ScriptCmdMap map[string]ScriptCmd
type ScriptFnMap  map[string]ScriptFn

const (
	CmdlineMaxRows = 10
)

func callPager(feedback string, pager string, pagerTitle string) string {
	var err error
	var shCall string
	var tempFile *os.File
	var tempFileContent string
	var tempFilePath string

	tempFile, err = os.CreateTemp("", "huiFeedback")
	if err != nil {
		panic("Could not create a temp file for feedback.")
	}
	defer os.Remove(tempFile.Name())
	tempFilePath = tempFile.Name()

	tempFileContent = feedback
	if tempFileContent[len(tempFileContent)-1] != '\n' {
		tempFileContent = fmt.Sprintf("%v\n", tempFileContent)
	}

	_, err = io.WriteString(tempFile, tempFileContent)
	if err != nil {
		panic("Could not write feedback to temp file.")
	}

	if pager == "./pkg/courier" || pager == "courier" {
		shCall = fmt.Sprintf(`%v %v -t "%v"`,
			pager,
			tempFilePath,
			pagerTitle)
	} else {
		shCall = fmt.Sprintf("%v %v", pager, tempFilePath)
	}

	return HandleShellSession(shCall)
}

func handleCommand(active *bool,
	cmdLine           string,
	cmdLineRows       []string,
	contentLineCount  int,
	cursor            *int,
	customCmds        map[string]ScriptCmd) string {
	var (
		cmdLineParts []string
		err          error
		fn           ScriptCmd
		num          uint64
		ret          string = ""
	)

	cmdLineParts = strings.SplitN(cmdLine, " ", 2)
	fn = customCmds[cmdLineParts[0]]
	if fn != nil {
		return fn(cmdLineParts[1])
	}

	switch cmdLine {
	case "q":
		fallthrough
	case "quit":
		fallthrough
	case "exit":
		*active = false

	default:
		num, err = strconv.ParseUint(cmdLine, 10, 32)

		if err != nil {
			ret = fmt.Sprintf("Command \"%v\" not recognised",
				cmdLine)
		} else {
			if int(num) < contentLineCount {
				*cursor = int(num)
			} else {
				*cursor = contentLineCount
			}
		}
	}

	for i := 0; i < len(cmdLineRows)-1; i++ {
		cmdLineRows[len(cmdLineRows)-1-i] =
			cmdLineRows[len(cmdLineRows)-1-i-1]
	}
	cmdLineRows[0] = cmdLine
	return ret
}

func HandleKeyCmdline(key string,
	active *bool,
	cmdLine *string,
	cmdLineCursor *int,
	cmdLineInsert *bool,
	cmdLineRowIdx *int,
	cmdLineRows []string,
	cmdMap ScriptCmdMap,
	cmdMode *bool,
	comCfg *ComConfig,
	contentLineCount int,
	cursor *int,
	feedback *string) {

	switch key {
	case comCfg.Keys.Cmdenter:
		*feedback = handleCommand(active,
			*cmdLine,
			cmdLineRows,
			contentLineCount,
			cursor,
			cmdMap)
		fallthrough
	case csi.SigInt:
		fallthrough
	case csi.SigTstp:
		*cmdLine = ""
		*cmdLineCursor = 0
		*cmdLineInsert = false
		*cmdLineRowIdx = -1
		*cmdMode = false
		fmt.Printf(csi.CursorHide)

	case csi.Backspace:
		if *cmdLineCursor > 0 {
			*cmdLine = (*cmdLine)[:*cmdLineCursor-1] +
				(*cmdLine)[*cmdLineCursor:]
			*cmdLineCursor--
		}

	case csi.CursorRight:
		if *cmdLineCursor < len(*cmdLine) {
			*cmdLineCursor++
		}

	case csi.CursorUp:
		if *cmdLineRowIdx < len(cmdLineRows)-1 {
			*cmdLineRowIdx++
			*cmdLine = cmdLineRows[*cmdLineRowIdx]
			*cmdLineCursor = len(*cmdLine)
		}

	case csi.CursorLeft:
		if *cmdLineCursor > 0 {
			*cmdLineCursor--
		}

	case csi.CursorDown:
		if *cmdLineRowIdx >= 0 {
			*cmdLineRowIdx--
		}
		if *cmdLineRowIdx >= 0 {
			*cmdLine = cmdLineRows[*cmdLineRowIdx]
		} else {
			*cmdLine = ""
		}
		*cmdLineCursor = len(*cmdLine)

	case csi.Home:
		*cmdLineCursor = 0

	case csi.Insert:
		*cmdLineInsert = !(*cmdLineInsert)

	case csi.Delete:
		if *cmdLineCursor < len(*cmdLine) {
			*cmdLine = (*cmdLine)[:*cmdLineCursor] +
				(*cmdLine)[*cmdLineCursor+1:]
		}

	case csi.End:
		*cmdLineCursor = len(*cmdLine)

	default:
		if len(key) == 1 {
			var insertReplace = 0

			if *cmdLineInsert == true &&
				*cmdLineCursor < len(*cmdLine) {
				insertReplace = 1
			}

			*cmdLine = (*cmdLine)[:*cmdLineCursor] +
				key +
				(*cmdLine)[*cmdLineCursor+insertReplace:]
			*cmdLineCursor++
		}
	}
}

func HandleShell(shell string) string {
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

func HandleShellSession(shell string) string {
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

	fmt.Printf("%v", csi.Clear)

	if len(strerr) > 0 {
		return string(strerr)
	}

	return ""
}
