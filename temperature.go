package weather

import (
	"fmt"
	"math"
	"strings"
	"time"
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
	var b strings.Builder
	fmt.Fprintln(&b, "Date:", t.day.Date)
	fmt.Fprintln(&b, "Until:", fmtDuration(time.Until(time.Unix(t.day.DateEpoch, 0)).Truncate(24*time.Hour)+24*time.Hour))
	diff1, diff2 := t.Difference()
	fmt.Fprintf(&b, "TempMaxDiff: %.1f\n", diff1)
	fmt.Fprintf(&b, "TempMinDiff: %.1f\n", diff2)
	fmt.Fprintln(&b, "Detail:")
	fmt.Fprintln(&b, "#0", t.previous)
	fmt.Fprintln(&b, "#1", t.day)
	return b.String()
}

func WillTempRiseFall(days []Day, difference float64) (res []TempRiseFall, err error) {
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
