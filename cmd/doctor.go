package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
	"sb-sync/pkg/sync"
)

var (
	doctorJSON bool
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check the health and configuration of sb-sync",
	Run: func(cmd *cobra.Command, args []string) {
		status := sync.HealthCheck()

		if doctorJSON {
			data, err := json.MarshalIndent(status, "", "  ")
			if err != nil {
				fmt.Printf("[ERROR] Failed to marshal status: %v\n", err)
				os.Exit(1)
			}
			fmt.Println(string(data))
			return
		}

		fmt.Println("=== sb-sync Health Check ===")
		fmt.Println()

		singboxStatus := "❌ Not installed"
		if status.SingBoxInstalled {
			singboxStatus = "✅ Installed"
		}
		fmt.Printf("sing-box binary: %s\n", singboxStatus)

		configStatus := "❌ Not found"
		if status.ConfigExists {
			configStatus = "✅ Found"
		}
		fmt.Printf("Config file:      %s (%s)\n", configStatus, config.GetConfigPath())

		webdavConfigured := "❌ Not configured"
		if status.WebDAVConfigured {
			webdavConfigured = "✅ Configured"
		}
		fmt.Printf("WebDAV:           %s\n", webdavConfigured)

		webdavReachable := "❌ Not reachable"
		if status.WebDAVReachable {
			webdavReachable = "✅ Reachable"
		}
		fmt.Printf("WebDAV server:    %s\n", webdavReachable)

		fmt.Println()

		if !status.WebDAVConfigured {
			fmt.Println("[WARN] WebDAV is not configured. Run 'sb-sync config set-dav' to set up synchronization.")
		}

		if status.WebDAVConfigured && !status.WebDAVReachable {
			fmt.Println("[WARN] WebDAV server is not reachable. Check your network connection and credentials.")
		}

		allHealthy := status.SingBoxInstalled &&
			status.ConfigExists &&
			status.WebDAVConfigured &&
			status.WebDAVReachable

		if allHealthy {
			fmt.Println("[INFO] All checks passed! sb-sync is healthy.")
			os.Exit(0)
		} else {
			fmt.Println("[WARN] Some checks failed. Please review the configuration.")
			os.Exit(1)
		}
	},
}

func init() {
	doctorCmd.Flags().BoolVar(&doctorJSON, "json", false, "Output in JSON format")
	rootCmd.AddCommand(doctorCmd)
}
