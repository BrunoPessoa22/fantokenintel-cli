package internal

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds persisted CLI settings.
type Config struct {
	APIKey  string `toml:"api_key"`
	APIURL  string `toml:"api_url"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".fti", "config.toml"), nil
}

// LoadConfig reads ~/.fti/config.toml. Missing file returns empty Config, no error.
func LoadConfig() (Config, error) {
	path, err := configPath()
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Config{}, nil
	}
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("reading config: %w", err)
	}
	return cfg, nil
}

// SaveConfig writes cfg to ~/.fti/config.toml, creating the directory if needed.
func SaveConfig(cfg Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("creating config dir: %w", err)
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("writing config: %w", err)
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(cfg)
}

// ResolveAPIKey returns the API key using precedence:
// 1. flagValue (from --api-key flag)
// 2. FTI_API_KEY env var
// 3. ~/.fti/config.toml
func ResolveAPIKey(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	if v := os.Getenv("FTI_API_KEY"); v != "" {
		return v, nil
	}
	cfg, err := LoadConfig()
	if err != nil {
		return "", err
	}
	return cfg.APIKey, nil
}

// ResolveBaseURL returns the API base URL, falling back to the default.
func ResolveBaseURL(defaultURL string) string {
	if v := os.Getenv("FTI_API_URL"); v != "" {
		return v
	}
	cfg, _ := LoadConfig()
	if cfg.APIURL != "" {
		return cfg.APIURL
	}
	return defaultURL
}
