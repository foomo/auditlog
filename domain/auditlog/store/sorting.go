package auditlogstore

type SortField string

const (
	SortFieldTimestamp SortField = "timestamp"
	SortFieldService   SortField = "service"
	SortFieldFunc      SortField = "func"
	SortFieldAction    SortField = "action"
	SortFieldUserID    SortField = "userId"
)

type Direction string

const (
	DirectionAscending  Direction = "ascending"
	DirectionDescending Direction = "descending"
)

type Sort struct {
	Field     SortField `json:"field"`
	Direction Direction `json:"direction"`
}

// GetSortValue maps the direction to mongo's sort int. Descending is the default
// because most audit log readers want newest entries first.
func (d Direction) GetSortValue() int {
	switch d {
	case DirectionAscending:
		return 1
	case DirectionDescending:
		return -1
	default:
		return -1
	}
}
