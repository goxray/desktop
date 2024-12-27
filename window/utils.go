package window

import (
	"fmt"

	"fyne.io/fyne/v2/lang"
)

// bytesToString returns short string representation of bytes, starting from Mb and ending in GB.
func bytesToString(bytes int) string {
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
