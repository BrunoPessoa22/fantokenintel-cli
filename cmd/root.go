package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags "-X cmd.Version=x.y.z"
var Version = "dev"

var (
	apiKey  string
	jsonOut bool
)

const defaultBaseURL = "https://web-production-ad7c4.up.railway.app"

var rootCmd = &cobra.Command{
	Use:     "fti",
	Short:   "Fan Token Intel CLI",
	Version: Version,
	Long: `fti — command-line interface for Fan Token Intel

Market data, whale tracking, and trading signals for sports fan tokens.
Powered by the Fan Token Intel API (https://fantokenintel.vercel.app).

Get an API key:
  fti auth register

Quick start:
  fti tokens list
  fti signals active --token PSG
  fti whales --all`,
	SilenceUsage: true,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "API key (overrides FTI_API_KEY env and ~/.fti/config.toml)")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output raw JSON")
}
