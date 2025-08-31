package pagination

type Meta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

func BuildMeta(page, pageSize int, total int64) Meta {
	totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
	return Meta{Page: page, PageSize: pageSize, Total: total, TotalPages: totalPages}
}
