package widget

import (
	"context"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"

	customtheme "github.com/goxray/ui/theme"
)

type NetStats interface {
	// BytesWritten should return total number of bytes written to the connection.
	BytesWritten() int
	// BytesRead should return total number of bytes read from the connection.
	BytesRead() int
	RecordInterval() time.Duration
}

func NewLiveNetworkStats(ctx context.Context, source NetStats) *fyne.Container {
	cnt := container.NewHBox(
		container.NewPadded(canvas.NewText("↑... "+lang.L("GB"), theme.Color(customtheme.ColorNameTextMuted))),
		container.NewPadded(canvas.NewText("↓... "+lang.L("GB"), theme.Color(customtheme.ColorNameTextMuted))),
	)

	go func() {
		prevReadBytes, prevWrittenBytes := "", ""

		// Small optimization to not rerender widget if no changes occurred to the values.
		shouldUpdateUI := func(readBytes, writtenBytes string) bool {
			noChange := readBytes == prevReadBytes && writtenBytes == prevWrittenBytes
			if noChange {
				return false
			}

			prevReadBytes, prevWrittenBytes = readBytes, writtenBytes

			return true
		}

		for {
			readBytes := fmt.Sprintf("↑%s", bytesToHumanFriendlyString(source.BytesRead()))
			writtenBytes := fmt.Sprintf("↓%s", bytesToHumanFriendlyString(source.BytesWritten()))
			if shouldUpdateUI(readBytes, writtenBytes) {
				readText := cnt.Objects[0].(*fyne.Container)
				writeText := cnt.Objects[1].(*fyne.Container)

				readText.Objects[0] = canvas.NewText(readBytes, theme.Color(customtheme.ColorNameTextMuted))
				writeText.Objects[0] = canvas.NewText(writtenBytes, theme.Color(customtheme.ColorNameTextMuted))
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(source.RecordInterval()):
			}
		}
	}()

	return cnt
}

// bytesToHumanFriendlyString returns short string representation of bytes, starting from Mb and ending in GB.
func bytesToHumanFriendlyString(bytes int) string {
	const bytesToMegabit = 125000
	const megaBitToGB = 8000
	postfix := lang.L("Mb")
	value := float64(bytes) / bytesToMegabit
	if value > 1000 { // threshold to turn into GB
		value = value / megaBitToGB
		postfix = lang.L("GB")
	}

	return fmt.Sprintf("%.2f %s", value, postfix)
}
