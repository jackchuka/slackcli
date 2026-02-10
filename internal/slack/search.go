package slack

// Search is handled by SearchMessages in messages.go
// This file provides additional search-specific types if needed.

type SearchFilesResult struct {
	Matches []File `json:"matches"`
	Total   int    `json:"total"`
}
