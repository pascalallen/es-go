package storage

type Projection interface {
	Name() string
	Script() string
}
