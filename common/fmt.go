// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package common

import (
	"fmt"
)

func Cprintf(fg FgColor, bg BgColor, format string, a ...any) (n int, err error) {
	var output string

	output = Csprintf(fg, bg, format, a...)

	return fmt.Printf(output)
}

func Csprintf(fg FgColor, bg BgColor, format string, a ...any) string {
	var ret string

	ret = fmt.Sprintf(format, a...)
	ret = fmt.Sprintf("%v%v%v%v%v",
	                  fg, bg, ret, SEQ_FG_DEFAULT, SEQ_BG_DEFAULT)

	return ret
}
