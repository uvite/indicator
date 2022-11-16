package indicator

import (
	"github.com/shopspring/decimal"
	"math"

	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/types"
)

// Refer: Arnaud Legoux Moving Average
// Refer: https://capital.com/arnaud-legoux-moving-average
// Also check https://github.com/DaveSkender/Stock.Indicators/blob/main/src/a-d/Alma/Alma.cs
// @param offset: Gaussian applied to the combo line. 1->ema, 0->sma
// @param sigma: the standard deviation applied to the combo line. This makes the combo line sharper
//
//go:generate callbackgen -type ALMA
type JMA struct {
	types.SeriesBase
	types.IntervalWindow         // required
	Phase                float64 // required: recommend to be 0.5
	Power                float64 // required: recommend to be 5

	E0              floats.Slice
	E1              floats.Slice
	E2              floats.Slice
	E3              floats.Slice
	PhaseRatio      float64
	alpha           float64
	beta            float64
	input           floats.Slice
	Values          floats.Slice
	UpdateCallbacks []func(value float64)
}

const MaxNumOfJMA = 5_000
const MaxNumOfJMATruncateSize = 100

func abc(i float64) float64 {
	re, _ := decimal.NewFromFloat(i).Round(2).Float64()
	return re
}

func (inc *JMA) Update(value float64) {

	if inc.Values.Length() == 0 {
		inc.SeriesBase.Series = inc

		if inc.Phase < (-100.) {
			inc.PhaseRatio = 0.5
		} else if inc.Phase > 100 {
			inc.PhaseRatio = 2.5
		} else {
			inc.PhaseRatio = inc.Phase/100 + 1.5
		}

		inc.beta = ((0.45 * float64(inc.Window-1)) / (0.45*float64(inc.Window-1) + 2))
		inc.alpha = (math.Pow(inc.beta, inc.Power))

		inc.Values.Update(value)
		inc.input.Push(value)
		return

	}
	inc.input.Push(value)
	if len(inc.input) > MaxNumOfJMA {
		inc.input = inc.input[MaxNumOfJMATruncateSize-1:]
	}
	inc.E0.Update((1-inc.alpha)*value + +inc.alpha*inc.E0.Index(0))
	inc.E1.Update((value-inc.E0.Last())*(1-inc.beta) + inc.beta*inc.E1.Last())

	inc.E2.Update((inc.E0.Last()+inc.PhaseRatio*inc.E1.Last()-inc.Values.Last())*math.Pow(1-inc.alpha, 2) + math.Pow(inc.alpha, 2)*inc.E2.Last())

	inc.Values.Push(inc.E2.Last() + inc.Values.Last())

}

func (inc *JMA) Last() float64 {
	if len(inc.Values) == 0 {
		return 0
	}
	return inc.Values[len(inc.Values)-1]
}

func (inc *JMA) Index(i int) float64 {
	if i >= len(inc.Values) {
		return 0
	}
	return inc.Values[len(inc.Values)-i-1]
}

func (inc *JMA) Length() int {
	return len(inc.Values)
}

var _ types.SeriesExtend = &JMA{}

func (inc *JMA) CalculateAndUpdate(allKLines []types.KLine) {
	if inc.input == nil {
		for _, k := range allKLines {
			inc.Update(k.Close.Float64())
			//inc.EmitUpdate(inc.Last())
		}
		return
	}
	inc.Update(allKLines[len(allKLines)-1].Close.Float64())
	//inc.EmitUpdate(inc.Last())
}

func (inc *JMA) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {
	if inc.Interval != interval {
		return
	}
	inc.CalculateAndUpdate(window)
}

func (inc *JMA) Bind(updater KLineWindowUpdater) {
	updater.OnKLineWindowUpdate(inc.handleKLineWindowUpdate)
}
