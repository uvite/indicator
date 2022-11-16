package indicator

import (
	"fmt"
	"github.com/c9s/bbgo/pkg/datatype/floats"
	"github.com/c9s/bbgo/pkg/fixedpoint"
	"github.com/c9s/bbgo/pkg/types"
	"math"
	"time"
)

const MaxNumOfEWMA = 5_000
const MaxNumOfEWMATruncateSize = 100

var three = fixedpoint.NewFromInt(3)

var zeroTime = time.Time{}

// Refer: Commodity Channel Index
// Refer URL: http://www.andrewshamlet.net/2017/07/08/python-tutorial-cci
// with modification of ddof=0 to let standard deviation to be divided by N instead of N-1
//
//go:generate callbackgen -type CCI
type CCI struct {
	types.SeriesBase
	types.IntervalWindow
	Input        floats.Slice
	TypicalPrice floats.Slice
	MA           floats.Slice
	Values       floats.Slice
	EndTime      time.Time

	UpdateCallbacks []func(value float64)
	Sqrt            bool
}

func (inc *CCI) Update(value float64) {
	if len(inc.TypicalPrice) == 0 {
		inc.SeriesBase.Series = inc
		inc.TypicalPrice.Push(value)
		inc.Input.Push(value)
		return
	} else if len(inc.TypicalPrice) > MaxNumOfEWMA {
		inc.TypicalPrice = inc.TypicalPrice[MaxNumOfEWMATruncateSize-1:]
		inc.Input = inc.Input[MaxNumOfEWMATruncateSize-1:]
	}

	inc.Input.Push(value)
	tp := inc.TypicalPrice.Last() - inc.Input.Index(inc.Window) + value
	inc.TypicalPrice.Push(tp)
	if len(inc.Input) < inc.Window {
		return
	}
	ma := tp / float64(inc.Window)
	inc.MA.Push(ma)
	if len(inc.MA) > MaxNumOfEWMA {
		inc.MA = inc.MA[MaxNumOfEWMATruncateSize-1:]
	}
	md := 0.

	if inc.Sqrt {
		for i := 0; i < inc.Window; i++ {
			diff := inc.Input.Index(i) - ma
			md += diff * diff
		}
		md = math.Sqrt(md / float64(inc.Window))
		fmt.Println("sqrt")
	} else {
		for i := 0; i < inc.Window; i++ {
			diff := inc.Input.Index(i) - ma
			md += math.Abs(diff)
		}
		md = (md / float64(inc.Window))
		//fmt.Println("simple")

	}

	cci := (value - ma) / (0.015 * md)

	inc.Values.Push(cci)
	if len(inc.Values) > MaxNumOfEWMA {
		inc.Values = inc.Values[MaxNumOfEWMATruncateSize-1:]
	}
}

func (inc *CCI) Update1(value float64) {
	if len(inc.TypicalPrice) == 0 {
		inc.SeriesBase.Series = inc
		inc.TypicalPrice.Push(value)
		inc.Input.Push(value)
		return
	} else if len(inc.TypicalPrice) > MaxNumOfEWMA {
		inc.TypicalPrice = inc.TypicalPrice[MaxNumOfEWMATruncateSize-1:]
		inc.Input = inc.Input[MaxNumOfEWMATruncateSize-1:]
	}

	inc.Input.Push(value)
	tp := inc.TypicalPrice.Last() - inc.Input.Index(inc.Window) + value
	inc.TypicalPrice.Push(tp)
	if len(inc.Input) < inc.Window {
		return
	}
	ma := tp / float64(inc.Window)
	inc.MA.Push(ma)
	if len(inc.MA) > MaxNumOfEWMA {
		inc.MA = inc.MA[MaxNumOfEWMATruncateSize-1:]
	}
	md := 0.
	//for i := 0; i < inc.Window; i++ {
	//	diff := inc.Input.Index(i) - ma
	//	md += diff * diff
	//}
	//md = math.Sqrt(md / float64(inc.Window))

	for i := 0; i < inc.Window; i++ {
		diff := inc.Input.Index(i) - ma
		md += math.Abs(diff)
	}
	md = (md / float64(inc.Window))

	cci := (value - ma) / (0.015 * md)
	//fmt.Printf("%.4f,%.4f,%.4f,%.4f,%.4f,%.4f \n", value, ma, md, 0.015, 0.015*md, cci)

	inc.Values.Push(cci)
	if len(inc.Values) > MaxNumOfEWMA {
		inc.Values = inc.Values[MaxNumOfEWMATruncateSize-1:]
	}
}

func (inc *CCI) Last() float64 {
	if len(inc.Values) == 0 {
		return 0
	}
	return inc.Values[len(inc.Values)-1]
}

func (inc *CCI) Index(i int) float64 {
	if i >= len(inc.Values) {
		return 0
	}
	return inc.Values[len(inc.Values)-1-i]
}

func (inc *CCI) Length() int {
	return len(inc.Values)
}

var _ types.SeriesExtend = &CCI{}

func (inc *CCI) PushK(k types.KLine) {

	//fmt.Println((k.High+k.Low+k.Close)/3.0, k.High.Add(k.Low).Add(k.Close).Div(three).Float64())
	inc.Update(k.High.Add(k.Low).Add(k.Close).Div(three).Float64())
}

func (inc *CCI) RePushK(k types.KLine) {
	//fmt.Println("len:", inc.Last(), inc.Index(1), inc.Index(2), inc.Length())
	//fmt.Println("\n kline:", k.Open, k.High, k.Low, k.Close, k.EndTime)

	if !k.Closed {
		inc.Input.Pop(int64(inc.Input.Length() - 1))
		inc.Values.Pop(int64(inc.Values.Length() - 1))
		inc.MA.Pop(int64(inc.MA.Length() - 1))
		inc.TypicalPrice.Pop(int64(inc.TypicalPrice.Length() - 1))

		//inc.Values[len(inc.Values)-1]
		//spew.Dump(inc.Values)
		inc.Update(k.High.Add(k.Low).Add(k.Close).Div(three).Float64())
		//fmt.Println("-2,-1,0,len:", inc.MA.Index(len(inc.Values)-3), inc.MA.Index(len(inc.Values)-2), inc.MA.Last(), inc.MA.Length())

		return
	} else {
		inc.Update(k.High.Add(k.Low).Add(k.Close).Div(three).Float64())

		inc.EndTime = k.EndTime.Time()
	}

}

func (inc *CCI) CalculateAndUpdate(allKLines []types.KLine) {
	if inc.TypicalPrice.Length() == 0 {
		for _, k := range allKLines {
			inc.PushK(k)

		}
	} else {
		k := allKLines[len(allKLines)-1]
		inc.PushK(k)

	}
}

func (inc *CCI) handleKLineWindowUpdate(interval types.Interval, window types.KLineWindow) {
	if inc.Interval != interval {
		return
	}

	inc.CalculateAndUpdate(window)
}
