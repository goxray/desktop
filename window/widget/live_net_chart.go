package widget

import (
	"context"
	"image/color"
	"log/slog"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/goxray/desktop/internal/netchart"
	customtheme "github.com/goxray/desktop/theme"
)

type Recorder interface {
	// Read should return values for uplink for each previous RecordInterval.
	// Number of values returned must match Written.
	Read() []float64
	// Written should return values for downlink for each previous RecordInterval.
	// Number of values returned must match Written.
	Written() []float64
	RecordInterval() time.Duration
}

// NewLiveNetworkChart creates new chart container with automatically updated net statistics from recorder.
//
// This method spawns a new goroutine to update the chart in background, so be sure to close the ctx when you are done.
func NewLiveNetworkChart(ctx context.Context, upLabel, downLabel string, size fyne.Size, recorder Recorder) *fyne.Container {
	data := map[string][]float64{upLabel: {}, downLabel: {}}
	r, g, b, _ := theme.Color(customtheme.ColorNameGraphGreen).RGBA()
	colorGreen := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	r, g, b, _ = theme.Color(customtheme.ColorNameGraphBlue).RGBA()
	colorBlue := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	colors := map[string]color.RGBA{upLabel: colorGreen, downLabel: colorBlue}

	chart := netchart.New(float64(size.Width), float64(size.Height), 0.6)

	ctx, cancel := context.WithCancel(ctx)
	updateChart := func() {
		if err := chart.UpdateNamed(data, colors, []string{upLabel, downLabel}); err != nil {
			slog.Error(err.Error())
			cancel()
		}
	}

	// Initialize the chart with initial data.
	data[upLabel] = recorder.Read()
	data[downLabel] = recorder.Written()
	updateChart()

	go func() {
		emptyDrawn := false
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(recorder.RecordInterval()):
			}

			data[upLabel] = recorder.Read()
			data[downLabel] = recorder.Written()

			if allMapZeroes(data) { // Optimization: no need to rerender stale zeroed chart.
				if !emptyDrawn { // draw empty chart just once
					updateChart()
					emptyDrawn = true
				}

				continue
			}

			updateChart()
			emptyDrawn = false
		}
	}()

	return chart.Container()
}

func allMapZeroes(data map[string][]float64) bool {
	for _, v := range data {
		if !allZeroes(v) {
			return false
		}
	}

	return true
}

func allZeroes(vals []float64) bool {
	for _, v := range vals {
		if v != 0 {
			return false
		}
	}

	return true
}
