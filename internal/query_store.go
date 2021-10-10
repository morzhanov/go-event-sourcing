package internal

type querystore struct {
}

type QueryStore interface {
	Get()
	Set()
}

func NewQueryStore() QueryStore {
	return &querystore{}
}
