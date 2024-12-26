package window

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/goxray/ui/icon"
	customtheme "github.com/goxray/ui/theme"
	"github.com/goxray/ui/window/form"
	customwidget "github.com/goxray/ui/window/widget"
)

type SettingsDraft[T ListItem] struct {
	window fyne.Window
	list   binding.DataList

	onAdd    func(data FormData) error
	onUpdate func(FormData, T) error
	onDelete func(T) error
}

func NewSettingsDraft[T ListItem](
	a fyne.App,
	list binding.DataList,
	onAdd func(data FormData) error,
	onUpdate func(FormData, T) error,
	onDelete func(T) error,
) *SettingsDraft[T] {
	w := a.NewWindow(settingsText)
	w.CenterOnScreen()
	w.RequestFocus()
	w.Resize(fyne.NewSize(750, 450))

	return &SettingsDraft[T]{
		window:   w,
		onAdd:    onAdd,
		onUpdate: onUpdate,
		onDelete: onDelete,
		list:     list,
	}
}

func (w *SettingsDraft[T]) Window() fyne.Window {
	return w.window
}

func (w *SettingsDraft[T]) Show() {
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon( // Connections list settings tab
			configsText,
			icon.Settings,
			w.createSettingsContainer(),
		),
		container.NewTabItemWithIcon( // About tab with static app info
			aboutText,
			theme.QuestionIcon(),
			w.createAboutContainer(),
		),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	w.window.SetContent(container.NewBorder(nil,
		container.NewVBox( // Footer with build info for the whole window
			widget.NewSeparator(),
			container.NewBorder(nil, nil, nil, widget.NewRichTextFromMarkdown(
				fmt.Sprintf("[%s](%s) *v%s build %d*",
					fyne.CurrentApp().Metadata().Name,
					repositoryLink,
					fyne.CurrentApp().Metadata().Version,
					fyne.CurrentApp().Metadata().Build),
			)),
		),
		nil, nil,
		tabs, // Actual window content (tabs)
	))

	w.window.RequestFocus()
	w.window.Show()
}

func (w *SettingsDraft[T]) createAboutContainer() *fyne.Container {
	return container.NewCenter(container.NewVBox(widget.NewRichTextFromMarkdown(string(mdAboutContent))))
}

func (w *SettingsDraft[T]) createSettingsContainer() *fyne.Container {
	return container.NewBorder(nil, nil,
		w.createAddForm(), nil, // Add form on the left
		container.NewBorder( // All other space is occupied by connections list
			nil, nil,
			widget.NewSeparator(), nil,
			container.NewBorder(
				container.NewVBox(widget.NewLabel(configsListHeaderText), widget.NewSeparator()),
				nil, nil, nil, w.createDynamicList(),
			),
		),
	)
}

func (w *SettingsDraft[T]) createAddForm() *fyne.Container {
	inputLink := &widget.Entry{PlaceHolder: linkPlaceholderText}
	inputLabel := &widget.Entry{PlaceHolder: linkNamePlaceholderText}
	errLabel := &widget.Label{Importance: widget.DangerImportance}
	errLabel.Hide()

	addBtn := &widget.Button{
		Icon: theme.ContentAddIcon(),
		Text: addText,
		OnTapped: func() {
			data := FormData{Label: inputLabel.Text, Link: inputLink.Text}
			if handleAddItem(data, errLabel, w.onAdd) {
				inputLabel.Text = ""
				inputLink.Text = ""
				inputLabel.Refresh()
				inputLink.Refresh()
			}
		},
		Importance: widget.HighImportance,
	}

	return container.NewVBox(
		widget.NewLabel(insertYourConnectionURLText),
		inputLabel,
		inputLink,
		errLabel,
		container.NewBorder(nil, nil, nil, addBtn), // Fit button to the right side
	)
}

func (w *SettingsDraft[T]) createDynamicList() *fyne.Container {
	updateForm := form.NewUpdateConfig(updateText, deleteText)
	configInfoText := customwidget.NewTextWithCopy(w.window.Clipboard())

	netStatsChart := container.NewWithoutLayout(&fyne.Container{})
	itemSettings := container.NewBorder(
		widget.NewSeparator(),
		updateForm.Container(),
		nil, nil,
		container.NewBorder(nil, nil, netStatsChart, nil, configInfoText.Container()),
	)
	itemSettings.Hidden = true

	list := widget.NewListWithData(w.list,
		func() fyne.CanvasObject { return nil },
		func(_ binding.DataItem, _ fyne.CanvasObject) {},
	)
	list.HideSeparators = true

	// Cache active charts
	activeCharts := map[widget.ListItemID]*fyne.Container{}
	listContainer := container.NewBorder(nil, itemSettings, nil, nil, list)

	list.CreateItem = func() fyne.CanvasObject {
		dataStats := container.NewHBox(
			container.NewPadded(canvas.NewText("↑0.0 GB", theme.Color(customtheme.ColorNameTextMuted))),
			container.NewPadded(canvas.NewText("↓0.0 GB", theme.Color(customtheme.ColorNameTextMuted))),
		)

		cnt := container.NewBorder(nil, nil,
			container.NewPadded(widget.NewIcon(nil)),
			dataStats, container.NewHBox(widget.NewLabel("template"), container.NewHBox()),
		)
		return cnt
	}
	list.UpdateItem = func(id widget.ListItemID, o fyne.CanvasObject) {
		val := getListItem(w.list, id)
		updateForm.Disable(val.Active())
		activeIcon := o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Icon)
		if val.Active() {
			activeIcon.SetResource(icon.ListActive)
		} else {
			activeIcon.SetResource(nil)
		}

		stats := o.(*fyne.Container).Objects[2].(*fyne.Container)
		readBytes := fmt.Sprintf("↑%s", bytesToString(val.Recorder().BytesRead()))
		writtenBytes := fmt.Sprintf("↓%s", bytesToString(val.Recorder().BytesWritten()))
		stats.Objects[0].(*fyne.Container).Objects[0] = canvas.NewText(readBytes, theme.Color(customtheme.ColorNameTextMuted))
		stats.Objects[1].(*fyne.Container).Objects[0] = canvas.NewText(writtenBytes, theme.Color(customtheme.ColorNameTextMuted))
		o.Refresh()

		o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf(
			"%s [%s]", val.Label(), val.XRayConfig()["Address"],
		))

		if _, ok := activeCharts[id]; !ok {
			activeCharts[id] = getNetStatChartsDemo(context.Background(), fyne.NewSize(250, 100), val.Recorder())
		}

		if len(o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects) == 0 {
			o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container).Objects = []fyne.CanvasObject{
				customwidget.NewBadge(val.XRayConfig()["Protocol"], theme.Color(customtheme.ColorNameTextMuted)),
				customwidget.NewBadge(val.XRayConfig()["Type"], theme.Color(customtheme.ColorNameTextMuted)),
				customwidget.NewBadge(val.XRayConfig()["TLS"], theme.Color(customtheme.ColorNameTextMuted)),
			}
		}
	}
	list.OnSelected = func(id widget.ListItemID) {
		itemSettings.Show()
		defer itemSettings.Refresh()

		val := getListItem(w.list, id)
		updateForm.Disable(val.Active())
		updateForm.SetInputs(val.Label(), val.Link())

		netStatsChart.Objects[0] = activeCharts[id]
		netStatsChart.Refresh()
		configInfoText.ParseMarkdown(xrayConfigToMd(val.XRayConfig()))

		// wrap err handling
		handleErr := func(err error) {
			if err != nil {
				updateForm.SetError(err)

				return
			}
			list.UnselectAll()
			updateForm.SetError(nil)
			itemSettings.Refresh()
			itemSettings.Hide()
			listContainer.Refresh()
		}

		updateForm.OnUpdate(func() {
			data := FormData{Label: updateForm.InputLabel(), Link: updateForm.InputLink()}
			handleErr(handleUpdateItem(data, val.(T), w.onUpdate))
		})
		updateForm.OnDelete(func() {
			if val.Active() {
				handleErr(errChangeActiveItem)

				return
			}

			handleErr(w.onDelete(val.(T)))
		})
	}

	return listContainer
}

func bytesToString(bytes int) string {
	const bytesToMegabit = 125000
	const megaBitToGB = 8000
	postfix := "Mb"
	value := float64(bytes) / bytesToMegabit
	if value > 1000 { // threshold to turn into GB
		value = value / megaBitToGB
		postfix = "GB"
	}

	return fmt.Sprintf("%.2f %s", value, postfix)
}

func handleAddItem(data FormData, errLabel *widget.Label, onAdd func(FormData) error) bool {
	if err := data.Validate(); err != nil {
		errLabel.SetText(err.Error())
		errLabel.Show()

		return false
	}

	if err := onAdd(data); err != nil {
		errLabel.SetText(err.Error())
		errLabel.Show()

		return false
	}

	return true
}

func handleUpdateItem[T ListItem](data FormData, val T, onUpdate func(FormData, T) error) error {
	if val.Active() {
		return errChangeActiveItem
	}

	if err := data.Validate(); err != nil {
		return err
	}

	return onUpdate(data, val)
}

func getListItem(list binding.DataList, id widget.ListItemID) ListItem {
	i, err := list.GetItem(id)
	if err != nil {
		// Unexpected undefined behaviour, better to just panic
		panic(fmt.Errorf("error getting data item: %w", err))
	}
	untyped, _ := i.(binding.Untyped).Get()

	return untyped.(ListItem)
}

type disableWidget interface {
	Disable()
	Enable()
	Refresh()
}

func disableAll(disable bool, widgets ...disableWidget) {
	for _, wdg := range widgets {
		if disable {
			wdg.Disable()
		} else {
			wdg.Enable()
		}
		wdg.Refresh()
	}
}

func xrayConfigToMd(x map[string]string) string {
	includeOrder := []string{
		"Address",
		"Type", "TLS", "Protocol", "Port",
		"ID", "Remark", "TlsFingerprint", "SNI",
		"Security", "Aid", "Host", "Network", "Path", "ALPN", "Authority", "ServiceName", "Mode",
	}

	str := ""
	for _, k := range includeOrder {
		if x[k] == "" {
			continue
		}

		str += fmt.Sprintf("**%s**: %s\n\n", k, x[k])
	}

	return str
}
