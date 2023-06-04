package data

import (
	"math"
	"strings"

	"greenlight.aenkas.org/internal/validator"
)

type ListParams struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

type Metadata struct {
	CurrentPage  int `json:"currentPage,omitempty"`
	PageSize     int `json:"pageSize,omitempty"`
	FirstPage    int `json:"firstPage,omitempty"`
	LastPage     int `json:"lastPage,omitempty"`
	TotalRecords int `json:"totalRecords,omitempty"`
}

func getMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

func (self ListParams) limit() int {
	return self.PageSize
}

func (self ListParams) offset() int {
	return (self.Page - 1) * self.PageSize
}

func (self ListParams) sortColumn() string {
	for _, safeValue := range self.SortSafelist {
		if self.Sort == safeValue {
			return strings.TrimPrefix(self.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + self.Sort)
}

func (self ListParams) sortDirection() string {
	if strings.HasPrefix(self.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func ValidateListParams(v *validator.Validator, p ListParams) {
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(p.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(p.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.In(p.Sort, p.SortSafelist...), "sort", "invalid sort value")
}
