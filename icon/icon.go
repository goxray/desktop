/*
Package icon provides compiled icons for the application.
*/
package icon

import (
	"embed"
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

//go:embed *
var assetsFS embed.FS

// Credits: https://www.svgrepo.com/collection/wolf-kit-solid-glyph-icons/3
// License: https://www.svgrepo.com/page/licensing/#CC%20Attribution
//
// All icons were modified, the color and some properties of SVG files were altered.
var (
	LogoPassive  fyne.Resource = &fyne.StaticResource{StaticName: "logo_passive", StaticContent: MustReadFile("icon_default.svg")}
	LogoActive   fyne.Resource = &fyne.StaticResource{StaticName: "logo_active", StaticContent: MustReadFile("icon_active.svg")}
	Warning      fyne.Resource = &fyne.StaticResource{StaticName: "warning", StaticContent: MustReadFile("warn.svg")}
	Settings     fyne.Resource = &fyne.StaticResource{StaticName: "settings", StaticContent: MustReadFile("settings.svg")}
	LinkOn       fyne.Resource = &fyne.StaticResource{StaticName: "link_on", StaticContent: MustReadFile("link.svg")}
	LinkProgress fyne.Resource = &fyne.StaticResource{StaticName: "link_in_progress", StaticContent: MustReadFile("loading.svg")}
	LinkOff      fyne.Resource = &fyne.StaticResource{StaticName: "link_off", StaticContent: MustReadFile("link_off.svg")}
	ListActive   fyne.Resource = &fyne.StaticResource{StaticName: "list_active", StaticContent: MustReadFile("list_active.svg")}
)

func init() {
	LinkProgress = theme.NewThemedResource(LinkProgress)
	Settings = theme.NewThemedResource(Settings)
}

func MustReadFile(path string) []byte {
	b, err := assetsFS.ReadFile("assets/" + path)
	if err != nil {
		panic(err)
	}

	return b
}
