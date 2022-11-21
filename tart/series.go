package tart

type Series []any

func NewSeries(a ...any) Series {
	return Series(a)
}

func (s *Series) Push(v any) {
	*s = append(*s, v)
}

func (s *Series) Update(v any) {
	*s = append(*s, v)
}

func (s *Series) Pop(i int64) (v any) {
	v = (*s)[i]
	*s = append((*s)[:i], (*s)[i+1:]...)
	return v
}

func (s Series) Tail(size int) Series {
	length := len(s)
	if length <= size {
		win := make(Series, length)
		copy(win, s)
		return win
	}

	win := make(Series, size)
	copy(win, s[length-size:])
	return win
}
func (s *Series) Last() any {
	length := len(*s)
	if length > 0 {
		return (*s)[length-1]
	}
	return ""
}

func (s *Series) Index(i int) any {
	length := len(*s)
	if length-i <= 0 || i < 0 {
		return ""
	}
	return (*s)[length-i-1]
}

func (s *Series) Length() int {
	return len(*s)
}

func (V Series) Filter(predicate func(item any, index int) bool) Series {
	result := Series{}

	for i, item := range V {
		if predicate(item, i) {
			result = append(result, item)
		}
	}

	return result
}

func (V Series) Map(iteratee func(item any, index int) any) Series {
	result := make(Series, len(V))

	for i, item := range V {
		result[i] = iteratee(item, i)
	}

	return result
}

func (V Series) Reduce(accumulator func(agg any, item any, index int) any, initial any) any {
	for i, item := range V {
		initial = accumulator(initial, item, i)
	}

	return initial
}

func (collection Series) ForEach(iteratee func(item any, index int)) {
	for i, item := range collection {
		iteratee(item, i)
	}
}

//func (s Series) Change(args ...any) any {
//	fmt.Println(args[2:])
//
//	for _, arg := range args {
//		switch arg.(type) {
//		case int:
//			fmt.Println(arg, "is an int value.")
//		case string:
//			fmt.Println(arg, "is a string value.")
//		case int64:
//			fmt.Println(arg, "is an int64 value.")
//		default:
//			fmt.Println(arg, "is an unknown type.")
//		}
//	}
//	fmt.Println(s.Last(), s.Index(1))
//	return s.Last() != s.Index(1)
//}
