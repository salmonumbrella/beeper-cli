package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/beeper-cli/internal/api"
	"github.com/salmonumbrella/beeper-cli/internal/outfmt"
)

func newMessagesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "messages",
		Short: "Manage messages",
	}

	cmd.AddCommand(newMessagesListCmd())
	cmd.AddCommand(newMessagesSearchCmd())
	cmd.AddCommand(newMessagesSendCmd())

	return cmd
}

func newMessagesListCmd() *cobra.Command {
	var (
		cursor    string
		direction string
		limit     int
		chat      string
	)

	cmd := &cobra.Command{
		Use:   "list [chat-id]",
		Short: "List messages in a chat",
		Long: `List messages in a chat.

You can specify the chat either by ID or by using --chat with a name:
  beeper messages list <chat-id>
  beeper messages list --chat "Kishan"`,
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

			params := url.Values{}
			if cursor != "" {
				params.Set("cursor", cursor)
			}
			if direction != "" {
				params.Set("direction", direction)
			}

			// Fetch chat info first to get the other person's name
			chatResp, err := client.Get(cmd.Context(), "/v1/chats/"+url.PathEscape(chatID))
			if err != nil {
				return api.UserFriendlyError(err)
			}
			var chat api.Chat
			if chatResp.StatusCode == 200 {
				_ = json.NewDecoder(chatResp.Body).Decode(&chat)
			}
			_ = chatResp.Body.Close()

			path := "/v1/chats/" + url.PathEscape(chatID) + "/messages"
			if len(params) > 0 {
				path += "?" + params.Encode()
			}

			resp, err := client.Get(cmd.Context(), path)
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, "Chat"); err != nil {
				return err
			}

			var result api.ListMessagesResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			messages := result.Items
			if limit > 0 && len(messages) > limit {
				messages = messages[:limit]
			}

			// Use chat title as the other person's name for DMs
			otherName := chat.Title
			if otherName == "" {
				otherName = "Them"
			}

			return outfmt.Output(cmd.Context(), result, func(w io.Writer) {
				for _, m := range messages {
					sender := senderName(m.SenderID)
					if sender == "Them" {
						sender = otherName
					}
					_, _ = fmt.Fprintf(w, "[%s] %s: %s\n", formatTime(m.Timestamp), sender, m.Text)
				}
				if result.HasMore {
					_, _ = fmt.Fprintf(w, "\n(more messages available, use --cursor %s)\n", result.Cursor)
				}
			})
		},
	}

	cmd.Flags().StringVar(&cursor, "cursor", "", "Pagination cursor")
	cmd.Flags().StringVar(&direction, "direction", "", "Direction: before or after")
	cmd.Flags().IntVar(&limit, "limit", 0, "Limit number of messages")
	cmd.Flags().StringVar(&chat, "chat", "", "Specify chat by name (searches for matching chat)")

	return cmd
}

func newMessagesSearchCmd() *cobra.Command {
	var (
		chatIDs   string
		dateAfter string
		limit     int
	)

	cmd := &cobra.Command{
		Use:   "search <query>",
		Short: "Search messages",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]

			client, err := getClient()
			if err != nil {
				return err
			}

			params := url.Values{}
			params.Set("query", query)
			if flags.Account != "" {
				for _, id := range strings.Split(flags.Account, ",") {
					params.Add("accountIDs", strings.TrimSpace(id))
				}
			}
			if chatIDs != "" {
				for _, id := range strings.Split(chatIDs, ",") {
					params.Add("chatIDs", strings.TrimSpace(id))
				}
			}
			if dateAfter != "" {
				params.Set("dateAfter", dateAfter)
			}

			resp, err := client.Get(cmd.Context(), "/v1/messages/search?"+params.Encode())
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, ""); err != nil {
				return err
			}

			var result api.SearchMessagesResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			messages := result.Messages
			if limit > 0 && len(messages) > limit {
				messages = messages[:limit]
			}

			return outfmt.Output(cmd.Context(), result, func(w io.Writer) {
				tw := outfmt.NewTableWriter(w)
				tw.SetHeader([]string{"Chat", "Sender", "Time", "Message"})
				for _, m := range messages {
					chatName := m.ChatID
					if chat, ok := result.Chats[m.ChatID]; ok {
						chatName = chat.Title
					}
					tw.Append([]string{
						truncate(chatName, 20),
						truncate(m.Sender, 15),
						formatTime(m.Timestamp),
						truncate(m.Text, 40),
					})
				}
				tw.Render()
				if result.HasMore {
					_, _ = fmt.Fprintf(w, "\n(%d+ results, showing first page)\n", len(messages))
				}
			})
		},
	}

	cmd.Flags().StringVar(&chatIDs, "chat", "", "Filter by chat ID(s), comma-separated")
	cmd.Flags().StringVar(&dateAfter, "after", "", "Messages after date (ISO format)")
	cmd.Flags().IntVar(&limit, "limit", 0, "Limit number of results")

	return cmd
}

func newMessagesSendCmd() *cobra.Command {
	var (
		text    string
		replyTo string
		to      string
	)

	cmd := &cobra.Command{
		Use:   "send [chat-id]",
		Short: "Send a message",
		Long: `Send a message to a chat.

You can specify the chat either by ID or by using --to with a name:
  beeper messages send <chat-id> --text "Hello"
  beeper messages send --to "Kishan" --text "Hello"

Using --to searches for a chat by name and uses the first match.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if text == "" {
				return fmt.Errorf("--text is required")
			}

			client, err := getClient()
			if err != nil {
				return err
			}

			// Resolve chat ID from --to flag or positional arg
			var chatID string
			if to != "" {
				resolved, err := resolveChatByName(cmd, client, to)
				if err != nil {
					return err
				}
				chatID = resolved
			} else if len(args) == 1 {
				chatID = args[0]
			} else {
				return fmt.Errorf("either <chat-id> argument or --to flag is required")
			}

			body := api.SendMessageRequest{
				Text:             text,
				ReplyToMessageID: replyTo,
			}

			resp, err := client.Post(cmd.Context(), "/v1/chats/"+url.PathEscape(chatID)+"/messages", body)
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, "Chat"); err != nil {
				return err
			}

			var result api.SendMessageResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			return outfmt.Output(cmd.Context(), result, func(w io.Writer) {
				_, _ = fmt.Fprintf(w, "Message sent (ID: %s)\n", result.MessageID)
			})
		},
	}

	cmd.Flags().StringVar(&text, "text", "", "Message text (required)")
	cmd.Flags().StringVar(&replyTo, "reply-to", "", "Reply to message ID")
	cmd.Flags().StringVar(&to, "to", "", "Send to chat by name (searches for matching chat)")
	_ = cmd.MarkFlagRequired("text")

	return cmd
}

// resolveChatByName searches for a chat by name and returns its ID.
// If multiple matches are found, it returns the first one.
// If no matches are found, it returns an error.
func resolveChatByName(cmd *cobra.Command, client *api.Client, name string) (string, error) {
	params := url.Values{}
	params.Set("query", name)

	resp, err := client.Get(cmd.Context(), "/v1/chats/search?"+params.Encode())
	if err != nil {
		return "", api.UserFriendlyError(err)
	}
	defer func() { _ = resp.Body.Close() }()

	if err := api.ParseErrorWithContext(resp, ""); err != nil {
		return "", err
	}

	var result api.ListChatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Items) == 0 {
		return "", fmt.Errorf("no chat found matching %q", name)
	}

	chat := result.Items[0]
	if len(result.Items) > 1 {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "Note: %d chats match %q, using %q (%s)\n",
			len(result.Items), name, chat.Title, chat.Network)
	}

	return chat.ID, nil
}

// looksLikeChatID returns true if the string looks like a chat ID rather than a name
func looksLikeChatID(s string) bool {
	// Chat IDs typically contain special characters like ! : @ or ##
	return strings.Contains(s, "!") || strings.Contains(s, ":") || strings.Contains(s, "@") || strings.Contains(s, "##")
}

// senderName returns a display name for a sender ID
func senderName(senderID string) string {
	if strings.Contains(senderID, ":beeper.com") {
		return "You"
	}
	// Extract username from Matrix-style ID like @username:server
	if strings.HasPrefix(senderID, "@") {
		parts := strings.SplitN(senderID[1:], ":", 2)
		if len(parts) > 0 && len(parts[0]) < 20 {
			return parts[0]
		}
	}
	return "Them"
}
