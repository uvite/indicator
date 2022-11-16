package indicator

import (
	"fmt"
	"math"

	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/types"
)

type JURIK struct {
	types.SeriesBase
	types.IntervalWindow         // required
	Phase                float64 // required: recommend to be 0.5
	Power                float64 // required: recommend to be 5

	Uband floats.Slice
	Lband floats.Slice
	//BsMin  floats.Slice
	//BsMax  floats.Slice
	Volty  floats.Slice
	Vsum   floats.Slice
	avolty floats.Slice
	dVolty floats.Slice
	Ma1    floats.Slice
	Ma2    floats.Slice
	det0   floats.Slice
	E2     floats.Slice

	PhaseRatio      float64
	alpha           float64
	beta            float64
	bet             float64
	pow1            float64
	len1            float64
	div             float64
	input           floats.Slice
	Values          floats.Slice
	UpdateCallbacks []func(value float64)
}

const MaxNumOfJURIK = 5_000
const MaxNumOfJURIKTruncateSize = 100

func (inc *JURIK) Update(value float64) {

	if inc.Values.Length() == 0 {
		inc.SeriesBase.Series = inc
		inc.len1 = math.Max(math.Log(math.Sqrt(0.5*float64(inc.Window-1)))/math.Log(2.0)+2.0, 0)
		len2 := math.Sqrt(0.5*float64(inc.Window-1)) * inc.len1
		inc.pow1 = math.Max(inc.len1-2.0, 0.5)
		inc.beta = ((0.45 * float64(inc.Window-1)) / (0.45*float64(inc.Window-1) + 2))

		inc.div = 1.0 / (10.0 + 10.0*(math.Min(math.Max(float64(inc.Window-10), 0), 100))/100)

		if inc.Phase < (-100.) {
			inc.PhaseRatio = 0.5
		} else if inc.Phase > 100 {
			inc.PhaseRatio = 2.5
		} else {
			inc.PhaseRatio = inc.Phase/100 + 1.5
		}

		inc.bet = len2 / (len2 + 1)

		inc.Uband.Update(value)

		inc.Lband.Update(value)

		inc.Vsum.Update(0)
		inc.Volty.Update(0)
		inc.avolty.Update(0)
		inc.Values.Update(0)
		inc.input.Push(value)
		return

	}

	inc.input.Push(value)
	if len(inc.input) < inc.Window {
		return
	}
	if len(inc.input) > MaxNumOfJURIK {
		inc.input = inc.input[MaxNumOfJURIKTruncateSize-1:]
	}

	//Price volatility
	del1 := value - inc.Uband.Last()
	del2 := value - inc.Lband.Last()
	//if math.Abs(del1) != math.Abs(del2) {
	//	inc.Volty.Update(math.Max(math.Abs(del1), math.Abs(del2)))
	//} else {
	//	inc.Volty.Update(0)
	//}

	if math.Abs(del1) > math.Abs(del2) {
		inc.Volty.Update(math.Abs(del1))
	} else {
		inc.Volty.Update(math.Abs(del2))

	}

	//Relative price volatility factor

	fmt.Println(del1, del2, inc.Volty.Last(), "\n\n")

	//Relative price volatility factor

	inc.Vsum.Update(inc.Vsum.Last() + inc.div*(inc.Volty.Last()-inc.Volty.Index(10)))

	inc.avolty.Update(inc.avolty.Last() + (2.0/(math.Max(float64(4.0*inc.Window), 30)+1.0))*(inc.Vsum.Last()-inc.avolty.Last()))

	dvolty1 := 0.
	if inc.avolty.Last() > 0 {
		dvolty1 = inc.Volty.Last() / inc.avolty.Last()
	}
	dvolty := (math.Max(1, math.Min(math.Pow(float64(inc.len1), 1.0/inc.pow1), dvolty1)))

	//fmt.Println(inc.Volty.Last(), inc.Vsum.Last(), inc.avolty.Last(), dvolty, "\n\n")
	//Jurik volatility bands
	pow2 := math.Pow(dvolty, inc.pow1)

	Kv := math.Pow(inc.bet, math.Sqrt(pow2))

	if del1 > 0 {
		inc.Uband.Update(value)
	} else {
		inc.Uband.Update(value - Kv*del1)
	}
	if del2 < 0 {
		inc.Lband.Update(value)
	} else {
		inc.Lband.Update(value - Kv*del2)
	}

	//Jurik Dynamic Factor
	alpha := math.Pow(inc.beta, pow2)

	inc.Ma1.Update((1-alpha)*value + alpha*inc.Ma1.Last())
	inc.det0.Update((value-inc.Ma1.Last())*(1-inc.beta) + inc.beta*inc.det0.Last())

	inc.Ma2.Update(inc.Ma1.Last() + inc.PhaseRatio*inc.det0.Last())

	inc.E2.Update((inc.Ma2.Last()-inc.Values.Last())*math.Pow(1-alpha, 2) + math.Pow(alpha, 2)*inc.E2.Last())

	inc.Values.Push(inc.E2.Last() + inc.Values.Last())

}

func (inc *JURIK) Last() float64 {
	if len(inc.Values) == 0 {
		return 0
	}
	return inc.Values[len(inc.Values)-1]
}

func (inc *JURIK) Index(i int) float64 {
	if i >= len(inc.Values) {
		return 0
	}
	return inc.Values[len(inc.Values)-i-1]
}

func (inc *JURIK) Length() int {
	return len(inc.Values)
}

var _ types.SeriesExtend = &JURIK{}

func (inc *JURIK) CalculateAndUpdate(allKLines []types.KLine) {
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

func (inc *JURIK) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {
	if inc.Interval != interval {
		return
	}
	inc.CalculateAndUpdate(window)
}

func (inc *JURIK) Bind(updater KLineWindowUpdater) {
	updater.OnKLineWindowUpdate(inc.handleKLineWindowUpdate)
}
