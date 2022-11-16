package indicator

import (
	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/indicator"
	"github.com/c9s/bbgo/pkg/types"
	"math"
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
type HMA struct {
	types.SeriesBase
	types.IntervalWindow         // required
	Offset               float64 // required: recommend to be 0.5
	Sigma                float64 // required: recommend to be 5
	weight               []float64
	sum                  float64
	input                []float64
	Values               floats.Slice
	half                 types.UpdatableSeriesExtend
	all                  types.UpdatableSeriesExtend
	sqrt                 types.UpdatableSeriesExtend
	UpdateCallbacks      []func(value float64)
}

const MaxNumOfHMA = 5_000
const MaxNumOfHMATruncateSize = 300

func (inc *HMA) Update(value float64) {

	if len(inc.Values) == 0 {
		inc.SeriesBase.Series = inc
		half := int(math.Round(float64(inc.Window / 2)))
		sqrt := int(math.Round(math.Sqrt(float64(inc.Window)) - 0.5))
		inc.half = &WMA{IntervalWindow: types.IntervalWindow{Window: half}}
		inc.all = &WMA{IntervalWindow: types.IntervalWindow{Window: inc.Window}}
		inc.sqrt = &WMA{IntervalWindow: types.IntervalWindow{Window: sqrt}}
		inc.Values.Push(value)
		return
	} else if len(inc.Values) > MaxNumOfHMA {

		inc.Values = inc.Values[MaxNumOfHMATruncateSize-1:]
	}
	inc.half.Update(value)

	inc.all.Update(value)

	inc.sqrt.Update(2*inc.half.Last() - inc.all.Last())
	ema := inc.sqrt.Last()
	inc.Values.Push(ema)

}

func (inc *HMA) Last() float64 {
	if len(inc.Values) == 0 {
		return 0
	}
	return inc.Values[len(inc.Values)-1]
}

func (inc *HMA) Index(i int) float64 {
	if i >= len(inc.Values) {
		return 0
	}
	return inc.Values[len(inc.Values)-i-1]
}

func (inc *HMA) Length() int {
	return len(inc.Values)
}

var _ types.SeriesExtend = &HMA{}

func (inc *HMA) CalculateAndUpdate(allKLines []types.KLine) {
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

func (inc *HMA) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {
	if inc.Interval != interval {
		return
	}
	inc.CalculateAndUpdate(window)
}

func (inc *HMA) Bind(updater indicator.KLineWindowUpdater) {
	updater.OnKLineWindowUpdate(inc.handleKLineWindowUpdate)
}
