// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package common

import (
	"fmt"
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

func DrawUpper(cfg ComCfg, header string, title string) {
	fmt.Printf("%v%v%v\n", cfg.HeaderFg, cfg.HeaderBg, header)
	fmt.Printf("%v%v%v\n", cfg.TitleFg, cfg.TitleBg, title)
}

func GenerateLower(cmdline  string,
                   cmdmode  bool,
                   comcfg   ComCfg,
                   feedback string,
                   termW    int)    string {
	var ret string
	
	if cmdmode == true {
		ret = fmt.Sprintf("%v%v%v%v",
			          comcfg.CmdlineFg,
			          comcfg.CmdlineBg,
			          comcfg.CmdlinePrefix,
			          cmdline)
	} else {
		feedback = strings.TrimSpace(feedback)
		ret = fmt.Sprintf("%v%v", comcfg.FeedbackPrefix, feedback)
		if strNeededLines(ret, termW) > 1 {
			// TODO will become a call to the pager later
			ret = comcfg.FeedbackPrefix
		}
		
		ret = fmt.Sprintf("%v%v%v",
		                  comcfg.FeedbackFg,
		                  comcfg.FeedbackBg,
		                  ret)
	}
	
	return ret
}

func SetCursor(x, y int) {
	fmt.Printf("\033[%v;%vH", y, x);
}

func strNeededLines(s string, termW int) uint {
	var ret uint = 0
	var line int = 0

	for i := 0; i < len(s); i++ {
		line++

		if line >= termW {
			line = 0
			ret++
		}
	}

	if line > 0 {
		ret++
	}

	return ret
}
