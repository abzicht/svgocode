package forgo

type Clonable[T any] interface {
	Clone() T
}

func Clone[T Clonable[T]](t []T) []T {
	t2 := make([]T, len(t))
	for i, _ := range t {
		t2[i] = t[i].Clone()
	}
	return t2
}
