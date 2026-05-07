// File: api/dictionary.go
// Description: Handles all interactions with the Merriam-Webster APIs and Datamuse.

package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

// DatamuseWord is a single result from the Datamuse API.
type DatamuseWord struct {
	Word string `json:"word"`
}

// DatamuseResult holds synonyms and antonyms fetched from Datamuse.
type DatamuseResult struct {
	Synonyms []string
	Antonyms []string
}

// FetchDictionaryData gets definitions from the dictionary API.
// When MW has no entry it returns a []string of suggestions; those are silently
// skipped so the caller receives an empty slice.
func FetchDictionaryData(word string) ([]DictionaryResponse, error) {
	if MW_DICTIONARY_API_KEY == "" {
		return nil, fmt.Errorf("dictionary API key is missing")
	}
	apiURL := fmt.Sprintf("%s%s?key=%s", DICTIONARY_API_URL, word, MW_DICTIONARY_API_KEY)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dictionary API returned status: %s", resp.Status)
	}

	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("could not decode dictionary JSON: %w", err)
	}

	var data []DictionaryResponse
	for _, item := range raw {
		var entry DictionaryResponse
		if err := json.Unmarshal(item, &entry); err != nil {
			continue
		}
		data = append(data, entry)
	}
	return data, nil
}

// FetchThesaurusData gets synonyms/antonyms from the thesaurus API.
// When MW has no entry for a word it returns a []string of suggestions instead of
// []ThesaurusResponse; in that case we return an empty slice so callers can fall
// back to another source.
func FetchThesaurusData(word string) ([]ThesaurusResponse, error) {
	if MW_THESAURUS_API_KEY == "" {
		return nil, fmt.Errorf("thesaurus API key is missing")
	}
	apiURL := fmt.Sprintf("%s%s?key=%s", THESAURUS_API_URL, word, MW_THESAURUS_API_KEY)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("thesaurus API returned status: %s", resp.Status)
	}

	var raw []json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("could not decode thesaurus JSON: %w", err)
	}

	var data []ThesaurusResponse
	for _, item := range raw {
		var entry ThesaurusResponse
		if err := json.Unmarshal(item, &entry); err != nil {
			// MW returned a suggestion string instead of an entry object — skip it.
			continue
		}
		data = append(data, entry)
	}
	return data, nil
}

// FetchDatamuseData fetches synonyms and antonyms from the Datamuse API (no key required).
func FetchDatamuseData(word string) (DatamuseResult, error) {
	encoded := url.QueryEscape(word)

	synResp, err := http.Get("https://api.datamuse.com/words?rel_syn=" + encoded)
	if err != nil {
		return DatamuseResult{}, err
	}
	defer synResp.Body.Close()

	var synWords []DatamuseWord
	if err := json.NewDecoder(synResp.Body).Decode(&synWords); err != nil {
		return DatamuseResult{}, fmt.Errorf("could not decode Datamuse synonyms: %w", err)
	}

	antResp, err := http.Get("https://api.datamuse.com/words?rel_ant=" + encoded)
	if err != nil {
		return DatamuseResult{}, err
	}
	defer antResp.Body.Close()

	var antWords []DatamuseWord
	if err := json.NewDecoder(antResp.Body).Decode(&antWords); err != nil {
		return DatamuseResult{}, fmt.Errorf("could not decode Datamuse antonyms: %w", err)
	}

	result := DatamuseResult{}
	for _, w := range synWords {
		result.Synonyms = append(result.Synonyms, w.Word)
	}
	for _, w := range antWords {
		result.Antonyms = append(result.Antonyms, w.Word)
	}
	return result, nil
}
