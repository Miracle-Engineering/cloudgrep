package model

type Filter interface {
	String() string
}

type NoFilter struct {
}

func (f NoFilter) String() string {
	return "NoFilter"
}
