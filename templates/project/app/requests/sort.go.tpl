package requests

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
