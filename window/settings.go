package window

import (
	_ "embed"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/goxray/ui/icon"
)

const repositoryLink = "https://github.com/goxray"

// Texts for future translations.
const (
	SettingsText                = "Settings"
	ConfigsText                 = "Configs"
	AboutText                   = "About"
	LinkPlaceholderText         = "vless://example.com..."
	LinkNamePlaceholderText     = "Display name"
	AddText                     = "Add"
	UpdateText                  = "Update"
	DeleteText                  = "Delete"
	InsertYourConnectionURLText = "Insert your connection URL"

	ErrChangeActiveText = "disconnect before editing"
	ErrLabelOrLinkEmpty = "label or link empty"
)

var (
	ErrChangeActiveItem     = errors.New(ErrChangeActiveText)
	ErrEmptyUpdateFormValue = errors.New(ErrLabelOrLinkEmpty)
)

type FormData struct {
	Label string
	Link  string
}

func (f *FormData) Validate() error {
	if f.Label == "" || f.Link == "" {
		return ErrEmptyUpdateFormValue
	}

	return nil
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
	w := a.NewWindow(SettingsText)
	w.CenterOnScreen()
	w.RequestFocus()
	w.Resize(fyne.NewSize(700, 580))

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon( // Connections list settings tab
			ConfigsText,
			icon.Settings,
			createSettingsContainer(list, onAdd, onUpdate, onDelete),
		),
		container.NewTabItemWithIcon( // About tab with static app info
			AboutText,
			theme.QuestionIcon(),
			createAboutContainer(),
		),
	)
	tabs.SetTabLocation(container.TabLocationTop)

	w.SetContent(container.NewBorder(nil,
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

	return w
}

func createAboutContainer() *fyne.Container {
	return container.NewCenter(container.NewVBox(widget.NewRichTextFromMarkdown(string(mdAboutContent))))
}

func createSettingsContainer[T ListItem](
	list binding.DataList,
	onAdd func(data FormData) error,
	onUpdate func(FormData, T) error,
	onDelete func(T) error,
) *fyne.Container {
	return container.NewBorder(nil, nil,
		createAddForm(onAdd), nil, // Add form on the left
		container.NewBorder( // All other space is occupied by connections list
			nil, nil,
			widget.NewSeparator(), nil,
			container.NewBorder(
				container.NewVBox(widget.NewLabel("Available connection links:"), widget.NewSeparator()),
				nil, nil, nil, createDynamicList(list, onDelete, onUpdate),
			),
		),
	)
}

func createAddForm(onAdd func(data FormData) error) *fyne.Container {
	inputLink := &widget.Entry{PlaceHolder: LinkPlaceholderText}
	inputLabel := &widget.Entry{PlaceHolder: LinkNamePlaceholderText}
	errLabel := &widget.Label{Importance: widget.DangerImportance}
	errLabel.Hide()

	handleError := func(err error) {
		if err != nil {
			errLabel.SetText(err.Error())
			errLabel.Show()

			return
		}
		errLabel.Hide()
	}

	addBtn := &widget.Button{
		Icon: theme.ContentAddIcon(),
		Text: AddText,
		OnTapped: func() {
			data := FormData{Label: inputLabel.Text, Link: inputLink.Text}
			handleError(handleAddItem(data, onAdd))
		},
		Importance: widget.HighImportance,
	}

	return container.NewVBox(
		widget.NewLabel(InsertYourConnectionURLText),
		inputLabel,
		inputLink,
		errLabel,
		container.NewBorder(nil, nil, nil, addBtn), // Fit button to the right side
	)
}

func createDynamicList[T ListItem](
	connectionsList binding.DataList,
	onDelete func(itm T) error,
	onUpdate func(FormData, T) error,
) *fyne.Container {
	saveBtn := &widget.Button{
		Text:       UpdateText,
		Icon:       theme.DocumentSaveIcon(),
		Importance: widget.HighImportance,
	}
	deleteBtn := &widget.Button{
		Text:       DeleteText,
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
			data := FormData{Label: newLabelInput.Text, Link: newLinkInput.Text}
			handleErr(handleUpdateItem(data, val.(T), onUpdate))
		}
		deleteBtn.OnTapped = func() {
			if val.Active() {
				handleErr(ErrChangeActiveItem)

				return
			}

			handleErr(onDelete(val.(T)))
		}
	}

	return listContainer
}

func handleAddItem(data FormData, onAdd func(FormData) error) error {
	if err := data.Validate(); err != nil {
		return err
	}

	return onAdd(data)
}

func handleUpdateItem[T ListItem](data FormData, val T, onUpdate func(FormData, T) error) error {
	if val.Active() {
		return ErrChangeActiveItem
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
}

func disableAll(disable bool, buttons ...disableWidget) {
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
