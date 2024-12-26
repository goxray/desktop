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
			return color.RGBA{R: 50, G: 199, B: 67, A: 255}
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
			return color.RGBA{R: 49, G: 181, B: 64, A: 255}
		case theme.ColorNameError:
			return color.RGBA{237, 72, 66, 255}
		case theme.ColorNamePrimary:
			return color.RGBA{84, 139, 235, 255}
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
			return color.RGBA{R: 89, G: 217, B: 110, A: 255}
		case ColorNameGraphBlue:
			return color.RGBA{R: 105, G: 150, B: 255, A: 255}
		}
	}
	switch c {
	case ColorNameTextMuted:
		return color.RGBA{R: 160, G: 160, B: 160, A: 255}
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
