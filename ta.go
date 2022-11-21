package tart

import (
	"fmt"
	"github.com/uvite/indicator/floats"
	"github.com/uvite/indicator/indicator"
	. "github.com/uvite/indicator/tart"
	"github.com/uvite/indicator/types"
)

type Talib struct {
}

func (ta Talib) Alma(in floats.Slice, n int64, offset float64, sigma int) floats.Slice {

	//out := make(floats.Slice, len(in))

	alma := indicator.ALMA{
		IntervalWindow: types.IntervalWindow{Window: int(n)},
		Offset:         offset,
		Sigma:          sigma,
	}
	//alma.CalculateAndUpdate(tt.kLines)
	//s := NewSma(n)
	fmt.Println(alma.Length(), "alma.length")
	for _, v := range in {
		alma.Update(v)
	}

	return alma.Values
}

// 以上为bbgo indicator 实现

func (ta Talib) Sma(in floats.Slice, n int64) floats.Slice {

	out := make(floats.Slice, len(in))

	s := NewSma(n)
	for i, v := range in {
		out[i] = s.Update(v)
	}

	return out
}
func (ta Talib) Atr(h, l, c floats.Slice, n int64) floats.Slice {

	out := make(floats.Slice, len(c))

	a := NewAtr(n)
	for i := 0; i < len(c); i++ {
		out[i] = a.Update(h[i], l[i], c[i])
	}

	return out
}
func (ta Talib) Ema(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	k := 2.0 / float64(n+1)
	e := NewEma(n, k)
	for i, v := range in {
		out[i] = e.Update(v)
	}

	return out
}
func (ta Talib) Rsi(in floats.Slice, n int64) floats.Slice {

	out := make(floats.Slice, len(in))

	r := NewRsi(n)
	for i, v := range in {
		out[i] = r.Update(v)
	}

	return out
}
func (ta Talib) Wma(in floats.Slice, n int64) floats.Slice {

	out := make(floats.Slice, len(in))

	w := NewWma(n)
	for i, v := range in {
		out[i] = w.Update(v)
	}

	return out
}
func (ta Talib) WillR(h, l, c floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(c))

	w := NewWillR(n)
	for i := 0; i < len(c); i++ {
		out[i] = w.Update(h[i], l[i], c[i])
	}

	return out
}
func (ta Talib) Var(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	s := NewVar(n)
	for i, v := range in {
		out[i] = s.Update(v)
	}

	return out
}

func (ta Talib) UltOsc(h, l, c floats.Slice, n1, n2, n3 int64) floats.Slice {
	out := make(floats.Slice, len(c))

	u := NewUltOsc(n1, n2, n3)
	for i := 0; i < len(c); i++ {
		out[i] = u.Update(h[i], l[i], c[i])
	}

	return out
}
func (ta Talib) Trix(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	t := NewTrix(n)
	for i, v := range in {
		out[i] = t.Update(v)
	}

	return out
}
func (ta Talib) Trima(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	t := NewTrima(n)
	for i, v := range in {
		out[i] = t.Update(v)
	}

	return out
}
func (ta Talib) TRange(h, l, c floats.Slice) floats.Slice {
	out := make(floats.Slice, len(c))

	t := NewTRange()
	for i := 0; i < len(c); i++ {
		out[i] = t.Update(h[i], l[i], c[i])
	}

	return out
}

func (ta Talib) StochRsi(in floats.Slice, n, kN int64, dt MaType, dN int64) (floats.Slice, floats.Slice) {
	k := make(floats.Slice, len(in))
	d := make(floats.Slice, len(in))

	s := NewStochRsi(n, kN, dt, dN)
	for i, v := range in {
		k[i], d[i] = s.Update(v)
	}

	return k, d
}

func (ta Talib) StdDev(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	s := NewStdDev(n)
	for i, v := range in {
		out[i] = s.Update(v)
	}

	return out
}
func (ta Talib) Roc(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	r := NewRoc(n)
	for i, v := range in {
		out[i] = r.Update(v)
	}

	return out
}

func (ta Talib) Obv(c, v floats.Slice) floats.Slice {
	out := make(floats.Slice, len(c))

	o := NewObv()
	for i := 0; i < len(c); i++ {
		out[i] = o.Update(c[i], v[i])
	}

	return out
}
func (ta Talib) Natr(h, l, c floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(c))

	a := NewNatr(n)
	for i := 0; i < len(c); i++ {
		out[i] = a.Update(h[i], l[i], c[i])
	}

	return out
}

func (ta Talib) Macd(in floats.Slice, fastN, slowN, signalN int64) (floats.Slice, floats.Slice, floats.Slice) {
	macd := make(floats.Slice, len(in))
	signal := make(floats.Slice, len(in))
	hist := make(floats.Slice, len(in))

	m := NewMacd(fastN, slowN, signalN)
	for i, v := range in {
		macd[i], signal[i], hist[i] = m.Update(v)
	}

	return macd, signal, hist
}
func (ta Talib) Kama(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	k := NewKama(n)
	for i, v := range in {
		out[i] = k.Update(v)
	}

	return out
}
func (ta Talib) Dx(h, l, c floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(c))

	d := NewDx(n)
	for i := 0; i < len(c); i++ {
		out[i] = d.Update(h[i], l[i], c[i])
	}

	return out
}

func (ta Talib) Diff(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	d := NewDiff(n)
	for i, v := range in {
		out[i] = d.Update(v)
	}

	return out
}

func (ta Talib) Dev(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	d := NewDev(n)
	for i, v := range in {
		out[i] = d.Update(v)
	}

	return out
}
func (ta Talib) Ppo(in floats.Slice, t MaType, fastN, slowN int64) floats.Slice {
	out := make(floats.Slice, len(in))

	p := NewPpo(t, fastN, slowN)
	for i, v := range in {
		out[i] = p.Update(v)
	}

	return out
}

func (ta Talib) Dema(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	k := 2.0 / float64(n+1)
	d := NewDema(n, k)
	for i, v := range in {
		out[i] = d.Update(v)
	}

	return out
}
func (ta Talib) Cmo(in floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(in))

	c := NewCmo(n)
	for i, v := range in {
		out[i] = c.Update(v)
	}

	return out
}

func (ta Talib) Cci(h, l, c floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(h))

	d := NewCci(n)
	for i := 0; i < len(h); i++ {
		out[i] = d.Update(h[i], l[i], c[i])
	}

	return out
}

func (ta Talib) BBands(t MaType, in floats.Slice, n int64, upNStdDev, dnNStdDev float64) (floats.Slice, floats.Slice, floats.Slice) {
	m := make(floats.Slice, len(in))
	u := make(floats.Slice, len(in))
	l := make(floats.Slice, len(in))

	b := NewBBands(t, n, upNStdDev, dnNStdDev)
	for i, v := range in {
		u[i], m[i], l[i] = b.Update(v)
	}

	return u, m, l
}
func (ta Talib) AroonOsc(h, l floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(h))

	a := NewAroonOsc(n)
	for i := 0; i < len(h); i++ {
		out[i] = a.Update(h[i], l[i])
	}

	return out
}
func (ta Talib) Aroon(h, l floats.Slice, n int64) (floats.Slice, floats.Slice) {
	dn := make(floats.Slice, len(h))
	up := make(floats.Slice, len(h))

	a := NewAroon(n)
	for i := 0; i < len(h); i++ {
		dn[i], up[i] = a.Update(h[i], l[i])
	}

	return dn, up
}
func (ta Talib) Apo(t MaType, in floats.Slice, fastN, slowN int64) floats.Slice {
	out := make(floats.Slice, len(in))

	a := NewApo(t, fastN, slowN)
	for i, v := range in {
		out[i] = a.Update(v)
	}

	return out
}

func (ta Talib) Adx(h, l, c floats.Slice, n int64) floats.Slice {
	out := make(floats.Slice, len(c))

	a := NewAdx(n)
	for i := 0; i < len(c); i++ {
		out[i] = a.Update(h[i], l[i], c[i])
	}

	return out
}

//{"Wma","WillR","Var","UltOsc","Trix","Trima","TRange","StochRsi","StdDev","Roc","Obv",
//"Natr","Macd","Kama","Dx","Diff","Dev","ppo"}
