package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
	"sb-sync/pkg/downloader"
)

var (
	dryRun bool
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install sing-box core",
	Run: func(cmd *cobra.Command, args []string) {
		if dryRun {
			fmt.Println("[DRY-RUN] Would perform the following actions:")
			fmt.Println("[DRY-RUN] 1. Check for latest sing-box version from GitHub")
			fmt.Println("[DRY-RUN] 2. Download the latest release archive")
			fmt.Println("[DRY-RUN] 3. Extract binary to:", config.AppConfig.InstallDir)
			return
		}

		fmt.Println("[INFO] Checking for latest sing-box version...")
		version, err := downloader.GetLatestVersion()
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			return
		}
		fmt.Printf("[INFO] Latest version: %s\n", version)

		fmt.Println("[INFO] Downloading...")
		tmpFile, err := downloader.DownloadSingBox(version)
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			return
		}

		fmt.Println("[INFO] Extracting...")
		err = downloader.ExtractBinary(tmpFile, config.AppConfig.InstallDir)
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			return
		}

		fmt.Printf("[INFO] Successfully installed sing-box to %s\n", config.AppConfig.InstallDir)
	},
}

func init() {
	installCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show what would be done without actually doing it")
	rootCmd.AddCommand(installCmd)
}
