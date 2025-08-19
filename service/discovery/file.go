package discovery

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

func getFilePath() string {
	if runtime.GOOS == "windows" {
		localAppData := os.Getenv("LOCALAPPDATA")
		return filepath.Join(localAppData, "UnitedToken", "service.json")
	}
	home := os.Getenv("HOME")
	return filepath.Join(home, ".local", "share", "uni-token", "service.json")
}

func SetupFileDiscovery(port int) error {
	filePath := getFilePath()

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	// Write initial service data
	if err := os.WriteFile(filePath, []byte(GetData(&port)), 0644); err != nil {
		return err
	}

	fmt.Printf("Service discovery file created at %s\n", filePath)

	// Setup cleanup on exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		// Write final service data with null port on exit
		os.WriteFile(filePath, []byte(GetData(nil)), 0644)
		os.Exit(0)
	}()

	return nil
}
