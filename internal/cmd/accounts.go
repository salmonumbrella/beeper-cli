package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/beeper-cli/internal/api"
	"github.com/salmonumbrella/beeper-cli/internal/outfmt"
)

func newAccountsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "accounts",
		Short: "List connected accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}

			resp, err := client.Get(cmd.Context(), "/v1/accounts")
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, ""); err != nil {
				return err
			}

			var accounts []api.Account
			if err := json.NewDecoder(resp.Body).Decode(&accounts); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			return outfmt.Output(cmd.Context(), accounts, func(w io.Writer) {
				tw := outfmt.NewTableWriter(w)
				tw.SetHeader([]string{"ID", "Network", "User"})
				for _, a := range accounts {
					name := a.ProfileName
					if name == "" {
						name = a.ProfileUsername
					}
					tw.Append([]string{a.ID, a.NetworkName, name})
				}
				tw.Render()
			})
		},
	}
}

func getClient() (*api.Client, error) {
	store, err := openSecretsStore()
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	accounts, err := store.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list accounts: %w", err)
	}

	if len(accounts) == 0 {
		return nil, fmt.Errorf("no tokens configured. Run: beeper auth add")
	}

	// Use first account (or could check flags.Account)
	creds, err := store.Get(accounts[0].Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}

	opts := []api.ClientOption{}
	if flags.Debug {
		opts = append(opts, api.WithDebug(true))
	}

	return api.NewClient(api.DefaultBaseURL, creds.Token, opts...), nil
}
