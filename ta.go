package tart

import "github.com/iamjinlei/go-tart/floats"

type Ta struct {
}

func (ta Ta) Sma(in floats.Slice, n int64) floats.Slice {

	out := make([]float64, len(in))

	s := NewSma(n)
	for i, v := range in {
		out[i] = s.Update(v)
	}

	return out
}
func (ta Ta) Atr(h, l, c floats.Slice, n int64) floats.Slice {

	out := make([]float64, len(c))

	a := NewAtr(n)
	for i := 0; i < len(c); i++ {
		out[i] = a.Update(h[i], l[i], c[i])
	}

	return out
}
func (ta Ta) Ema(in floats.Slice, n int64) floats.Slice {
	out := make([]float64, len(in))

	k := 2.0 / float64(n+1)
	e := NewEma(n, k)
	for i, v := range in {
		out[i] = e.Update(v)
	}

	return out
}
func (ta Ta) Rsi(in floats.Slice, n int64) floats.Slice {

	out := make([]float64, len(in))

	r := NewRsi(n)
	for i, v := range in {
		out[i] = r.Update(v)
	}

	return out
}
