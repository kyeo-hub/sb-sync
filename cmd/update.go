package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
	"sb-sync/pkg/downloader"
)

var updateDryRun bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check and update sing-box to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		if updateDryRun {
			fmt.Println("[DRY-RUN] Would perform the following actions:")
			fmt.Println("[DRY-RUN] 1. Check for latest sing-box version from GitHub")
			fmt.Println("[DRY-RUN] 2. Compare with current installed version")
			fmt.Println("[DRY-RUN] 3. Download and extract if newer")
			return
		}

		fmt.Println("[INFO] Checking for latest sing-box version...")
		latest, err := downloader.GetLatestVersion()
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			return
		}

		current, err := getCurrentVersion()
		if err != nil {
			fmt.Printf("[WARN] Could not determine current version: %v\n", err)
		} else {
			fmt.Printf("[INFO] Current version: %s\n", current)
			fmt.Printf("[INFO] Latest version:  %s\n", latest)

			if current == latest {
				fmt.Println("[INFO] sing-box is already up to date.")
				return
			}
		}

		fmt.Printf("[INFO] Updating to %s...\n", latest)
		tmpFile, err := downloader.DownloadSingBox(latest)
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			return
		}

		err = downloader.ExtractBinary(tmpFile, config.AppConfig.InstallDir)
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			return
		}

		fmt.Println("[INFO] Successfully updated sing-box.")
	},
}

func init() {
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "Show what would be done without actually doing it")
	rootCmd.AddCommand(updateCmd)
}

func getCurrentVersion() (string, error) {
	binPath := config.GetSingBoxBinary()

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return "", fmt.Errorf("sing-box not installed")
	}

	out, err := exec.Command(binPath, "version").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get version: %w", err)
	}

	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 3 {
			return "v" + parts[2], nil
		}
	}

	return "", fmt.Errorf("failed to parse version output")
}

func GetCurrentVersion() string {
	version, err := getCurrentVersion()
	if err != nil {
		return "unknown"
	}
	return version
}
