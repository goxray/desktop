//go:build !linux && !darwin

package root

import "log/slog"

func PromptRootAccess() {
	slog.Warn("PromptRootAccess not implemented on this platform, run the program as root manually")
}

func Init() {}

func OnStart() {}
