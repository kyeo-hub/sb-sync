package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	svc "sb-sync/pkg/service"
)

func newService() (svc.ServiceInterface, error) {
	return svc.NewService()
}

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage the sing-box service",
}

var installServiceCmd = &cobra.Command{
	Use:   "install",
	Short: "Install sing-box as a system service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := newService()
		if err != nil {
			fmt.Printf("[ERROR] Failed to create service: %v\n", err)
			return
		}
		err = s.Install()
		if err != nil {
			fmt.Printf("[ERROR] Failed to install service: %v\n", err)
			return
		}
		fmt.Println("[INFO] Service installed successfully.")
	},
}

var uninstallServiceCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := newService()
		if err != nil {
			fmt.Printf("[ERROR] Failed to create service: %v\n", err)
			return
		}
		err = s.Uninstall()
		if err != nil {
			fmt.Printf("[ERROR] Failed to uninstall service: %v\n", err)
			return
		}
		fmt.Println("[INFO] Service uninstalled successfully.")
	},
}

var startServiceCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := newService()
		if err != nil {
			fmt.Printf("[ERROR] Failed to create service: %v\n", err)
			return
		}
		err = s.Start()
		if err != nil {
			fmt.Printf("[ERROR] Failed to start service: %v\n", err)
			return
		}
		fmt.Println("[INFO] Service started successfully.")
	},
}

var stopServiceCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := newService()
		if err != nil {
			fmt.Printf("[ERROR] Failed to create service: %v\n", err)
			return
		}
		err = s.Stop()
		if err != nil {
			fmt.Printf("[ERROR] Failed to stop service: %v\n", err)
			return
		}
		fmt.Println("[INFO] Service stopped successfully.")
	},
}

var restartServiceCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := newService()
		if err != nil {
			fmt.Printf("[ERROR] Failed to create service: %v\n", err)
			return
		}
		fmt.Println("[INFO] Stopping service...")
		if err := s.Stop(); err != nil {
			fmt.Printf("[WARN] Failed to stop service: %v\n", err)
		}
		fmt.Println("[INFO] Starting service...")
		if err := s.Start(); err != nil {
			fmt.Printf("[ERROR] Failed to start service: %v\n", err)
			return
		}
		fmt.Println("[INFO] Service restarted successfully.")
	},
}

var statusServiceCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := newService()
		if err != nil {
			fmt.Printf("[ERROR] Failed to create service: %v\n", err)
			return
		}
		status, err := s.Status()
		if err != nil {
			fmt.Printf("[ERROR] Failed to get service status: %v\n", err)
			return
		}
		fmt.Printf("[INFO] Service status: %v\n", status)
	},
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
