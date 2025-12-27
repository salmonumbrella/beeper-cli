package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for beeper.

To load completions:

Bash:
  $ source <(beeper completion bash)
  # To load completions for each session, execute once:
  # Linux:
  $ beeper completion bash > /etc/bash_completion.d/beeper
  # macOS:
  $ beeper completion bash > $(brew --prefix)/etc/bash_completion.d/beeper

Zsh:
  # If shell completion is not already enabled in your environment,
  # you will need to enable it. You can execute the following once:
  $ echo "autoload -U compinit; compinit" >> ~/.zshrc
  # To load completions for each session, execute once:
  $ beeper completion zsh > "${fpath[1]}/_beeper"
  # You will need to start a new shell for this setup to take effect.

Fish:
  $ beeper completion fish | source
  # To load completions for each session, execute once:
  $ beeper completion fish > ~/.config/fish/completions/beeper.fish

PowerShell:
  PS> beeper completion powershell | Out-String | Invoke-Expression
  # To load completions for every new session, run:
  PS> beeper completion powershell > beeper.ps1
  # and source this file from your PowerShell profile.
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
	return cmd
}
