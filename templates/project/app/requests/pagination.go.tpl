package requests

import "github.com/samber/lo"

type Pager struct {
	Pagination `json:",inline"`
	Total      int64       `json:"total"`
	Items      interface{} `json:"items"`
}

type Pagination struct {
	Page  int `json:"page" form:"page" query:"page"`
	Limit int `json:"limit" form:"limit" query:"limit"`
}

func (filter *Pagination) Offset() int {
	return (filter.Page - 1) * filter.Limit
}

func (filter *Pagination) Format() *Pagination {
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if !lo.Contains([]int{10, 20, 50, 100}, filter.Limit) {
		filter.Limit = 10
	}

	return filter
}
