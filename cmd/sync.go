package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"sb-sync/pkg/sync"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync config from WebDAV",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Syncing from WebDAV...")
		err := sync.SyncFromWebDAV()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		fmt.Println("Successfully synced configuration.")
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
