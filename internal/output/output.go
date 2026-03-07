package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

// DisplayMode defines how results are rendered to the terminal.
type DisplayMode string

const (
	// ModeTable renders aligned columns using tab-separated layout.
	ModeTable DisplayMode = "table"
	// ModeJSON renders data as pretty-printed JSON.
	ModeJSON DisplayMode = "json"
	// ModeCSV renders comma-separated values.
	ModeCSV DisplayMode = "csv"
)

// Render writes the provided records to stdout in the specified display mode.
// columns controls the ordering of fields; records holds the data as a slice of maps.
func Render(mode DisplayMode, columns []string, records []map[string]any) error {
	switch mode {
	case ModeJSON:
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(records)

	case ModeCSV:
		writer := csv.NewWriter(os.Stdout)
		if writeErr := writer.Write(columns); writeErr != nil {
			return writeErr
		}
		for _, rec := range records {
			cells := make([]string, len(columns))
			for idx, col := range columns {
				cells[idx] = fmt.Sprint(rec[col])
			}
			if writeErr := writer.Write(cells); writeErr != nil {
				return writeErr
			}
		}
		writer.Flush()
		return writer.Error()

	default:
		// Fall back to a human-readable tabular layout.
		tw := tabwriter.NewWriter(os.Stdout, 0, 4, 2, ' ', 0)
		fmt.Fprintln(tw, strings.Join(columns, "\t"))
		for _, rec := range records {
			entries := make([]string, len(columns))
			for idx, col := range columns {
				entries[idx] = fmt.Sprint(rec[col])
			}
			fmt.Fprintln(tw, strings.Join(entries, "\t"))
		}
		return tw.Flush()
	}
}

// SelectColumns narrows each record to only the specified keys.
// When no keys are provided the original records are returned unchanged.
func SelectColumns(records []map[string]any, keys []string) []map[string]any {
	if len(keys) == 0 {
		return records
	}
	filtered := make([]map[string]any, 0, len(records))
	for _, rec := range records {
		subset := make(map[string]any, len(keys))
		for _, k := range keys {
			if val, exists := rec[k]; exists {
				subset[k] = val
			}
		}
		filtered = append(filtered, subset)
	}
	return filtered
}
