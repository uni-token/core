//go:build !linux && !windows

package urlScheme

import "fmt"

func registerURLSchemeWindows(options UrlSchemeRegisterOption) error {
	return fmt.Errorf("Windows URL scheme registration not implemented on this platform")
}

func registerURLSchemeLinux(options UrlSchemeRegisterOption) error {
	return fmt.Errorf("Linux URL scheme registration not implemented on this platform")
}
