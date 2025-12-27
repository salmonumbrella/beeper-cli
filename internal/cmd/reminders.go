package cmd

import (
	"fmt"
	"net/url"
	"time"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/beeper-cli/internal/api"
)

func newRemindersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reminders",
		Short: "Manage chat reminders",
	}

	cmd.AddCommand(newRemindersSetCmd())
	cmd.AddCommand(newRemindersClearCmd())

	return cmd
}

func newRemindersSetCmd() *cobra.Command {
	var (
		at   string
		chat string
	)

	cmd := &cobra.Command{
		Use:   "set [chat-id]",
		Short: "Set a reminder for a chat",
		Long: `Set a reminder for a chat.

You can specify the chat by ID or by using --chat with a name:
  beeper reminders set <chat-id> --at "2024-12-26 10:00"
  beeper reminders set --chat "Kishan" --at "2024-12-26 10:00"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if at == "" {
				return fmt.Errorf("--at is required")
			}

			reminderTime, err := time.Parse(time.RFC3339, at)
			if err != nil {
				// Try simpler format
				reminderTime, err = time.Parse("2006-01-02 15:04", at)
				if err != nil {
					return fmt.Errorf("invalid time format. Use: 2006-01-02 15:04 or RFC3339")
				}
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			// Resolve chat ID from --chat flag or positional arg
			var chatID string
			if chat != "" {
				resolved, err := resolveChatByName(cmd, client, chat)
				if err != nil {
					return err
				}
				chatID = resolved
			} else if len(args) == 1 {
				chatID = args[0]
			} else {
				return fmt.Errorf("either <chat-id> argument or --chat flag is required")
			}

			body := api.ReminderRequest{Reminder: api.NewReminderTime(reminderTime)}
			resp, err := client.Post(cmd.Context(), "/v1/chats/"+url.PathEscape(chatID)+"/reminders", body)
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, "Chat"); err != nil {
				return err
			}

			fmt.Printf("Reminder set for %s\n", reminderTime.Format("Jan 2, 2006 at 3:04 PM"))
			return nil
		},
	}

	cmd.Flags().StringVar(&at, "at", "", "Reminder time (required, e.g., '2024-12-25 10:00' or RFC3339)")
	cmd.Flags().StringVar(&chat, "chat", "", "Specify chat by name (searches for matching chat)")
	_ = cmd.MarkFlagRequired("at")

	return cmd
}

func newRemindersClearCmd() *cobra.Command {
	var chat string

	cmd := &cobra.Command{
		Use:   "clear [chat-id]",
		Short: "Clear a chat reminder",
		Long: `Clear a reminder from a chat.

You can specify the chat by ID or by using --chat with a name:
  beeper reminders clear <chat-id>
  beeper reminders clear --chat "Kishan"`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}

			// Resolve chat ID from --chat flag or positional arg
			var chatID string
			if chat != "" {
				resolved, err := resolveChatByName(cmd, client, chat)
				if err != nil {
					return err
				}
				chatID = resolved
			} else if len(args) == 1 {
				chatID = args[0]
			} else {
				return fmt.Errorf("either <chat-id> argument or --chat flag is required")
			}

			resp, err := client.Delete(cmd.Context(), "/v1/chats/"+url.PathEscape(chatID)+"/reminders")
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, "Chat"); err != nil {
				return err
			}

			fmt.Println("Reminder cleared")
			return nil
		},
	}

	cmd.Flags().StringVar(&chat, "chat", "", "Specify chat by name (searches for matching chat)")

	return cmd
}
