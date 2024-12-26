package widget

import (
	"context"
	"image/color"
	"log/slog"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/goxray/ui/internal/netchart"
	customtheme "github.com/goxray/ui/theme"
)

type Recorder interface {
	Read() []float64
	Written() []float64
	BytesRead() int
	BytesWritten() int
	RecordInterval() time.Duration
}

// NewLiveNetworkChart creates new chart container with automatically updated net statistics from recorder.
//
// This method spawns a new goroutine to update the chart in background, so be sure to close the ctx when you are done.
func NewLiveNetworkChart(ctx context.Context, size fyne.Size, recorder Recorder) *fyne.Container {
	const uplinkName, downlinkName = " ● upload", "● download"
	data := map[string][]float64{uplinkName: {}, downlinkName: {}}

	r, g, b, _ := theme.Color(customtheme.ColorNameGraphGreen).RGBA()
	colorGreen := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	r, g, b, _ = theme.Color(customtheme.ColorNameGraphBlue).RGBA()
	colorBlue := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	colors := map[string]color.RGBA{uplinkName: colorGreen, downlinkName: colorBlue}

	chart := netchart.New(float64(size.Width), float64(size.Height), 0.6)

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			data[uplinkName] = recorder.Read()
			data[downlinkName] = recorder.Written()

			if err := chart.UpdateNamed(data, colors, []string{uplinkName, downlinkName}); err != nil {
				slog.Error(err.Error())
				cancel()
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(recorder.RecordInterval()):
				continue
			}
		}
	}()

	return chart.Container()
}
