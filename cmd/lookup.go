// File: cmd/lookup.go
// Description: Contains the core logic for the lookup command.

package cmd

import (
	"fmt"
	"os"

	"godef/api"
	formatter "godef/int"
)

// executeLookup handles fetching data from both dictionary and thesaurus APIs.
func executeLookup(word string) {
	// Create channels to receive data from concurrent API calls.
	dictChan := make(chan []api.DictionaryResponse)
	thesChan := make(chan []api.ThesaurusResponse)
	errChan := make(chan error, 2) // Buffer for two potential errors

	// Fetch dictionary and thesaurus data concurrently.
	go func() {
		data, err := api.FetchDictionaryData(word)
		if err != nil {
			errChan <- fmt.Errorf("dictionary error: %w", err)
			return
		}
		dictChan <- data
	}()

	go func() {
		data, err := api.FetchThesaurusData(word)
		if err != nil {
			errChan <- fmt.Errorf("thesaurus error: %w", err)
			return
		}
		thesChan <- data
	}()

	var dictData []api.DictionaryResponse
	var thesData []api.ThesaurusResponse

	// Wait for both API calls to complete.
	for i := 0; i < 2; i++ {
		select {
		case data := <-dictChan:
			dictData = data
		case data := <-thesChan:
			thesData = data
		case err := <-errChan:
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	// Pass both results to the formatter for printing.
	formatter.PrintResults(dictData, thesData)
}
