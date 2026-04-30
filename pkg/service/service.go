package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
	"sb-sync/pkg/config"
	"sb-sync/pkg/sync"
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
	configPath := config.GetConfigPath()
	binPath := filepath.Join(config.AppConfig.InstallDir, "sing-box")
	if filepath.Separator == '\\' {
		binPath += ".exe"
	}

	// Initial sync
	fmt.Println("Performing initial sync...")
	if err := sync.SyncFromWebDAV(); err != nil {
		fmt.Printf("Initial sync failed: %v\n", err)
	}

	interval := time.Duration(config.AppConfig.SyncInterval) * time.Minute
	if interval < time.Minute {
		interval = time.Hour // Default to 1 hour if interval is invalid
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		p.cmd = exec.Command(binPath, "run", "-c", configPath)
		p.cmd.Stdout = os.Stdout
		p.cmd.Stderr = os.Stderr
		
		if err := p.cmd.Start(); err != nil {
			fmt.Printf("Failed to start sing-box: %v\n", err)
			time.Sleep(10 * time.Second)
			continue
		}

		done := make(chan error, 1)
		go func() {
			done <- p.cmd.Wait()
		}()

		select {
		case <-p.exit:
			if p.cmd != nil && p.cmd.Process != nil {
				p.cmd.Process.Kill()
			}
			return
		case err := <-done:
			if err != nil {
				fmt.Printf("sing-box exited with error: %v\n", err)
			}
			fmt.Println("sing-box stopped. Restarting in 5 seconds...")
			time.Sleep(5 * time.Second)
		case <-ticker.C:
			fmt.Println("Checking for config updates...")
			if _, updated, err := sync.SyncFromWebDAVWithStatus(); err == nil {
				if updated {
					fmt.Println("Config updated, restarting sing-box...")
					if p.cmd != nil && p.cmd.Process != nil {
						p.cmd.Process.Kill()
					}
				} else {
					fmt.Println("Config is up to date.")
				}
			} else {
				fmt.Printf("Sync failed: %v\n", err)
			}
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
