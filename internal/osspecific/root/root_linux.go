package root

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"syscall"
)

func PromptRootAccess() {}

// Init is used on Linux to temporarily downgrade root privileges to user.
// It is required to properly initialize tray menu dbus.
func Init() {
	downgradePrivileges()
}

func OnStart() {
	resetPrivileges()
}

func downgradePrivileges() {
	if runtime.GOOS != "linux" {
		return
	}

	// sudo adds SUDO_UID env variable with original user id
	if os.Getenv("SUDO_UID") == "" {
		log.Fatal("compatibility issue: unable to detect sudo")
	}

	usrUID, err := strconv.Atoi(os.Getenv("SUDO_UID"))
	if err != nil {
		log.Fatal(fmt.Errorf("sudo_uid is not int: %q: %w", os.Getenv("SUDO_UID"), err))
	}

	if err := syscall.Setuid(usrUID); err != nil {
		log.Fatal(err)
	}
}

func resetPrivileges() {
	if runtime.GOOS != "linux" {
		return
	}

	if err := syscall.Setuid(0); err != nil {
		log.Fatal(err)
	}
}
