package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/salmonumbrella/beeper-cli/internal/api"
	"github.com/salmonumbrella/beeper-cli/internal/auth"
	"github.com/salmonumbrella/beeper-cli/internal/outfmt"
	"github.com/salmonumbrella/beeper-cli/internal/secrets"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}

	cmd.AddCommand(newAuthAddCmd())
	cmd.AddCommand(newAuthListCmd())
	cmd.AddCommand(newAuthRemoveCmd())
	cmd.AddCommand(newAuthLoginCmd())
	cmd.AddCommand(newAuthTestCmd())

	return cmd
}

func newAuthAddCmd() *cobra.Command {
	var tokenFlag string

	cmd := &cobra.Command{
		Use:   "add [name]",
		Short: "Add a new token",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "default"
			if len(args) > 0 {
				name = args[0]
			}

			token := tokenFlag
			if token == "" {
				fmt.Print("Enter token: ")
				tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
				if err != nil {
					return fmt.Errorf("failed to read token: %w", err)
				}
				fmt.Println()
				token = strings.TrimSpace(string(tokenBytes))
			}

			if token == "" {
				return fmt.Errorf("token cannot be empty")
			}

			// Validate token
			client := api.NewClient(api.DefaultBaseURL, token)
			resp, err := client.Get(cmd.Context(), "/v1/accounts")
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, ""); err != nil {
				return err
			}

			// Save to keyring
			store, err := secrets.NewStore()
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}

			if err := store.Set(name, secrets.Credentials{
				Token:     token,
				CreatedAt: time.Now(),
			}); err != nil {
				return fmt.Errorf("failed to save token: %w", err)
			}

			fmt.Printf("Token saved as '%s'\n", name)
			return nil
		},
	}

	cmd.Flags().StringVar(&tokenFlag, "token", "", "Token (or enter interactively)")

	return cmd
}

func newAuthListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List configured tokens",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := secrets.NewStore()
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}

			accounts, err := store.List()
			if err != nil {
				return fmt.Errorf("failed to list accounts: %w", err)
			}

			if len(accounts) == 0 {
				fmt.Println("No tokens configured. Run: beeper auth add")
				return nil
			}

			return outfmt.Output(cmd.Context(), accounts, func(w io.Writer) {
				tw := outfmt.NewTableWriter(w)
				tw.SetHeader([]string{"Name", "Created"})
				for _, a := range accounts {
					tw.Append([]string{a.Name, a.CreatedAt.Format("2006-01-02 15:04")})
				}
				tw.Render()
			})
		},
	}
}

func newAuthRemoveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a token",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			store, err := secrets.NewStore()
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}

			if err := store.Delete(name); err != nil {
				return fmt.Errorf("failed to remove token: %w", err)
			}

			fmt.Printf("Token '%s' removed\n", name)
			return nil
		},
	}
}

func newAuthTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test [name]",
		Short: "Test a token",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "default"
			if len(args) > 0 {
				name = args[0]
			}

			store, err := secrets.NewStore()
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}

			creds, err := store.Get(name)
			if err != nil {
				return fmt.Errorf("token '%s' not found. Run: beeper auth add", name)
			}

			client := api.NewClient(api.DefaultBaseURL, creds.Token)
			resp, err := client.Get(cmd.Context(), "/v1/accounts")
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, ""); err != nil {
				return err
			}

			fmt.Printf("Token '%s' is valid\n", name)
			return nil
		},
	}
}

func openSecretsStore() (secrets.Store, error) {
	return secrets.NewStore()
}

func newAuthLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login via browser (interactive token setup)",
		RunE: func(cmd *cobra.Command, args []string) error {
			server, err := auth.NewSetupServer()
			if err != nil {
				return fmt.Errorf("failed to create setup server: %w", err)
			}

			ctx := context.Background()
			result, err := server.Start(ctx)
			if err != nil {
				return fmt.Errorf("setup failed: %w", err)
			}

			if result.Error != nil {
				return result.Error
			}

			fmt.Printf("\nAuthentication successful! Token saved as '%s'\n", result.Name)
			fmt.Println("You can now use the Beeper CLI. Try: beeper chats list")
			return nil
		},
	}
}
