package main

import (
	"fmt"
	tart "github.com/uvite/indicator"

	"reflect"
)

func Change(args ...any) any {
	//fmt.Println(args[2:])
	var res tart.Series
	var len int = 1
	for _, arg := range args {
		//fmt.Println(reflect.TypeOf(arg))
		switch arg.(type) {
		case tart.Series:
			res = arg.(tart.Series)
		case int:
			len = arg.(int)

		default:
			fmt.Println(arg, "is an unknown type.")
		}
	}
	if res.Length() > 0 {
		fmt.Println(reflect.TypeOf(res.Index(0)).String())
		if reflect.TypeOf(res.Index(0)).String() == "float64" {
			diff := res.Index(0).(float64) - res.Index(len).(float64)
			return diff
		}
		if reflect.TypeOf(res.Index(0)).String() == "string" {
			diff := res.Index(len).(string) != res.Index(0).(string)
			return diff
		}

	}
	return ""
}
func main() {
	res := tart.NewSeries()
	//for i := 0; i < 10; i++ {
	//	res.Push(i)
	//}
	res.Push("hold")
	res.Push("buy")
	res.Push("buy")
	res.Push("sell")

	fmt.Println(Change(res))

	//f1 := res.Filter(func(item any, index int) bool {
	//	return item.(int)%2 == 0
	//})
	//f2 := res.Map(func(item any, index int) any {
	//	return item.(int) * 2
	//})
	//
	//f3 := res.Reduce(func(agg any, item any, index int) any {
	//	fmt.Println(agg)
	//	return item.(int) + agg.(int)
	//}, 3)
	//fmt.Println(f1, f2, f3)

}
