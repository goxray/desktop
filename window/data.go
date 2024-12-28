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
)

var (
	errChangeActiveItem     = errors.New("disconnect before editing")
	errEmptyUpdateFormValue = errors.New("label or link empty")
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
	// Read should return values for uplink for each previous RecordInterval.
	// Number of values returned must match Written.
	Read() []float64
	// Written should return values for downlink for each previous RecordInterval.
	// Number of values returned must match Written.
	Written() []float64
	// BytesRead should return the total number of bytes for uplink.
	BytesRead() int
	// BytesWritten should return the total number of bytes for downlink.
	BytesWritten() int
	RecordInterval() time.Duration
}

type ListItem interface {
	NetworkRecorder

	Label() string
	Link() string
	XRayConfig() map[string]string
	Active() bool
}
