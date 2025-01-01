package theme

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var _ fyne.Theme = (*AppTheme)(nil)

const (
	ColorNameGraphGreen     = "goxray_theme_graph_green"
	ColorNameGraphBlue      = "goxray_theme_graph_blue"
	ColorNameTextMuted      = "goxray_theme_text_muted"
	ColorNameTextErrorMuted = "goxray_theme_text_error_muted"
)

// AppTheme represents small styling additions to fyne default theme.
type AppTheme struct {
	Variant fyne.ThemeVariant
}

func (m AppTheme) Color(c fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	v := m.Variant

	if v == theme.VariantDark {
		switch c {
		case theme.ColorNameBackground:
			return color.RGBA{R: 30, G: 31, B: 34, A: 255}
		case theme.ColorNameSeparator:
			return color.RGBA{R: 50, G: 51, B: 54, A: 255}
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
			return color.RGBA{R: 237, G: 72, B: 66, A: 255}
		case theme.ColorNamePrimary:
			return color.RGBA{R: 84, G: 139, B: 235, A: 255}
		}
	}

	return m.customColor(c, v)
}

func (m AppTheme) customColor(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case ColorNameTextMuted:
		switch v {
		case theme.VariantDark:
			return color.RGBA{R: 140, G: 140, B: 140, A: 255}
		case theme.VariantLight:
			return color.RGBA{R: 160, G: 160, B: 160, A: 255}
		}
	case ColorNameTextErrorMuted:
		switch v {
		case theme.VariantDark:
			return color.RGBA{194, 118, 118, 255}
		case theme.VariantLight:
			return color.RGBA{247, 139, 139, 255}
		}
	case ColorNameGraphGreen:
		switch v {
		case theme.VariantDark:
			return color.RGBA{R: 62, G: 194, B: 84, A: 255}
		case theme.VariantLight:
			return color.RGBA{R: 89, G: 217, B: 110, A: 255}
		}
	case ColorNameGraphBlue:
		switch v {
		case theme.VariantDark:
			return color.RGBA{R: 62, G: 104, B: 240, A: 255}
		case theme.VariantLight:
			return color.RGBA{R: 105, G: 150, B: 255, A: 255}
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
