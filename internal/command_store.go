package internal

type commandstore struct {
}

type CommandStore interface {
	Get()
	Set()
}

func NewCommandStore() CommandStore {
	return &commandstore{}
}
