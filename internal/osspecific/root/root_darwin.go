package root

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
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
	// Keep theme settings of the user.
	themeFlag := fmt.Sprintf("-theme_variant=%d", fyne.CurrentApp().Settings().ThemeVariant())

	p, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("could not get executable path: %v", err))
	}
	cmd := elevate.WithIcon(filepath.Dir(p)+appPackageIcons).WithPrompt("GoXRay requires admin privileges").
		Command(p, append(os.Args[1:], appendFlag, themeFlag)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()
}
