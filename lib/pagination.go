package lib

type PaginationData struct {
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPage  int               `json:"total_page"`
	TotalItems int64             `json:"total_items"`
	Links      map[string]string `json:"links"`
}

func Pagination(page, limit, totalPage int, totalItems int64, links map[string]string) *PaginationData {
	return &PaginationData{
		Page:       page,
		Limit:      limit,
		TotalPage:  totalPage,
		TotalItems: totalItems,
		Links:      links,
	}
}
