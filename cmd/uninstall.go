package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"sb-sync/pkg/config"
)

var (
	uninstallKeepConfig bool
	uninstallForce      bool
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall sb-sync and optionally sing-box",
	Run: func(cmd *cobra.Command, args []string) {
		runUninstall()
	},
}

func init() {
	uninstallCmd.Flags().BoolVar(&uninstallKeepConfig, "keep-config", false, "Keep configuration files")
	uninstallCmd.Flags().BoolVar(&uninstallForce, "force", false, "Force uninstall without confirmation")
	rootCmd.AddCommand(uninstallCmd)
}

func runUninstall() {
	fmt.Println("=== sb-sync Uninstall ===")
	fmt.Println()

	installDir := config.GetInstallDir()
	sbSyncBinary := "/usr/local/bin/sb-sync"

	if !uninstallForce {
		fmt.Println("警告: 此操作将删除以下内容:")
		fmt.Printf("  - sb-sync 二进制文件: %s\n", sbSyncBinary)
		fmt.Printf("  - 安装目录: %s\n", installDir)
		if !uninstallKeepConfig {
			fmt.Println("  - 所有配置文件")
		}
		fmt.Println()
		fmt.Print("是否继续? (y/N): ")

		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("取消卸载.")
			os.Exit(0)
		}
	}

	fmt.Println("[INFO] Stopping service...")
	stopService()

	fmt.Println("[INFO] Removing service files...")
	removeManualService()

	if !uninstallKeepConfig {
		fmt.Println("[INFO] Removing installation directory...")
		if err := os.RemoveAll(installDir); err != nil {
			fmt.Printf("[WARN] Failed to remove install dir: %v\n", err)
		}
	} else {
		fmt.Println("[INFO] Keeping installation directory (--keep-config)")
	}

	fmt.Println("[INFO] Removing sb-sync binary...")
	if err := os.Remove(sbSyncBinary); err != nil {
		fmt.Printf("[WARN] Failed to remove binary: %v\n", err)
	}

	fmt.Println()
	fmt.Println("[INFO] Checking for shell completion files...")
	completionPaths := []string{
		"/etc/bash_completion.d/sb-sync",
		"/usr/local/share/zsh/site-functions/_sb-sync",
		"/usr/share/zsh/site-functions/_sb-sync",
	}

	for _, path := range completionPaths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("[INFO] Removing completion file: %s\n", path)
			os.Remove(path)
		}
	}

	fmt.Println()
	fmt.Println("[INFO] sb-sync uninstall completed!")
	fmt.Println()
	fmt.Println("注意: 如果之前手动安装了 sing-box，请使用以下命令完全删除:")
	fmt.Println("  rm -f /usr/bin/sing-box /usr/local/bin/sing-box")
	fmt.Println("  rm -f /etc/sing-box/config.json")
}
