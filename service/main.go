package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/kardianos/service"

	"uni-token-service/discovery"
	"uni-token-service/logic"
	urlScheme "uni-token-service/logic/url_scheme"
	"uni-token-service/server"
	"uni-token-service/store"
)

const (
	serviceName        = "UniTokenService"
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
	store.Init("data.db")

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
	svcConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceDisplayName,
		Description: serviceDescription,
		Arguments:   []string{},
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
		err := s.Run()
		if err != nil {
			panic(err)
		}
		return
	}
	command := os.Args[1]

	if command == "version" {
		fmt.Printf("Service Version: %d\n", logic.GetVersion())
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

	if command == "setup-in-sudo" {
		handleSetupInSudo(&s)
		return
	}

	prg.logger.Infof("Executing command: %s\n", command)

	err = service.Control(s, command)
	if err != nil {
		panic(err)
	}
}

func handleSetup() {
	// Register URL scheme for the application
	err := urlScheme.RegisterURLScheme(urlScheme.UrlSchemeRegisterOption{
		Scheme:  "uni-token",
		AppName: "UniToken",
	})
	if err != nil {
		panic(err)
	}

	// Install and start the service
	execPath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	res, err := logic.SudoExec(
		execPath+" setup-in-sudo",
		&logic.SudoOptions{
			Name: serviceDisplayName,
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

func handleSetupInSudo(s *service.Service) {
	err := service.Control(*s, "install")
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		fmt.Println("Failed to install service:", err)
	}
	err = service.Control(*s, "start")
	if err != nil {
		panic(err)
	}
}
