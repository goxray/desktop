package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Theme = (*AppTheme)(nil)

// AppTheme represents small styling additions to fyne default theme.
type AppTheme struct{}

func (m AppTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if v == theme.VariantDark {
		switch c {
		case theme.ColorNameBackground:
			return color.RGBA{R: 0x29, G: 0x29, B: 0x29, A: 0xff}
		case theme.ColorNameSeparator:
			return color.RGBA{R: 0x19, G: 0x19, B: 0x19, A: 0x55}
		}
	}

	return theme.DefaultTheme().Color(c, v)
}

func (m AppTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m AppTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (m AppTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameInnerPadding {
		return theme.DefaultTheme().Size(name) * 0.7
	}

	return theme.DefaultTheme().Size(name) * 0.9
}
