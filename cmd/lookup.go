// File: cmd/lookup.go
// Description: Contains the core logic for the lookup command.

package cmd

import (
	"fmt"
	"os"

	"godef/api"
	formatter "godef/int"
)

// executeLookup handles fetching data from the dictionary (and optionally thesaurus) APIs.
func executeLookup(word string, verbose bool) {
	if verbose {
		executeLookupVerbose(word)
	} else {
		executeLookupDefinition(word)
	}
}

func executeLookupDefinition(word string) {
	dictChan := make(chan []api.DictionaryResponse)
	errChan := make(chan error, 1)

	go func() {
		data, err := api.FetchDictionaryData(word)
		if err != nil {
			errChan <- fmt.Errorf("dictionary error: %w", err)
			return
		}
		dictChan <- data
	}()

	select {
	case data := <-dictChan:
		formatter.PrintResults(data, nil, api.DatamuseResult{})
	case err := <-errChan:
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func executeLookupVerbose(word string) {
	dictChan := make(chan []api.DictionaryResponse)
	thesChan := make(chan []api.ThesaurusResponse)
	datuChan := make(chan api.DatamuseResult)
	errChan := make(chan error, 3)

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

	go func() {
		data, err := api.FetchDatamuseData(word)
		if err != nil {
			errChan <- fmt.Errorf("datamuse error: %w", err)
			return
		}
		datuChan <- data
	}()

	var dictData []api.DictionaryResponse
	var thesData []api.ThesaurusResponse
	var datuData api.DatamuseResult

	for i := 0; i < 3; i++ {
		select {
		case data := <-dictChan:
			dictData = data
		case data := <-thesChan:
			thesData = data
		case data := <-datuChan:
			datuData = data
		case err := <-errChan:
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}

	formatter.PrintResults(dictData, thesData, datuData)
}
