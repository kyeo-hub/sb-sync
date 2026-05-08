package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sb-sync/pkg/sync"
)

var syncDryRun bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync config from WebDAV",
	Run: func(cmd *cobra.Command, args []string) {
		if syncDryRun {
			fmt.Println("[DRY-RUN] Would perform the following actions:")
			fmt.Println("[DRY-RUN] 1. Connect to WebDAV server")
			fmt.Println("[DRY-RUN] 2. Download config from remote path")
			fmt.Println("[DRY-RUN] 3. Save to local config file")
			return
		}

		fmt.Println("[INFO] Syncing from WebDAV...")
		err := sync.SyncFromWebDAV()
		if err != nil {
			fmt.Printf("[ERROR] %v\n", err)
			return
		}
		fmt.Println("[INFO] Successfully synced configuration.")
	},
}

func init() {
	syncCmd.Flags().BoolVar(&syncDryRun, "dry-run", false, "Show what would be done without actually doing it")
	rootCmd.AddCommand(syncCmd)
}
