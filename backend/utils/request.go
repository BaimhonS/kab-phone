package utils

type QueryPagination struct {
	Search   string `query:"search"`
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
}
