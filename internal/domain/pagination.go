package domain

const (
	PageLimitMax     = 120
	PageLimitMin     = 1
	PageLimitDefault = 12
)

type PaginationQuery struct {
	Limit int `json:"limit" validate:"min=1,max=100"`
	Page  int `json:"page" validate:"min=1"`
}

type PaginatedResults struct {
	PaginationQuery
	Total int            `json:"total"`
	Data  *[]interface{} `json:"data"`
}

func (pq *PaginationQuery) GetOffset() int {
	return pq.Limit * (pq.Page - 1)
}

func (pi *PaginationQuery) CreateResults() (PaginatedResults, int) {
	var results PaginatedResults
	results.PaginationQuery = *pi
	offset := pi.GetOffset()
	return results, offset
}

type IPaginationService interface {
	Paginate(model any, results *PaginatedResults, offset int, filters ...Filter) error
}
