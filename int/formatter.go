// File: internal/formatter.go
// Description: Formats and prints the API data to the console with colors.

package formatter

import (
	"fmt"
	"strings"

	"godef/api"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorBold   = "\033[1m"
)

// PrintResults formats and prints the combined dictionary and thesaurus data.
func PrintResults(dictData []api.DictionaryResponse, thesData []api.ThesaurusResponse) {
	if len(dictData) == 0 {
		fmt.Println("No definition results found.")
		return
	}

	// --- Print Header ---
	// The word might have variations (e.g., initial:1), clean it up.
	word := strings.Split(dictData[0].Meta.ID, ":")[0]
	fmt.Printf("%s%s%s%s\n", ColorBold, ColorYellow, word, ColorReset)

	// --- Print Definitions ---
	fmt.Printf("\n%sDEFINITIONS%s\n", ColorBold, ColorReset)
	fmt.Println(strings.Repeat("─", 40))
	for _, entry := range dictData {
		// Check if entry is a valid definition (not just a suggestion list)
		if len(entry.Shortdef) > 0 {
			fmt.Printf("%s(%s)%s\n", ColorRed, entry.Fl, ColorReset)
			for i, def := range entry.Shortdef {
				fmt.Printf("  %d. %s\n", i+1, def)
			}
			fmt.Println()
		}
	}

	// --- Print Thesaurus Info ---
	if len(thesData) > 0 && len(thesData[0].Meta.Syns) > 0 {
		fmt.Printf("\n%sTHESAURUS%s\n", ColorBold, ColorReset)
		fmt.Println(strings.Repeat("─", 40))
		for _, entry := range thesData {
			if len(entry.Meta.Syns) > 0 || len(entry.Meta.Ants) > 0 {
				fmt.Printf("%s(%s)%s\n", ColorRed, entry.Fl, ColorReset)
				if len(entry.Meta.Syns) > 0 {
					fmt.Printf("  %sSynonyms:%s %s\n", ColorGreen, ColorReset, strings.Join(entry.Meta.Syns[0], ", "))
				}
				if len(entry.Meta.Ants) > 0 {
					fmt.Printf("  %sAntonyms:%s %s\n", ColorBlue, ColorReset, strings.Join(entry.Meta.Ants[0], ", "))
				}
				fmt.Println()
			}
		}
	}
}
