// File: api/dictionary.go
// Description: Handles all interactions with the Merriam-Webster APIs.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// API Keys provided by the user.
const MW_DICTIONARY_API_KEY = "bf1709d5-50fc-4eb0-bc63-0622d7112d1a"
const MW_THESAURUS_API_KEY = "39a7934d-4ebd-4d3f-8b6d-9d9a7590979a"

// Corrected API URLs.
const DICTIONARY_API_URL = "https://www.dictionaryapi.com/api/v3/references/collegiate/json/"
const THESAURUS_API_URL = "https://www.dictionaryapi.com/api/v3/references/thesaurus/json/"

// DictionaryResponse defines the structs for the Collegiate Dictionary API.
type DictionaryResponse struct {
	Meta     Meta     `json:"meta"`
	Fl       string   `json:"fl"` // Functional Label (part of speech)
	Shortdef []string `json:"shortdef"`
}

// ThesaurusResponse defines the structs for the Collegiate Thesaurus API.
type ThesaurusResponse struct {
	Meta Meta   `json:"meta"`
	Fl   string `json:"fl"`
	Def  []Def  `json:"def"`
}

type Meta struct {
	ID    string     `json:"id"`
	UUID  string     `json:"uuid"`
	Src   string     `json:"src"`
	Stems []string   `json:"stems"`
	Ants  [][]string `json:"ants"`
	Syns  [][]string `json:"syns"`
}

type Def struct {
	Sseq [][][]interface{} `json:"sseq"`
}

// FetchDictionaryData gets definitions from the dictionary API.
func FetchDictionaryData(word string) ([]DictionaryResponse, error) {
	if MW_DICTIONARY_API_KEY == "" {
		return nil, fmt.Errorf("dictionary API key is missing")
	}
	url := fmt.Sprintf("%s%s?key=%s", DICTIONARY_API_URL, word, MW_DICTIONARY_API_KEY)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dictionary API returned status: %s", resp.Status)
	}

	var data []DictionaryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("could not decode dictionary JSON: %w", err)
	}
	return data, nil
}

// FetchThesaurusData gets synonyms/antonyms from the thesaurus API.
func FetchThesaurusData(word string) ([]ThesaurusResponse, error) {
	if MW_THESAURUS_API_KEY == "" {
		return nil, fmt.Errorf("thesaurus API key is missing")
	}
	url := fmt.Sprintf("%s%s?key=%s", THESAURUS_API_URL, word, MW_THESAURUS_API_KEY)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("thesaurus API returned status: %s", resp.Status)
	}

	var data []ThesaurusResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("could not decode thesaurus JSON: %w", err)
	}
	return data, nil
}
