# ygm - You've Got Marketing CLI

Command-line interface for [You've Got Marketing](https://youvegotmarketing.com).

## Installation

### From Source

```bash
go install github.com/CromulentConsulting/ygm-cli/cmd/ygm@latest
```

### From Releases

Download the latest binary from [Releases](https://github.com/CromulentConsulting/ygm-cli/releases).

## Usage

### Authentication

```bash
ygm login
```

This opens your browser for authentication. Enter the code shown in your terminal.

### View Brand DNA

```bash
ygm brand          # Human-readable output
ygm brand --json   # JSON output for scripts
```

### Tasks

```bash
# List tasks
ygm tasks                      # All tasks
ygm tasks --status=pending     # Filter by status
ygm tasks --platform=instagram # Filter by platform
ygm tasks --json               # JSON output

# Create a task
ygm tasks create --title "Post on Reddit" --platform reddit
ygm tasks create --title "Launch tweet" --description "Announce v2" --platform twitter --date 2026-02-11

# Update a task
ygm tasks update 42 --status completed
ygm tasks update 42 --title "New title" --description "Updated copy"

# Discard a task (soft-delete)
ygm tasks discard 42
```

### Get Context for AI Prompts

```bash
ygm context
```

Returns a JSON dump of your brand DNA, marketing plan, and pending tasks - perfect for including in AI prompts.

### Multi-Organization Support

```bash
ygm --org=other-org brand   # Use a specific organization
```

## Configuration

Config is stored in `~/.config/ygm/config.yml`:

```yaml
version: 1
default_org: acme-corp
api_url: https://youvegotmarketing.com

accounts:
  acme-corp:
    token: ygm_abc123...
    user_email: user@example.com
    org_id: 1
    org_name: Acme Corp
```

## Development

```bash
# Build
make build

# Install locally
make install

# Run tests
make test

# Cross-compile for all platforms
make build-all
```

## License

MIT
