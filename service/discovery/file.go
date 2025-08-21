package discovery

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

func getServiceRootPath() string {
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
		return filepath.Join(getServiceRootPath(), "service.exe")
	} else {
		return filepath.Join(getServiceRootPath(), "service")
	}
}

func getServiceJsonPath() string {
	return filepath.Join(getServiceRootPath(), "service.json")
}

func GetDbPath() string {
	return filepath.Join(getServiceRootPath(), "data.db")
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

	// Setup cleanup on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		// Write final service data with null port on exit
		os.WriteFile(filePath, []byte(GetServiceInfo(nil)), 0644)
		os.Exit(0)
	}()

	return nil
}
