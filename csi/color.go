// SPDX-License-Identifier: GPL-2.0-or-later
// Copyright (C) 2024 - 2025  Andy Frank Schoknecht

package csi

import (
	"fmt"
)

type FgColor struct {
	Active bool
	R, G, B uint
}

func (c FgColor) String(
) string {
	if c.Active {
		return fmt.Sprintf("\x1b[38;2;%v;%v;%vm", c.R, c.G, c.B)
	} else {
		return FgDefault
	}
}

type BgColor struct {
	Active bool
	R, G, B uint
}

func (c BgColor) String(
) string {
	if c.Active {
		return fmt.Sprintf("\x1b[48;2;%v;%v;%vm", c.R, c.G, c.B)
	} else {
		return BgDefault
	}
}
