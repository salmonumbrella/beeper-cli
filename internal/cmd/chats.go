package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/beeper-cli/internal/api"
	"github.com/salmonumbrella/beeper-cli/internal/outfmt"
)

func newChatsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chats",
		Short: "Manage chats",
	}

	cmd.AddCommand(newChatsListCmd())
	cmd.AddCommand(newChatsGetCmd())
	cmd.AddCommand(newChatsSearchCmd())
	cmd.AddCommand(newChatsArchiveCmd())
	cmd.AddCommand(newChatsArchiveReadCmd())

	return cmd
}

func newChatsListCmd() *cobra.Command {
	var (
		unreadOnly bool
		inbox      string
		limit      int
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List chats",
		Long: `List recent chats from your inbox.

Note: This command returns only the ~25 most recent chats due to API limitations.
To find ALL chats, use 'beeper chats search <query>' which searches your entire history.

To archive all read chats (searching through ALL chats), use:
  beeper chats archive-read`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}

			params := url.Values{}
			if inbox != "" {
				params.Set("inbox", inbox)
			}
			if flags.Account != "" {
				for _, id := range strings.Split(flags.Account, ",") {
					params.Add("accountIDs", strings.TrimSpace(id))
				}
			}

			path := "/v1/chats"
			if len(params) > 0 {
				path += "?" + params.Encode()
			}

			resp, err := client.Get(cmd.Context(), path)
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, ""); err != nil {
				return err
			}

			var result api.ListChatsResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			chats := result.Items
			if unreadOnly {
				filtered := make([]api.Chat, 0)
				for _, c := range chats {
					if c.UnreadCount > 0 {
						filtered = append(filtered, c)
					}
				}
				chats = filtered
			}

			if limit > 0 && len(chats) > limit {
				chats = chats[:limit]
			}

			colorEnabled := outfmt.ShouldColorize(outfmt.GetColor(cmd.Context()))
			return outfmt.Output(cmd.Context(), chats, func(w io.Writer) {
				tw := outfmt.NewTableWriter(w)
				headers := []string{"ID", "Name", "Network", "Unread", "Last Activity"}
				if colorEnabled {
					for i, h := range headers {
						headers[i] = outfmt.Colorize(h, outfmt.Bold, true)
					}
				}
				tw.SetHeader(headers)
				for _, c := range chats {
					unread := ""
					if c.UnreadCount > 0 {
						unread = fmt.Sprintf("%d", c.UnreadCount)
					}
					tw.Append([]string{
						truncate(c.ID, 20),
						truncate(c.Title, 30),
						c.Network,
						unread,
						formatTime(c.LastActivity),
					})
				}
				tw.Render()
			})
		},
	}

	cmd.Flags().BoolVar(&unreadOnly, "unread", false, "Show only unread chats")
	cmd.Flags().StringVar(&inbox, "inbox", "", "Filter by inbox: primary, low-priority, archive")
	cmd.Flags().IntVar(&limit, "limit", 0, "Limit number of results")

	return cmd
}

func newChatsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <chat-id>",
		Short: "Get chat details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			chatID := args[0]

			client, err := getClient()
			if err != nil {
				return err
			}

			resp, err := client.Get(cmd.Context(), "/v1/chats/"+url.PathEscape(chatID))
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, "Chat"); err != nil {
				return err
			}

			var chat api.Chat
			if err := json.NewDecoder(resp.Body).Decode(&chat); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			return outfmt.Output(cmd.Context(), chat, func(w io.Writer) {
				_, _ = fmt.Fprintf(w, "ID:          %s\n", chat.ID)
				_, _ = fmt.Fprintf(w, "Name:        %s\n", chat.Title)
				_, _ = fmt.Fprintf(w, "Network:     %s\n", chat.Network)
				_, _ = fmt.Fprintf(w, "Account:     %s\n", chat.AccountID)
				_, _ = fmt.Fprintf(w, "Type:        %s\n", chat.Type)
				_, _ = fmt.Fprintf(w, "Unread:      %d\n", chat.UnreadCount)
				_, _ = fmt.Fprintf(w, "Archived:    %t\n", chat.IsArchived)
				_, _ = fmt.Fprintf(w, "Last Active: %s\n", formatTime(chat.LastActivity))
				if chat.Participants != nil && len(chat.Participants.Items) > 0 {
					_, _ = fmt.Fprintf(w, "\nParticipants:\n")
					for _, p := range chat.Participants.Items {
						me := ""
						if p.IsSelf {
							me = " (me)"
						}
						_, _ = fmt.Fprintf(w, "  - %s%s\n", p.FullName, me)
					}
				}
			})
		},
	}
}

func newChatsSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search chats",
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

			resp, err := client.Get(cmd.Context(), "/v1/chats/search?"+params.Encode())
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, ""); err != nil {
				return err
			}

			var result api.ListChatsResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}

			colorEnabled := outfmt.ShouldColorize(outfmt.GetColor(cmd.Context()))
			return outfmt.Output(cmd.Context(), result.Items, func(w io.Writer) {
				tw := outfmt.NewTableWriter(w)
				headers := []string{"ID", "Name", "Network", "Type"}
				if colorEnabled {
					for i, h := range headers {
						headers[i] = outfmt.Colorize(h, outfmt.Bold, true)
					}
				}
				tw.SetHeader(headers)
				for _, c := range result.Items {
					tw.Append([]string{
						truncate(c.ID, 20),
						truncate(c.Title, 30),
						c.Network,
						c.Type,
					})
				}
				tw.Render()
			})
		},
	}
}

func newChatsArchiveCmd() *cobra.Command {
	var (
		unarchive bool
		chat      string
	)

	cmd := &cobra.Command{
		Use:   "archive [chat-id]",
		Short: "Archive or unarchive a chat",
		Long: `Archive or unarchive a chat.

You can specify the chat either by ID or by using --chat with a name:
  beeper chats archive <chat-id>
  beeper chats archive --chat "Kishan"`,
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

			body := map[string]bool{"archived": !unarchive}
			resp, err := client.Post(cmd.Context(), "/v1/chats/"+url.PathEscape(chatID)+"/archive", body)
			if err != nil {
				return api.UserFriendlyError(err)
			}
			defer func() { _ = resp.Body.Close() }()

			if err := api.ParseErrorWithContext(resp, "Chat"); err != nil {
				return err
			}

			action := "archived"
			if unarchive {
				action = "unarchived"
			}
			fmt.Printf("Chat %s\n", action)
			return nil
		},
	}

	cmd.Flags().BoolVar(&unarchive, "unarchive", false, "Unarchive instead of archive")
	cmd.Flags().StringVar(&chat, "chat", "", "Specify chat by name (searches for matching chat)")

	return cmd
}

func newChatsArchiveReadCmd() *cobra.Command {
	var dryRun bool

	cmd := &cobra.Command{
		Use:   "archive-read",
		Short: "Archive all read chats",
		Long: `Archive all chats that have zero unread messages.

This command searches through ALL your chats (not just the recent 25) and
archives any chat where unreadCount is 0.

Use --dry-run to preview which chats would be archived without making changes.

Note: The regular 'chats list' only returns ~25 recent chats. This command
uses search to find ALL chats across your entire history.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient()
			if err != nil {
				return err
			}

			// Search patterns to find all chats (vowels + common consonants cover most names)
			patterns := []string{"a", "e", "i", "o", "u", "s", "t", "n", "r", "l", "m", "c", "d", "p", "b", "k", "g", "h", "j", "w", "y", "1", "2", "3"}

			// Collect unique chat IDs to archive
			toArchive := make(map[string]string) // id -> title

			for _, pattern := range patterns {
				params := url.Values{}
				params.Set("query", pattern)

				resp, err := client.Get(cmd.Context(), "/v1/chats/search?"+params.Encode())
				if err != nil {
					continue // Skip failed searches
				}

				var result api.ListChatsResponse
				if resp.StatusCode == 200 {
					_ = json.NewDecoder(resp.Body).Decode(&result)
				}
				_ = resp.Body.Close()

				for _, chat := range result.Items {
					if chat.UnreadCount == 0 && !chat.IsArchived {
						toArchive[chat.ID] = chat.Title
					}
				}
			}

			if len(toArchive) == 0 {
				fmt.Println("No read chats to archive")
				return nil
			}

			if dryRun {
				fmt.Printf("Would archive %d chats:\n", len(toArchive))
				for _, title := range toArchive {
					fmt.Printf("  - %s\n", title)
				}
				return nil
			}

			fmt.Printf("Archiving %d read chats...\n", len(toArchive))
			archived := 0
			for id, title := range toArchive {
				body := map[string]bool{"archived": true}
				resp, err := client.Post(cmd.Context(), "/v1/chats/"+url.PathEscape(id)+"/archive", body)
				if err != nil {
					fmt.Printf("  ✗ %s (error)\n", title)
					continue
				}
				_ = resp.Body.Close()
				if resp.StatusCode == 200 {
					fmt.Printf("  ✓ %s\n", title)
					archived++
				}
			}
			fmt.Printf("\nArchived %d/%d chats\n", archived, len(toArchive))
			return nil
		},
	}

	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview which chats would be archived")

	return cmd
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	local := t.Local()
	now := time.Now()
	if local.Year() == now.Year() && local.YearDay() == now.YearDay() {
		return local.Format("3:04 PM")
	}
	return local.Format("Jan 2, 3:04 PM")
}
