package cmd

import (
	"fmt"
	"strings"

	"github.com/BrunoPessoa22/fantokenintel-cli/internal"
	"github.com/spf13/cobra"
)

var tokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "Fan token market data",
}

// ── tokens list ──────────────────────────────────────────────────────────────

var (
	tokensSortBy string
	tokensOrder  string
)

var tokensListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all fan tokens with market metrics",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL := internal.ResolveBaseURL(defaultBaseURL)
		c := internal.NewClient(baseURL, "")

		params := map[string]string{
			"sort_by": tokensSortBy,
			"order":   tokensOrder,
		}
		q := buildQuery(params)

		var tokens []struct {
			Symbol          string  `json:"symbol"`
			Name            string  `json:"name"`
			Team            string  `json:"team"`
			Price           float64 `json:"price"`
			PriceChange1h   float64 `json:"price_change_1h"`
			PriceChange24h  float64 `json:"price_change_24h"`
			Volume24h       float64 `json:"volume_24h"`
			MarketCap       float64 `json:"market_cap"`
			HealthGrade     string  `json:"health_grade"`
			HealthScore     float64 `json:"health_score"`
		}

		raw, err := c.Get("/api/tokens", q, &tokens)
		if err != nil {
			return err
		}

		if jsonOut {
			internal.PrintJSON(raw)
			return nil
		}

		t := internal.NewTable("SYMBOL", "NAME", "PRICE", "1H%", "24H%", "VOLUME", "MCAP", "HEALTH")
		t.Header()
		for _, tk := range tokens {
			t.Row(
				internal.Cyan.Sprint(tk.Symbol),
				internal.TruncStr(tk.Name, 22),
				internal.FormatPrice(tk.Price),
				internal.FormatChange(tk.PriceChange1h),
				internal.FormatChange(tk.PriceChange24h),
				internal.FormatVolume(tk.Volume24h),
				internal.FormatVolume(tk.MarketCap),
				gradeColor(tk.HealthGrade, tk.HealthScore),
			)
		}
		t.Flush()
		fmt.Printf("\n%d tokens\n", len(tokens))
		return nil
	},
}

// ── tokens get ───────────────────────────────────────────────────────────────

var tokensGetCmd = &cobra.Command{
	Use:   "get <SYMBOL>",
	Short: "Get detailed info for a specific token",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		symbol := strings.ToUpper(args[0])
		baseURL := internal.ResolveBaseURL(defaultBaseURL)
		c := internal.NewClient(baseURL, "")

		var resp struct {
			Token struct {
				ID                 int    `json:"id"`
				Symbol             string `json:"symbol"`
				Name               string `json:"name"`
				Team               string `json:"team"`
				League             string `json:"league"`
				Country            string `json:"country"`
				TotalSupply        int64  `json:"total_supply"`
				CirculatingSupply  int64  `json:"circulating_supply"`
				LaunchDate         string `json:"launch_date"`
			} `json:"token"`
			Metrics struct {
				Price           float64 `json:"price"`
				PriceChange1h   float64 `json:"price_change_1h"`
				PriceChange24h  float64 `json:"price_change_24h"`
				PriceChange7d   float64 `json:"price_change_7d"`
				Volume24h       float64 `json:"volume_24h"`
				MarketCap       float64 `json:"market_cap"`
				TotalHolders    int     `json:"total_holders"`
				HolderChange24h int     `json:"holder_change_24h"`
				HealthScore     float64 `json:"health_score"`
				HealthGrade     string  `json:"health_grade"`
				Liquidity1pct   float64 `json:"liquidity_1pct"`
				SpreadBps       float64 `json:"spread_bps"`
			} `json:"metrics"`
			Exchanges []struct {
				Name      string  `json:"name"`
				Price     float64 `json:"price"`
				Volume24h float64 `json:"volume_24h"`
				SpreadBps float64 `json:"spread_bps"`
				BestBid   float64 `json:"best_bid"`
				BestAsk   float64 `json:"best_ask"`
			} `json:"exchanges"`
		}

		raw, err := c.Get("/api/tokens/"+symbol, nil, &resp)
		if err != nil {
			return err
		}

		if jsonOut {
			internal.PrintJSON(raw)
			return nil
		}

		tk := resp.Token
		m := resp.Metrics

		internal.Bold.Printf("\n%s — %s\n", tk.Symbol, tk.Name)
		fmt.Printf("  Team:      %s\n", tk.Team)
		fmt.Printf("  League:    %s\n", tk.League)
		fmt.Printf("  Country:   %s\n", tk.Country)
		if tk.LaunchDate != "" {
			fmt.Printf("  Launch:    %s\n", internal.Dim.Sprint(tk.LaunchDate))
		}

		fmt.Println()
		internal.Bold.Println("Market")
		fmt.Printf("  Price:       %s\n", internal.FormatPrice(m.Price))
		fmt.Printf("  1h / 24h:    %s / %s\n", internal.FormatChange(m.PriceChange1h), internal.FormatChange(m.PriceChange24h))
		fmt.Printf("  7d:          %s\n", internal.FormatChange(m.PriceChange7d))
		fmt.Printf("  Volume 24h:  %s\n", internal.FormatVolume(m.Volume24h))
		fmt.Printf("  Market cap:  %s\n", internal.FormatVolume(m.MarketCap))
		fmt.Printf("  Holders:     %d (%s 24h)\n", m.TotalHolders, formatHolderDelta(m.HolderChange24h))
		fmt.Printf("  Health:      %s\n", gradeColor(m.HealthGrade, m.HealthScore))
		fmt.Printf("  Liquidity:   %s  Spread: %.1f bps\n", internal.FormatVolume(m.Liquidity1pct), m.SpreadBps)

		if len(resp.Exchanges) > 0 {
			fmt.Println()
			internal.Bold.Println("Exchanges")
			t := internal.NewTable("EXCHANGE", "PRICE", "VOLUME", "BID", "ASK", "SPREAD")
			t.Header()
			for _, ex := range resp.Exchanges {
				t.Row(
					ex.Name,
					internal.FormatPrice(ex.Price),
					internal.FormatVolume(ex.Volume24h),
					internal.FormatPrice(ex.BestBid),
					internal.FormatPrice(ex.BestAsk),
					fmt.Sprintf("%.1f bps", ex.SpreadBps),
				)
			}
			t.Flush()
		}
		fmt.Println()
		return nil
	},
}

func gradeColor(grade string, score float64) string {
	label := fmt.Sprintf("%s (%.0f)", grade, score)
	switch grade {
	case "A":
		return internal.Green.Sprint(label)
	case "B":
		return internal.Cyan.Sprint(label)
	case "C":
		return internal.Yellow.Sprint(label)
	default:
		return internal.Red.Sprint(label)
	}
}

func formatHolderDelta(d int) string {
	if d > 0 {
		return internal.Green.Sprintf("+%d", d)
	}
	if d < 0 {
		return internal.Red.Sprintf("%d", d)
	}
	return internal.Dim.Sprint("0")
}

func init() {
	tokensListCmd.Flags().StringVar(&tokensSortBy, "sort-by", "volume_24h", "Sort field (volume_24h, price_change_24h, market_cap, health_score)")
	tokensListCmd.Flags().StringVar(&tokensOrder, "order", "desc", "Sort order (asc, desc)")

	tokensCmd.AddCommand(tokensListCmd)
	tokensCmd.AddCommand(tokensGetCmd)
	rootCmd.AddCommand(tokensCmd)
}
