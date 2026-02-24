package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/fatih/color"
)

// Colors
var (
	Bold    = color.New(color.Bold)
	Dim     = color.New(color.Faint)
	Green   = color.New(color.FgGreen)
	Red     = color.New(color.FgRed)
	Yellow  = color.New(color.FgYellow)
	Cyan    = color.New(color.FgCyan)
	Magenta = color.New(color.FgMagenta)
	White   = color.New(color.FgWhite)
)

// PrintJSON pretty-prints raw JSON bytes.
func PrintJSON(data []byte) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		os.Stdout.Write(data)
		return
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(v) //nolint:errcheck
}

// Table is a simple tab-aligned table writer.
type Table struct {
	w       *tabwriter.Writer
	headers []string
}

// NewTable creates a Table writing to stdout with the given column headers.
func NewTable(headers ...string) *Table {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	return &Table{w: w, headers: headers}
}

// Header prints the header row with a separator line.
func (t *Table) Header() {
	Bold.Fprintln(t.w, strings.Join(t.headers, "\t"))
	seps := make([]string, len(t.headers))
	for i, h := range t.headers {
		seps[i] = strings.Repeat("─", len(h))
	}
	Dim.Fprintln(t.w, strings.Join(seps, "\t"))
}

// Row prints a data row. Values are tab-separated.
func (t *Table) Row(cols ...string) {
	fmt.Fprintln(t.w, strings.Join(cols, "\t"))
}

// Flush flushes the tabwriter buffer.
func (t *Table) Flush() {
	t.w.Flush()
}

// FormatChange formats a percentage change with color (+ green, - red).
func FormatChange(pct float64) string {
	s := fmt.Sprintf("%+.2f%%", pct)
	if pct > 0 {
		return Green.Sprint(s)
	}
	if pct < 0 {
		return Red.Sprint(s)
	}
	return s
}

// FormatPrice formats a USD price.
func FormatPrice(p float64) string {
	if p == 0 {
		return Dim.Sprint("—")
	}
	if p < 0.01 {
		return fmt.Sprintf("$%.6f", p)
	}
	if p < 1 {
		return fmt.Sprintf("$%.4f", p)
	}
	return fmt.Sprintf("$%.3f", p)
}

// FormatVolume formats a large USD volume with K/M suffix.
func FormatVolume(v float64) string {
	if v == 0 {
		return Dim.Sprint("—")
	}
	switch {
	case v >= 1_000_000:
		return fmt.Sprintf("$%.1fM", v/1_000_000)
	case v >= 1_000:
		return fmt.Sprintf("$%.1fK", v/1_000)
	default:
		return fmt.Sprintf("$%.0f", v)
	}
}

// FormatConfidence formats a 0-1 confidence as a coloured percentage.
func FormatConfidence(c float64) string {
	pct := fmt.Sprintf("%.0f%%", c*100)
	switch {
	case c >= 0.85:
		return Green.Sprint(pct)
	case c >= 0.70:
		return Yellow.Sprint(pct)
	default:
		return Red.Sprint(pct)
	}
}

// FormatDirection formats "short"/"long" with colour.
func FormatDirection(d string) string {
	switch strings.ToLower(d) {
	case "short":
		return Red.Sprint("short")
	case "long":
		return Green.Sprint("long")
	default:
		return d
	}
}

// FormatSide formats "buy"/"sell" with colour.
func FormatSide(s string) string {
	switch strings.ToLower(s) {
	case "sell":
		return Red.Sprint("sell")
	case "buy":
		return Green.Sprint("buy")
	default:
		return s
	}
}

// FormatOutcome formats signal outcomes with colour.
func FormatOutcome(o string) string {
	switch o {
	case "target_hit":
		return Green.Sprint("target_hit")
	case "stopped_out":
		return Red.Sprint("stopped_out")
	case "expired":
		return Dim.Sprint("expired")
	default:
		return Dim.Sprint(o)
	}
}

// TruncStr truncates a string to max length with ellipsis.
func TruncStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

// Fatal prints an error message to stderr and exits.
func Fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, Red.Sprint("error")+": "+format+"\n", args...)
	os.Exit(1)
}
