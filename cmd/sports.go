package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/BrunoPessoa22/fantokenintel-cli/internal"
	"github.com/spf13/cobra"
)

var (
	sportsToken string
	sportsDays  int
)

var sportsCmd = &cobra.Command{
	Use:   "sports",
	Short: "Upcoming sports matches for fan token teams",
}

var sportsUpcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "List upcoming matches with token context",
	RunE: func(cmd *cobra.Command, args []string) error {
		baseURL := internal.ResolveBaseURL(defaultBaseURL)
		c := internal.NewClient(baseURL, "")

		params := map[string]string{
			"days":  strconv.Itoa(sportsDays),
			"limit": "100",
		}
		if sportsToken != "" {
			params["token"] = strings.ToUpper(sportsToken)
		}
		q := buildQuery(params)

		var resp struct {
			Count       int    `json:"count"`
			Days        int    `json:"days"`
			TokenFilter string `json:"token_filter"`
			Matches     []struct {
				MatchID         string  `json:"match_id"`
				HomeTeam        string  `json:"home_team"`
				AwayTeam        string  `json:"away_team"`
				MatchDate       string  `json:"match_date"`
				Competition     string  `json:"competition"`
				Status          string  `json:"status"`
				ImportanceScore float64 `json:"importance_score"`
				HomeToken       string  `json:"home_token"`
				AwayToken       string  `json:"away_token"`
			} `json:"matches"`
		}

		raw, err := c.Get("/api/matches/upcoming", q, &resp)
		if err != nil {
			return err
		}

		if jsonOut {
			internal.PrintJSON(raw)
			return nil
		}

		filter := "all tokens"
		if sportsToken != "" {
			filter = strings.ToUpper(sportsToken)
		}

		internal.Bold.Printf("\nUpcoming matches — %s  (next %d days)\n\n", filter, sportsDays)

		if resp.Count == 0 {
			internal.Dim.Println("No upcoming matches found.")
			return nil
		}

		t := internal.NewTable("DATE", "HOME", "AWAY", "COMPETITION", "TOKENS", "IMP")
		t.Header()
		for _, m := range resp.Matches {
			tokens := tokenPair(m.HomeToken, m.AwayToken)
			t.Row(
				shortTime(m.MatchDate),
				internal.TruncStr(m.HomeTeam, 18),
				internal.TruncStr(m.AwayTeam, 18),
				internal.TruncStr(m.Competition, 18),
				tokens,
				importanceBar(m.ImportanceScore),
			)
		}
		t.Flush()
		fmt.Printf("\n%d match(es)\n", resp.Count)
		return nil
	},
}

func tokenPair(home, away string) string {
	parts := []string{}
	if home != "" {
		parts = append(parts, internal.Cyan.Sprint(home))
	}
	if away != "" {
		parts = append(parts, internal.Cyan.Sprint(away))
	}
	if len(parts) == 0 {
		return internal.Dim.Sprint("—")
	}
	return strings.Join(parts, " / ")
}

func importanceBar(score float64) string {
	switch {
	case score >= 80:
		return internal.Green.Sprintf("%.0f ●●●", score)
	case score >= 50:
		return internal.Yellow.Sprintf("%.0f ●●○", score)
	default:
		return internal.Dim.Sprintf("%.0f ●○○", score)
	}
}

func init() {
	sportsUpcomingCmd.Flags().StringVar(&sportsToken, "token", "", "Filter by token symbol (e.g. PSG)")
	sportsUpcomingCmd.Flags().IntVar(&sportsDays, "days", 14, "Look-ahead window in days")

	sportsCmd.AddCommand(sportsUpcomingCmd)
	rootCmd.AddCommand(sportsCmd)
}
