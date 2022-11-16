package indicator

import (
	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/indicator"
	"github.com/c9s/bbgo/pkg/types"
)

// Refer: Arnaud Legoux Moving Average
// Refer: https://capital.com/arnaud-legoux-moving-average
// Also check https://github.com/DaveSkender/Stock.Indicators/blob/main/src/a-d/Alma/Alma.cs
// @param offset: Gaussian applied to the combo line. 1->ema, 0->sma
// @param sigma: the standard deviation applied to the combo line. This makes the combo line sharper
//
//  hullma = ta.wma(2*ta.wma(src, length/2)-ta.wma(src, length), math.floor(math.sqrt(length)))
//
//wma 实现
//norm = 0.0
//sum = 0.0
//for i = 0 to y - 1
//	weight = (y - i) * y
//	norm := norm + weight
//	sum := sum + x[i] * weight
//	sum / norm

//go:generate callbackgen -type ALMA
type WMA struct {
	types.SeriesBase
	types.IntervalWindow         // required
	Offset               float64 // required: recommend to be 0.5
	Sigma                float64 // required: recommend to be 5

	weight          []float64
	sum             float64
	input           []float64
	Values          floats.Slice
	UpdateCallbacks []func(value float64)
}

const MaxNumOfWMA = 5_000
const MaxNumOfWMATruncateSize = 300

func (inc *WMA) Update(value float64) {
	if inc.weight == nil {
		inc.SeriesBase.Series = inc
		inc.weight = make([]float64, inc.Window)
		//m := inc.Offset * (float64(inc.Window) - 1.)
		//s := float64(inc.Window) / float64(inc.Sigma)
		inc.sum = 0.
		for i := 0; i < inc.Window; i++ {
			//diff := float64(i) - m
			//wt := math.Exp(-diff * diff / 2. / s / s)
			wt := float64(inc.Window-i) * float64(inc.Window)
			//norm := norm + weight
			inc.sum += wt
			inc.weight[i] = wt
		}
	}
	inc.input = append(inc.input, value)
	if len(inc.input) >= inc.Window {
		weightedSum := 0.0
		inc.input = inc.input[len(inc.input)-inc.Window:]
		for i := 0; i < inc.Window; i++ {
			weightedSum += inc.weight[inc.Window-i-1] * inc.input[i]
		}
		inc.Values.Push(weightedSum / inc.sum)
		if len(inc.Values) > MaxNumOfWMA {
			inc.Values = inc.Values[MaxNumOfWMATruncateSize-1:]
		}
	}
}

func (inc *WMA) Last() float64 {
	if len(inc.Values) == 0 {
		return 0
	}
	return inc.Values[len(inc.Values)-1]
}

func (inc *WMA) Index(i int) float64 {
	if i >= len(inc.Values) {
		return 0
	}
	return inc.Values[len(inc.Values)-i-1]
}

func (inc *WMA) Length() int {
	return len(inc.Values)
}

var _ types.SeriesExtend = &WMA{}

func (inc *WMA) CalculateAndUpdate(allKLines []types.KLine) {
	if inc.input == nil {
		for _, k := range allKLines {
			inc.Update(k.Close.Float64())
			inc.EmitUpdate(inc.Last())
		}
		return
	}
	inc.Update(allKLines[len(allKLines)-1].Close.Float64())
	inc.EmitUpdate(inc.Last())
}

func (inc *WMA) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {
	if inc.Interval != interval {
		return
	}
	inc.CalculateAndUpdate(window)
}

func (inc *WMA) Bind(updater indicator.KLineWindowUpdater) {
	updater.OnKLineWindowUpdate(inc.handleKLineWindowUpdate)
}
