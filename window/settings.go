package window

import (
	"context"
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/goxray/ui/icon"
	customtheme "github.com/goxray/ui/theme"
	"github.com/goxray/ui/window/form"
	customwidget "github.com/goxray/ui/window/widget"
)

type Settings[T ListItem] struct {
	window fyne.Window
	list   binding.DataList

	onAdd    func(data FormData) error
	onUpdate func(FormData, T) error
	onDelete func(T) error

	ctx       context.Context
	ctxCancel context.CancelFunc
}

func NewSettings[T ListItem](
	a fyne.App,
	list binding.DataList,
	onAdd func(data FormData) error,
	onUpdate func(FormData, T) error,
	onDelete func(T) error,
) *Settings[T] {
	w := a.NewWindow(lang.L("Settings"))
	w.CenterOnScreen()
	w.RequestFocus()
	w.Resize(fyne.NewSize(750, 450))

	ctx, cancel := context.WithCancel(context.Background())

	s := &Settings[T]{
		window:    w,
		onAdd:     onAdd,
		onUpdate:  onUpdate,
		onDelete:  onDelete,
		list:      list,
		ctx:       ctx,
		ctxCancel: cancel,
	}
	s.init()

	return s
}

func (w *Settings[T]) Refresh() {
	if w.window.Canvas() != nil {
		w.window.Content().Refresh()
	}
}

func (w *Settings[T]) OnClosed(fn func()) {
	w.window.SetOnClosed(func() {
		w.ctxCancel()
		fn()
	})
}

func (w *Settings[T]) init() {
	content := container.NewBorder(nil,
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
	)
	w.window.SetContent(content)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon( // Connections list settings tab
			lang.L("Configs"),
			icon.Settings,
			w.createSettingsContainer(),
		),
		container.NewTabItemWithIcon( // About tab with static app info
			lang.L("About"),
			theme.QuestionIcon(),
			w.createAboutContainer(),
		),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	content.Objects = append(content.Objects, tabs)

	w.Refresh()
}

func (w *Settings[T]) Show() {
	w.Refresh()
	w.window.RequestFocus()
	w.window.Show()
}

func (w *Settings[T]) createAboutContainer() *fyne.Container {
	return container.NewCenter(container.NewVBox(widget.NewRichTextFromMarkdown(string(mdAboutContent))))
}

func (w *Settings[T]) createSettingsContainer() *fyne.Container {
	return container.NewBorder(nil, nil,
		w.createAddForm(), nil, // Add form on the left
		container.NewBorder( // All other space is occupied by connections list
			nil, nil,
			widget.NewSeparator(), nil,
			container.NewBorder(
				container.NewVBox(widget.NewLabel(lang.L("Available connection configurations:")), widget.NewSeparator()),
				nil, nil, nil, w.createDynamicList(),
			),
		),
	)
}

func (w *Settings[T]) createAddForm() *fyne.Container {
	inputLink := &widget.Entry{PlaceHolder: "vless://example.com..."}
	inputLabel := &widget.Entry{PlaceHolder: lang.L("Display name")}
	errLabel := &widget.Label{Importance: widget.DangerImportance}
	errLabel.Hide()

	addBtn := &widget.Button{
		Icon: theme.ContentAddIcon(),
		Text: lang.L("Add"),
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
		widget.NewLabel(lang.L("Insert your connection URL")),
		inputLabel,
		inputLink,
		errLabel,
		container.NewBorder(nil, nil, nil, addBtn), // Fit button to the right side
	)
}

func (w *Settings[T]) createDynamicList() *fyne.Container {
	updateForm := form.NewUpdateConfig(lang.L("Update"), lang.L("Delete"))
	configInfoText := customwidget.NewTextWithCopy(w.window.Clipboard())

	netStatsChart := container.NewWithoutLayout(&fyne.Container{})
	itemSettings := container.NewBorder(
		widget.NewSeparator(),
		updateForm.Container(),
		nil, nil,
		container.NewBorder(nil, nil, netStatsChart, nil, configInfoText.Container()),
	)
	itemSettings.Hidden = true

	list := widget.NewListWithData(w.list, nil, nil)
	list.HideSeparators = true

	// Small caches to reuse sensitive widgets.
	activeCharts := map[widget.ListItemID]*fyne.Container{}       // Cache for active live charts
	renderedBadges := map[widget.ListItemID][]fyne.CanvasObject{} // Cache for badges
	activeNetStats := map[widget.ListItemID]*fyne.Container{}     // Cache for net stats counters
	var selectedItem widget.ListItemID = -1

	list.CreateItem = func() fyne.CanvasObject {
		dataStats := container.NewHBox(
			container.NewPadded(canvas.NewText("↑0.0 "+lang.L("GB"), theme.Color(customtheme.ColorNameTextMuted))),
			container.NewPadded(canvas.NewText("↓0.0 "+lang.L("GB"), theme.Color(customtheme.ColorNameTextMuted))),
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

		val := getListItem(w.list, id)

		if val.Active() {
			activeIcon.SetResource(icon.ListActive)
		} else {
			activeIcon.SetResource(nil)
		}
		if id == selectedItem {
			updateForm.ToggleHide(val.Active())
		}

		label.SetText(fmt.Sprintf("%s [%s]", val.Label(), val.XRayConfig()["Address"]))

		if _, ok := activeNetStats[id]; !ok {
			activeNetStats[id] = customwidget.NewLiveNetworkStats(w.ctx, val)
			o.(*fyne.Container).Objects[2].(*fyne.Container).Objects = activeNetStats[id].Objects
		}

		if _, ok := activeCharts[id]; !ok {
			activeCharts[id] = customwidget.NewLiveNetworkChart(w.ctx, " ● "+lang.L("upload"), "● "+lang.L("download"),
				fyne.NewSize(250, 100), val)
		}

		if _, ok := renderedBadges[id]; !ok {
			renderedBadges[id] = createBadgesForVal(val)
		}

		badges.Objects = renderedBadges[id]
	}
	list.OnUnselected = func(id widget.ListItemID) {
		itemSettings.Hide()
	}
	list.OnSelected = func(id widget.ListItemID) {
		selectedItem = id
		defer itemSettings.Show()
		defer itemSettings.Refresh()

		val := getListItem(w.list, id)

		netStatsChart.Objects[0] = activeCharts[id]
		configInfoText.ParseMarkdown(xrayConfigToStrings(val.XRayConfig()))

		updateForm.ToggleHide(val.Active())
		updateForm.SetInputs(val.Label(), val.Link())
		updateForm.OnUpdate(func() error {
			// Update badges to reflect config changes in update.
			defer func() { renderedBadges[id] = createBadgesForVal(val) }()

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

// createBadgesForVal generates badges set for list value.
func createBadgesForVal(val ListItem) []fyne.CanvasObject {
	showTagsFor := []string{"Protocol", "TLS", "Flow"}
	// Specify specific key:values that should be marked with different badge color.
	specialColors := map[string]map[string]color.Color{
		// TLS none is a terrible security issue, mark it red.
		"TLS": {"none": theme.Color(customtheme.ColorNameTextErrorMuted)},
	}

	badges := make([]fyne.CanvasObject, 0, len(showTagsFor))
	for _, tag := range showTagsFor {
		value, ok := val.XRayConfig()[tag]
		if !ok || value == "" {
			continue
		}
		clr := theme.Color(customtheme.ColorNameTextMuted)
		if colorsVal, ok := specialColors[tag]; ok && colorsVal[value] != nil {
			clr = colorsVal[value]
		}

		badges = append(badges, customwidget.NewBadge(val.XRayConfig()[tag], clr))
	}

	return badges
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

func xrayConfigToStrings(x map[string]string) (md string, toCopy string) {
	const separator = "separator"
	includeOrder := []string{
		"Address",
		"Type", "TLS", "Protocol", "Port",
		"ID", "Remark", "TlsFingerprint", "SNI",
		"Security", "Aid", "Host", "Network", "Path", "ALPN", "Authority", "ServiceName", "Mode",
		separator, // separate base info from protocol-specific info
		"Flow", "Pbk", "Sid", "Spx", "Fp",
	}

	for _, k := range includeOrder {
		if k == separator {
			md += fmt.Sprintf("---\n") // draw horizontal line in MD.
			continue
		}
		if x[k] == "" {
			continue
		}

		md += fmt.Sprintf("**%s**: %s\n\n", lang.L(k), x[k])
		toCopy += fmt.Sprintf("%s: %s\n", lang.L(k), x[k])
	}

	return md, toCopy
}
