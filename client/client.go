package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"sync"

	"github.com/diverged/tavily-go/models"
	"github.com/diverged/tavily-go/utils"
)

// TavilyClient represents a client for the Tavily API
type TavilyClient struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewTavilyClient creates a new TavilyClient
func NewTavilyClient(apiKey string) *TavilyClient {
	return &TavilyClient{
		BaseURL:    "https://api.tavily.com/search",
		APIKey:     apiKey,
		HTTPClient: &http.Client{},
	}
}

// Search performs a search query
func (c *TavilyClient) Search(req models.SearchRequest) (*models.SearchResponse, error) {
	req.APIKey = c.APIKey

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var searchResp models.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return &searchResp, nil
}

// GetSearchContext retrieves the search context for a query
func (c *TavilyClient) GetSearchContext(query string, searchDepth string, maxTokens int) (string, error) {
	req := models.SearchRequest{
		Query:       query,
		SearchDepth: searchDepth,
	}

	resp, err := c.Search(req)
	if err != nil {
		return "", err
	}

	context := make([]map[string]string, len(resp.Results))
	for i, result := range resp.Results {
		context[i] = map[string]string{
			"url":     result.URL,
			"content": result.Content,
		}
	}

	jsonContext, err := utils.GetMaxItemsFromList(interfaceSlice(context), maxTokens)
	if err != nil {
		return "", fmt.Errorf("error getting max items from list: %w", err)
	}

	return jsonContext, nil
}

// QnASearch performs a Q&A search
func (c *TavilyClient) QnASearch(query string, searchDepth string) (string, error) {
	req := models.SearchRequest{
		Query:         query,
		SearchDepth:   searchDepth,
		IncludeAnswer: true,
	}

	resp, err := c.Search(req)
	if err != nil {
		return "", err
	}

	return resp.Answer, nil
}

// GetCompanyInfo retrieves company information
func (c *TavilyClient) GetCompanyInfo(query string, searchDepth string, maxResults int) ([]models.SearchResult, error) {
	topics := []string{"news", "general", "finance"}
	var wg sync.WaitGroup
	resultsChan := make(chan []models.SearchResult, len(topics))
	errorsChan := make(chan error, len(topics))

	for _, topic := range topics {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			req := models.SearchRequest{
				Query:       query,
				SearchDepth: searchDepth,
				Topic:       t,
				MaxResults:  maxResults,
			}
			resp, err := c.Search(req)
			if err != nil {
				errorsChan <- err
				return
			}
			resultsChan <- resp.Results
		}(topic)
	}

	wg.Wait()
	close(resultsChan)
	close(errorsChan)

	if len(errorsChan) > 0 {
		return nil, <-errorsChan
	}

	var allResults []models.SearchResult
	seenURLs := make(map[string]bool)

	for results := range resultsChan {
		for _, result := range results {
			if !seenURLs[result.URL] {
				allResults = append(allResults, result)
				seenURLs[result.URL] = true
			}
		}
	}

	// Sort results by score (descending)
	sort.Slice(allResults, func(i, j int) bool {
		return allResults[i].Score > allResults[j].Score
	})

	// Limit to maxResults
	if len(allResults) > maxResults {
		allResults = allResults[:maxResults]
	}

	return allResults, nil
}

// Helper function to convert a slice of maps to a slice of interface{}
func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
