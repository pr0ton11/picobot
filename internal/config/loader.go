package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// LoadConfig loads config from ~/.picobot/config.json if present, then applies any environment variable overrides on top.
// If no config file exists, it automatically runs onboarding to create one.
func LoadConfig() (Config, error) {
	cfgPath, _, err := ResolveDefaultPaths()
	if err != nil {
		return Config{}, fmt.Errorf("resolving config path: %w", err)
	}

	// Auto-onboard if config doesn't exist yet.
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		fmt.Println("First run detected — running onboard...")
		cfgOut, wsOut, err := Onboard()
		if err != nil {
			return Config{}, fmt.Errorf("auto-onboard failed: %w", err)
		}
		fmt.Printf("Wrote config to %s\nInitialized workspace at %s\n", cfgOut, wsOut)
	}

	var cfg Config
	f, err := os.Open(cfgPath)
	if err == nil {
		defer func() { _ = f.Close() }()
		if err := json.NewDecoder(f).Decode(&cfg); err != nil {
			return Config{}, err
		}
	}
	// env vars always take precedence over the config file, enabling runtime overrides without editing config.json.
	applyEnvOverrides(&cfg)
	return cfg, nil
}

// applyEnvOverrides updates config fields from all environment variables.
func applyEnvOverrides(cfg *Config) {
	// --- Provider settings ---
	if v := os.Getenv("OPENAI_API_KEY"); v != "" {
		if cfg.Providers.OpenAI == nil {
			cfg.Providers.OpenAI = &ProviderConfig{}
		}
		cfg.Providers.OpenAI.APIKey = v
	}
	if v := os.Getenv("OPENAI_API_BASE"); v != "" {
		if cfg.Providers.OpenAI == nil {
			cfg.Providers.OpenAI = &ProviderConfig{}
		}
		cfg.Providers.OpenAI.APIBase = v
	}

	// --- Agent defaults ---
	if v := os.Getenv("PICOBOT_MODEL"); v != "" {
		cfg.Agents.Defaults.Model = v
	}
	if v := os.Getenv("PICOBOT_MAX_TOKENS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Agents.Defaults.MaxTokens = n
		}
	}
	if v := os.Getenv("PICOBOT_MAX_TOOL_ITERATIONS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			cfg.Agents.Defaults.MaxToolIterations = n
		}
	}

	// --- Telegram ---
	if v := os.Getenv("TELEGRAM_BOT_TOKEN"); v != "" {
		cfg.Channels.Telegram.Enabled = true
		cfg.Channels.Telegram.Token = v
	}
	if v := os.Getenv("TELEGRAM_ALLOW_FROM"); v != "" {
		cfg.Channels.Telegram.AllowFrom = splitCSV(v)
	}

	// --- Discord ---
	if v := os.Getenv("DISCORD_BOT_TOKEN"); v != "" {
		cfg.Channels.Discord.Enabled = true
		cfg.Channels.Discord.Token = v
	}
	if v := os.Getenv("DISCORD_ALLOW_FROM"); v != "" {
		cfg.Channels.Discord.AllowFrom = splitCSV(v)
	}

	// --- Slack ---
	if v := os.Getenv("SLACK_APP_TOKEN"); v != "" {
		cfg.Channels.Slack.Enabled = true
		cfg.Channels.Slack.AppToken = v
	}
	if v := os.Getenv("SLACK_BOT_TOKEN"); v != "" {
		cfg.Channels.Slack.Enabled = true
		cfg.Channels.Slack.BotToken = v
	}
	if v := os.Getenv("SLACK_ALLOW_USERS"); v != "" {
		cfg.Channels.Slack.AllowUsers = splitCSV(v)
	}
	if v := os.Getenv("SLACK_ALLOW_CHANNELS"); v != "" {
		cfg.Channels.Slack.AllowChannels = splitCSV(v)
	}

	// --- Signal ---
	if v := os.Getenv("SIGNAL_API_URL"); v != "" {
		cfg.Channels.Signal.Enabled = true
		cfg.Channels.Signal.APIURL = v
	}
	if v := os.Getenv("SIGNAL_API_TOKEN"); v != "" {
		cfg.Channels.Signal.APIToken = v
	}
	if v := os.Getenv("SIGNAL_NUMBER"); v != "" {
		cfg.Channels.Signal.Number = v
	}
	if v := os.Getenv("SIGNAL_ALLOW_FROM"); v != "" {
		cfg.Channels.Signal.AllowFrom = splitCSV(v)
	}
}

// splitCSV splits a comma-separated string into trimmed, non-empty parts.
func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
