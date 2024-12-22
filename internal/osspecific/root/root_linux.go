package root

import "log/slog"

func PromptRootAccess() {
	if !hasPermissions() {
		slog.Warn("PromptRootAccess not implemented in Linux yet, run the program as root")
	}
}
