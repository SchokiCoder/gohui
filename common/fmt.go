// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024  Andy Frank Schoknecht

package common

import (
	"fmt"
	"strings"
)

func Cprinta(alignment string,
             fg FgColor,
             bg BgColor,
             termW int,
             str string) (n int, err error) {
	return fmt.Printf("%v", Csprinta(alignment, fg, bg, termW, str))
}

func Cprintf(fg FgColor, bg BgColor, format string, a ...any) (n int, err error) {
	var output string

	output = Csprintf(fg, bg, format, a...)

	return fmt.Printf(output)
}

func Cprintfa(alignment string,
              fg FgColor,
              bg BgColor,
              termW int,
              format string,
              a ...any)      (n int, err error) {
	var str = Csprintfa(alignment, fg, bg, termW, format, a...)
	return fmt.Printf("%v", str)
}

func Csprinta(alignment string,
              fg FgColor,
              bg BgColor,
              termW int,
              str string) string {
	var ret = Csprintfa(alignment, fg, bg, termW, "%v", str)

	switch alignment {
	case "left":
		return fmt.Sprintf("%v\n", ret)

	case "center":
		fallthrough
	case "centered":
		return fmt.Sprintf("%v\n", ret)

	case "right":
		return fmt.Sprintf("%v", ret)

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

func Csprintfa(alignment string,
               fg FgColor,
               bg BgColor,
               termW int,
               format string,
               a ...any)      string {
	var str string
	var strlen int
	
	str = fmt.Sprintf(format, a...)
	strlen = len(str)
	str = Csprintf(fg, bg, "%v", str)

	switch alignment {
	case "left":
		return str

	case "center":
		fallthrough
	case "centered":
		return fmt.Sprintf("%v%v",
		                   strings.Repeat(" ", (termW - strlen) / 2),
		                   str)

	case "right":
		return fmt.Sprintf("%v%v",
		                   strings.Repeat(" ", termW - strlen),
		                   str)

	default:
		panic(fmt.Sprintf(`Unknown alignment "%v".`, alignment))
	}
}
