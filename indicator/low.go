package indicator

import (
	"time"

	"github.com/uvite/indicator/floats"
	"github.com/uvite/indicator/types"
)

//go:generate callbackgen -type Low
type Low struct {
	types.IntervalWindow
	types.SeriesBase

	Values  floats.Slice
	EndTime time.Time

	updateCallbacks []func(value float64)
}

func (inc *Low) Update(value float64) {
	if len(inc.Values) == 0 {
		inc.SeriesBase.Series = inc
	}

	inc.Values.Push(value)
}

func (inc *Low) PushK(k types.KLine) {
	if k.EndTime.Before(inc.EndTime) {
		return
	}

	inc.Update(k.Low.Float64())
	inc.EndTime = k.EndTime.Time()
	inc.EmitUpdate(inc.Last())
}
