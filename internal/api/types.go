package api

import "time"

// Account represents a connected chat network
type Account struct {
	ID              string `json:"id"`
	NetworkName     string `json:"networkName"`
	NetworkIcon     string `json:"networkIcon,omitempty"`
	ProfileName     string `json:"profileName,omitempty"`
	ProfileEmail    string `json:"profileEmail,omitempty"`
	ProfilePhone    string `json:"profilePhone,omitempty"`
	ProfileUsername string `json:"profileUsername,omitempty"`
}

type GetAccountsResponse struct {
	Accounts []Account `json:"accounts"`
}

// Chat represents a conversation
type Chat struct {
	ID           string           `json:"id"`
	LocalChatID  string           `json:"localChatID,omitempty"`
	AccountID    string           `json:"accountID"`
	Network      string           `json:"network,omitempty"`
	Title        string           `json:"title"`
	Type         string           `json:"type"` // "dm" or "group"
	UnreadCount  int              `json:"unreadCount"`
	IsMuted      bool             `json:"isMuted"`
	IsArchived   bool             `json:"isArchived"`
	IsPinned     bool             `json:"isPinned"`
	LastActivity time.Time        `json:"lastActivity"`
	Preview      *MessagePreview  `json:"preview,omitempty"`
	Participants *ParticipantList `json:"participants,omitempty"`
	ReminderAt   *time.Time       `json:"reminderAt,omitempty"`
}

// MessagePreview is the preview shown in chat list
type MessagePreview struct {
	Text      string    `json:"text"`
	Sender    string    `json:"sender,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// ParticipantList wraps the participants array
type ParticipantList struct {
	Items []Participant `json:"items"`
}

type Participant struct {
	ID          string `json:"id"`
	FullName    string `json:"fullName"`
	Username    string `json:"username,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	ImgURL      string `json:"imgURL,omitempty"`
	IsSelf      bool   `json:"isSelf"`
	IsAdmin     bool   `json:"isAdmin"`
	IsVerified  bool   `json:"isVerified"`
}

type ListChatsResponse struct {
	Items   []Chat `json:"items"`
	Total   int    `json:"total"`
	HasMore bool   `json:"hasMore"`
	Cursor  string `json:"cursor,omitempty"`
}

// Message represents a chat message
type Message struct {
	ID        string    `json:"id"`
	ChatID    string    `json:"chatID"`
	AccountID string    `json:"accountID"`
	SenderID  string    `json:"senderID"`
	Sender    string    `json:"sender"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	IsMe      bool      `json:"isMe"`
	ReplyTo   *Message  `json:"replyTo,omitempty"`
}

type ListMessagesResponse struct {
	Items   []Message `json:"items"`
	HasMore bool      `json:"hasMore"`
	Cursor  string    `json:"cursor,omitempty"`
}

type SearchMessagesResponse struct {
	Messages []Message       `json:"messages"`
	Chats    map[string]Chat `json:"chats"`
	HasMore  bool            `json:"hasMore"`
	Cursor   string          `json:"cursor,omitempty"`
}

// SendMessageRequest for sending a message
type SendMessageRequest struct {
	Text             string `json:"text,omitempty"`
	ReplyToMessageID string `json:"replyToMessageID,omitempty"`
}

type SendMessageResponse struct {
	MessageID string `json:"messageID"`
	Deeplink  string `json:"deeplink"`
}

// ReminderRequest for setting a reminder
type ReminderRequest struct {
	Reminder ReminderTime `json:"reminder"`
}

// ReminderTime represents the time format expected by Beeper API
type ReminderTime struct {
	RemindAtMs int64 `json:"remindAtMs"`
}

// NewReminderTime creates a ReminderTime from a time.Time
func NewReminderTime(t time.Time) ReminderTime {
	return ReminderTime{
		RemindAtMs: t.UnixMilli(),
	}
}

// FocusRequest for focusing Beeper window
type FocusRequest struct {
	ChatID              string `json:"chatID,omitempty"`
	MessageID           string `json:"messageID,omitempty"`
	DraftText           string `json:"draftText,omitempty"`
	DraftAttachmentPath string `json:"draftAttachmentPath,omitempty"`
}

// SearchResponse for unified search
type SearchResponse struct {
	Chats    []Chat    `json:"chats"`
	Messages []Message `json:"messages"`
	HasMore  bool      `json:"hasMore"`
}
