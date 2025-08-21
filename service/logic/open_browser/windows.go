//go:build windows

package openBrowser

import "golang.org/x/sys/windows"

func OpenBrowser(targetUser string, url string) error {
	return windows.ShellExecute(0, nil, windows.StringToUTF16Ptr(url), nil, nil, windows.SW_SHOWNORMAL)
}
