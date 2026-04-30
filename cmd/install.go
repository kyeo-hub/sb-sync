package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
	"sb-sync/pkg/downloader"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install sing-box core",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking for latest sing-box version...")
		version, err := downloader.GetLatestVersion()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Printf("Latest version: %s\n", version)

		fmt.Println("Downloading...")
		tmpFile, err := downloader.DownloadSingBox(version)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("Extracting...")
		err = downloader.ExtractBinary(tmpFile, config.AppConfig.InstallDir)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("Successfully installed sing-box to %s\n", config.AppConfig.InstallDir)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
