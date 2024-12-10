package functional

// Types
type MapFunc[I, O any] func(I) O
type FilterFunc[I any] func(I) bool
type ReduceFunc[I any] func(value, acc I) I


// Map
func Map[S ~[]I, I, O any](s S, fn MapFunc[I, O]) <-chan O {
    ch := make(chan O, len(s))
    go func() {
        for _, value := range s {
            ch <- fn(value)
        }
        close(ch)
    }()
    return ch
}

func MapSat[S ~[]I, I any](s S, condition func(I) bool) <-chan int {
    return Map(s, func(value I) int {
        return map[bool]int{true: 1, false: 0}[condition(value)]
    })
}


// Filter
func Filter[S ~[]I, I any](s S, fn FilterFunc[I]) S {
    filtered := S{}
    for _, value := range s {
        if fn(value) {
            filtered = append(filtered, value)
        }
    }
    return filtered
}

func CountSat[S ~[]I, I any](s S, coniditon func(I) bool) int {
    return len(Filter(s, coniditon))
}

func AnySat[S ~[]I, I any](s S, condition func(I) bool) bool {
    return CountSat(s, condition) > 0
}


// Reduce
func Reduce[I any](ch <-chan I, initial I, fn ReduceFunc[I]) I {
	var acc I = initial
	for value := range ch {
		acc = fn(acc, value)
	}
	return acc
}

func ReduceSum(ch <-chan int) int {
    return Reduce(ch, 0, func(acc, value int) int {
        return acc + value
    })
}
