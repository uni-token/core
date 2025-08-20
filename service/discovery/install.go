package discovery

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
)

func InstallExecutable() error {
	targetPath := GetServiceExecutablePath()
	selfPath := os.Args[0]
	if runtime.GOOS != "windows" {
		var err error
		selfPath, err = os.Executable()
		if err != nil {
			return err
		}
	}

	if targetPath == selfPath {
		return nil
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		if err := copyFile(selfPath, targetPath); err != nil {
			return err
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	err = os.Chmod(dst, 0755)
	if err != nil {
		return err
	}

	return nil
}
