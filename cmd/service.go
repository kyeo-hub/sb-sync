package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage the sing-box service",
}

var installServiceCmd = &cobra.Command{
	Use:   "install",
	Short: "Install sing-box as a system service",
	Run: func(cmd *cobra.Command, args []string) {
		installManualService()
	},
}

var uninstallServiceCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		stopService()
		removeManualService()
	},
}

var startServiceCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		startService()
	},
}

var stopServiceCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		stopService()
	},
}

var restartServiceCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		stopService()
		startService()
	},
}

var statusServiceCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		checkStatus()
	},
}

func getPIDFile() string {
	return filepath.Join(config.GetInstallDir(), "sing-box.pid")
}

func getLogFile() string {
	return filepath.Join(config.GetInstallDir(), "sing-box.log")
}

func getLockFile() string {
	return filepath.Join(config.GetInstallDir(), "sing-box.lock")
}

func installManualService() {
	fmt.Println("[INFO] Systemd not available, using manual service management.")
	
	pidFile := getPIDFile()
	if _, err := os.Stat(pidFile); err == nil {
		fmt.Println("[WARN] Service may already be running. Use 'sb-sync service stop' first.")
	}
	
	scriptPath := filepath.Join(config.GetInstallDir(), "sing-box.sh")
	script := fmt.Sprintf(`#!/bin/bash
exec %s run -c %s >> %s 2>&1
`, config.GetSingBoxBinary(), config.GetConfigPath(), getLogFile())
	
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		fmt.Printf("[ERROR] Failed to create service script: %v\n", err)
		return
	}
	
	fmt.Println("[INFO] Manual service script created.")
	fmt.Printf("[INFO] To start: sb-sync service start\n")
	fmt.Printf("[INFO] Logs: %s\n", getLogFile())
}

func removeManualService() {
	scriptPath := filepath.Join(config.GetInstallDir(), "sing-box.sh")
	os.Remove(scriptPath)
	os.Remove(getPIDFile())
	os.Remove(getLogFile())
	os.Remove(getLockFile())
	fmt.Println("[INFO] Manual service files removed.")
}

func startService() {
	pidFile := getPIDFile()
	
	if _, err := os.Stat(pidFile); err == nil {
		if pid, err := os.ReadFile(pidFile); err == nil {
			if proc, err := os.FindProcess(int(pid[0])); err == nil && proc != nil {
				if err := proc.Signal(os.Signal(nil)); err == nil {
					fmt.Println("[INFO] Service is already running.")
					return
				}
			}
		}
	}
	
	binPath := config.GetSingBoxBinary()
	configPath := config.GetConfigPath()
	logFile := getLogFile()
	
	fmt.Printf("[INFO] Starting sing-box...\n")
	fmt.Printf("[INFO] Binary: %s\n", binPath)
	fmt.Printf("[INFO] Config: %s\n", configPath)
	fmt.Printf("[INFO] Log: %s\n", logFile)
	
	outFile, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("[ERROR] Failed to open log file: %v\n", err)
		return
	}
	defer outFile.Close()
	
	cmd := exec.Command(binPath, "run", "-c", configPath)
	cmd.Stdout = outFile
	cmd.Stderr = outFile
	
	if err := cmd.Start(); err != nil {
		fmt.Printf("[ERROR] Failed to start sing-box: %v\n", err)
		return
	}
	
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", cmd.Process.Pid)), 0644); err != nil {
		fmt.Printf("[WARN] Failed to write PID file: %v\n", err)
	}
	
	fmt.Printf("[INFO] Service started successfully. PID: %d\n", cmd.Process.Pid)
}

func stopService() {
	pidFile := getPIDFile()
	
	if _, err := os.Stat(pidFile); err != nil {
		fmt.Println("[INFO] Service is not running (no PID file).")
		return
	}
	
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Printf("[ERROR] Failed to read PID file: %v\n", err)
		return
	}
	
	var pid int
	if _, err := fmt.Sscanf(string(pidData), "%d", &pid); err != nil {
		fmt.Printf("[ERROR] Invalid PID file: %v\n", err)
		os.Remove(pidFile)
		return
	}
	
	proc, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("[INFO] Process not found: %d\n", pid)
		os.Remove(pidFile)
		return
	}
	
	if err := proc.Kill(); err != nil {
		fmt.Printf("[WARN] Failed to kill process: %v\n", err)
	}
	
	proc.Wait()
	os.Remove(pidFile)
	fmt.Println("[INFO] Service stopped.")
}

func checkStatus() {
	pidFile := getPIDFile()
	
	if _, err := os.Stat(pidFile); err != nil {
		fmt.Println("[INFO] Service status: stopped")
		return
	}
	
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		fmt.Println("[INFO] Service status: unknown (no PID)")
		return
	}
	
	var pid int
	if _, err := fmt.Sscan(string(pidData), &pid); err != nil {
		fmt.Println("[INFO] Service status: unknown (invalid PID)")
		return
	}
	
	if proc, err := os.FindProcess(pid); err == nil {
		if err := proc.Signal(os.Signal(nil)); err == nil {
			fmt.Printf("[INFO] Service status: running (PID: %d)\n", pid)
			return
		}
	}
	
	fmt.Println("[INFO] Service status: stopped (stale PID file)")
	os.Remove(pidFile)
}

func init() {
	serviceCmd.AddCommand(installServiceCmd)
	serviceCmd.AddCommand(uninstallServiceCmd)
	serviceCmd.AddCommand(startServiceCmd)
	serviceCmd.AddCommand(stopServiceCmd)
	serviceCmd.AddCommand(restartServiceCmd)
	serviceCmd.AddCommand(statusServiceCmd)
	rootCmd.AddCommand(serviceCmd)
}
