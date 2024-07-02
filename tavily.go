package tavilygo

import (
	"github.com/diverged/tavily-go/client"
	"github.com/diverged/tavily-go/models"
)

// NewClient creates a new TavilyClient with the given API key.
func NewClient(apiKey string) *client.TavilyClient {
	return client.NewTavilyClient(apiKey)
}

// Search performs a search query using the Tavily API and returns the search response.
func Search(c *client.TavilyClient, req models.SearchRequest) (*models.SearchResponse, error) {
	return c.Search(req)
}

// GetSearchContext retrieves the search context for a query, limiting the response to the specified maximum number of tokens.
func GetSearchContext(c *client.TavilyClient, query string, searchDepth string, maxTokens int) (string, error) {
	return c.GetSearchContext(query, searchDepth, maxTokens)
}

// QnASearch performs a Q&A search using the Tavily API and returns the answer.
func QnASearch(c *client.TavilyClient, query string, searchDepth string) (string, error) {
	return c.QnASearch(query, searchDepth)
}

// GetCompanyInfo retrieves company information by performing searches across multiple topics and returns a sorted list of unique search results.
func GetCompanyInfo(c *client.TavilyClient, query string, searchDepth string, maxResults int) ([]models.SearchResult, error) {
	return c.GetCompanyInfo(query, searchDepth, maxResults)
}
