package forgo

type Equals[T any] interface {
	Equal(T) bool
}

func Equal[T Equals[T]](t ...T) bool {
	if len(t) == 0 {
		return false
	}
	t0 := t[0]
	for i, _ := range t[1:] {
		if !t[i].Equal(t0) {
			return false
		}
	}
	return true
}
