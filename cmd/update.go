package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
	"sb-sync/pkg/downloader"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check and update sing-box to latest version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking for latest sing-box version...")
		latest, err := downloader.GetLatestVersion()
		if err != nil {
			fmt.Printf("Error getting latest version: %v\n", err)
			return
		}

		current, err := getCurrentVersion()
		if err != nil {
			fmt.Printf("Could not determine current version: %v\n", err)
			// Continue to install anyway
		} else {
			fmt.Printf("Current version: %s\n", current)
			fmt.Printf("Latest version:  %s\n", latest)

			if current == latest {
				fmt.Println("sing-box is already up to date.")
				return
			}
		}

		fmt.Printf("Updating to %s...\n", latest)
		tmpFile, err := downloader.DownloadSingBox(latest)
		if err != nil {
			fmt.Printf("Error downloading: %v\n", err)
			return
		}

		err = downloader.ExtractBinary(tmpFile, config.AppConfig.InstallDir)
		if err != nil {
			fmt.Printf("Error extracting: %v\n", err)
			return
		}

		fmt.Println("Successfully updated sing-box.")
	},
}

func getCurrentVersion() (string, error) {
	binPath := filepath.Join(config.AppConfig.InstallDir, "sing-box")
	if filepath.Separator == '\\' {
		binPath += ".exe"
	}

	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		return "", fmt.Errorf("sing-box not installed")
	}

	out, err := exec.Command(binPath, "version").Output()
	if err != nil {
		return "", err
	}

	// Output format is usually "sing-box version 1.2.3"
	lines := strings.Split(string(out), "\n")
	if len(lines) > 0 {
		parts := strings.Fields(lines[0])
		if len(parts) >= 3 {
			return "v" + parts[2], nil
		}
	}

	return "", fmt.Errorf("failed to parse version output")
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
