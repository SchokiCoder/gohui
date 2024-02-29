// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package common

import (
	"fmt"
	"strings"
)

func Cprintf(fg FgColor, bg BgColor, format string, a ...any) (n int, err error) {
	var output string

	output = Csprintf(fg, bg, format, a...)

	return fmt.Printf(output)
}

func Cprinta(alignment string,
             fg FgColor,
             bg BgColor,
             termW int,
             str string) (n int, err error) {
	var strlen = len(str)

	str = Csprintf(fg, bg, "%v", str)

	switch alignment {
	case "left":
		return fmt.Printf("%v\n", str)

	case "center":
		fallthrough
	case "centered":
		return fmt.Printf("%v%v\n",
		                  strings.Repeat(" ", (termW - strlen) / 2),
		                  str)

	case "right":
		return fmt.Printf("%v%v",
		                  strings.Repeat(" ", termW - strlen),
		                  str)

	default:
		panic(fmt.Sprintf(`Unknown alignment "%v".`, alignment))
	}
}

func Csprintf(fg FgColor, bg BgColor, format string, a ...any) string {
	var ret string

	ret = fmt.Sprintf(format, a...)
	ret = fmt.Sprintf("%v%v%v%v%v",
	                  fg, bg, ret, SEQ_FG_DEFAULT, SEQ_BG_DEFAULT)

	return ret
}
