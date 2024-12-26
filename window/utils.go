package window

import (
	"fmt"
)

// bytesToString returns short string representation of bytes, starting from Mb and ending in GB.
func bytesToString(bytes int) string {
	const bytesToMegabit = 125000
	const megaBitToGB = 8000
	postfix := "Mb"
	value := float64(bytes) / bytesToMegabit
	if value > 1000 { // threshold to turn into GB
		value = value / megaBitToGB
		postfix = "GB"
	}

	return fmt.Sprintf("%.2f %s", value, postfix)
}
