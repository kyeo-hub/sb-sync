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

func init() {
	setDavCmd.Flags().String("url", "", "WebDAV server URL")
	setDavCmd.Flags().String("user", "", "WebDAV username")
	setDavCmd.Flags().String("pass", "", "WebDAV password")
	setDavCmd.Flags().String("path", "", "Remote path to config.json")

	setProxyCmd.Flags().String("url", "", "GitHub proxy URL (e.g. https://ghproxy.com/)")

	configCmd.AddCommand(setDavCmd)
	configCmd.AddCommand(setProxyCmd)
	rootCmd.AddCommand(configCmd)
}
