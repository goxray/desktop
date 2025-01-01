package main

import (
	"embed"
	"errors"
	"flag"
	"log/slog"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/lang"
	"fyne.io/systray"
	"github.com/lilendian0x00/xray-knife/v2/xray"

	"github.com/goxray/ui/icon"
	"github.com/goxray/ui/internal/connlist"
	"github.com/goxray/ui/internal/osspecific/dock"
	"github.com/goxray/ui/internal/osspecific/root"
	"github.com/goxray/ui/internal/traylist"
	"github.com/goxray/ui/theme"
	"github.com/goxray/ui/window"
)

const (
	AppTitleName = "Go XRay VPN Client"
)

var MenuIcons = &traylist.IconSet{
	LogoActive:  icon.LogoActive,
	LogoPassive: icon.LogoPassive,
	NotSelected: icon.LinkOff,
	InProgress:  icon.LinkProgress,
	Selected:    icon.LinkOn,
	Settings:    icon.Settings,
	Warning:     icon.Warning,
}

//go:embed translation
var translations embed.FS

func initialize() {
	debug.SetGCPercent(10)
	root.PromptRootAccess()
}

func onstart() {
	systray.SetTooltip(AppTitleName)
	dock.HideIconInDock()
}

func main() {
	a := app.New()
	initialize()
	a.Settings().SetTheme(&theme.AppTheme{Variant: getThemeVariant()})
	a.Lifecycle().SetOnStarted(onstart)
	if err := lang.AddTranslationsFS(translations, "translation"); err != nil {
		slog.Warn("failed to init translations", "error", err)
	}

	items := connlist.New()
	list := binding.BindUntypedList(items.AllUntyped())
	trayMenu := traylist.NewDefault[*connlist.Item](lang.L(AppTitleName), toDesktopApp(a), MenuIcons)
	settingsLoader := NewSaveFile(a.Preferences())

	// Tray menu setup.
	var settingsWindow *window.Settings[*connlist.Item]
	trayMenu.OnSettingsClick(func() {
		if settingsWindow == nil {
			settingsWindow = window.NewSettings[*connlist.Item](a, list, AddFormH(items), UpdateFormH(), DeleteItemH(items))
			settingsWindow.OnClosed(func() { settingsWindow = nil })
		}
		settingsWindow.Show()
	})
	trayMenu.OnItemClick(ConnectHandler(trayMenu))
	trayMenu.Show()

	// Update all UI elements when items are updated.
	items.OnAdd(func(item *connlist.Item) {
		trayMenu.Add(item.Label(), item)
		if err := list.Append(item); err != nil {
			slog.Warn(err.Error())
		}
	})
	items.OnDelete(func(item *connlist.Item) {
		err := errors.Join(trayMenu.Remove(item), list.Remove(item))
		if err != nil {
			slog.Error(err.Error())
		}
	})
	items.OnChange(func() {
		trayMenu.Refresh()
		settingsLoader.Update(items)
		if settingsWindow != nil {
			settingsWindow.Refresh()
		}
	})
	// Disconnect any active connections on quitting/panic.
	defer func() {
		if trayMenu.HasActive() {
			if err := trayMenu.GetActive().Disconnect(); err != nil {
				slog.Error(err.Error())
			}
		}
	}()

	settingsLoader.Load(items) // Initialize items from savefile and update windows/tray with new items.
	a.Run()
}

func DeleteItemH(list *connlist.Collection) func(itm *connlist.Item) error {
	return func(itm *connlist.Item) error {
		list.RemoveItem(itm)

		return nil
	}
}

func UpdateFormH() func(data window.FormData, itm *connlist.Item) error {
	return func(updated window.FormData, item *connlist.Item) error {
		return item.Update(updated.Link, updated.Label)
	}
}

func AddFormH(list *connlist.Collection) func(data window.FormData) error {
	return func(new window.FormData) error {
		proto, err := xray.ParseXrayConfig(new.Link)
		if err != nil {
			return err
		}

		if len(proto.ConvertToGeneralConfig().Remark) > 32 {
			return errors.New("remark value is too long")
		}

		return list.AddItem(new.Label, new.Link)
	}
}

func ConnectHandler(trayItems *traylist.List[*connlist.Item]) func(id int) error {
	return func(id int) error {
		// If clicked item is connected - just disconnect and return.
		if trayItems.IsActive(id) {
			return trayItems.Get(id).Disconnect()
		}

		// Disconnect active connections before connecting the clicked one.
		if trayItems.HasActive() {
			err := trayItems.GetActive().Disconnect()
			if err != nil {
				return err
			}
		}

		return trayItems.Get(id).Connect()
	}
}

func toDesktopApp(a fyne.App) desktop.App {
	desk, ok := a.(desktop.App)
	if !ok {
		panic("Only desktop mode supported")
	}

	return desk
}

func getThemeVariant() fyne.ThemeVariant {
	var fThemeVariant = flag.Int("theme_variant", 0, "set app theme variant")
	flag.Bool("xroot", false, "is root")
	flag.Parse()

	if fThemeVariant != nil {
		return fyne.ThemeVariant(*fThemeVariant)
	}

	return fyne.CurrentApp().Settings().ThemeVariant()
}
