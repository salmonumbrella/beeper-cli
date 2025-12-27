# 💬 Beeper CLI — Messaging in your terminal.

Manage chats, send messages, and control Beeper Desktop from the command line.

## Features

- **Authentication** - browser-based login with secure keyring storage
- **Chat management** - list, search, archive, and organize conversations
- **Desktop control** - focus Beeper window, navigate to chats, pre-fill drafts
- **Messaging** - send messages, search history, and view conversations
- **Multi-network support** - manage chats across all connected networks (iMessage, WhatsApp, Telegram, Signal, etc.)
- **Reminders** - set and clear chat reminders

## Installation

### Homebrew

```bash
brew install salmonumbrella/tap/beeper-cli
```

### Go Install

```bash
go install github.com/salmonumbrella/beeper-cli/cmd/beeper@latest
```

### Build from Source

```bash
git clone https://github.com/salmonumbrella/beeper-cli
cd beeper-cli
make build && make install
```

## Quick Start

### 1. Enable Local API

Open Beeper Desktop → Settings → Developers → Enable Local API

### 2. Authenticate

```bash
beeper auth add
```

This opens your browser to authenticate with Beeper and securely stores the token in your system keyring.

### 3. Test Connection

```bash
beeper auth test
```

### 4. Start Using

```bash
beeper accounts              # List connected networks
beeper chats list            # Show recent chats
beeper messages search "hi"  # Search message history
```

## Configuration

### Account Filtering

Filter commands by specific chat networks:

```bash
# Via flag
beeper chats list --account whatsapp

# Multiple networks
beeper messages search "invoice" --account whatsapp,telegram
```

Get account IDs from `beeper accounts`.

### Environment Variables

- `BEEPER_OUTPUT` - Output format: `text` (default) or `json`
- `BEEPER_COLOR` - Color mode: `auto` (default), `always`, or `never`
- `NO_COLOR` - Set to any value to disable colors (standard convention)

## Security

### Credential Storage

Tokens are stored securely in your system's keychain:
- **macOS**: Keychain Access
- **Linux**: Secret Service (GNOME Keyring, KWallet)
- **Windows**: Credential Manager

## Commands

### Authentication

```bash
beeper auth add              # Authenticate via browser (opens browser)
beeper auth list             # List configured tokens
beeper auth remove <name>    # Remove a token
beeper auth test             # Test current token
```

### Accounts

```bash
beeper accounts              # List connected chat networks
```

### Chats

```bash
beeper chats list                           # List recent chats
beeper chats list --unread                  # Only unread chats
beeper chats list --inbox primary           # Filter by inbox
beeper chats list --inbox archive           # Show archived chats
beeper chats list --account whatsapp        # Filter by network
beeper chats get <chat-id>                  # Get chat details
beeper chats search <query>                 # Search for chats by name
beeper chats archive <chat-id>              # Archive a chat
beeper chats archive --chat "John"          # Archive by name
beeper chats archive --unarchive <chat-id>  # Unarchive
beeper chats archive-read                   # Archive all read chats
beeper chats archive-read --dry-run         # Preview what would be archived
```

### Messages

```bash
beeper messages list <chat-id>              # List messages in a chat
beeper messages list --chat "John"          # List by chat name
beeper messages list <chat-id> --limit 50   # Limit results
beeper messages search <query>              # Search all messages
beeper messages search "invoice" --account telegram
beeper messages send <chat-id> --text "Hello!"
beeper messages send --chat "John" --text "Meeting at 3pm"
```

### Reminders

```bash
beeper reminders set <chat-id> --at "2024-12-25 10:00"
beeper reminders set --chat "John" --at "tomorrow 9am"
beeper reminders clear <chat-id>
```

### Focus (Desktop Control)

```bash
beeper focus                                # Bring Beeper to foreground
beeper focus --chat "John"                  # Open specific chat
beeper focus --chat "John" --draft "Hi!"    # Open with pre-filled draft
```

## Output Formats

### Text

Human-readable tables with colors and formatting:

```bash
$ beeper chats list
ID                    NAME              NETWORK     UNREAD  LAST ACTIVITY
!abc123...            John Smith        WhatsApp            3:42 PM
!def456...            Team Chat         Telegram    5       2:15 PM
!ghi789...            Mom               iMessage    2       Yesterday
```

### JSON

Machine-readable output:

```bash
$ beeper chats list -o json
[
  {
    "id": "!abc123...",
    "title": "John Smith",
    "network": "WhatsApp",
    "unreadCount": 0,
    "lastActivity": "2024-12-25T15:42:00Z"
  }
]
```

Data goes to stdout, errors and progress to stderr for clean piping.

## Examples

### Archive all read chats

```bash
# Preview first
beeper chats archive-read --dry-run

# Then execute
beeper chats archive-read
```

### Send a quick message

```bash
beeper messages send --chat "John" --text "Running 5 mins late!"
```

### Search messages across networks

```bash
beeper messages search "pdf" --account whatsapp -o json | jq '.messages[]'
```

### Focus chat with draft from clipboard

```bash
beeper focus --chat "Team" --draft "$(pbpaste)"
```

### Automation

Use `-o json` for scripting and `--dry-run` to preview changes:

```bash
# Count unread chats
beeper chats list -o json | jq '[.[] | select(.unreadCount > 0)] | length'

# Get all Telegram chat IDs
beeper chats list --account telegram -o json --query '.[].id'

# Notify on unread count
unread=$(beeper chats list -o json | jq '[.[] | select(.unreadCount > 0)] | length')
[ "$unread" -gt 10 ] && osascript -e "display notification \"$unread unread\" with title \"Beeper\""
```

### Debug Mode

Enable verbose output for troubleshooting:

```bash
beeper --debug chats list
# Shows: api request method=GET url=http://localhost:9988/...
# Shows: api response status=200 content_length=1234
```

### JQ Filtering

Filter JSON output with JQ expressions:

```bash
# Get only unread chats
beeper chats list -o json --query '.[] | select(.unreadCount > 0)'

# Extract chat IDs
beeper chats list -o json --query '.[].id'

# Get chats from specific network
beeper chats list -o json --query '.[] | select(.network == "Telegram")'
```

## Global Flags

All commands support these flags:

- `--account <id>` - Filter by account/network ID
- `-o, --output <format>` - Output format: `text` or `json` (default: text)
- `--color <mode>` - Color mode: `auto`, `always`, or `never` (default: auto)
- `--debug` - Enable debug output (shows API requests/responses)
- `--query <expr>` - JQ filter expression for JSON output
- `--help` - Show help for any command
- `--version` - Show version information

## Shell Completions

Generate shell completions for your preferred shell:

### Bash

```bash
# macOS (Homebrew):
beeper completion bash > $(brew --prefix)/etc/bash_completion.d/beeper

# Linux:
beeper completion bash > /etc/bash_completion.d/beeper

# Or source directly:
source <(beeper completion bash)
```

### Zsh

```zsh
beeper completion zsh > "${fpath[1]}/_beeper"
```

### Fish

```fish
beeper completion fish > ~/.config/fish/completions/beeper.fish
```

### PowerShell

```powershell
beeper completion powershell | Out-String | Invoke-Expression
```

## Development

After cloning, install git hooks:

```bash
make setup
```

This installs [lefthook](https://github.com/evilmartians/lefthook) pre-commit and pre-push hooks for linting and testing.

## License

MIT

## Links

- [Beeper](https://beeper.com)
- [GitHub Repository](https://github.com/salmonumbrella/beeper-cli)
