package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/beeper-cli/internal/api"
)

func newFocusCmd() *cobra.Command {
	var (
		chat       string
		messageID  string
		draftText  string
		attachment string
	)

	cmd := &cobra.Command{
		Use:   "focus",
		Short: "Focus Beeper Desktop window",
		Long: `Focus Beeper Desktop window, optionally navigating to a specific chat.

You can specify the chat by ID or by name:
  beeper focus --chat "Kishan" --draft "Hello!"
  beeper focus --chat "!abc123:beeper.com"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}

			// Resolve chat name to ID if needed
			chatID := chat
			if chat != "" && !looksLikeChatID(chat) {
				resolved, err := resolveChatByName(cmd, client, chat)
				if err != nil {
					return err
				}
				chatID = resolved
			}

			// If we have both chatID and draftText, use two-step approach
			// to avoid draft being applied to wrong chat
			if chatID != "" && draftText != "" {
				// Step 1: Focus on chat first (no draft)
				focusBody := api.FocusRequest{
					ChatID:    chatID,
					MessageID: messageID,
				}
				resp, err := client.Post(cmd.Context(), "/v1/focus", focusBody)
				if err != nil {
					return api.UserFriendlyError(err)
				}
				_ = resp.Body.Close()
				if err := api.ParseErrorWithContext(resp, ""); err != nil {
					return err
				}

				// Brief pause to let Beeper switch chats
				time.Sleep(100 * time.Millisecond)

				// Step 2: Now set the draft
				draftBody := api.FocusRequest{
					ChatID:              chatID,
					DraftText:           draftText,
					DraftAttachmentPath: attachment,
				}
				resp, err = client.Post(cmd.Context(), "/v1/focus", draftBody)
				if err != nil {
					return api.UserFriendlyError(err)
				}
				defer func() { _ = resp.Body.Close() }()
				if err := api.ParseErrorWithContext(resp, ""); err != nil {
					return err
				}

				fmt.Println("Focused Beeper on chat with draft")
				return nil
			}

			// Single request if no draft or no chatID
			body := api.FocusRequest{
				ChatID:              chatID,
				MessageID:           messageID,
				DraftText:           draftText,
				DraftAttachmentPath: attachment,
			}

			resp, err := client.Post(cmd.Context(), "/v1/focus", body)
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, ""); err != nil {
				return err
			}

			if chatID != "" {
				fmt.Println("Focused Beeper on chat")
			} else {
				fmt.Println("Focused Beeper")
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&chat, "chat", "", "Navigate to specific chat (by name or ID)")
	cmd.Flags().StringVar(&messageID, "message", "", "Navigate to specific message")
	cmd.Flags().StringVar(&draftText, "draft", "", "Pre-fill draft text")
	cmd.Flags().StringVar(&attachment, "attachment", "", "Pre-fill draft attachment path")

	return cmd
}
