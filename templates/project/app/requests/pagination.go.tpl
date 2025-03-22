package requests

import "github.com/samber/lo"

type Pager struct {
	Pagination `json:",inline"`
	Total      int64 `json:"total"`
	Items      any   `json:"items"`
}

type Pagination struct {
	Page  int64 `json:"page" form:"page" query:"page"`
	Limit int64 `json:"limit" form:"limit" query:"limit"`
}

func (filter *Pagination) Offset() int64 {
	return (filter.Page - 1) * filter.Limit
}

func (filter *Pagination) Format() *Pagination {
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if !lo.Contains([]int64{10, 20, 50, 100}, filter.Limit) {
		filter.Limit = 10
	}

	return filter
}
