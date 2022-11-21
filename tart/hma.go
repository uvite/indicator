package tart

import "math"

// The Triple Exponential Moving Average (Hma) reduces the lag of traditional
// EMAs, making it more responsive and better-suited for short-term trading.
// Shortly after developing the Double Exponential Moving Average (DEMA) in 1994,
// Patrick Mulloy took the concept a step further and created the Triple
// Exponential Moving Average (Hma). Like its predecessor DEMA, the Hma overlay
// uses the lag difference between different EMAs to adjust a traditional EMA.
// However, Hma's formula uses a triple-smoothed EMA in addition to the single-
// and double-smoothed EMAs employed in the formula for DEMA. The offset created
// using these three EMAs produces a moving average that stays even closer to the
// price bars than DEMA.
//
//	https://school.stockcharts.com/doku.php?id=technical_indicators:Hma
//	https://www.investopedia.com/terms/t/triple-exponential-moving-average.asp
type Hma struct {
	n    int64
	sz   int64
	ema1 *Wma
	ema2 *Wma
	ema3 *Wma
}

func NewHma(n int64) *Hma {
	return &Hma{
		n:    n,
		sz:   0,
		ema1: NewWma(n),
		ema2: NewWma(n / 2),
		ema3: NewWma(int64(math.Floor(math.Sqrt(float64(n))))),
	}
}

//hullma = ta.wma(2*ta.wma(src, length/2)-ta.wma(src, length), math.floor(math.sqrt(length)))

func (t *Hma) Update(v float64) float64 {
	t.sz++

	e1 := t.ema1.Update(v)
	e2 := t.ema2.Update(v)

	if t.sz > t.n-1 {

		e3 := t.ema3.Update(2*e2 - e1)
		return e3
	}

	return 0
}

func (t *Hma) InitPeriod() int64 {
	return t.n*3 - 3
}

func (t *Hma) Valid() bool {
	return t.sz > t.InitPeriod()
}

// The Triple Exponential Moving Average (Hma) reduces the lag of traditional
// EMAs, making it more responsive and better-suited for short-term trading.
// Shortly after developing the Double Exponential Moving Average (DEMA) in 1994,
// Patrick Mulloy took the concept a step further and created the Triple
// Exponential Moving Average (Hma). Like its predecessor DEMA, the Hma overlay
// uses the lag difference between different EMAs to adjust a traditional EMA.
// However, Hma's formula uses a triple-smoothed EMA in addition to the single-
// and double-smoothed EMAs employed in the formula for DEMA. The offset created
// using these three EMAs produces a moving average that stays even closer to the
// price bars than DEMA.
//
//	https://school.stockcharts.com/doku.php?id=technical_indicators:Hma
//	https://www.investopedia.com/terms/t/triple-exponential-moving-average.asp
func HmaArr(in []float64, n int64) []float64 {
	out := make([]float64, len(in))

	///k := 2.0 / float64(n+1)
	t := NewHma(n)
	for i, v := range in {
		out[i] = t.Update(v)
	}

	return out
}
