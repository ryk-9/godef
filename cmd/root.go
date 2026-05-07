// File: cmd/root.go
// Description: Sets up the root command for the CLI using Cobra.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "godef [word]",
	Short: "A dictionary and thesaurus tool using the Merriam-Webster API.",
	Long: `GoDef provides definitions and parts of speech for a given word.
Use -v to also show synonyms and antonyms.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		word := args[0]
		executeLookup(word, verbose)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show synonyms and antonyms")
}
