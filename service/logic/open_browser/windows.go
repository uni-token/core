//go:build windows

package openBrowser

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"uni-token-service/constants"
)

func OpenURLInUserSession(url string) error {
	sessionId := windows.WTSGetActiveConsoleSessionId()
	if sessionId == 0xFFFFFFFF { // 0xFFFFFFFF for no active session
		return fmt.Errorf("could not get active console session ID: no user is logged in")
	}

	var userToken windows.Token
	err := windows.WTSQueryUserToken(sessionId, &userToken)
	if err != nil {
		return fmt.Errorf("could not get user token for session %d: %v", sessionId, err)
	}
	defer userToken.Close()

	var envBlock *uint16
	err = windows.CreateEnvironmentBlock(&envBlock, userToken, false)
	if err != nil {
		return fmt.Errorf("could not create environment block: %v", err)
	}
	defer windows.DestroyEnvironmentBlock(envBlock)

	// chromePath := "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"
	// if _, err := os.Stat(chromePath); os.IsNotExist(err) {
	// 	chromePath = "C:\\Program Files (x86)\\Google\\Chrome\\Application\\chrome.exe"
	// 	if _, err := os.Stat(chromePath); os.IsNotExist(err) {
	// 		return fmt.Errorf("google Chrome not found at default paths")
	// 	}
	// }

	cmdPath, err := syscall.UTF16PtrFromString("C:\\Windows\\System32\\cmd.exe")
	if err != nil {
		return fmt.Errorf("failed to create cmd path string pointer: %v", err)
	}

	cmdLine := fmt.Sprintf(`/C start "" "%s"`, url)
	cmdLinePtr, err := syscall.UTF16PtrFromString(cmdLine)
	if err != nil {
		return fmt.Errorf("failed to create command line string pointer: %v", err)
	}

	startupInfo := &windows.StartupInfo{
		Cb:    uint32(unsafe.Sizeof(windows.StartupInfo{})),
		Flags: windows.STARTF_USESHOWWINDOW,
	}
	processInfo := &windows.ProcessInformation{}

	err = windows.CreateProcessAsUser(
		userToken,
		cmdPath,
		cmdLinePtr,
		nil,
		nil,
		false,
		windows.CREATE_UNICODE_ENVIRONMENT,
		envBlock,
		nil,
		startupInfo,
		processInfo,
	)
	if err != nil {
		return fmt.Errorf("CreateProcessAsUser failed: %v", err)
	}

	windows.CloseHandle(processInfo.Process)
	windows.CloseHandle(processInfo.Thread)

	return nil
}

func OpenBrowser(targetUser string, url string) error {
	if constants.ShouldChangeUser() {
		return OpenURLInUserSession(url)
	} else {
		return windows.ShellExecute(0, nil, windows.StringToUTF16Ptr(url), nil, nil, windows.SW_SHOWNORMAL)
	}
}
