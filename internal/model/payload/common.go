package payload

import "math"

type ResponseData struct {
	Data any `json:"data"`
}

type ResponseList struct {
	Total  int64          `json:"total"`
	Params any            `json:"params"`
	Meta   PaginationMeta `json:"meta"`
	Items  any            `json:"items"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	NextPage    int `json:"next_page"`
	TotalPage   int `json:"total_page"`
}

type ResponseError[T any] struct {
	Error T `json:"error"`
}

func Ok(data any) ResponseData {
	return ResponseData{Data: data}
}

func Paginated(data any, count int64) ResponseList {
	return ResponseList{
		Items: data,
	}
}

const (
	defaultLimit = 25
	maxLimit     = 100
)

type PaginationFilter struct {
	Page   int    `json:"page" query:"page"`
	Limit  int    `json:"limit" query:"limit"`
	Search string `json:"search" query:"search"`
}

func (m *PaginationFilter) Normalize() {
	if m.Limit > maxLimit {
		m.Limit = maxLimit
	}

	if m.Limit == 0 {
		m.Limit = defaultLimit
	}

	if m.Page == 0 {
		m.Page = 1
	}
}

func (m *PaginationFilter) Paginate(items any, total int64) ResponseList {
	return ResponseList{
		Items:  items,
		Total:  total,
		Params: m,
		Meta: PaginationMeta{
			CurrentPage: m.Page,
			NextPage:    m.Page + 1,
			TotalPage: func() int {
				if total == 0 {
					return 1
				}
				return int(math.Ceil(float64(total) / float64(m.Limit)))
			}(),
		},
	}
}
