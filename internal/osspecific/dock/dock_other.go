//go:build !linux && !darwin

package dock

import "log/slog"

func HideIconInDock() {
	slog.Warn("hiding dock icon not implemented on this platform")

	return
}
