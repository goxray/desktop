/*
Package root implements OS specific functionality to gain root privileges for the application.
*/
package root

import (
	"os"
	"slices"
)

const appendFlag = "-xroot=true"

func hasPermissions() bool {
	return slices.Contains(os.Args[1:], appendFlag)
}
