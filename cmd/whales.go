package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/BrunoPessoa22/fantokenintel-cli/internal"
	"github.com/spf13/cobra"
)

var (
	whalesAll      bool
	whalesHours    int
	whalesLimit    int
	whalesMinValue float64
	whalesWatch    bool
	whalesInterval int
)

var whalesCmd = &cobra.Command{
	Use:   "whales [SYMBOL]",
	Short: "CEX + DEX whale trade activity",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL := internal.ResolveBaseURL(defaultBaseURL)
		c := internal.NewClient(baseURL, "")

		symbol := ""
		if !whalesAll && len(args) > 0 {
			symbol = strings.ToUpper(args[0])
		}

		if !whalesWatch {
			return whalesCombined(c, symbol)
		}

		// Watch mode: poll on a ticker, clear screen between updates.
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

		ticker := time.NewTicker(time.Duration(whalesInterval) * time.Second)
		defer ticker.Stop()

		clearScreen()
		if err := whalesCombined(c, symbol); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
		}
		internal.Dim.Printf("\n  Refreshing every %ds — Ctrl+C to stop\n", whalesInterval)

		for {
			select {
			case <-sig:
				fmt.Println()
				return nil
			case <-ticker.C:
				clearScreen()
				if err := whalesCombined(c, symbol); err != nil {
					fmt.Fprintln(os.Stderr, "error:", err)
				}
				internal.Dim.Printf("\n  Refreshing every %ds — Ctrl+C to stop\n", whalesInterval)
			}
		}
	},
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

type whaleTrade struct {
	Time         string  `json:"time"`
	Venue        string  `json:"venue"`
	Symbol       string  `json:"symbol"`
	Exchange     string  `json:"exchange"`
	Side         string  `json:"side"`
	Price        float64 `json:"price"`
	Quantity     float64 `json:"quantity"`
	ValueUSD     float64 `json:"value_usd"`
	IsAggressive bool    `json:"is_aggressive"`
	TxHash       string  `json:"tx_hash"`
}

func whalesCombined(c *internal.Client, symbol string) error {
	params := map[string]string{
		"limit":     strconv.Itoa(whalesLimit),
		"min_value": fmt.Sprintf("%.0f", whalesMinValue),
		"hours":     strconv.Itoa(whalesHours),
	}
	if symbol != "" {
		params["symbol"] = symbol
	}
	q := buildQuery(params)

	var resp struct {
		Transactions []whaleTrade `json:"transactions"`
		Count        int          `json:"count"`
		CexCount     int          `json:"cex_count"`
		DexCount     int          `json:"dex_count"`
		Threshold    float64      `json:"threshold_usd"`
		Timestamp    string       `json:"timestamp"`
	}

	raw, err := c.Get("/api/whales/combined", q, &resp)
	if err != nil {
		return err
	}

	if jsonOut {
		internal.PrintJSON(raw)
		return nil
	}

	filter := "all tokens"
	if symbol != "" {
		filter = symbol
	}

	internal.Bold.Printf("\nWhale trades — %s  (min %s, last %dh)\n\n",
		filter, internal.FormatVolume(whalesMinValue), whalesHours)

	if resp.Count == 0 {
		internal.Dim.Println("No whale trades found.")
		return nil
	}

	t := internal.NewTable("TIME", "VENUE", "TOKEN", "EXCHANGE", "SIDE", "PRICE", "QTY", "VALUE")
	t.Header()
	for _, tr := range resp.Transactions {
		aggressiveFlag := ""
		if tr.IsAggressive {
			aggressiveFlag = internal.Yellow.Sprint(" *")
		}
		t.Row(
			internal.Dim.Sprint(shortTime(tr.Time)),
			strings.ToUpper(tr.Venue),
			internal.Cyan.Sprint(tr.Symbol),
			tr.Exchange,
			internal.FormatSide(tr.Side)+aggressiveFlag,
			internal.FormatPrice(tr.Price),
			formatQty(tr.Quantity),
			internal.FormatVolume(tr.ValueUSD),
		)
	}
	t.Flush()
	fmt.Printf("\n%d trades  CEX:%d  DEX:%d  (*)=aggressive\n", resp.Count, resp.CexCount, resp.DexCount)
	return nil
}

func formatQty(q float64) string {
	if q >= 1_000_000 {
		return fmt.Sprintf("%.1fM", q/1_000_000)
	}
	if q >= 1_000 {
		return fmt.Sprintf("%.1fK", q/1_000)
	}
	return fmt.Sprintf("%.0f", q)
}

func init() {
	whalesCmd.Flags().BoolVar(&whalesAll, "all", false, "Show whales for all tokens")
	whalesCmd.Flags().IntVar(&whalesHours, "hours", 24, "Look-back window in hours")
	whalesCmd.Flags().IntVar(&whalesLimit, "limit", 50, "Max trades to show")
	whalesCmd.Flags().Float64Var(&whalesMinValue, "min-value", 50000, "Minimum trade value in USD")
	whalesCmd.Flags().BoolVar(&whalesWatch, "watch", false, "Poll and refresh continuously")
	whalesCmd.Flags().IntVar(&whalesInterval, "interval", 30, "Refresh interval in seconds (with --watch)")

	rootCmd.AddCommand(whalesCmd)
}
