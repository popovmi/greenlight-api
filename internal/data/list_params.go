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

func (lp ListParams) limit() int {
	return lp.PageSize
}

func (lp ListParams) offset() int {
	return (lp.Page - 1) * lp.PageSize
}

func (lp ListParams) sortColumn() string {
	for _, safeValue := range lp.SortSafelist {
		if lp.Sort == safeValue {
			return strings.TrimPrefix(lp.Sort, "-")
		}
	}
	panic("unsafe sort parameter: " + lp.Sort)
}

func (lp ListParams) sortDirection() string {
	if strings.HasPrefix(lp.Sort, "-") {
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
