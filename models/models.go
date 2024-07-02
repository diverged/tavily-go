package models

// SearchRequest represents the structure for a search request
type SearchRequest struct {
	Query             string   `json:"query"`
	SearchDepth       string   `json:"search_depth,omitempty"`
	Topic             string   `json:"topic,omitempty"`
	IncludeAnswer     bool     `json:"include_answer,omitempty"`
	IncludeRawContent bool     `json:"include_raw_content,omitempty"`
	MaxResults        int      `json:"max_results,omitempty"`
	IncludeDomains    []string `json:"include_domains,omitempty"`
	ExcludeDomains    []string `json:"exclude_domains,omitempty"`
	IncludeImages     bool     `json:"include_images,omitempty"`
	APIKey            string   `json:"api_key"`
	UseCache          bool     `json:"use_cache,omitempty"`
}

// SearchResult represents the structure for a search result
type SearchResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Content string  `json:"content"`
	Score   float64 `json:"score"`
}

// SearchResponse represents the structure for a search response
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Answer  string         `json:"answer,omitempty"`
}
