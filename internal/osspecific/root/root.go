/*
Package root implements OS specific functionality to gain root privileges for the application.
*/
package root

import (
	"os"
)

const appendFlag = "-xroot"

func hasPermissions() bool {
	return os.Args[len(os.Args)-1] == appendFlag
}
