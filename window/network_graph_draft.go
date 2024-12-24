package window

import (
	"image/color"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"github.com/goxray/ui/internal/netchart"
	customtheme "github.com/goxray/ui/theme"
)

type RecorderI interface {
	Read() []float64
	Written() []float64
	RecordInterval() time.Duration
}

// Recorder
// TODO: just to debug graph, remove and pass properly via item! maybe we should get it from the item?
var Recorder RecorderI

// TODO: move to netchart pkg
func getNetStatChartsDemo(size fyne.Size) *fyne.Container {
	const uplinkName = " ● upload"
	const downlinkName = "● download"
	data := map[string][]float64{
		uplinkName:   {},
		downlinkName: {},
	}

	chart := netchart.New(float64(size.Width), float64(size.Height), 0.6)

	r, g, b, _ := theme.Color(customtheme.ColorNameGraphGreen).RGBA()
	colorGreen := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}
	r, g, b, _ = theme.Color(customtheme.ColorNameGraphBlue).RGBA()
	colorBlue := color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: 255}

	colors := map[string]color.RGBA{
		uplinkName:   colorGreen,
		downlinkName: colorBlue,
	}

	data[uplinkName] = Recorder.Read()
	data[downlinkName] = Recorder.Written()

	if err := chart.UpdateNamed(data, colors, []string{uplinkName, downlinkName}); err != nil {
		log.Fatal(err)
	}

	// go func() {
	// 	for {
	// 		data[uplinkName] = recorder.Read()
	// 		data[downlinkName] = recorder.Written()
	//
	// 		if err := chart.UpdateNamed(data, colors, []string{uplinkName, downlinkName}); err != nil {
	// 			log.Fatal(err)
	// 		}
	// 		<-time.After(recorder.RecordInterval())
	// 	}
	// }()

	return chart.Container()
}
