// Package cmd provides shared helper functions used across subcommands.
package cmd

import "sort"

// sortedMapKeys returns the keys of a string-keyed map in alphabetical order.
func sortedMapKeys(m map[string]any) []string {
	result := make([]string, 0, len(m))
	for k := range m {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}
