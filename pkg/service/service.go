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

type ServiceInterface interface {
	Install() error
	Uninstall() error
	Start() error
	Stop() error
	Status() (string, error)
}

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
	binPath := config.GetSingBoxBinary()

	fmt.Printf("%s Performing initial sync...\n", config.LogPrefixInfo)
	if err := sync.SyncFromWebDAV(); err != nil {
		fmt.Printf("%s Initial sync failed: %v\n", config.LogPrefixWarn, err)
	}

	interval := time.Duration(config.AppConfig.SyncInterval) * time.Minute
	if interval < config.MinimumSyncInterval {
		interval = config.DefaultSyncInterval
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		p.cmd = exec.Command(binPath, "run", "-c", configPath)
		p.cmd.Stdout = os.Stdout
		p.cmd.Stderr = os.Stderr

		if err := p.cmd.Start(); err != nil {
			fmt.Printf("%s Failed to start sing-box: %v\n", config.LogPrefixError, err)
			time.Sleep(config.RestartDelay)
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
				fmt.Printf("%s sing-box exited with error: %v\n", config.LogPrefixWarn, err)
			}
			fmt.Printf("%s sing-box stopped. Restarting in %v...\n", config.LogPrefixInfo, config.RestartDelay)
			time.Sleep(config.RestartDelay)
		case <-ticker.C:
			fmt.Printf("%s Checking for config updates...\n", config.LogPrefixInfo)
			if _, updated, err := sync.SyncFromWebDAVWithStatus(); err == nil {
				if updated {
					fmt.Printf("%s Config updated, restarting sing-box...\n", config.LogPrefixInfo)
					if p.cmd != nil && p.cmd.Process != nil {
						p.cmd.Process.Kill()
					}
				} else {
					fmt.Printf("%s Config is up to date.\n", config.LogPrefixInfo)
				}
			} else {
				fmt.Printf("%s Sync failed: %v\n", config.LogPrefixError, err)
			}
		}
	}
}

func (p *program) Stop(s service.Service) error {
	close(p.exit)
	time.Sleep(config.StopDelay)
	if p.cmd != nil && p.cmd.Process != nil {
		p.cmd.Process.Kill()
	}
	return nil
}

func GetServiceConfig() *service.Config {
	return &service.Config{
		Name:        config.ServiceName,
		DisplayName: config.ServiceDisplayName,
		Description: config.ServiceDescription,
	}
}

func NewService() (ServiceInterface, error) {
	prg := &program{}
	s, err := service.New(prg, GetServiceConfig())
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	return &serviceWrapper{s}, nil
}

type serviceWrapper struct {
	s service.Service
}

func (w *serviceWrapper) Install() error {
	return w.s.Install()
}

func (w *serviceWrapper) Uninstall() error {
	return w.s.Uninstall()
}

func (w *serviceWrapper) Start() error {
	return w.s.Start()
}

func (w *serviceWrapper) Stop() error {
	return w.s.Stop()
}

func (w *serviceWrapper) Status() (string, error) {
	status, err := w.s.Status()
	if err != nil {
		return "", fmt.Errorf("failed to get status: %w", err)
	}
	return fmt.Sprintf("%v", status), nil
}

func GetSingBoxBinaryPath() string {
	binPath := filepath.Join(config.GetInstallDir(), "sing-box")
	if filepath.Separator == '\\' {
		binPath += ".exe"
	}
	return binPath
}
