package window

import (
	_ "embed"
	"errors"
	"time"
)

//go:embed about_static.md
var mdAboutContent []byte

// Texts for future translations.
const (
	repositoryLink = "https://github.com/goxray"

	settingsText                = "Settings"
	configsText                 = "Configs"
	aboutText                   = "About"
	linkPlaceholderText         = "vless://example.com..."
	linkNamePlaceholderText     = "Display name"
	addText                     = "Add"
	updateText                  = "Update"
	deleteText                  = "Delete"
	insertYourConnectionURLText = "Insert your connection URL"
	configsListHeaderText       = "Available connection configurations:"

	errChangeActiveText = "disconnect before editing"
	errLabelOrLinkEmpty = "label or link empty"
)

var (
	errChangeActiveItem     = errors.New(errChangeActiveText)
	errEmptyUpdateFormValue = errors.New(errLabelOrLinkEmpty)
)

type FormData struct {
	Label string
	Link  string
}

func (f *FormData) Validate() error {
	if f.Label == "" || f.Link == "" {
		return errEmptyUpdateFormValue
	}

	return nil
}

type NetworkRecorder interface {
	Read() []float64
	Written() []float64
	BytesRead() int
	BytesWritten() int
	RecordInterval() time.Duration
}

type ListItem interface {
	Label() string
	Link() string
	XRayConfig() map[string]string
	Active() bool

	Recorder() NetworkRecorder
}
