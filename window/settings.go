package window

import (
	_ "embed"
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/goxray/ui/icon"
)

type FormData struct {
	Label string
	Link  string
}

type ListItem interface {
	Label() string
	Link() string
	XRayConfig() map[string]string
	Active() bool
}

//go:embed about_static.md
var mdAboutContent []byte

func NewSettings[T ListItem](
	a fyne.App,
	onAdd func(data FormData) error,
	onUpdate func(FormData, T) error,
	onDelete func(T) error,
	list binding.DataList,
) fyne.Window {
	w := a.NewWindow("Settings")
	w.CenterOnScreen()
	w.RequestFocus()

	addForm := createAddForm(onAdd)
	itemsList := container.NewBorder(
		container.NewVBox(widget.NewLabel("Available connection links:"), widget.NewSeparator()),
		nil, nil, nil, createDynamicList(list, onDelete, onUpdate),
	)

	configsTab := container.NewBorder(nil, nil, addForm, nil, container.NewBorder(
		nil, nil,
		widget.NewSeparator(), nil,
		itemsList,
	),
	)
	aboutTab := container.NewCenter(container.NewVBox(widget.NewRichTextFromMarkdown(string(mdAboutContent))))

	metadata := fyne.CurrentApp().Metadata()
	repoLink, _ := url.Parse("https://github.com/goxray")
	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Configs", icon.Settings, configsTab),
		container.NewTabItemWithIcon("About", theme.QuestionIcon(), aboutTab),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(container.NewBorder(nil,
		container.NewVBox(
			widget.NewSeparator(),
			container.NewBorder(nil, nil, nil, widget.NewRichTextFromMarkdown(
				fmt.Sprintf("[%s](%s) *v%s build %d*",
					metadata.Name, repoLink.String(),
					metadata.Version, metadata.Build),
			)),
		),
		nil, nil, tabs,
	))
	w.Resize(fyne.NewSize(700, 580))

	return w
}

func createAddForm(onAdd func(data FormData) error) *fyne.Container {
	inputLink := &widget.Entry{PlaceHolder: "vless://example.com..."}
	inputLabel := &widget.Entry{PlaceHolder: "Display name"}
	errLabel := &widget.Label{Importance: widget.DangerImportance}
	errLabel.Hide()

	addBtn := &widget.Button{
		Icon: theme.ContentAddIcon(),
		Text: "Add",
		OnTapped: func() {
			if err := onAdd(FormData{Label: inputLabel.Text, Link: inputLink.Text}); err != nil {
				errLabel.SetText(err.Error())
				errLabel.Show()
				return
			}
			errLabel.Hide()
		},
		Importance: widget.HighImportance,
	}

	return container.NewBorder(nil, nil, nil, container.NewVBox(
		widget.NewLabel("Insert your connection URL"),
		inputLabel,
		inputLink,
		errLabel,
		container.NewBorder(nil, nil, nil, addBtn),
	))
}

func createDynamicList[T ListItem](
	connectionsList binding.DataList,
	onDelete func(itm T) error,
	onUpdate func(FormData, T) error,
) *fyne.Container {
	saveBtn := &widget.Button{
		Text:       "Update",
		Icon:       theme.DocumentSaveIcon(),
		Importance: widget.HighImportance,
	}
	deleteBtn := &widget.Button{
		Text:       "Delete",
		Icon:       theme.DeleteIcon(),
		Importance: widget.DangerImportance,
	}

	configInfo := widget.NewRichTextFromMarkdown("configuration info")
	errLabel := &widget.Label{Text: "error", Importance: widget.DangerImportance}
	errLabel.Hide()
	newLabelInput := widget.NewEntry()
	newLinkInput := widget.NewEntry()
	itemSettings := container.NewBorder(
		widget.NewSeparator(),
		container.NewVBox(
			widget.NewSeparator(),
			container.NewVBox(errLabel, newLabelInput, newLinkInput),
			container.NewBorder(nil, nil, nil, deleteBtn, saveBtn),
		), nil, nil,

		configInfo,
	)
	itemSettings.Hidden = true

	list := widget.NewListWithData(connectionsList,
		func() fyne.CanvasObject { return nil },
		func(_ binding.DataItem, _ fyne.CanvasObject) {},
	)
	list.HideSeparators = true
	list.CreateItem = func() fyne.CanvasObject {
		cnt := container.NewBorder(nil, nil,
			container.NewPadded(widget.NewIcon(nil)),
			widget.NewRichTextFromMarkdown(""), widget.NewLabel("template"),
		)
		return cnt
	}

	listContainer := container.NewBorder(nil, itemSettings, nil, nil, list)
	list.UpdateItem = func(id widget.ListItemID, o fyne.CanvasObject) {
		val := getListItem(connectionsList, id)
		disableAll(val.Active(), saveBtn, deleteBtn)
		activeIcon := o.(*fyne.Container).Objects[1].(*fyne.Container).Objects[0].(*widget.Icon)
		if val.Active() {
			activeIcon.SetResource(icon.ListActive)
		} else {
			activeIcon.SetResource(nil)
		}
		o.Refresh()

		o.(*fyne.Container).Objects[0].(*widget.Label).SetText(fmt.Sprintf(
			"%s [%s]", val.Label(), val.XRayConfig()["Address"],
		))
	}
	list.OnSelected = func(id widget.ListItemID) {
		itemSettings.Show()
		defer itemSettings.Refresh()

		val := getListItem(connectionsList, id)
		disableAll(val.Active(), saveBtn, deleteBtn)

		configInfo.ParseMarkdown(xrayConfigToMd(val.XRayConfig()))
		newLabelInput.SetText(val.Label())
		newLinkInput.SetText(val.Link())

		// wrap err handling
		handleErr := func(err error) {
			if err != nil {
				errLabel.SetText(err.Error())
				errLabel.Show()

				return
			}
			list.UnselectAll()
			errLabel.Hide()
			itemSettings.Refresh()
			itemSettings.Hide()
			listContainer.Refresh()
		}

		saveBtn.OnTapped = func() {
			if val.Active() {
				handleErr(fmt.Errorf("disconnect before editing"))

				return
			}

			handleErr(onUpdate(FormData{Label: newLabelInput.Text, Link: newLinkInput.Text}, val.(T)))
		}
		deleteBtn.OnTapped = func() {
			if val.Active() {
				handleErr(fmt.Errorf("disconnect before editing"))

				return
			}

			handleErr(onDelete(val.(T)))
		}
	}

	return listContainer
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

func disableAll(disable bool, buttons ...*widget.Button) {
	for _, btn := range buttons {
		if disable {
			btn.Disable()
		} else {
			btn.Enable()
		}
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
