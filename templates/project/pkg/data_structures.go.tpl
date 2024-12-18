package pkg

import (
	"strings"

	"github.com/samber/lo"
)

type SortQueryFilter struct {
	Asc  *string `json:"asc" form:"asc"`
	Desc *string `json:"desc" form:"desc"`
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

func (s *SortQueryFilter) DescID() *SortQueryFilter {
	if s.Desc == nil {
		s.Desc = lo.ToPtr("id")
	}

	items := s.DescFields()
	if lo.Contains(items, "id") {
		return s
	}

	items = append(items, "id")
	s.Desc = lo.ToPtr(strings.Join(items, ","))
	return s
}

type PageDataResponse struct {
	PageQueryFilter `json:",inline"`
	Total           int64       `json:"total"`
	Items           interface{} `json:"items"`
}

type PageQueryFilter struct {
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
