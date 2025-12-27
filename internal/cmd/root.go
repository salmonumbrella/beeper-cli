package cmd

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/beeper-cli/internal/outfmt"
)

type rootFlags struct {
	Account string
	Output  string
	Query   string
	Color   string
	Debug   bool
}

var flags rootFlags

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "beeper",
		Short:        "CLI for Beeper Desktop",
		Long:         "A command-line interface for Beeper Desktop's local API.",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ctx = outfmt.WithFormat(ctx, flags.Output)
			ctx = outfmt.WithQuery(ctx, flags.Query)
			ctx = outfmt.WithColor(ctx, flags.Color)
			cmd.SetContext(ctx)
			return nil
		},
	}

	cmd.PersistentFlags().StringVarP(&flags.Account, "account", "a", "", "Filter by account ID(s), comma-separated")
	cmd.PersistentFlags().StringVarP(&flags.Output, "output", "o", "text", "Output format: text|json")
	cmd.PersistentFlags().StringVar(&flags.Query, "query", "", "JQ filter for JSON output")
	cmd.PersistentFlags().StringVar(&flags.Color, "color", "auto", "Color mode: auto|always|never")
	cmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "Enable debug logging for API requests")

	cmd.AddCommand(newAuthCmd())
	cmd.AddCommand(newAccountsCmd())
	cmd.AddCommand(newChatsCmd())
	cmd.AddCommand(newMessagesCmd())
	cmd.AddCommand(newRemindersCmd())
	cmd.AddCommand(newFocusCmd())
	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newCompletionCmd())

	return cmd
}

func Execute(args []string) error {
	cmd := NewRootCmd()
	cmd.SetArgs(args)
	return cmd.ExecuteContext(context.Background())
}
