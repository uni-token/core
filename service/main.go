package main

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/kardianos/service"

	"uni-token-service/constants"
	"uni-token-service/discovery"
	"uni-token-service/logic"
	"uni-token-service/logic/url_scheme"
	"uni-token-service/server"
	"uni-token-service/store"
)

const (
	serviceNamePrefix  = "UniTokenService"
	serviceDisplayName = "UniToken Service"
	serviceDescription = "UniToken Service"
)

type program struct {
	exit   chan struct{}
	logger service.Logger
	port   int
}

func (p *program) Start(s service.Service) error {
	p.logger.Infof("Service '%s' is starting...", serviceDisplayName)
	p.exit = make(chan struct{})

	go p.run()
	return nil
}

func (p *program) run() {
	store.Init(discovery.GetDbPath())

	p.logger.Info("Service is running. Starting main logic...")

	port, err := server.SetupAPIServer()
	if err != nil {
		p.logger.Errorf("Failed to setup API server: %v", err)
		return
	}
	p.port = port

	time.Sleep(100 * time.Millisecond)

	if err := discovery.SetupFileDiscovery(port); err != nil {
		p.logger.Errorf("Failed to setup file discovery: %v", err)
		return
	}

	p.logger.Infof("Service started successfully on port %d", port)

	<-p.exit
	p.logger.Info("Service main logic stopped.")
}

func (p *program) Stop(s service.Service) error {
	p.logger.Info("Service is stopping...")
	close(p.exit)
	return nil
}

func main() {
	username := constants.GetUserName()
	serviceName := serviceNamePrefix + "-" + url.PathEscape(username)

	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceDisplayName + " - " + username,
		Description: serviceDescription,
		Arguments:   []string{"run"},
		EnvVars: map[string]string{
			"UNI_TOKEN_SERVICE_USER": username,
			"UNI_TOKEN_SERVICE_ROOT": discovery.GetServiceRootPath(),
		},
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		panic(err)
	}

	prg.logger, err = s.Logger(nil)
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		// Directly run the executable
		handleSetup()
		logic.OpenUI(url.Values{}, false)
		return
	}

	command := os.Args[1]

	if command == "run" {
		err := s.Run()
		if err != nil {
			panic(err)
		}
		return
	}

	if command == "version" {
		fmt.Printf("Service Version: %d\n", constants.GetVersion())
		return
	}

	if command == "url" {
		handleUrlScheme(os.Args[2])
		return
	}

	if command == "setup" {
		handleSetup()
		return
	}

	if command == "install-and-start" {
		handleInstallAndStart(&s, serviceName)
		return
	}

	if command == "uninstall" {
		handleSudo(false, []string{"uninstall-impl"})
		return
	}

	if command == "uninstall-impl" {
		handleUninstall(s, serviceName)
		return
	}

	if command == "sudo" {
		if len(os.Args) < 3 {
			fmt.Println("Usage: sudo <command>")
			return
		}
		handleSudo(false, os.Args[2:])
		return
	}

	prg.logger.Infof("Executing command: %s\n", command)

	err = service.Control(s, command)
	if err != nil {
		panic(err)
	}
}

func handleSetup() {
	err := discovery.InstallExecutable()
	if err != nil {
		panic(err)
	}

	if discovery.IsServiceRunning() {
		return
	}

	// Register URL scheme for the application
	execPath := discovery.GetServiceExecutablePath()
	err = urlScheme.RegisterURLScheme(urlScheme.UrlSchemeRegisterOption{
		Scheme:         "uni-token",
		AppName:        "UniToken",
		ExecutablePath: execPath,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sudo is required to install and start the service.\n")

	// Install and start the service in sudo mode
	handleSudo(true, []string{"install-and-start"})
}

func handleSudo(useInstalled bool, args []string) {
	var execPath string
	if useInstalled {
		execPath = discovery.GetServiceExecutablePath()
	} else {
		var err error
		execPath, err = os.Executable()
		if err != nil {
			panic(err)
		}
	}
	command := execPath + " " + strings.Join(args, " ")

	fmt.Printf("Executing command with sudo: %s\n", command)

	res, err := logic.SudoExec(
		command,
		&logic.SudoOptions{
			Name: serviceDisplayName,
			Env: map[string]string{
				"UNI_TOKEN_SERVICE_USER": constants.GetUserName(),
				"UNI_TOKEN_SERVICE_ROOT": discovery.GetServiceRootPath(),
			},
		},
	)
	if err != nil {
		panic(err)
	}
	if res.Stdout != "" {
		fmt.Println(res.Stdout)
	}
	if res.Stderr != "" {
		fmt.Println(res.Stderr)
	}
}

func handleUrlScheme(url string) {
	if !strings.HasPrefix(url, "uni-token://") {
		fmt.Println("Invalid URL scheme. Expected 'uni-token://'.")
		return
	}
	url = strings.TrimPrefix(url, "uni-token://")

	switch url {
	case "start":
		handleSetup()
	default:
		fmt.Printf("Unknown URL action: %s\n", url)
	}
}

func handleInstallAndStart(s *service.Service, serviceName string) {
	err := service.Control(*s, "install")
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		fmt.Println("Failed to install service:", err)
	}
	fmt.Printf("Installed service \"%s\".\n", serviceName)

	err = service.Control(*s, "start")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Started service \"%s\".\n", serviceName)
}

func handleUninstall(s service.Service, serviceName string) {
	err := service.Control(s, "stop")
	if err != nil {
		fmt.Printf("Failed to stop service: %v\n", err)
	} else {
		fmt.Printf("Stopped service \"%s\".\n", serviceName)
	}

	err = service.Control(s, "uninstall")
	if err != nil {
		fmt.Printf("Failed to uninstall service: %v\n", err)
	} else {
		fmt.Printf("Uninstalled service \"%s\".\n", serviceName)
	}

	// err = urlScheme.UnregisterURLScheme()
	// if err != nil {
	// 	fmt.Printf("Failed to unregister URL scheme: %v\n", err)
	// }

	rootPath := discovery.GetServiceRootPath()
	err = os.RemoveAll(rootPath)
	if err != nil {
		fmt.Printf("Failed to remove service root path %s: %v\n", rootPath, err)
	} else {
		fmt.Printf("Removed service root path %s.\n", rootPath)
	}
}
