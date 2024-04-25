// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package common

import (
	"github.com/SchokiCoder/gohui/csi"

	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

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

func DrawUpper(comcfg ComConfig, header []string, termW int, title []string) {
	for _, v := range header {
		Cprinta(comcfg.Header.Alignment,
			comcfg.Header.Fg,
			comcfg.Header.Bg,
			termW,
			v)
	}

	for _, v := range title {
		Cprinta(comcfg.Title.Alignment,
			comcfg.Title.Fg,
			comcfg.Title.Bg,
			termW,
			v)
	}
}

func GenerateLower(cmdline string,
	cmdmode bool,
	comcfg ComConfig,
	feedback *string,
	pagerTitle string,
	termW int) string {
	var (
		fits bool
		ret  string
	)

	if cmdmode == true {
		ret = Csprintfa(comcfg.CmdLine.Alignment,
			comcfg.CmdLine.Fg,
			comcfg.CmdLine.Bg,
			termW,
			"%v%v",
			comcfg.CmdLine.Prefix,
			cmdline)
	} else {
		ret, fits = tryFitFeedback(*feedback, comcfg.Feedback.Prefix, termW)
		if fits == false {
			ret = callPager(*feedback, comcfg.Pager.Name, pagerTitle)
			*feedback = ""
			ret, _ = tryFitFeedback(ret, comcfg.Feedback.Prefix, termW)
		}

		ret = Csprintfa(comcfg.Feedback.Alignment,
			comcfg.Feedback.Fg,
			comcfg.Feedback.Bg,
			termW,
			"%v",
			ret)
	}

	return ret
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

	fmt.Printf("%v", csi.CLEAR)

	if len(strerr) > 0 {
		return string(strerr)
	}

	return ""
}

func PrintAbout(appLicense,
	appLicenseUrl,
	appName,
	appNameFormal,
	appRepo,
	appVersion string) {
	fmt.Printf("The source code of \"%v\" aka %v %v is available, "+
		`licensed under the %v at:
%v

If you did not receive a copy of the license, see below:
%v
`,
		appNameFormal, appName, appVersion, appLicense,
		appRepo,
		appLicenseUrl)
}

func PrintVersion(appName, appVersion string) {
	fmt.Printf("%v: version %v\n", appName, appVersion)
}

func SplitByLines(maxLineLen int, str string) []string {
	var step1 []string
	var step2 []string
	var lastCut int

	step1 = strings.Split(str, "\n")

	for _, v := range step1 {
		if len(v) <= maxLineLen {
			step2 = append(step2, v)
			continue
		}

		lastCut = 0
		for len(v[lastCut:]) > maxLineLen {
			step2 = append(step2, v[lastCut:lastCut+maxLineLen])
			lastCut += maxLineLen
		}
		step2 = append(step2, v[lastCut:])
	}

	return step2
}

func tryFitFeedback(feedback string,
	feedbackPrefix string,
	termW int) (string, bool) {
	var (
		retStr  string
		retFits bool
	)

	retStr = strings.TrimSpace(feedback)
	retStr = fmt.Sprintf("%v%v", feedbackPrefix, retStr)

	if len(SplitByLines(termW, retStr)) > 1 {
		retStr = feedbackPrefix
		retFits = false
	} else {
		retFits = true
	}

	return retStr, retFits
}
