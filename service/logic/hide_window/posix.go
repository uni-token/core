//go:build !windows

package hideWindow

import "os/exec"

func HideWindow(cmd *exec.Cmd) {
	// No-op on Unix-like systems
}
