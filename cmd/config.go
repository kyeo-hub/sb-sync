package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
}

var setDavCmd = &cobra.Command{
	Use:   "set-dav",
	Short: "Set WebDAV configuration",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		user, _ := cmd.Flags().GetString("user")
		pass, _ := cmd.Flags().GetString("pass")
		path, _ := cmd.Flags().GetString("path")

		if url != "" {
			config.AppConfig.WebDAV.URL = url
		}
		if user != "" {
			config.AppConfig.WebDAV.Username = user
		}
		if pass != "" {
			config.AppConfig.WebDAV.Password = pass
		}
		if path != "" {
			config.AppConfig.WebDAV.FilePath = path
		}

		if err := config.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}
		fmt.Println("WebDAV configuration updated.")
	},
}

var setProxyCmd = &cobra.Command{
	Use:   "set-proxy",
	Short: "Set GitHub proxy",
	Run: func(cmd *cobra.Command, args []string) {
		proxy, _ := cmd.Flags().GetString("url")
		config.AppConfig.GithubProxy = proxy
		if err := config.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}
		fmt.Println("GitHub proxy updated.")
	},
}

var setIntervalCmd = &cobra.Command{
	Use:   "set-interval",
	Short: "Set sync interval in minutes",
	Run: func(cmd *cobra.Command, args []string) {
		interval, _ := cmd.Flags().GetInt("minutes")
		if interval < 1 {
			fmt.Println("Interval must be at least 1 minute.")
			return
		}
		config.AppConfig.SyncInterval = interval
		if err := config.Save(); err != nil {
			fmt.Printf("Error saving config: %v\n", err)
			return
		}
		fmt.Printf("Sync interval set to %d minutes.\n", interval)
	},
}

var showConfigCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("WebDAV URL:      %s\n", config.AppConfig.WebDAV.URL)
		fmt.Printf("WebDAV User:     %s\n", config.AppConfig.WebDAV.Username)
		pass := "********"
		if config.AppConfig.WebDAV.Password == "" {
			pass = "(not set)"
		}
		fmt.Printf("WebDAV Pass:     %s\n", pass)
		fmt.Printf("WebDAV Path:     %s\n", config.AppConfig.WebDAV.FilePath)
		fmt.Printf("GitHub Proxy:    %s\n", config.AppConfig.GithubProxy)
		fmt.Printf("Install Dir:     %s\n", config.AppConfig.InstallDir)
		fmt.Printf("Sync Interval:   %d minutes\n", config.AppConfig.SyncInterval)
		fmt.Printf("Config Path:     %s\n", config.GetConfigPath())
	},
}

func init() {
	setDavCmd.Flags().String("url", "", "WebDAV server URL")
	setDavCmd.Flags().String("user", "", "WebDAV username")
	setDavCmd.Flags().String("pass", "", "WebDAV password")
	setDavCmd.Flags().String("path", "", "Remote path to config.json")

	setProxyCmd.Flags().String("url", "", "GitHub proxy URL (e.g. https://gh-proxy.com/)")

	setIntervalCmd.Flags().Int("minutes", 60, "Sync interval in minutes")

	configCmd.AddCommand(setDavCmd)
	configCmd.AddCommand(setProxyCmd)
	configCmd.AddCommand(setIntervalCmd)
	configCmd.AddCommand(showConfigCmd)
	rootCmd.AddCommand(configCmd)
}
