package auditlogstore

type PagedResult[T any] struct {
	Results  []*T `json:"results"`
	Total    int  `json:"total"`
	Page     int  `json:"page"`
	PageSize int  `json:"pageSize"`
}
