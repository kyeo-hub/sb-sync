package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kardianos/service"
	"sb-sync/pkg/config"
)

type program struct {
	exit chan struct{}
	cmd  *exec.Cmd
}

func (p *program) Start(s service.Service) error {
	p.exit = make(chan struct{})
	go p.run()
	return nil
}

func (p *program) run() {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".sb-sync", "config.json")
	binPath := filepath.Join(config.AppConfig.InstallDir, "sing-box")
	if filepath.Separator == '\\' {
		binPath += ".exe"
	}

	for {
		p.cmd = exec.Command(binPath, "run", "-c", configPath)
		p.cmd.Stdout = os.Stdout
		p.cmd.Stderr = os.Stderr
		
		err := p.cmd.Run()
		if err != nil {
			fmt.Printf("sing-box exited with error: %v\n", err)
		}

		select {
		case <-p.exit:
			return
		default:
			// Restart on crash
			fmt.Println("Restarting sing-box...")
		}
	}
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}
	return nil
}

func GetServiceConfig() *service.Config {
	return &service.Config{
		Name:        "sb-sync-singbox",
		DisplayName: "Sing-Box (sb-sync managed)",
		Description: "Sing-Box proxy service managed by sb-sync",
	}
}

func NewService() (service.Service, error) {
	prg := &program{}
	s, err := service.New(prg, GetServiceConfig())
	if err != nil {
		return nil, err
	}
	return s, nil
}
