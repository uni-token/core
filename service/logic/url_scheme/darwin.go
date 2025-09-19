//go:build darwin

package urlScheme

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func RegisterURLScheme(options UrlSchemeRegisterOption) error {
	// On macOS, URL schemes are typically registered through:
	// 1. Creating an app bundle with Info.plist containing URL scheme registration
	// 2. Using Launch Services to register the scheme handler

	// For command-line applications, we'll create a minimal app bundle
	// and register it with the system

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create app bundle structure in ~/Applications/
	appName := options.AppName + ".app"
	appBundlePath := filepath.Join(homeDir, "Applications", appName)
	contentsPath := filepath.Join(appBundlePath, "Contents")
	macOSPath := filepath.Join(contentsPath, "MacOS")

	// Create directories
	if err := os.MkdirAll(macOSPath, 0755); err != nil {
		return fmt.Errorf("failed to create app bundle directories: %w", err)
	}

	// Create Info.plist with URL scheme registration
	infoPlistPath := filepath.Join(contentsPath, "Info.plist")
	infoPlistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>%s</string>
	<key>CFBundleIdentifier</key>
	<string>com.%s.%s</string>
	<key>CFBundleName</key>
	<string>%s</string>
	<key>CFBundleVersion</key>
	<string>1.0</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleURLTypes</key>
	<array>
		<dict>
			<key>CFBundleURLName</key>
			<string>%s URL Handler</string>
			<key>CFBundleURLSchemes</key>
			<array>
				<string>%s</string>
			</array>
		</dict>
	</array>
</dict>
</plist>`,
		options.AppName,
		strings.ToLower(options.AppName),
		strings.ToLower(options.AppName),
		options.AppName,
		options.AppName,
		options.Scheme)

	if err := os.WriteFile(infoPlistPath, []byte(infoPlistContent), 0644); err != nil {
		return fmt.Errorf("failed to write Info.plist: %w", err)
	}

	// Create executable wrapper script in the app bundle
	executablePath := filepath.Join(macOSPath, options.AppName)
	wrapperScript := fmt.Sprintf(`#!/bin/bash
# URL scheme handler wrapper for %s
exec "%s" url "$@"
`, options.AppName, options.ExecutablePath)

	if err := os.WriteFile(executablePath, []byte(wrapperScript), 0755); err != nil {
		return fmt.Errorf("failed to write wrapper executable: %w", err)
	}

	// Register the app bundle with Launch Services
	// This makes the system aware of the URL scheme handler
	cmd := exec.Command("/System/Library/Frameworks/CoreServices.framework/Frameworks/LaunchServices.framework/Support/lsregister",
		"-f", appBundlePath)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to register app bundle with Launch Services: %w\nOutput: %s", err, output)
	}

	// Optional: Set as default handler for the scheme
	// This requires the user to approve the change, so we'll skip it for now
	// cmd = exec.Command("open", "-a", appBundlePath, fmt.Sprintf("%s://test", options.Scheme))

	return nil
}
