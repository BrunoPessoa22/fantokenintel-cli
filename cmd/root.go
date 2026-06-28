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

The descriptive fan-token intelligence layer for agents: prices, whale-distribution
flows, on-chain data and match-impact across fan tokens. Data, not advice.
Powered by the Fan Token Intel API (https://fantokenintel.com).

Quick start:
  fti tokens list
  fti whales --all
  fti sports upcoming

Premium — the event-impact moat (match-event → market-adjusted price reactions)
is metered per call on the x402 gateway: https://x402.brunopessoa.com
Free higher limits with a key:  fti auth register`,
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
