package logic

// Credit: https://github.com/jorangreef/sudo-prompt

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	PermissionDenied = "user did not grant permission"
	NoPolkitAgent    = "no polkit authentication agent found"
	MaxBuffer        = 134217728
)

var (
	envVarRegex        = regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]*$")
	sudoRegex          = regexp.MustCompile("(?i)^sudo")
	nameRegex          = regexp.MustCompile("^[a-z0-9 ]+$")
	windowsEscapeRegex = regexp.MustCompile(`([<>\\|&^])`)
)

// SudoOptions represents the options for executing a command with elevated privileges
type SudoOptions struct {
	Name string            `json:"name"`
	Icns string            `json:"icns"`
	Env  map[string]string `json:"env"`
}

// Instance represents an instance of a sudo command execution
type Instance struct {
	Command string
	Options SudoOptions
	UUID    string
	Path    string
	// Windows specific paths
	PathElevate string
	PathExecute string
	PathCommand string
	PathStdout  string
	PathStderr  string
	PathStatus  string
}

// Result represents the result of command execution
type Result struct {
	Stdout string
	Stderr string
	Error  error
}

// SudoExec executes a command with elevated privileges
func SudoExec(command string, options *SudoOptions) (*Result, error) {
	if command == "" {
		return nil, errors.New("command should be a string")
	}

	if sudoRegex.MatchString(command) {
		return nil, errors.New("command should not be prefixed with \"sudo\"")
	}

	// Set default options
	if options == nil {
		options = &SudoOptions{}
	}

	// Set default name from process title if not provided
	if options.Name == "" {
		title := filepath.Base(os.Args[0])
		if ValidName(title) {
			options.Name = title
		} else {
			return nil, errors.New("process title cannot be used as a valid name")
		}
	} else if !ValidName(options.Name) {
		return nil, errors.New("options.name must be alphanumeric only (spaces are allowed) and <= 70 characters")
	}

	// Validate icns
	if options.Icns != "" {
		if strings.TrimSpace(options.Icns) == "" {
			return nil, errors.New("options.icns must not be empty if provided")
		}
	}

	// Validate environment variables
	if options.Env != nil {
		if len(options.Env) == 0 {
			return nil, errors.New("options.env must not be empty if provided")
		}
		for key, value := range options.Env {
			if !envVarRegex.MatchString(key) {
				return nil, fmt.Errorf("options.env has an invalid environment variable name: %q", key)
			}
			if strings.ContainsAny(value, "\r\n") {
				return nil, fmt.Errorf("options.env has an invalid environment variable value: %q", value)
			}
		}
	}

	platform := runtime.GOOS
	if platform != "darwin" && platform != "linux" && platform != "windows" {
		return nil, errors.New("platform not yet supported")
	}

	instance := &Instance{
		Command: command,
		Options: *options,
	}

	return attempt(instance)
}

// attempt attempts to execute the command based on the platform
func attempt(instance *Instance) (*Result, error) {
	platform := runtime.GOOS
	switch platform {
	case "darwin":
		return macExec(instance)
	case "linux":
		return linuxExec(instance)
	case "windows":
		return windowsExec(instance)
	default:
		return nil, errors.New("platform not yet supported")
	}
}

// escapeDoubleQuotes escapes double quotes in a string
func escapeDoubleQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

// ValidName validates if a name is valid
func ValidName(name string) bool {
	// We use 70 characters as a limit to side-step any issues with Unicode
	// normalization form causing a 255 character string to exceed the fs limit.
	if !nameRegex.MatchString(strings.ToLower(name)) {
		return false
	}
	if strings.TrimSpace(name) == "" {
		return false
	}
	if len(name) > 70 {
		return false
	}
	return true
}

// generateUUID generates a UUID for the instance
func generateUUID(instance *Instance) (string, error) {
	random := make([]byte, 256)
	_, err := rand.Read(random)
	if err != nil {
		// Fallback to timestamp and random
		random = []byte(fmt.Sprintf("%d%f", time.Now().UnixNano(), time.Since(time.Unix(0, 0)).Seconds()))
	}

	hash := sha256.New()
	hash.Write([]byte("sudo-prompt-3"))
	hash.Write([]byte(instance.Options.Name))
	hash.Write([]byte(instance.Command))
	hash.Write(random)

	uuid := hex.EncodeToString(hash.Sum(nil))
	if len(uuid) < 32 {
		return "", errors.New("expected a valid UUID")
	}
	uuid = uuid[len(uuid)-32:]

	if len(uuid) != 32 {
		return "", errors.New("expected a valid UUID")
	}

	return uuid, nil
}

// remove removes a file or directory recursively
func remove(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("argument path not defined")
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		if strings.Contains(path, "\"") {
			return errors.New("argument path cannot contain double-quotes")
		}
		cmd = exec.Command("cmd", "/c", fmt.Sprintf(`rmdir /s /q "%s"`, path))
	} else {
		cmd = exec.Command("/bin/rm", "-rf", escapeDoubleQuotes(filepath.Clean(path)))
	}

	return cmd.Run()
}

// linuxExec executes a command on Linux with elevated privileges
func linuxExec(instance *Instance) (*Result, error) {
	binary, err := linuxBinary()
	if err != nil {
		return nil, err
	}

	var command []string

	// Preserve current working directory
	cwd, _ := os.Getwd()
	command = append(command, fmt.Sprintf(`cd "%s";`, escapeDoubleQuotes(cwd)))

	// Export environment variables
	for key, value := range instance.Options.Env {
		command = append(command, fmt.Sprintf(`export %s="%s";`, key, escapeDoubleQuotes(value)))
	}

	command = append(command, fmt.Sprintf(`"%s"`, escapeDoubleQuotes(binary)))

	if strings.Contains(strings.ToLower(binary), "kdesudo") {
		command = append(command, "--comment")
		command = append(command, fmt.Sprintf(`"%s wants to make changes. Enter your password to allow this."`, instance.Options.Name))
		command = append(command, "-d") // Do not show the command to be run in the dialog
		command = append(command, "--")
	} else if strings.Contains(strings.ToLower(binary), "pkexec") {
		command = append(command, "--disable-internal-agent")
	}

	magic := "SUDOPROMPT\n"
	bashCommand := fmt.Sprintf(`/bin/bash -c "echo %s; %s"`,
		escapeDoubleQuotes(strings.TrimSpace(magic)),
		escapeDoubleQuotes(instance.Command))
	command = append(command, bashCommand)

	cmdStr := strings.Join(command, " ")
	cmd := exec.Command("/bin/sh", "-c", cmdStr)

	output, err := cmd.CombinedOutput()
	stdout := string(output)
	stderr := ""

	// Check if elevation succeeded using magic marker
	elevated := strings.HasPrefix(stdout, magic)
	if elevated {
		stdout = stdout[len(magic):]
	}

	// Only normalize the error if it is definitely not a command error
	if err != nil && !elevated {
		if strings.Contains(stderr, "No authentication agent found") {
			return &Result{
				Stdout: stdout,
				Stderr: stderr,
				Error:  errors.New(NoPolkitAgent),
			}, nil
		} else {
			return &Result{
				Stdout: stdout,
				Stderr: stderr,
				Error:  errors.New(PermissionDenied),
			}, nil
		}
	}

	result := &Result{
		Stdout: stdout,
		Stderr: stderr,
		Error:  err,
	}
	return result, nil
}

// linuxBinary finds the appropriate binary for Linux elevation
func linuxBinary() (string, error) {
	paths := []string{"/usr/bin/kdesudo", "/usr/bin/pkexec"}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", errors.New("unable to find pkexec or kdesudo")
}

// macExec executes a command on macOS with elevated privileges
func macExec(instance *Instance) (*Result, error) {
	temp := os.TempDir()
	if temp == "" {
		return nil, errors.New("os.TempDir() not defined")
	}

	user := os.Getenv("USER")
	if user == "" {
		return nil, errors.New("env['USER'] not defined")
	}

	uuid, err := generateUUID(instance)
	if err != nil {
		return nil, err
	}

	instance.UUID = uuid
	instance.Path = filepath.Join(temp, instance.UUID, instance.Options.Name+".app")

	// Create cleanup function
	cleanup := func() error {
		return remove(filepath.Dir(instance.Path))
	}

	// Execute the macOS-specific logic
	err = macApplet(instance)
	if err != nil {
		cleanup()
		return nil, err
	}

	err = macIcon(instance)
	if err != nil {
		cleanup()
		return nil, err
	}

	err = macPropertyList(instance)
	if err != nil {
		cleanup()
		return nil, err
	}

	err = macCommand(instance)
	if err != nil {
		cleanup()
		return nil, err
	}

	stdout, stderr, err := macOpen(instance)
	if err != nil {
		cleanup()
		return &Result{
			Stdout: stdout,
			Stderr: stderr,
			Error:  err,
		}, nil
	}

	result, err := macResult(instance)
	cleanup()
	return result, err
}

// macApplet creates the macOS applet
func macApplet(instance *Instance) error {
	parent := filepath.Dir(instance.Path)
	err := os.MkdirAll(parent, 0755)
	if err != nil {
		return err
	}

	zip := filepath.Join(parent, "sudo-prompt-applet.zip")

	const APPLET_DATA = "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"

	data, err := base64.StdEncoding.DecodeString(APPLET_DATA)
	if err != nil {
		return err
	}

	err = os.WriteFile(zip, data, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("/usr/bin/unzip", "-o", zip, "-d", instance.Path)
	return cmd.Run()
}

// macIcon sets up the icon for the macOS applet
func macIcon(instance *Instance) error {
	if instance.Options.Icns == "" {
		return nil
	}

	data, err := os.ReadFile(instance.Options.Icns)
	if err != nil {
		return err
	}

	icns := filepath.Join(instance.Path, "Contents", "Resources", "applet.icns")
	return os.WriteFile(icns, data, 0644)
}

// macCommand creates the command script for macOS
func macCommand(instance *Instance) error {
	path := filepath.Join(instance.Path, "Contents", "MacOS", "sudo-prompt-command")

	var script []string

	// Preserve current working directory
	cwd, _ := os.Getwd()
	script = append(script, fmt.Sprintf(`cd "%s"`, escapeDoubleQuotes(cwd)))

	// Export environment variables
	for key, value := range instance.Options.Env {
		script = append(script, fmt.Sprintf(`export %s="%s"`, key, escapeDoubleQuotes(value)))
	}

	script = append(script, instance.Command)

	return os.WriteFile(path, []byte(strings.Join(script, "\n")), 0755)
}

// macPropertyList sets up the property list for macOS
func macPropertyList(instance *Instance) error {
	plist := filepath.Join(instance.Path, "Contents", "Info.plist")
	value := instance.Options.Name + " Password Prompt"

	if strings.Contains(value, "'") {
		return errors.New("value should not contain single quotes")
	}

	cmd := exec.Command("/usr/bin/defaults", "write", plist, "CFBundleName", value)
	return cmd.Run()
}

// macOpen opens the macOS applet
func macOpen(instance *Instance) (string, string, error) {
	binary := filepath.Join(instance.Path, "Contents", "MacOS", "applet")
	cwd := filepath.Dir(binary)

	cmd := exec.Command("./" + filepath.Base(binary))
	cmd.Dir = cwd

	output, err := cmd.CombinedOutput()
	return string(output), "", err
}

// macResult processes the result from macOS execution
func macResult(instance *Instance) (*Result, error) {
	cwd := filepath.Join(instance.Path, "Contents", "MacOS")

	codeData, err := os.ReadFile(filepath.Join(cwd, "code"))
	if err != nil {
		if os.IsNotExist(err) {
			return &Result{
				Stdout: "",
				Stderr: "",
				Error:  errors.New(PermissionDenied),
			}, nil
		}
		return nil, err
	}

	stdoutData, err := os.ReadFile(filepath.Join(cwd, "stdout"))
	if err != nil {
		return nil, err
	}

	stderrData, err := os.ReadFile(filepath.Join(cwd, "stderr"))
	if err != nil {
		return nil, err
	}

	code, err := strconv.Atoi(strings.TrimSpace(string(codeData)))
	if err != nil {
		return nil, err
	}

	stdout := string(stdoutData)
	stderr := string(stderrData)

	result := &Result{
		Stdout: stdout,
		Stderr: stderr,
	}

	if code != 0 {
		result.Error = fmt.Errorf("Command failed: %s\n%s", instance.Command, stderr)
	}

	return result, nil
}

// windowsExec executes a command on Windows with elevated privileges
func windowsExec(instance *Instance) (*Result, error) {
	temp := os.TempDir()
	if temp == "" {
		return nil, errors.New("os.TempDir() not defined")
	}

	uuid, err := generateUUID(instance)
	if err != nil {
		return nil, err
	}

	instance.UUID = uuid
	instance.Path = filepath.Join(temp, instance.UUID)

	if strings.Contains(instance.Path, "\"") {
		return nil, errors.New("instance.path cannot contain double-quotes")
	}

	instance.PathElevate = filepath.Join(instance.Path, "elevate.vbs")
	instance.PathExecute = filepath.Join(instance.Path, "execute.bat")
	instance.PathCommand = filepath.Join(instance.Path, "command.bat")
	instance.PathStdout = filepath.Join(instance.Path, "stdout")
	instance.PathStderr = filepath.Join(instance.Path, "stderr")
	instance.PathStatus = filepath.Join(instance.Path, "status")

	err = os.MkdirAll(instance.Path, 0755)
	if err != nil {
		return nil, err
	}

	// Cleanup function
	cleanup := func() error {
		return remove(instance.Path)
	}

	err = windowsWriteExecuteScript(instance)
	if err != nil {
		cleanup()
		return nil, err
	}

	err = windowsWriteCommandScript(instance)
	if err != nil {
		cleanup()
		return nil, err
	}

	stdout, stderr, err := windowsElevate(instance)
	if err != nil {
		cleanup()
		return &Result{
			Stdout: stdout,
			Stderr: stderr,
			Error:  err,
		}, nil
	}

	err = windowsWaitForStatus(instance)
	if err != nil {
		cleanup()
		return nil, err
	}

	result, err := windowsResult(instance)
	cleanup()
	return result, err
}

// windowsElevate elevates the process on Windows
func windowsElevate(instance *Instance) (string, string, error) {
	escapedPath := strings.ReplaceAll(instance.PathExecute, "'", "`'")
	command := fmt.Sprintf(`powershell.exe Start-Process -FilePath "'%s'" -WindowStyle hidden -Verb runAs`, escapedPath)

	cmd := exec.Command("cmd", "/c", command)

	// Close stdin to prevent PowerShell from waiting indefinitely on Windows 7
	stdin, err := cmd.StdinPipe()
	if err == nil {
		stdin.Close()
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), "", errors.New(PermissionDenied)
	}

	return string(output), "", nil
}

// windowsWriteExecuteScript writes the execute script for Windows
func windowsWriteExecuteScript(instance *Instance) error {
	script := []string{
		"@echo off",
		fmt.Sprintf(`call "%s" > "%s" 2> "%s"`, instance.PathCommand, instance.PathStdout, instance.PathStderr),
		fmt.Sprintf(`(echo %%ERRORLEVEL%%) > "%s"`, instance.PathStatus),
	}

	return os.WriteFile(instance.PathExecute, []byte(strings.Join(script, "\r\n")), 0644)
}

// windowsWriteCommandScript writes the command script for Windows
func windowsWriteCommandScript(instance *Instance) error {
	cwd, _ := os.Getwd()
	if strings.Contains(cwd, "\"") {
		return errors.New("process.cwd() cannot contain double-quotes")
	}

	script := []string{
		"@echo off",
		"chcp 65001>nul",               // Set code page to UTF-8
		fmt.Sprintf(`cd /d "%s"`, cwd), // Preserve current working directory
	}

	// Export environment variables
	for key, value := range instance.Options.Env {
		escaped := windowsEscapeRegex.ReplaceAllString(value, "^$1")
		script = append(script, fmt.Sprintf("set %s=%s", key, escaped))
	}

	script = append(script, instance.Command)

	return os.WriteFile(instance.PathCommand, []byte(strings.Join(script, "\r\n")), 0644)
}

// windowsWaitForStatus waits for the status file to be created on Windows
func windowsWaitForStatus(instance *Instance) error {
	for {
		stat, err := os.Stat(instance.PathStatus)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		if err == nil && stat.Size() >= 2 {
			// File exists and has content
			break
		}

		// Wait and retry
		time.Sleep(1 * time.Second)

		// Check if stdout file exists to ensure command was executed
		if _, err := os.Stat(instance.PathStdout); err != nil {
			return errors.New(PermissionDenied)
		}
	}

	return nil
}

// windowsResult processes the result from Windows execution
func windowsResult(instance *Instance) (*Result, error) {
	codeData, err := os.ReadFile(instance.PathStatus)
	if err != nil {
		return nil, err
	}

	stdoutData, err := os.ReadFile(instance.PathStdout)
	if err != nil {
		return nil, err
	}

	stderrData, err := os.ReadFile(instance.PathStderr)
	if err != nil {
		return nil, err
	}

	code, err := strconv.Atoi(strings.TrimSpace(string(codeData)))
	if err != nil {
		return nil, err
	}

	stdout := string(stdoutData)
	stderr := string(stderrData)

	result := &Result{
		Stdout: stdout,
		Stderr: stderr,
	}

	if code != 0 {
		result.Error = fmt.Errorf("Command failed: %s\r\n%s", instance.Command, stderr)
	}

	return result, nil
}
