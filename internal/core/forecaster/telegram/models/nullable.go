package models

type Nullable[T any] struct {
	Value   T
	Defined bool
}

func NewNullable[T any](value T) Nullable[T] {
	return Nullable[T]{Value: value, Defined: true}
}
