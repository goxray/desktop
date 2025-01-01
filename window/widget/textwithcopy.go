package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	customtheme "github.com/goxray/ui/theme"
)

type hoverable struct {
	fyne.CanvasObject
	onMouseIn  func()
	onMouseOut func()
}

func (h *hoverable) MouseIn(*desktop.MouseEvent) {
	h.onMouseIn()
}

func (h *hoverable) MouseOut() {
	h.onMouseOut()
}

func (h *hoverable) MouseMoved(*desktop.MouseEvent) {}

// TextWithCopy represents a widget.RichText with copy button attached to the top right corner.
type TextWithCopy struct {
	content   *widget.RichText
	copyBtn   *widget.Button
	container *fyne.Container
	clipboard fyne.Clipboard
}

func NewTextWithCopy(clipboard fyne.Clipboard) *TextWithCopy {
	richText := widget.NewRichTextFromMarkdown("configuration info")

	copyBtn := widget.NewButtonWithIcon(
		"", theme.NewColoredResource(theme.ContentCopyIcon(), customtheme.ColorNameTextMuted), func() {
			clipboard.SetContent(richText.String())
		},
	)
	copyBtn.Hidden = true

	hv := &hoverable{
		CanvasObject: container.NewStack(),
		onMouseIn:    copyBtn.Show,
		onMouseOut:   copyBtn.Hide,
	}

	// Push copy btn to the right top corner
	cnt := container.NewStack(container.NewBorder(container.NewBorder(nil, nil, nil,
		container.NewPadded(container.NewPadded(
			copyBtn,
		)),
	), nil, nil, nil), container.NewScroll(richText), hv)

	return &TextWithCopy{
		content:   richText,
		container: cnt,
		copyBtn:   copyBtn,
		clipboard: clipboard,
	}
}

func (t *TextWithCopy) Container() *fyne.Container {
	return t.container
}

// ParseMarkdown updates the TextWithCopy RichText content and sets the text that will be copied to clipboard on copy button press.
func (t *TextWithCopy) ParseMarkdown(markdown string, toBeCopied string) {
	t.copyBtn.OnTapped = func() {
		if toBeCopied == "" {
			toBeCopied = markdown
		}

		t.clipboard.SetContent(toBeCopied)
	}

	t.content.ParseMarkdown(markdown)
}
