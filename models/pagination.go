package models

type PaginationParams struct {
	OrderBy string `form:"order_by"`
	Offset  int    `form:"offset,default=0"`
	Limit   int    `form:"limit,default=50"`
}
