//go:build windows

package urlScheme

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func RegisterURLScheme(options UrlSchemeRegisterOption) error {
	exePath := filepath.ToSlash(options.ExecutablePath)

	// Build registry path
	// HKEY_CURRENT_USER is user-specific, can be modified without admin privileges
	baseKey := "Software\\Classes\\" + options.Scheme
	commandKey := baseKey + "\\shell\\open\\command"

	// 1. Create or open main key (e.g.: HKEY_CURRENT_USER\Software\Classes\my-app)
	k, _, err := registry.CreateKey(registry.CURRENT_USER, baseKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create registry key %s: %w", baseKey, err)
	}
	defer k.Close()

	// Set default value and "URL Protocol" flag
	// @="URL: My App Protocol"
	if err := k.SetStringValue("", "URL: "+strings.ToUpper(options.Scheme)+" Protocol"); err != nil {
		return fmt.Errorf("failed to set default value for key %s: %w", baseKey, err)
	}
	// "URL Protocol"="" (empty string value, indicates this is a URL Protocol Handler)
	if err := k.SetStringValue("URL Protocol", ""); err != nil {
		return fmt.Errorf("failed to set URL Protocol value for key %s: %w", baseKey, err)
	}

	// 2. Create or open command key (e.g.: ...\my-app\shell\open\command)
	cmdK, _, err := registry.CreateKey(registry.CURRENT_USER, commandKey, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to create registry key %s: %w", commandKey, err)
	}
	defer cmdK.Close()

	// Set command key value to executable path and %1 (full URL placeholder)
	// @="\"C:\\Path\\To\\Your\\MyApp.exe\" \"%1\""
	commandValue := fmt.Sprintf(`"%s" "%%1"`, exePath) // Note: %% before %1 is for escaping
	if err := cmdK.SetStringValue("", commandValue); err != nil {
		return fmt.Errorf("failed to set command value for key %s: %w", commandKey, err)
	}

	return nil
}
