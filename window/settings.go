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

	ctx       context.Context
	ctxCancel context.CancelFunc
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

	ctx, cancel := context.WithCancel(context.Background())
	return &SettingsDraft[T]{
		window:    w,
		onAdd:     onAdd,
		onUpdate:  onUpdate,
		onDelete:  onDelete,
		list:      list,
		ctx:       ctx,
		ctxCancel: cancel,
	}
}

func (w *SettingsDraft[T]) OnClosed(fn func()) {
	w.window.SetOnClosed(func() {
		w.ctxCancel()
		fn()
	})
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
	renderedBadges := map[widget.ListItemID][]fyne.CanvasObject{}
	updateRenderedBadges := func(id widget.ListItemID, val ListItem) {
		renderedBadges[id] = []fyne.CanvasObject{
			customwidget.NewBadge(val.XRayConfig()["Protocol"], theme.Color(customtheme.ColorNameTextMuted)),
			customwidget.NewBadge(val.XRayConfig()["Type"], theme.Color(customtheme.ColorNameTextMuted)),
			customwidget.NewBadge(val.XRayConfig()["TLS"], theme.Color(customtheme.ColorNameTextMuted)),
		}
	}

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
		defer itemSettings.Refresh()
		activeIcon := o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Icon)
		label := o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[0].(*widget.Label)
		badges := o.(*fyne.Container).Objects[0].(*fyne.Container).Objects[1].(*fyne.Container)
		netStatsRead := o.(*fyne.Container).Objects[2].(*fyne.Container).Objects[0].(*fyne.Container)
		netStatsWritten := o.(*fyne.Container).Objects[2].(*fyne.Container).Objects[1].(*fyne.Container)

		val := getListItem(w.list, id)

		if val.Active() {
			activeIcon.SetResource(icon.ListActive)
		} else {
			activeIcon.SetResource(nil)
		}

		readBytes := fmt.Sprintf("↑%s", bytesToString(val.Recorder().BytesRead()))
		writtenBytes := fmt.Sprintf("↓%s", bytesToString(val.Recorder().BytesWritten()))
		netStatsRead.Objects[0] = canvas.NewText(readBytes, theme.Color(customtheme.ColorNameTextMuted))
		netStatsWritten.Objects[0] = canvas.NewText(writtenBytes, theme.Color(customtheme.ColorNameTextMuted))

		label.SetText(fmt.Sprintf("%s [%s]", val.Label(), val.XRayConfig()["Address"]))

		if _, ok := activeCharts[id]; !ok {
			activeCharts[id] = customwidget.NewLiveNetworkChart(w.ctx, fyne.NewSize(250, 100), val.Recorder())
		}

		if _, ok := renderedBadges[id]; !ok {
			updateRenderedBadges(id, val)
		}

		badges.Objects = renderedBadges[id]
	}
	list.OnUnselected = func(id widget.ListItemID) {
		itemSettings.Hide()
	}
	list.OnSelected = func(id widget.ListItemID) {
		defer itemSettings.Show()
		defer itemSettings.Refresh()

		val := getListItem(w.list, id)

		netStatsChart.Objects[0] = activeCharts[id]
		netStatsChart.Refresh()
		configInfoText.ParseMarkdown(xrayConfigToMd(val.XRayConfig()))

		updateForm.Disable(val.Active())
		updateForm.SetInputs(val.Label(), val.Link())
		updateForm.OnUpdate(func() error {
			defer updateRenderedBadges(id, val) // update badges only on config change

			data := FormData{Label: updateForm.InputLabel(), Link: updateForm.InputLink()}

			if val.Active() {
				return errChangeActiveItem
			}

			if err := data.Validate(); err != nil {
				return err
			}

			return w.onUpdate(data, val.(T))
		})
		updateForm.OnDelete(func() error {
			if val.Active() {
				return errChangeActiveItem
			}

			return w.onDelete(val.(T))
		})
		updateForm.OnSubmit(func() {
			list.UnselectAll()
		})
	}

	return container.NewBorder(nil, itemSettings, nil, nil, list)
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

func getListItem(list binding.DataList, id widget.ListItemID) ListItem {
	i, err := list.GetItem(id)
	if err != nil {
		// Unexpected undefined behaviour, better to just panic
		panic(fmt.Errorf("error getting data item: %w", err))
	}
	untyped, _ := i.(binding.Untyped).Get()

	return untyped.(ListItem)
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
