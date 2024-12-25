package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Theme = (*AppTheme)(nil)

const (
	ColorNameGraphGreen = "goxray_theme_graph_green"
	ColorNameGraphBlue  = "goxray_theme_graph_blue"
	ColorNameTextMuted  = "goxray_theme_text_muted"
)

// AppTheme represents small styling additions to fyne default theme.
type AppTheme struct{}

func (m AppTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if v == theme.VariantDark {
		switch c {
		case theme.ColorNameBackground:
			return color.RGBA{R: 0x29, G: 0x29, B: 0x29, A: 0xff}
		case theme.ColorNameSeparator:
			return color.RGBA{R: 0x19, G: 0x19, B: 0x19, A: 0x55}
		case theme.ColorNamePrimary:
			return color.RGBA{R: 50, G: 107, B: 199, A: 255}
		case theme.ColorNameSuccess:
			return color.RGBA{50, 199, 67, 255}
		case theme.ColorNameError:
			return color.RGBA{R: 199, G: 50, B: 50, A: 255}
		case theme.ColorNameForegroundOnPrimary, // Buttons text color
			theme.ColorNameForegroundOnError,
			theme.ColorNameForegroundOnWarning,
			theme.ColorNameForegroundOnSuccess:

			return color.RGBA{R: 240, G: 240, B: 240, A: 255}
		}
	}

	if v == theme.VariantLight {
		switch c {
		case theme.ColorNameSuccess:
			return color.RGBA{49, 181, 64, 255}
		}
	}

	// Custom colors:
	switch v {
	case theme.VariantDark:
		switch c {
		case ColorNameGraphGreen:
			return color.RGBA{62, 194, 84, 255}
		case ColorNameGraphBlue:
			return color.RGBA{62, 104, 240, 255}
		}

	case theme.VariantLight:
		switch c {
		case ColorNameGraphGreen:
			return color.RGBA{89, 217, 110, 255}
		case ColorNameGraphBlue:
			return color.RGBA{105, 150, 255, 255}
		}
	}
	switch c {
	case ColorNameTextMuted:
		return color.RGBA{255, 255, 255, 180}
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
