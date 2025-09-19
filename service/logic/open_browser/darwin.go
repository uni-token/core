//go:build darwin

package openBrowser

import (
	"os"
	"os/exec"
	"uni-token-service/constants"
)

func OpenBrowser(targetUser string, url string) error {
	var cmd *exec.Cmd
	if constants.ShouldChangeUser {
		// Use sudo to run as the target user on macOS
		cmd = exec.Command("sudo", "-u", targetUser, "open", url)
	} else {
		// Use the native macOS 'open' command
		cmd = exec.Command("open", url)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
