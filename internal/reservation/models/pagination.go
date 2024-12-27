package models

type PaginationParams struct {
	Page int `query:"page" validate:"required,gt=0"`
	Size int `query:"size" validate:"required,gt=0"`
}
