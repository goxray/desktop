package widget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	customtheme "github.com/goxray/ui/theme"
)

// TextWithCopy represents a widget.RichText with copy button attached to the top right corner.
type TextWithCopy struct {
	content   *widget.RichText
	container *fyne.Container
	clipboard fyne.Clipboard
}

func NewTextWithCopy(clipboard fyne.Clipboard) *TextWithCopy {
	richText := widget.NewRichTextFromMarkdown("configuration info")

	// Push copy btn to the right top corner
	copyConfigBtn := container.NewBorder(container.NewBorder(nil, nil, nil, container.NewPadded(container.NewPadded(
		widget.NewButtonWithIcon(
			"", theme.NewColoredResource(theme.ContentCopyIcon(), customtheme.ColorNameTextMuted), func() {
				clipboard.SetContent(richText.String())
			},
		),
	))), nil, nil, nil)

	return &TextWithCopy{
		content:   richText,
		container: container.NewStack(copyConfigBtn, container.NewVScroll(richText)),
		clipboard: clipboard,
	}
}

func (t *TextWithCopy) Container() *fyne.Container {
	return t.container
}

// ParseMarkdown updates the TextWithCopy RichText content.
func (t *TextWithCopy) ParseMarkdown(text string) {
	t.content.ParseMarkdown(text)
}
