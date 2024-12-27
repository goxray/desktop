package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/getlantern/elevate"
)

const appPackageIcons = "/../Resources/icon.icns"

func PromptRootAccess() {
	// A really hacky way to get admin password prompt for macOS
	if !hasPermissions() {
		runItselfAsRoot()
		os.Exit(0)
	}
}

func runItselfAsRoot() {
	p, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("could not get executable path: %v", err))
	}
	cmd := elevate.WithIcon(filepath.Dir(p)+appPackageIcons).WithPrompt("Go XRay Client requires admin privileges").
		Command(p, append(os.Args[1:], appendFlag)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
