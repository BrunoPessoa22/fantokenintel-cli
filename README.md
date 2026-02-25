# fti — Fan Token Intel CLI

The fastest way for AI agents to access fan token intelligence.

Your agent runs:

```bash
fti signals active --json | jq '.[0]'
fti whales --all --json | jq '.transactions | sort_by(.value_usd) | reverse | .[0]'
fti tokens list --json | jq '[.[] | select(.health_grade == "A")]'
fti prices PSG --history --json | jq '[.prices[] | {time: .time, price: .price}]'
```

One binary. Structured JSON output. No HTTP client, no SDK, no MCP config. If your agent can run a bash command, it can call `fti`.

Built for agents running in bash environments — Claude Code, Codex, shell agents, and automation pipelines. Also works great as a human CLI for quick terminal lookups.

---

## Install

### Homebrew (macOS / Linux)

```bash
brew tap BrunoPessoa22/fantokenintel
brew install fantokenintel
```

### Download binary

Grab the latest release from [GitHub Releases](https://github.com/BrunoPessoa22/fantokenintel-cli/releases) for your platform (darwin/linux/windows, amd64/arm64).

### Build from source

```bash
git clone https://github.com/BrunoPessoa22/fantokenintel-cli
cd fantokenintel-cli
go build -o fti .
```

---

## Auth

Get an API key first:

```bash
fti auth register       # interactive prompt
fti auth login          # save your key to ~/.fti/config.toml
fti auth me             # verify key + show tier/rate-limit
```

API key lookup order: `--api-key` flag → `FTI_API_KEY` env → `~/.fti/config.toml`

---

## Commands

### Tokens

```bash
fti tokens list                        # all fan tokens, sorted by volume
fti tokens list --sort-by health_score # sort options: volume_24h, price_change_24h, market_cap, health_score
fti tokens get PSG                     # full detail: market, exchanges, holders
```

### Prices

```bash
fti prices PSG                                    # current price snapshot
fti prices PSG --history                          # 7-day hourly history
fti prices PSG --history --days 14 --interval 4h
fti prices PSG --history --limit 20               # show last 20 rows
```

### Signals  *(API key required)*

```bash
fti signals active                      # all active signals
fti signals active --token PSG --min-confidence 0.8
fti signals history                     # last 30 days
fti signals history --token BAR --days 90 --outcome target_hit
```

### Whales

```bash
fti whales PSG                            # PSG whale trades (last 24h)
fti whales PSG --hours 1                  # last 1 hour
fti whales --all                          # all tokens
fti whales --all --min-value 100000       # filter by trade size
fti whales --all --watch                  # stream, refresh every 30s
fti whales --all --watch --interval 10    # refresh every 10s
```

### Sports

```bash
fti sports upcoming                     # all upcoming matches (14 days)
fti sports upcoming --token PSG --days 30
```

---

## JSON output

Every command accepts `--json` for machine-readable output:

```bash
fti signals active --json | jq '.[0].token'
fti whales --all --json | jq '.transactions | sort_by(.value_usd) | reverse | .[0]'
fti tokens list --json | jq '[.[] | select(.health_grade == "A")]'
```

---

## Config

`~/.fti/config.toml`:

```toml
api_key = "ti_live_..."
api_url = "https://web-production-ad7c4.up.railway.app"   # optional override
```

---

## Shell completions

```bash
# zsh
echo 'eval "$(fti completion zsh)"' >> ~/.zshrc

# bash
echo 'eval "$(fti completion bash)"' >> ~/.bashrc

# fish
fti completion fish | source
```

---

## Agent usage

`fti` is designed as a first-class data source for AI agents. Every command outputs structured JSON with `--json`:

```bash
# Signal intelligence
fti signals active --json                              # all active signals
fti signals active --token PSG --json                  # single token
fti whales --all --json | jq '.transactions | sort_by(.value_usd) | reverse | .[0]'

# Pre-match alpha
fti sports upcoming --json | jq '[.[] | select(.days_until < 2)]'

# Pipe into your reasoning step
fti signals active --json | jq '.[0]' | your-agent-script

# Token screening
fti tokens list --json | jq '[.[] | select(.health_grade == "A")]'
```

No HTTP client setup. No SDK. No MCP config required — if your agent has bash access, it has fan token data.

---

## Distribution

| Platform | Format |
|---|---|
| macOS arm64 / amd64 | `fti_darwin_arm64.tar.gz` |
| Linux arm64 / amd64 | `fti_linux_arm64.tar.gz` |
| Windows amd64 | `fti_windows_amd64.zip` |

Cross-compiled via [GoReleaser](https://goreleaser.com) on every `v*` tag.
