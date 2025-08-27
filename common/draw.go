// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024 - 2025  Andy Frank Schoknecht

package common

import (
	"fmt"
	"strings"
)

func DrawUpper(
	comCfg ComConfig,
	header []string,
	termW int,
	title []string,
) {
	for _, v := range header {
		Cprinta(comCfg.Header.Alignment,
			comCfg.Header.Fg,
			comCfg.Header.Bg,
			termW,
			v)
	}

	for _, v := range title {
		Cprinta(comCfg.Title.Alignment,
			comCfg.Title.Fg,
			comCfg.Title.Bg,
			termW,
			v)
	}
}

func GenerateLower(
	cmdline    string,
	cmdMode    bool,
	comCfg     ComConfig,
	fb         *Feedback,
	pagerTitle string,
	termW      int,
) string {
	var (
		fits bool
		ret  string
	)

	if cmdMode == true {
		ret = Csprintfa(comCfg.CmdLine.Alignment,
			comCfg.CmdLine.Fg,
			comCfg.CmdLine.Bg,
			termW,
			"%v%v",
			comCfg.CmdLine.Prefix,
			cmdline)
	} else {
		ret, fits = tryFitFeedback(*fb, comCfg.Feedback.Prefix, termW)
		if fits == false {
			ret = string(callPager(*fb, comCfg.Pager.Name, pagerTitle))
			*fb = ""
			ret, _ = tryFitFeedback(
				Feedback(ret),
				comCfg.Feedback.Prefix,
				termW)
		}

		ret = Csprintfa(comCfg.Feedback.Alignment,
			comCfg.Feedback.Fg,
			comCfg.Feedback.Bg,
			termW,
			"%v",
			ret)
	}

	return ret
}

func PrintAbout(
	appLicense,
	appLicenseUrl,
	appName,
	appNameFormal,
	appRepo,
	appVersion string,
) {
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

func PrintVersion(
	appName string,
	appVersion string,
) {
	fmt.Printf("%v: version %v\n", appName, appVersion)
}

func SplitByLines(
	maxLineLen int,
	str string,
) []string {
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

func tryFitFeedback(
	fb       Feedback,
	fbPrefix string,
	termW    int,
) (string, bool) {
	var (
		retStr  string
		retFits bool
	)

	retStr = strings.TrimSpace(string(fb))
	retStr = fmt.Sprintf("%v%v", fbPrefix, retStr)

	if len(SplitByLines(termW, retStr)) > 1 {
		retStr = fbPrefix
		retFits = false
	} else {
		retFits = true
	}

	return retStr, retFits
}
