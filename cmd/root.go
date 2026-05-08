package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
)

var (
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "sb-sync",
	Short:   "A cross-platform sing-box installer and config synchronizer",
	Version: Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := config.Init(); err != nil {
			fmt.Printf("[ERROR] Failed to initialize config: %v\n", err)
			os.Exit(1)
		}
		config.LoadEnvOverrides()
		if err := config.ValidateConfig(); err != nil {
			fmt.Printf("[WARN] Config validation warning: %v\n", err)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}
