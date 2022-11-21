package types

import "github.com/uvite/indicator/floats"

var _ Series = floats.Slice([]float64{}).Addr()
