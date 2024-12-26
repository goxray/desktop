package form

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UpdateConfig struct {
	errLabel           *widget.Label
	onSubmit           func()
	saveBtn, deleteBtn *widget.Button
	newLabel, newLink  *widget.Entry
	container          *fyne.Container
}

func NewUpdateConfig(updateBtnTitle, deleteBtnTitle string) *UpdateConfig {
	errLabel := &widget.Label{Text: "error", Importance: widget.DangerImportance}
	errLabel.Hide()
	newLabelInput := widget.NewEntry()
	newLinkInput := widget.NewEntry()

	saveBtn := &widget.Button{Text: updateBtnTitle, Icon: theme.DocumentCreateIcon(), Importance: widget.HighImportance}
	deleteBtn := &widget.Button{Text: deleteBtnTitle, Icon: theme.DeleteIcon(), Importance: widget.DangerImportance}

	return &UpdateConfig{
		errLabel:  errLabel,
		saveBtn:   saveBtn,
		deleteBtn: deleteBtn,
		newLabel:  newLabelInput,
		newLink:   newLinkInput,
		onSubmit:  func() {},
		container: container.NewVBox(
			widget.NewSeparator(),
			container.NewVBox(errLabel, newLabelInput, newLinkInput),
			container.NewBorder(nil, nil, nil, deleteBtn, saveBtn),
		),
	}
}

func (f *UpdateConfig) Container() *fyne.Container {
	return f.container
}

func (f *UpdateConfig) Disable(disable bool) {
	if disable {
		f.saveBtn.Disable()
		f.deleteBtn.Disable()
		f.newLabel.Disable()
		f.newLink.Disable()
	} else {
		f.saveBtn.Enable()
		f.deleteBtn.Enable()
		f.newLabel.Enable()
		f.newLink.Enable()
	}
}

func (f *UpdateConfig) SetInputs(label, link string) {
	f.newLink.SetText(link)
	f.newLabel.SetText(label)
}

func (f *UpdateConfig) OnSubmit(fn func()) {
	f.onSubmit = fn
}

func (f *UpdateConfig) SetError(err error) {
	if err != nil {
		f.errLabel.SetText(err.Error())
		f.errLabel.Show()
	} else {
		f.errLabel.SetText("")
		f.errLabel.Hide()
		f.onSubmit()
	}
}

func (f *UpdateConfig) InputLabel() string {
	return f.newLabel.Text
}

func (f *UpdateConfig) InputLink() string {
	return f.newLink.Text
}

func (f *UpdateConfig) OnUpdate(fn func() error) {
	f.saveBtn.OnTapped = func() {
		f.SetError(fn())
	}
}

func (f *UpdateConfig) OnDelete(fn func() error) {
	f.deleteBtn.OnTapped = func() {
		f.SetError(fn())
	}
}
