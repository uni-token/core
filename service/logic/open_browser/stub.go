//go:build !linux && !windows

package openBrowser

import "fmt"

func OpenBrowser(targetUser string, url string) error {
	return fmt.Errorf("OpenBrowser not implemented on this platform")
}
