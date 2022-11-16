package indicator

import (
	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/types"
	"math"
)

//go:generate callbackgen -type RSX
type RSX struct {
	types.SeriesBase
	types.IntervalWindow // required
	input                floats.Slice
	Values               floats.Slice
	F28                  floats.Slice
	F38                  floats.Slice
	F48                  floats.Slice
	F58                  floats.Slice
	F68                  floats.Slice
	F78                  floats.Slice
	F30                  floats.Slice
	F40                  floats.Slice
	F50                  floats.Slice
	F60                  floats.Slice
	F70                  floats.Slice
	F80                  floats.Slice

	UpdateCallbacks []func(value float64)
}

const MaxNumOfRSX = 5_000
const MaxNumOfRSXTruncateSize = 100

func (inc *RSX) Update(value float64) {

	if inc.Values.Length() == 0 {
		inc.SeriesBase.Series = inc
		inc.Values.Push(value)
		inc.input.Push(value * 100)

		return
	}
	inc.input.Push(value * 100)
	if len(inc.input) > MaxNumOfRSX {
		inc.input = inc.input[MaxNumOfRSXTruncateSize-1:]
	}
	mom0 := types.Change(&inc.input)
	moa0 := math.Abs(mom0.Last())
	Kg := 3. / float64(inc.Window+2.0)

	Hg := 1 - Kg

	//mom
	inc.F28.Update(Kg*mom0.Last() + Hg*inc.F28.Index(0))
	inc.F30.Update(Hg*inc.F30.Index(0) + Kg*inc.F28.Last())

	mom1 := inc.F28.Last()*1.5 - inc.F30.Last()*0.5

	inc.F38.Update(Hg*inc.F38.Index(0) + Kg*mom1)
	inc.F40.Update(Kg*inc.F38.Last() + Hg*inc.F40.Index(0))

	mom2 := inc.F38.Last()*1.5 - inc.F40.Last()*0.5

	inc.F48.Update(Hg*inc.F48.Index(0) + Kg*mom2)
	inc.F50.Update(Kg*inc.F48.Last() + Hg*inc.F50.Index(0))

	mom_out := inc.F48.Last()*1.5 - inc.F50.Last()*0.5

	inc.F58.Update(Hg*inc.F58.Index(0) + Kg*moa0)
	inc.F60.Update(Kg*inc.F58.Last() + Hg*inc.F60.Last())

	moa1 := inc.F58.Last()*1.5 - inc.F60.Last()*0.5

	inc.F68.Update(Hg*inc.F68.Index(0) + Kg*moa1)
	inc.F70.Update(Kg*inc.F68.Last() + Hg*inc.F70.Index(0))

	moa2 := inc.F68.Last()*1.5 - inc.F70.Last()*0.5

	inc.F78.Update(Hg*inc.F78.Index(0) + Kg*moa2)
	inc.F80.Update(Kg*inc.F78.Last() + Hg*inc.F80.Index(0))

	moa_out := inc.F78.Last()*1.5 - inc.F80.Last()*0.5

	rsx := math.Max(math.Min((mom_out/moa_out+1.0)*50.0, 100.00), 0.00)
	inc.Values.Push(rsx)

}

func (inc *RSX) Last() float64 {
	if len(inc.Values) == 0 {
		return 0
	}
	return inc.Values[len(inc.Values)-1]
}

func (inc *RSX) Index(i int) float64 {
	if i >= len(inc.Values) {
		return 0
	}
	return inc.Values[len(inc.Values)-i-1]
}

func (inc *RSX) Length() int {
	return len(inc.Values)
}

var _ types.SeriesExtend = &RSX{}

func (inc *RSX) CalculateAndUpdate(allKLines []types.KLine) {
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

func (inc *RSX) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {
	if inc.Interval != interval {
		return
	}
	inc.CalculateAndUpdate(window)
}

func (inc *RSX) Bind(updater KLineWindowUpdater) {
	updater.OnKLineWindowUpdate(inc.handleKLineWindowUpdate)
}
