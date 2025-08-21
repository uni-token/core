package discovery

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetServiceRootPath() string {
	overrideRoot := os.Getenv("UNI_TOKEN_SERVICE_ROOT")
	if overrideRoot != "" {
		return overrideRoot
	}
	if runtime.GOOS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		return filepath.Join(localAppData, "UniToken")
	} else {
		home := os.Getenv("HOME")
		return filepath.Join(home, ".local", "share", "uni-token")
	}
}

func GetServiceExecutablePath() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(GetServiceRootPath(), "service.exe")
	} else {
		return filepath.Join(GetServiceRootPath(), "service")
	}
}

func getServiceJsonPath() string {
	return filepath.Join(GetServiceRootPath(), "service.json")
}

func GetDbPath() string {
	return filepath.Join(GetServiceRootPath(), "data.db")
}

func SetupFileDiscovery(port int) error {
	filePath := getServiceJsonPath()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// Write initial service data
	if err := os.WriteFile(filePath, []byte(GetServiceInfo(&port)), 0644); err != nil {
		return err
	}

	fmt.Printf("Service discovery file created at %s\n", filePath)

	return nil
}
