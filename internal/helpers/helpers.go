package helpers

func Ptr[T any](v T) *T {
	return &v
}

type Number interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 |
		uint32 | int64 | uint64 | float32 | float64
}

func NilIfZero[T Number](v T) *T {
	if v == 0 {
		return nil
	}

	return &v
}
