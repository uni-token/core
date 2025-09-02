//go:build linux

package openBrowser

import (
	"os"
	"os/exec"
	"strings"
	"uni-token-service/constants"
)

func OpenBrowser(targetUser string, url string) error {
	providers := []string{"xdg-open", "x-www-browser", "www-browser"}

	for _, provider := range providers {
		if _, err := exec.LookPath(provider); err == nil {
			var cmd *exec.Cmd
			if constants.ShouldChangeUser {
				cmd = exec.Command("sudo", "-i", "-u", targetUser, provider, url)
			} else {
				cmd = exec.Command(provider, url)
			}
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			return cmd.Run()
		}
	}

	return &exec.Error{Name: strings.Join(providers, ","), Err: exec.ErrNotFound}
}
