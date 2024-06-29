package tavilygo

import (
	"github.com/diverged/tavily-go/client"
	"github.com/diverged/tavily-go/models"
)

// NewClient creates a new TavilyClient
func NewClient(apiKey string) *client.TavilyClient {
	return client.NewTavilyClient(apiKey)
}

// Search performs a search query
func Search(c *client.TavilyClient, req models.SearchRequest) (*models.SearchResponse, error) {
	return c.Search(req)
}

// GetSearchContext retrieves the search context for a query
func GetSearchContext(c *client.TavilyClient, query string, searchDepth string, maxTokens int) (string, error) {
	return c.GetSearchContext(query, searchDepth, maxTokens)
}

// QnASearch performs a Q&A search
func QnASearch(c *client.TavilyClient, query string, searchDepth string) (string, error) {
	return c.QnASearch(query, searchDepth)
}

// GetCompanyInfo retrieves company information
func GetCompanyInfo(c *client.TavilyClient, query string, searchDepth string, maxResults int) ([]models.SearchResult, error) {
	return c.GetCompanyInfo(query, searchDepth, maxResults)
}
