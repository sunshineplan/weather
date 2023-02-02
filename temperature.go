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
	fmt.Fprintf(&b, "Date: %s %s (%s later)\n", t.day.Date, t.day.Weekday(),
		fmtDuration(time.Until(t.day.Time()).Truncate(24*time.Hour)+24*time.Hour))
	diff1, diff2 := t.Difference()
	fmt.Fprintf(&b, "TempMaxDiff: %.1f°C, TempMinDiff: %.1f°C\n", diff1, diff2)
	fmt.Fprintln(&b, "Forecast:")
	fmt.Fprintln(&b, "#0", t.previous.Temperature())
	fmt.Fprintln(&b, "#1", t.day.Temperature())
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
