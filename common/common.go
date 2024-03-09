// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package common

import (
	"os/exec"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	SIGINT  = "\003"
	SIGTSTP = "\004"

	SEQ_CLEAR      = "\033[H\033[2J"
	SEQ_FG_DEFAULT = "\033[39m"
	SEQ_BG_DEFAULT = "\033[49m"
	SEQ_CRSR_HIDE  = "\033[?25l"
	SEQ_CRSR_SHOW  = "\033[?25h"
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
	if tempFileContent[len(tempFileContent) - 1] != '\n' {
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

func DrawUpper(cfg ComCfg, header []string, termW int, title []string) {
	for _, v := range header {
		Cprinta(cfg.HeaderAlignment,
		        cfg.HeaderFg,
		        cfg.HeaderBg,
		        termW,
		        v)
	}

	for _, v := range title {
		Cprinta(cfg.TitleAlignment,
		        cfg.TitleFg,
		        cfg.TitleBg,
		        termW,
		        v)
	}
}

func GenerateLower(cmdline    string,
                   cmdmode    bool,
                   comcfg     ComCfg,
                   feedback   *string,
                   pagerTitle string,
                   termW      int)    string {
	var fits bool
	var ret string
	
	if cmdmode == true {
		ret = Csprintfa(comcfg.CmdlineAlignment,
		                comcfg.CmdlineFg,
			        comcfg.CmdlineBg,
			        termW,
			        "%v%v",
			        comcfg.CmdlinePrefix,
			        cmdline)
	} else {
		ret, fits = tryFitFeedback(*feedback, comcfg.FeedbackPrefix, termW)
		if fits == false {
			ret = callPager(*feedback, comcfg.AppPager, pagerTitle)
			*feedback = ""
			ret, _ = tryFitFeedback(ret, comcfg.FeedbackPrefix, termW)
		}

		ret = Csprintfa(comcfg.FeedbackAlignment,
		                comcfg.FeedbackFg,
		                comcfg.FeedbackBg,
		                termW,
		                "%v",
		                ret)
	}

	return ret
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

	fmt.Printf("%v", SEQ_CLEAR)

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
	fmt.Printf(
`The source code of "%v" aka %v %v is available, licensed under the %v at:
%v

If you did not receive a copy of the license, see below:
%v
`,
	           appNameFormal, appName, appVersion, appLicense,
	           appRepo,
	           appLicenseUrl);
}

func PrintVersion(appName, appVersion string) {
	fmt.Printf("%v: version %v\n", appName, appVersion)
}

func SetCursor(x, y int) {
	fmt.Printf("\033[%v;%vH", y, x);
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
			step2 = append(step2, v[lastCut:lastCut + maxLineLen])
			lastCut += maxLineLen
		}
		step2 = append(step2, v[lastCut:])
	}

	return step2
}

func tryFitFeedback(feedback       string,
                    feedbackPrefix string,
                    termW          int)    (string, bool) {
	var retStr string
	var retFits bool

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

type CouRuntime struct {
	AcceptInput bool
	Active bool
	CmdLine string
	CmdMode bool
	Scroll int
	Feedback string
}

func NewCouRuntime() CouRuntime {
	return CouRuntime {
		AcceptInput: true,
		Active: true,
		CmdLine: "",
		CmdMode: false,
		Scroll: 0,
		Feedback: "",
	}
}

type HuiRuntime struct {
	AcceptInput bool
	Active bool
	CmdLine string
	CmdMode bool
	Cursor int
	Feedback string
}

func NewHuiRuntime() HuiRuntime {
	return HuiRuntime {
		AcceptInput: true,
		Active: true,
		CmdLine: "",
		CmdMode: false,
		Cursor: 0,
		Feedback: "",
	}
}
