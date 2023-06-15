package common

type SortQueryFilter struct {
	Asc  *string `json:"asc"`
	Desc *string `json:"desc"`
}

func (s *SortQueryFilter) AscFields() []string {
	if s.Asc == nil {
		return nil
	}
	return strings.Split(*s.Asc, ",")
}

func (s *SortQueryFilter) DescFields() []string {
	if s.Desc == nil {
		return nil
	}
	return strings.Split(*s.Desc, ",")
}

type PageDataResponse struct {
	PageQueryFilter `json:",inline"`
	Total           int64       `json:"total"`
	Items           interface{} `json:"items"`
}

type PageQueryFilter struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
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
