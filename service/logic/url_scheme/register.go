package urlScheme

import (
	"fmt"
	"runtime"
)

type UrlSchemeRegisterOption struct {
	Scheme  string
	AppName string
}

func RegisterURLScheme(options UrlSchemeRegisterOption) error {
	switch runtime.GOOS {
	case "windows":
		return registerURLSchemeWindows(options)
	case "linux":
		return registerURLSchemeLinux(options)
	// case "darwin":
	// 	return registerURLSchemeMacOS(options)
	default:
		return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}
}
