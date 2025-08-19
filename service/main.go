package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kardianos/service"

	"uni-token-service/discovery"
	"uni-token-service/logic"
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

	if command == "sudo" {
		command, err := filepath.Abs(os.Args[0])
		if err != nil {
			panic(err)
		}
		for _, arg := range os.Args[2:] {
			command += " " + arg
		}

		res, err := logic.SudoExec(command, &logic.SudoOptions{
			Name: serviceDisplayName,
		})
		if err != nil {
			panic(err)
		}
		if res.Stdout != "" {
			fmt.Println(res.Stdout)
		}
		if res.Stderr != "" {
			fmt.Println(res.Stderr)
		}
		return
	}

	if command == "version" {
		fmt.Printf("Service Version: %d\n", logic.GetVersion())
		return
	}

	if command == "setup" {
		err = service.Control(s, "install")
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			fmt.Println("Failed to install service:", err)
		}
		err = service.Control(s, "start")
		if err != nil {
			panic(err)
		}
		return
	}

	prg.logger.Infof("Executing command: %s\n", command)

	err = service.Control(s, command)
	if err != nil {
		panic(err)
	}
}
