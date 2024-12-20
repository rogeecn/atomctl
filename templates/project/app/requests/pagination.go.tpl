package requests

type Pager struct {
	Pagination `json:",inline"`
	Total      int64       `json:"total"`
	Items      interface{} `json:"items"`
}

type Pagination struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

func (filter *PageQueryFilter) Offset() int {
	return (filter.Page - 1) * filter.Limit
}

func (filter *PageQueryFilter) Format() *PageQueryFilter {
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	if filter.Limit > 50 {
		filter.Limit = 50
	}
	return filter
}
