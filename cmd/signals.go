package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Trading signals were retired in the descriptive-data pivot — Fan Token Intel is
// now a descriptive intelligence layer (data, not advice). This command is kept as
// a deprecation stub so existing scripts get a clear pointer instead of a silent
// empty result, and so old flags (--token, --days, ...) don't error.

var (
	signalsToken   string
	signalsMinConf float64
	signalsDays    int
	signalsOutcome string
	signalsLimit   int
)

func printSignalsDeprecation() {
	fmt.Println("Trading signals were retired. Fan Token Intel is now a descriptive")
	fmt.Println("data layer — it describes what's happening, it doesn't recommend trades.")
	fmt.Println()
	fmt.Println("Use instead:")
	fmt.Println("  fti tokens list        token health grades + prices")
	fmt.Println("  fti whales --all       whale-distribution flows")
	fmt.Println("  fti sports upcoming    matchday calendar")
	fmt.Println()
	fmt.Println("The premium event-impact moat (match-event → market-adjusted price")
	fmt.Println("reactions, with sample sizes and significance) is metered per call on")
	fmt.Println("the x402 gateway: https://x402.brunopessoa.com")
}

var signalsCmd = &cobra.Command{
	Use:   "signals",
	Short: "(deprecated) trading signals were retired — see descriptive tools",
	Run:   func(cmd *cobra.Command, args []string) { printSignalsDeprecation() },
}

var signalsActiveCmd = &cobra.Command{
	Use:    "active",
	Short:  "(deprecated)",
	Hidden: true,
	Run:    func(cmd *cobra.Command, args []string) { printSignalsDeprecation() },
}

var signalsHistoryCmd = &cobra.Command{
	Use:    "history",
	Short:  "(deprecated)",
	Hidden: true,
	Run:    func(cmd *cobra.Command, args []string) { printSignalsDeprecation() },
}

func init() {
	// Keep the old flags so legacy invocations parse cleanly before printing the notice.
	signalsActiveCmd.Flags().StringVar(&signalsToken, "token", "", "(deprecated)")
	signalsActiveCmd.Flags().Float64Var(&signalsMinConf, "min-confidence", 0.65, "(deprecated)")
	signalsHistoryCmd.Flags().StringVar(&signalsToken, "token", "", "(deprecated)")
	signalsHistoryCmd.Flags().IntVar(&signalsDays, "days", 30, "(deprecated)")
	signalsHistoryCmd.Flags().StringVar(&signalsOutcome, "outcome", "", "(deprecated)")
	signalsHistoryCmd.Flags().IntVar(&signalsLimit, "limit", 50, "(deprecated)")

	signalsCmd.AddCommand(signalsActiveCmd)
	signalsCmd.AddCommand(signalsHistoryCmd)
	rootCmd.AddCommand(signalsCmd)
}
