package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BrunoPessoa22/fantokenintel-cli/internal"
	"github.com/spf13/cobra"
)

var (
	pricesHistory  bool
	pricesDays     int
	pricesInterval string
	pricesLimit    int
)

var pricesCmd = &cobra.Command{
	Use:   "prices <SYMBOL>",
	Short: "Current price or historical price data",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		symbol := strings.ToUpper(args[0])
		baseURL := internal.ResolveBaseURL(defaultBaseURL)
		c := internal.NewClient(baseURL, "")

		if !pricesHistory {
			return currentPrice(c, symbol)
		}
		return priceHistory(c, symbol)
	},
}

func currentPrice(c *internal.Client, symbol string) error {
	var resp struct {
		Token struct {
			Symbol string `json:"symbol"`
			Name   string `json:"name"`
		} `json:"token"`
		Metrics struct {
			Price          float64 `json:"price"`
			PriceChange1h  float64 `json:"price_change_1h"`
			PriceChange24h float64 `json:"price_change_24h"`
			PriceChange7d  float64 `json:"price_change_7d"`
			Volume24h      float64 `json:"volume_24h"`
		} `json:"metrics"`
	}

	raw, err := c.Get("/api/tokens/"+symbol, nil, &resp)
	if err != nil {
		return err
	}

	if jsonOut {
		internal.PrintJSON(raw)
		return nil
	}

	internal.Bold.Printf("\n%s  %s\n\n", resp.Token.Symbol, resp.Token.Name)
	fmt.Printf("  Price:   %s\n", internal.FormatPrice(resp.Metrics.Price))
	fmt.Printf("  1h:      %s\n", internal.FormatChange(resp.Metrics.PriceChange1h))
	fmt.Printf("  24h:     %s\n", internal.FormatChange(resp.Metrics.PriceChange24h))
	fmt.Printf("  7d:      %s\n", internal.FormatChange(resp.Metrics.PriceChange7d))
	fmt.Printf("  Vol 24h: %s\n", internal.FormatVolume(resp.Metrics.Volume24h))
	fmt.Println()
	return nil
}

func priceHistory(c *internal.Client, symbol string) error {
	q := buildQuery(map[string]string{
		"interval": pricesInterval,
		"days":     strconv.Itoa(pricesDays),
	})

	var resp struct {
		Symbol      string `json:"symbol"`
		PeriodHours int    `json:"period_hours"`
		DataPoints  int    `json:"data_points"`
		Prices      []struct {
			Time      string  `json:"time"`
			Price     float64 `json:"price"`
			Volume    float64 `json:"volume"`
			Spread    float64 `json:"spread"`
			Liquidity float64 `json:"liquidity"`
		} `json:"prices"`
	}

	raw, err := c.Get("/api/tokens/"+symbol+"/history", q, &resp)
	if err != nil {
		return err
	}

	if jsonOut {
		internal.PrintJSON(raw)
		return nil
	}

	internal.Bold.Printf("\n%s price history — last %d days (%s interval)\n\n", symbol, pricesDays, pricesInterval)

	t := internal.NewTable("TIME", "PRICE", "VOLUME", "SPREAD")
	t.Header()
	rows := resp.Prices
	if pricesLimit > 0 && len(rows) > pricesLimit {
		rows = rows[len(rows)-pricesLimit:]
	}
	for _, p := range rows {
		t.Row(
			internal.Dim.Sprint(shortTime(p.Time)),
			internal.FormatPrice(p.Price),
			internal.FormatVolume(p.Volume),
			fmt.Sprintf("%.1f bps", p.Spread),
		)
	}
	t.Flush()
	fmt.Printf("\n%d data points\n", resp.DataPoints)
	return nil
}

// shortTime trims the seconds from an ISO timestamp.
func shortTime(ts string) string {
	if len(ts) >= 16 {
		return ts[:16]
	}
	return ts
}

func init() {
	pricesCmd.Flags().BoolVar(&pricesHistory, "history", false, "Show historical price data")
	pricesCmd.Flags().IntVar(&pricesDays, "days", 7, "Number of days of history")
	pricesCmd.Flags().StringVar(&pricesInterval, "interval", "1h", "Candle interval (1h, 4h, 1d)")
	pricesCmd.Flags().IntVar(&pricesLimit, "limit", 0, "Max rows to display (0 = all)")

	rootCmd.AddCommand(pricesCmd)
}
