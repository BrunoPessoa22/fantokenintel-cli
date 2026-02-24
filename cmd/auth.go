package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/BrunoPessoa22/fantokenintel-cli/internal"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage API keys and account",
}

// ── auth register ────────────────────────────────────────────────────────────

var authRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new API key interactively",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Name: ")
		name, _ := reader.ReadString('\n')
		name = strings.TrimSpace(name)

		fmt.Print("Email: ")
		email, _ := reader.ReadString('\n')
		email = strings.TrimSpace(email)

		fmt.Print("Description (optional): ")
		desc, _ := reader.ReadString('\n')
		desc = strings.TrimSpace(desc)

		fmt.Print("Scope [read/full] (default: read): ")
		scope, _ := reader.ReadString('\n')
		scope = strings.TrimSpace(scope)
		if scope == "" {
			scope = "read"
		}

		baseURL := internal.ResolveBaseURL(defaultBaseURL)
		c := internal.NewClient(baseURL, "")

		payload := map[string]string{
			"name":        name,
			"email":       email,
			"description": desc,
			"scope":       scope,
		}

		var resp struct {
			AgentID           string   `json:"agent_id"`
			APIKey            string   `json:"api_key"`
			Name              string   `json:"name"`
			Tier              string   `json:"tier"`
			RateLimitPerMin   int      `json:"rate_limit_per_minute"`
			EmailVerified     bool     `json:"email_verified"`
			Capabilities      []string `json:"capabilities"`
			Message           string   `json:"message"`
		}

		raw, err := c.Post("/api/v1/auth/register", payload, &resp)
		if err != nil {
			return err
		}

		if jsonOut {
			internal.PrintJSON(raw)
			return nil
		}

		color.New(color.Bold, color.FgGreen).Println("\nRegistration successful!")
		fmt.Printf("  API Key:    %s\n", color.New(color.Bold).Sprint(resp.APIKey))
		fmt.Printf("  Agent ID:   %s\n", internal.Dim.Sprint(resp.AgentID))
		fmt.Printf("  Tier:       %s\n", resp.Tier)
		fmt.Printf("  Rate limit: %d req/min\n", resp.RateLimitPerMin)
		fmt.Printf("  Scope:      %s\n", strings.Join(resp.Capabilities, ", "))
		if resp.Message != "" {
			fmt.Printf("\n  %s\n", internal.Dim.Sprint(resp.Message))
		}
		fmt.Printf("\nSave your key with:\n  fti auth login\n\n")
		return nil
	},
}

// ── auth login ───────────────────────────────────────────────────────────────

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save an API key to ~/.fti/config.toml",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Paste your API key (ti_live_...): ")
		key, _ := reader.ReadString('\n')
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("no API key provided")
		}

		cfg, err := internal.LoadConfig()
		if err != nil {
			return err
		}
		cfg.APIKey = key
		if err := internal.SaveConfig(cfg); err != nil {
			return err
		}

		internal.Green.Println("API key saved to ~/.fti/config.toml")
		return nil
	},
}

// ── auth me ──────────────────────────────────────────────────────────────────

var authMeCmd = &cobra.Command{
	Use:   "me",
	Short: "Show current API key info",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := internal.ResolveAPIKey(apiKey)
		if err != nil {
			return err
		}
		if key == "" {
			return fmt.Errorf("no API key found — run: fti auth login")
		}

		baseURL := internal.ResolveBaseURL(defaultBaseURL)
		c := internal.NewClient(baseURL, key)

		var resp struct {
			AgentID         string   `json:"agent_id"`
			Name            string   `json:"name"`
			Description     string   `json:"description"`
			Tier            string   `json:"tier"`
			Capabilities    []string `json:"capabilities"`
			RateLimitPerMin int      `json:"rate_limit_per_minute"`
			TotalRequests   int      `json:"total_requests"`
			CreatedAt       string   `json:"created_at"`
		}

		raw, err := c.Get("/api/v1/auth/me", nil, &resp)
		if err != nil {
			return err
		}

		if jsonOut {
			internal.PrintJSON(raw)
			return nil
		}

		internal.Bold.Printf("\n%s\n", resp.Name)
		fmt.Printf("  Agent ID:      %s\n", internal.Dim.Sprint(resp.AgentID))
		fmt.Printf("  Tier:          %s\n", internal.Cyan.Sprint(resp.Tier))
		fmt.Printf("  Scope:         %s\n", strings.Join(resp.Capabilities, ", "))
		fmt.Printf("  Rate limit:    %d req/min\n", resp.RateLimitPerMin)
		fmt.Printf("  Total calls:   %d\n", resp.TotalRequests)
		if resp.Description != "" {
			fmt.Printf("  Description:   %s\n", resp.Description)
		}
		fmt.Printf("  Created:       %s\n", internal.Dim.Sprint(resp.CreatedAt))
		fmt.Println()
		return nil
	},
}

func init() {
	authCmd.AddCommand(authRegisterCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authMeCmd)

	rootCmd.AddCommand(authCmd)
}
