package types

import (
	"time"

	"github.com/uvite/indicator/fixedpoint"
)

type PremiumIndex struct {
	Symbol          string           `json:"symbol"`
	MarkPrice       fixedpoint.Value `json:"markPrice"`
	LastFundingRate fixedpoint.Value `json:"lastFundingRate"`
	NextFundingTime time.Time        `json:"nextFundingTime"`
	Time            time.Time        `json:"time"`
}
