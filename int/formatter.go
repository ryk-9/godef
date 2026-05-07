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
func PrintResults(dictData []api.DictionaryResponse, thesData []api.ThesaurusResponse, datuData api.DatamuseResult) {
	if len(dictData) == 0 {
		fmt.Println("No definition results found.")
		return
	}

	word := strings.Split(dictData[0].Meta.ID, ":")[0]
	fmt.Printf("%s%s%s%s\n", ColorBold, ColorYellow, word, ColorReset)

	// --- Print Definitions ---
	fmt.Printf("\n%sDEFINITIONS%s\n", ColorBold, ColorReset)
	fmt.Println(strings.Repeat("─", 40))
	for _, entry := range dictData {
		if len(entry.Shortdef) > 0 {
			fmt.Printf("%s(%s)%s\n", ColorRed, entry.Fl, ColorReset)
			for i, def := range entry.Shortdef {
				fmt.Printf("  %d. %s\n", i+1, def)
			}
			fmt.Println()
		}
	}

	// --- Print Merriam-Webster Thesaurus ---
	mwHasThesaurus := false
	for _, entry := range thesData {
		if len(entry.Meta.Syns) > 0 || len(entry.Meta.Ants) > 0 {
			mwHasThesaurus = true
			break
		}
	}

	if mwHasThesaurus {
		fmt.Printf("\n%sTHESAURUS%s\n", ColorBold, ColorReset)
		fmt.Println(strings.Repeat("─", 40))
		for _, entry := range thesData {
			if len(entry.Meta.Syns) == 0 && len(entry.Meta.Ants) == 0 {
				continue
			}
			fmt.Printf("%s(%s)%s\n", ColorRed, entry.Fl, ColorReset)
			if len(entry.Meta.Syns) == 1 {
				fmt.Printf("  %sSynonyms:%s %s\n", ColorGreen, ColorReset, strings.Join(entry.Meta.Syns[0], ", "))
			} else {
				for i, group := range entry.Meta.Syns {
					fmt.Printf("  %sSynonyms (sense %d):%s %s\n", ColorGreen, i+1, ColorReset, strings.Join(group, ", "))
				}
			}
			if len(entry.Meta.Ants) == 1 {
				fmt.Printf("  %sAntonyms:%s %s\n", ColorBlue, ColorReset, strings.Join(entry.Meta.Ants[0], ", "))
			} else {
				for i, group := range entry.Meta.Ants {
					fmt.Printf("  %sAntonyms (sense %d):%s %s\n", ColorBlue, i+1, ColorReset, strings.Join(group, ", "))
				}
			}
			fmt.Println()
		}
	} else if len(datuData.Synonyms) > 0 || len(datuData.Antonyms) > 0 {
		// Fall back to Datamuse when MW thesaurus has nothing.
		fmt.Printf("\n%sTHESAURUS%s %s(via Datamuse)%s\n", ColorBold, ColorReset, ColorCyan, ColorReset)
		fmt.Println(strings.Repeat("─", 40))
		if len(datuData.Synonyms) > 0 {
			fmt.Printf("  %sSynonyms:%s %s\n", ColorGreen, ColorReset, strings.Join(datuData.Synonyms, ", "))
		}
		if len(datuData.Antonyms) > 0 {
			fmt.Printf("  %sAntonyms:%s %s\n", ColorBlue, ColorReset, strings.Join(datuData.Antonyms, ", "))
		}
		fmt.Println()
	}
}
