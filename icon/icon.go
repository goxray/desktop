/*
Package icon provides compiled icons for the application.
*/
package icon

import (
	"embed"
	_ "embed"
	"runtime"

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
	LogoPassive  = PrepareResource("icon_default.svg")
	LogoActive   = PrepareResource("icon_active.svg")
	Warning      = PrepareResource("warn.svg")
	Settings     = PrepareResource("settings.svg")
	LinkOn       = PrepareResource("link.svg")
	LinkProgress = PrepareResource("loading.svg")
	LinkOff      = PrepareResource("link_off.svg")
	ListActive   = PrepareResource("list_active.svg")
)

func init() {
	// Use svg icons for tray icon only on darwin.
	if runtime.GOOS != "darwin" {
		LogoPassive = PrepareResource("icon_default.png")
		LogoActive = PrepareResource("icon_active.png")
		Warning = PrepareResource("warn.png")
	}

	LinkProgress = theme.NewThemedResource(LinkProgress)
	Settings = theme.NewThemedResource(Settings)
}

func PrepareResource(path string) fyne.Resource {
	b, err := assetsFS.ReadFile("assets/" + path)
	if err != nil {
		panic(err)
	}

	return &fyne.StaticResource{StaticName: path, StaticContent: b}
}
