package data

import "greenlight.aenkas.org/internal/validator"

type ListParams struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

func ValidateListParams(v *validator.Validator, p ListParams) {
	v.Check(p.Page > 0, "page", "must be greater than zero")
	v.Check(p.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(p.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(p.PageSize <= 100, "page_size", "must be a maximum of 100")
	v.Check(validator.In(p.Sort, p.SortSafelist...), "sort", "invalid sort value")
}
