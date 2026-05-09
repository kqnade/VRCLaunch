package ui

import (
	"image/color"

	"gioui.org/font/gofont"
	"gioui.org/text"
	"gioui.org/widget/material"
)

func NewDarkTheme() *material.Theme {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
	th.Palette = material.Palette{
		Bg:         color.NRGBA{R: 0x14, G: 0x16, B: 0x1c, A: 0xff},
		Fg:         color.NRGBA{R: 0xea, G: 0xea, B: 0xee, A: 0xff},
		ContrastBg: color.NRGBA{R: 0x52, G: 0x6b, B: 0xff, A: 0xff},
		ContrastFg: color.NRGBA{R: 0xff, G: 0xff, B: 0xff, A: 0xff},
	}
	return th
}
