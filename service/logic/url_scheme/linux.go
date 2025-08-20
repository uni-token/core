//go:build linux

package urlScheme

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Stub for Windows function when building on Linux
func registerURLSchemeWindows(options UrlSchemeRegisterOption) error {
	return fmt.Errorf("windows URL scheme registration not available on linux")
}

// registerURLSchemeForLinux registers URL scheme on Linux
func registerURLSchemeLinux(options UrlSchemeRegisterOption) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// .desktop files are usually placed in ~/.local/share/applications/
	desktopDir := filepath.Join(homeDir, ".local", "share", "applications")
	if _, err := os.Stat(desktopDir); os.IsNotExist(err) {
		if err := os.MkdirAll(desktopDir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", desktopDir, err)
		}
	}

	// Define .desktop file name and path
	desktopFileName := fmt.Sprintf("%s.desktop", options.AppName)
	desktopFilePath := filepath.Join(desktopDir, desktopFileName)

	// Write .desktop file content
	content := fmt.Sprintf(`[Desktop Entry]
Name=%s
Comment=Opens %s:// links
Exec=%s url %%u
Terminal=false
Type=Application
MimeType=x-scheme-handler/%s;
`, options.AppName, options.Scheme, options.ExecutablePath, options.Scheme)

	if err := os.WriteFile(desktopFilePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write .desktop file: %w", err)
	}
	// Use xdg-mime command to register MIME type
	// This tells the system which application handles 'x-scheme-handler/my-app'
	cmd := exec.Command("xdg-mime", "default", desktopFileName, fmt.Sprintf("x-scheme-handler/%s", options.Scheme))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("xdg-mime execution failed: %w\nOutput: %s", err, output)
	}

	// Some desktop environments may need to update desktop database, but xdg-mime default is usually sufficient
	// cmd = exec.Command("update-desktop-database", desktopDir)
	// if output, err := cmd.CombinedOutput(); err != nil {
	// 	fmt.Printf("⚠️ Warning: update-desktop-database may fail or be unnecessary: %w\nOutput: %s\n", err, output)
	// }

	return nil
}
