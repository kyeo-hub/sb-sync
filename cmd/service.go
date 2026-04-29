package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sb-sync/pkg/service"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "Manage the sing-box service",
}

var installServiceCmd = &cobra.Command{
	Use:   "install",
	Short: "Install sing-box as a system service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := service.NewService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		err = s.Install()
		if err != nil {
			fmt.Printf("Error installing service: %v\n", err)
			return
		}
		fmt.Println("Service installed successfully.")
	},
}

var uninstallServiceCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := service.NewService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		err = s.Uninstall()
		if err != nil {
			fmt.Printf("Error uninstalling service: %v\n", err)
			return
		}
		fmt.Println("Service uninstalled successfully.")
	},
}

var startServiceCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := service.NewService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		err = s.Start()
		if err != nil {
			fmt.Printf("Error starting service: %v\n", err)
			return
		}
		fmt.Println("Service started successfully.")
	},
}

var stopServiceCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := service.NewService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		err = s.Stop()
		if err != nil {
			fmt.Printf("Error stopping service: %v\n", err)
			return
		}
		fmt.Println("Service stopped successfully.")
	},
}

var statusServiceCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the sing-box service",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := service.NewService()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		status, err := s.Status()
		if err != nil {
			fmt.Printf("Error checking status: %v\n", err)
			return
		}
		fmt.Printf("Service status: %v\n", status)
	},
}

func init() {
	serviceCmd.AddCommand(installServiceCmd)
	serviceCmd.AddCommand(uninstallServiceCmd)
	serviceCmd.AddCommand(startServiceCmd)
	serviceCmd.AddCommand(stopServiceCmd)
	serviceCmd.AddCommand(statusServiceCmd)
	rootCmd.AddCommand(serviceCmd)
}
