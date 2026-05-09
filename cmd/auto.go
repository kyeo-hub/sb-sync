package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
)

var (
	autoImport  bool
	autoMigrate bool
	autoKill    bool
	autoAll     bool
)

var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Auto-detect and migrate existing sing-box installation",
	Long: `Auto-detect existing sing-box installation and optionally:
  - Stop running sing-box processes
  - Import configuration
  - Setup sb-sync service

This command helps migrate from manual sing-box installation to sb-sync managed service.`,
	Run: func(cmd *cobra.Command, args []string) {
		runAutoDetect()
	},
}

func runAutoDetect() {
	fmt.Println("=== Auto-detecting existing sing-box installation ===")
	fmt.Println()

	detected := false
	hasRunningProcess := false

	if binary, version := detectSingBoxBinary(); binary != "" {
		detected = true
		fmt.Printf("[FOUND] sing-box binary: %s\n", binary)
		fmt.Printf("[FOUND] Version: %s\n", version)

		if isProcessRunning(binary) {
			hasRunningProcess = true
			fmt.Printf("[FOUND] Process is running\n")
		}
		fmt.Println()
	}

	configPath, configDir := detectConfigPath()
	if configPath != "" {
		detected = true
		fmt.Printf("[FOUND] Configuration file: %s\n", configPath)
		if configDir != "" {
			fmt.Printf("[FOUND] Config directory: %s\n", configDir)
		}
		fmt.Println()
	}

	serviceName, serviceStatus := detectSystemdService()
	if serviceName != "" {
		detected = true
		fmt.Printf("[FOUND] Systemd service: %s (%s)\n", serviceName, serviceStatus)
		if serviceStatus == "active" {
			hasRunningProcess = true
		}
		fmt.Println()
	}

	if !detected {
		fmt.Println("[INFO] No existing sing-box installation detected.")
		fmt.Println("[INFO] Run 'sb-sync install' to install sing-box.")
		fmt.Println()
		return
	}

	if hasRunningProcess && (autoKill || autoAll) {
		fmt.Println("[ACTION] Stopping existing sing-box processes...")
		killSingBoxProcesses()
		fmt.Println()
	}

	if autoImport || autoAll {
		fmt.Println("[ACTION] Importing configuration...")

		if binary, _ := detectSingBoxBinary(); binary != "" {
			binDir := filepath.Dir(binary)
			if binDir != "." && binDir != "/" {
				config.AppConfig.InstallDir = binDir
			}
		}

		if configPath != "" {
			config.AppConfig.InstallDir = filepath.Dir(configPath)
		}

		config.Save()
		fmt.Printf("[CONFIG] Install directory set to: %s\n", config.AppConfig.InstallDir)
		fmt.Printf("[CONFIG] Config path set to: %s\n", configPath)
		fmt.Println()
	}

	if autoMigrate || autoAll {
		if serviceName != "" && serviceStatus == "active" {
			fmt.Println("[ACTION] Disabling systemd service...")
			exec.Command("systemctl", "stop", serviceName).Run()
			exec.Command("systemctl", "disable", serviceName).Run()
			fmt.Printf("[INFO] Stopped and disabled %s\n", serviceName)
			fmt.Println()
		}
	}

	fmt.Println("=== Auto-migration completed ===")
	fmt.Println()
	fmt.Println("Current configuration:")
	fmt.Printf("  Install Dir: %s\n", config.AppConfig.InstallDir)
	fmt.Printf("  Config Path: %s\n", config.GetConfigPath())
	fmt.Println()
	fmt.Println("Next steps:")
	if config.AppConfig.WebDAV.URL == "" {
		fmt.Println("  1. Configure WebDAV: sb-sync config set-dav --url <url> --user <user> --pass <pass> --path <path>")
	}
	fmt.Println("  1. Sync config: sb-sync sync")
	fmt.Println("  2. Start service: sb-sync service start")
	fmt.Println("  3. Check status: sb-sync service status")
}

func isProcessRunning(binaryPath string) bool {
	cmd := exec.Command("pgrep", "-f", filepath.Base(binaryPath))
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

func killSingBoxProcesses() {
	patterns := []string{
		"sing-box",
	}

	for _, pattern := range patterns {
		cmd := exec.Command("pkill", "-f", pattern)
		output, err := cmd.CombinedOutput()

		if err == nil {
			fmt.Printf("[INFO] Killed processes matching: %s\n", pattern)
		} else if len(output) > 0 {
			errStr := string(output)
			if !strings.Contains(errStr, "no process found") && !strings.Contains(errStr, "No matching processes") {
				fmt.Printf("[WARN] pkill output: %s\n", errStr)
			}
		}
	}

	pidFile := filepath.Join(config.GetInstallDir(), "sing-box.pid")
	if _, err := os.Stat(pidFile); err == nil {
		pidData, _ := os.ReadFile(pidFile)
		var pid int
		if _, err := fmt.Sscanf(string(pidData), "%d", &pid); err == nil {
			proc, err := os.FindProcess(pid)
			if err == nil {
				proc.Kill()
				proc.Wait()
			}
		}
		os.Remove(pidFile)
		fmt.Println("[INFO] Removed sb-sync PID file")
	}
}

func detectSingBoxBinary() (string, string) {
	commonPaths := []string{
		"/usr/bin/sing-box",
		"/usr/local/bin/sing-box",
		"/opt/sing-box/sing-box",
		"/root/sing-box",
		"/home/*/sing-box",
	}

	for _, path := range commonPaths {
		if strings.Contains(path, "*") {
			if matches, _ := filepath.Glob(path); len(matches) > 0 {
				path = matches[0]
			} else {
				continue
			}
		}

		if _, err := os.Stat(path); err == nil {
			if version := getSingBoxVersion(path); version != "" {
				return path, version
			}
		}
	}

	if path, err := exec.LookPath("sing-box"); err == nil {
		if version := getSingBoxVersion(path); version != "" {
			return path, version
		}
	}

	return "", ""
}

func getSingBoxVersion(binaryPath string) string {
	cmd := exec.Command(binaryPath, "version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "version") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				return strings.TrimSpace(parts[1])
			}
			return strings.TrimSpace(line)
		}
	}

	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}

	return ""
}

func detectConfigPath() (string, string) {
	commonPaths := []string{
		"/etc/sing-box/config.json",
		"/etc/sing-box/config.yaml",
		"/etc/sing-box/config.yml",
		"/var/lib/sing-box/config.json",
		"/root/sing-box/config.json",
		"/opt/sing-box/config.json",
		"/usr/local/etc/sing-box/config.json",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			dir := filepath.Dir(path)
			return path, dir
		}
	}

	if systemdPath := detectConfigFromSystemd(); systemdPath != "" {
		return systemdPath, filepath.Dir(systemdPath)
	}

	return "", ""
}

func detectConfigFromSystemd() string {
	serviceFilePaths := []string{
		"/lib/systemd/system/sing-box.service",
		"/etc/systemd/system/sing-box.service",
		"/usr/lib/systemd/system/sing-box.service",
	}

	for _, servicePath := range serviceFilePaths {
		if content, err := os.ReadFile(servicePath); err == nil {
			serviceContent := string(content)

			lines := strings.Split(serviceContent, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)

				if strings.HasPrefix(line, "ExecStart=") {
					parts := strings.Fields(line)
					for _, part := range parts {
						if strings.HasPrefix(part, "-c=") {
							configPath := strings.TrimPrefix(part, "-c=")
							if _, err := os.Stat(configPath); err == nil {
								return configPath
							}
						} else if strings.HasPrefix(part, "--config=") {
							configPath := strings.TrimPrefix(part, "--config=")
							if _, err := os.Stat(configPath); err == nil {
								return configPath
							}
						}
					}

					for _, arg := range parts[1:] {
						if strings.HasSuffix(arg, ".json") || strings.HasSuffix(arg, ".yaml") || strings.HasSuffix(arg, ".yml") {
							if _, err := os.Stat(arg); err == nil {
								return arg
							}
						}
					}
				}

				if strings.HasPrefix(line, "Environment=") {
					envLine := strings.TrimPrefix(line, "Environment=")
					if strings.Contains(envLine, "CONFIG=") {
						parts := strings.Split(envLine, "\"")
						for _, part := range parts {
							if strings.HasSuffix(part, ".json") || strings.HasSuffix(part, ".yaml") {
								if _, err := os.Stat(part); err == nil {
									return part
								}
							}
						}
					}
				}
			}
		}
	}

	return ""
}

func detectSystemdService() (string, string) {
	serviceFiles := []string{
		"/lib/systemd/system/sing-box.service",
		"/etc/systemd/system/sing-box.service",
		"/usr/lib/systemd/system/sing-box.service",
	}

	for _, path := range serviceFiles {
		if _, err := os.Stat(path); err == nil {
			cmd := exec.Command("systemctl", "is-active", "sing-box.service")
			output, _ := cmd.Output()
			status := strings.TrimSpace(string(output))
			if status == "unknown" {
				status = "inactive"
			}
			return "sing-box.service", status
		}
	}

	cmd := exec.Command("systemctl", "list-units", "--all", "--no-legend")
	output, _ := cmd.Output()

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "sing-box") && strings.Contains(line, ".service") {
			parts := strings.Fields(line)
			if len(parts) >= 4 {
				serviceName := strings.TrimSuffix(parts[0], ".service") + ".service"
				status := parts[2]
				return serviceName, status
			}
		}
	}

	return "", ""
}

func init() {
	autoCmd.Flags().BoolVar(&autoImport, "import", false, "Import detected configuration")
	autoCmd.Flags().BoolVar(&autoMigrate, "migrate", false, "Stop and disable systemd service")
	autoCmd.Flags().BoolVar(&autoKill, "kill", false, "Kill running sing-box processes")
	autoCmd.Flags().BoolVar(&autoAll, "all", false, "Run all auto-migration steps (kill, import, migrate)")
	autoCmd.MarkFlagsMutuallyExclusive("all", "import", "migrate", "kill")
	rootCmd.AddCommand(autoCmd)
}
