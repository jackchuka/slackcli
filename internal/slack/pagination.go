package slack

type PaginationParams struct {
	Cursor string
	Limit  int
	All    bool
}

func (p PaginationParams) EffectiveLimit() int {
	if p.Limit > 0 {
		return p.Limit
	}
	return 100
}

type PaginatedResult[T any] struct {
	Items      []T    `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
}
