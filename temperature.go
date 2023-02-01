package weather

import (
	"fmt"
	"math"
)

type TempRiseFall struct {
	day, previous Day
}

func NewTempRiseFall(day, previous Day) TempRiseFall {
	return TempRiseFall{day, previous}
}

func (t *TempRiseFall) Day() *Day {
	return &t.day
}

func (t *TempRiseFall) Previous() *Day {
	return &t.previous
}

func (t *TempRiseFall) Difference() (float64, float64) {
	return t.day.TempMax - t.previous.TempMax, t.day.TempMin - t.previous.TempMin
}

func (t *TempRiseFall) IsRise() bool {
	if t.day.Temp == t.previous.Temp {
		return t.day.TempMax > t.previous.TempMax
	}
	return t.day.Temp > t.previous.Temp
}

func (t TempRiseFall) String() string {
	diff1, diff2 := t.Difference()
	return fmt.Sprintf(`
Date: %s
TempMaxDiff: %.1f
TempMinDiff: %.1f
Detail: %s
`, t.day.Date, diff1, diff2, t.day)
}

func WillTempRiseFall(api API, difference float64, query string, n int) (res []TempRiseFall, err error) {
	_, days, err := api.Forecast(query, n)
	if err != nil {
		return
	}

	var day, previous Day
	for _, i := range days {
		day = i
		if previous.Date != "" {
			if math.Abs(day.TempMax-previous.TempMax) >= difference || math.Abs(day.TempMin-previous.TempMin) >= difference {
				res = append(res, NewTempRiseFall(day, previous))
			}
		}
		previous = i
	}
	return
}

func (api Weather) WillTempRiseFall(difference float64, query string, n int) ([]TempRiseFall, error) {
	return WillTempRiseFall(api, difference, query, n)
}
