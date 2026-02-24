package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BrunoPessoa22/fantokenintel-cli/internal"
	"github.com/spf13/cobra"
)

var signalsCmd = &cobra.Command{
	Use:   "signals",
	Short: "Trading signals (active and historical)",
}

// ── signals active ───────────────────────────────────────────────────────────

var (
	signalsToken     string
	signalsMinConf   float64
	signalsDays      int
	signalsOutcome   string
	signalsLimit     int
)

var signalsActiveCmd = &cobra.Command{
	Use:   "active",
	Short: "Show currently active trading signals",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := internal.ResolveAPIKey(apiKey)
		if err != nil {
			return err
		}
		if key == "" {
			return fmt.Errorf("API key required — run: fti auth login")
		}

		baseURL := internal.ResolveBaseURL(defaultBaseURL)
		c := internal.NewClient(baseURL, key)

		params := map[string]string{
			"min_confidence": fmt.Sprintf("%.2f", signalsMinConf),
		}
		if signalsToken != "" {
			params["token"] = strings.ToUpper(signalsToken)
		}
		q := buildQuery(params)

		var resp struct {
			ActiveSignals int `json:"active_signals"`
			Signals       []struct {
				ID                 string  `json:"id"`
				Token              string  `json:"token"`
				Direction          string  `json:"direction"`
				Tier               string  `json:"tier"`
				SellRatio          float64 `json:"sell_ratio"`
				ConfidenceScore    float64 `json:"confidence_score"`
				EntryPrice         float64 `json:"entry_price"`
				TargetPrice        float64 `json:"target_price"`
				StopPrice          float64 `json:"stop_price"`
				CreatedAt          string  `json:"created_at"`
				ExpiresAt          string  `json:"expires_at"`
				PrimaryReason      string  `json:"primary_reason"`
				MaxProfitPct       float64 `json:"max_profit_pct"`
				TrailingStopStatus string  `json:"trailing_stop_status"`
			} `json:"signals"`
		}

		raw, err := c.Get("/api/v1/signals/active", q, &resp)
		if err != nil {
			return err
		}

		if jsonOut {
			internal.PrintJSON(raw)
			return nil
		}

		if resp.ActiveSignals == 0 {
			internal.Dim.Println("\nNo active signals.")
			return nil
		}

		internal.Bold.Printf("\n%d active signal(s)\n\n", resp.ActiveSignals)
		t := internal.NewTable("TOKEN", "DIR", "TIER", "CONF", "ENTRY", "TARGET", "STOP", "MAX%", "EXPIRES")
		t.Header()
		for _, s := range resp.Signals {
			t.Row(
				internal.Cyan.Sprint(s.Token),
				internal.FormatDirection(s.Direction),
				tierColor(s.Tier),
				internal.FormatConfidence(s.ConfidenceScore),
				internal.FormatPrice(s.EntryPrice),
				internal.FormatPrice(s.TargetPrice),
				internal.FormatPrice(s.StopPrice),
				fmt.Sprintf("%.1f%%", s.MaxProfitPct),
				shortTime(s.ExpiresAt),
			)
		}
		t.Flush()
		fmt.Println()
		return nil
	},
}

// ── signals history ──────────────────────────────────────────────────────────

var signalsHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "Show historical signal performance",
	RunE: func(cmd *cobra.Command, args []string) error {
		key, err := internal.ResolveAPIKey(apiKey)
		if err != nil {
			return err
		}
		if key == "" {
			return fmt.Errorf("API key required — run: fti auth login")
		}

		baseURL := internal.ResolveBaseURL(defaultBaseURL)
		c := internal.NewClient(baseURL, key)

		params := map[string]string{
			"days":  strconv.Itoa(signalsDays),
			"limit": strconv.Itoa(signalsLimit),
		}
		if signalsToken != "" {
			params["token"] = strings.ToUpper(signalsToken)
		}
		if signalsOutcome != "" {
			params["outcome"] = signalsOutcome
		}
		q := buildQuery(params)

		var resp struct {
			Signals []struct {
				ID              string  `json:"id"`
				Token           string  `json:"token"`
				Tier            string  `json:"tier"`
				SellRatio       float64 `json:"sell_ratio"`
				ConfidenceScore float64 `json:"confidence_score"`
				EntryPrice      float64 `json:"entry_price"`
				TargetPrice     float64 `json:"target_price"`
				StopPrice       float64 `json:"stop_price"`
				OutcomeStatus   string  `json:"outcome_status"`
				PnlPct          float64 `json:"pnl_pct"`
				MaxProfitPct    float64 `json:"max_profit_pct"`
				CreatedAt       string  `json:"created_at"`
				ExitTime        string  `json:"exit_time"`
				ExitPrice       float64 `json:"exit_price"`
			} `json:"signals"`
		}

		raw, err := c.Get("/api/v1/signals/history", q, &resp)
		if err != nil {
			return err
		}

		if jsonOut {
			internal.PrintJSON(raw)
			return nil
		}

		if len(resp.Signals) == 0 {
			internal.Dim.Println("\nNo signal history found.")
			return nil
		}

		internal.Bold.Printf("\n%d signal(s) — last %d days\n\n", len(resp.Signals), signalsDays)
		t := internal.NewTable("TOKEN", "TIER", "CONF", "ENTRY", "EXIT", "PNL%", "OUTCOME", "DATE")
		t.Header()
		for _, s := range resp.Signals {
			t.Row(
				internal.Cyan.Sprint(s.Token),
				tierColor(s.Tier),
				internal.FormatConfidence(s.ConfidenceScore),
				internal.FormatPrice(s.EntryPrice),
				internal.FormatPrice(s.ExitPrice),
				formatPnl(s.PnlPct),
				internal.FormatOutcome(s.OutcomeStatus),
				shortTime(s.CreatedAt),
			)
		}
		t.Flush()
		fmt.Println()
		return nil
	},
}

func tierColor(tier string) string {
	switch tier {
	case "high":
		return internal.Green.Sprint(tier)
	case "medium":
		return internal.Yellow.Sprint(tier)
	default:
		return internal.Dim.Sprint(tier)
	}
}

func formatPnl(pct float64) string {
	if pct == 0 {
		return internal.Dim.Sprint("—")
	}
	s := fmt.Sprintf("%+.1f%%", pct)
	if pct > 0 {
		return internal.Green.Sprint(s)
	}
	return internal.Red.Sprint(s)
}

func init() {
	signalsActiveCmd.Flags().StringVar(&signalsToken, "token", "", "Filter by token symbol")
	signalsActiveCmd.Flags().Float64Var(&signalsMinConf, "min-confidence", 0.65, "Minimum confidence (0-1)")

	signalsHistoryCmd.Flags().StringVar(&signalsToken, "token", "", "Filter by token symbol")
	signalsHistoryCmd.Flags().IntVar(&signalsDays, "days", 30, "Look-back period in days")
	signalsHistoryCmd.Flags().StringVar(&signalsOutcome, "outcome", "", "Filter by outcome (target_hit, stopped_out, expired)")
	signalsHistoryCmd.Flags().IntVar(&signalsLimit, "limit", 50, "Max results")

	signalsCmd.AddCommand(signalsActiveCmd)
	signalsCmd.AddCommand(signalsHistoryCmd)
	rootCmd.AddCommand(signalsCmd)
}
