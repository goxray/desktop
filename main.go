package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	_ "unsafe"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/systray"
	vpn "github.com/goxray/tun/pkg/client"
	"github.com/lilendian0x00/xray-knife/xray"

	"github.com/goxray/ui/icon"
	"github.com/goxray/ui/internal/osspecific/dock"
	"github.com/goxray/ui/internal/osspecific/root"
	"github.com/goxray/ui/internal/traylist"
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

func init() {
	debug.SetGCPercent(10)
	root.PromptRootAccess()
}

func onstart() {
	systray.SetTooltip(AppTitleName)
	dock.HideIconInDock()
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&AppTheme{})
	a.Lifecycle().SetOnStarted(onstart)

	client, err := vpn.NewClientWithOpts(vpn.Config{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})),
	})
	if err != nil {
		panic(fmt.Errorf("create vpn client: %v", err))
	}

	items := NewItemsList()
	list := binding.BindUntypedList(items.AllUntyped())
	trayMenu := traylist.NewDefault[*Item](AppTitleName, toDesktopApp(a), MenuIcons)
	settingsLoader := NewSaveFile(a.Preferences())

	// Tray menu setup.
	var settingsWindow fyne.Window
	trayMenu.OnSettingsClick(func() {
		if settingsWindow == nil {
			settingsWindow = window.NewSettings[*Item](a, AddFormH(items), UpdateFormH(items), DeleteItemH(items), list)
			settingsWindow.SetOnClosed(func() {
				settingsWindow = nil
			})
		}
		settingsWindow.Show()
		settingsWindow.RequestFocus()
	})
	trayMenu.OnItemClick(ConnectHandler(trayMenu, client))

	// Update all UI elements when items are updated.
	items.OnAdd(func(item *Item) {
		trayMenu.Add(item.Label(), item)
		if err := list.Append(item); err != nil {
			slog.Warn(err.Error())
		}
	})
	items.OnDelete(func(item *Item) {
		err := errors.Join(trayMenu.Remove(item), list.Remove(item))
		if err != nil {
			slog.Error(err.Error())
		}
	})
	items.OnChange(func() {
		trayMenu.Refresh()
		settingsLoader.Update(items)
	})

	settingsLoader.Load(items) // Initialize items from savefile and update windows/tray with new items.

	trayMenu.Show()
	a.Run()
}

func DeleteItemH(list *ItemsList) func(itm *Item) error {
	return func(itm *Item) error {
		list.RemoveItem(itm)

		return nil
	}
}

func UpdateFormH(list *ItemsList) func(data window.FormData, itm *Item) error {
	return func(updated window.FormData, item *Item) error {
		proto, err := xray.ParseXrayConfig(updated.Link)
		if err != nil {
			return err
		}

		item.SetXRayConfig(proto.ConvertToGeneralConfig())
		item.LinkVal = updated.Link
		item.LabelVal = updated.Label
		list.UpdateItem()

		return nil
	}
}

func AddFormH(list *ItemsList) func(data window.FormData) error {
	return func(new window.FormData) error {
		proto, err := xray.ParseXrayConfig(new.Link)
		if err != nil {
			return err
		}

		list.Add(NewItem(new.Label, new.Link, proto.ConvertToGeneralConfig()))

		return nil
	}
}

func ConnectHandler(trayItems *traylist.List[*Item], client *vpn.Client) func(id int) error {
	return func(id int) error {
		// If clicked item is connected - just disconnect and return.
		if trayItems.IsActive(id) {
			return client.Disconnect(context.Background())
		}

		// Disconnect active connections before connecting the clicked one.
		if trayItems.HasActive() {
			err := client.Disconnect(context.Background())
			if err != nil {
				return err
			}
		}

		return client.Connect(trayItems.Get(id).Link())
	}
}

func toDesktopApp(a fyne.App) desktop.App {
	desk, ok := a.(desktop.App)
	if !ok {
		panic("Only desktop mode supported")
	}

	return desk
}
