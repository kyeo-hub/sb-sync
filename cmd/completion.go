package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: `Generate shell completion scripts for sb-sync.

To load completions:

Bash:

  $ source <(sb-sync completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ sb-sync completion bash > /etc/bash_completion.d/sb-sync
  # macOS:
  $ sb-sync completion bash > /usr/local/etc/bash_completion.d/sb-sync

Zsh:

  # If shell completion is not enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ sb-sync completion zsh > "${fpath[1]}/_sb-sync"

  # You will need to start a new shell for this setup to take effect.

Fish:

  $ sb-sync completion fish | source

  # To load completions for each session, execute once:
  $ sb-sync completion fish > ~/.config/fish/completions/sb-sync.fish

PowerShell:

  PS> sb-sync completion powershell | Out-String | Invoke-Expression

  # To load completions for each session, execute once:
  PS> sb-sync completion powershell > sb-sync.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return rootCmd.GenBashCompletion(os.Stdout)
		case "zsh":
			return rootCmd.GenZshCompletion(os.Stdout)
		case "fish":
			return rootCmd.GenFishCompletion(os.Stdout, true)
		case "powershell":
			return rootCmd.GenPowerShellCompletionWithDesc(os.Stdout)
		}
		return fmt.Errorf("unsupported shell: %s", args[0])
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
