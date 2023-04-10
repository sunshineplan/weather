package weather

import (
	"fmt"
	"strings"
)

type TempRiseFall struct {
	day, previous Day
}

func NewTempRiseFall(day, previous Day) TempRiseFall {
	return TempRiseFall{day, previous}
}

func (t TempRiseFall) Day() Day {
	return t.day
}

func (t TempRiseFall) Previous() Day {
	return t.previous
}

func (t TempRiseFall) Difference() [][]Temperature {
	return [][]Temperature{
		{t.day.TempMax - t.previous.TempMax, t.day.TempMin - t.previous.TempMin, t.day.Temp - t.previous.Temp},
		{t.day.FeelsLikeMax - t.previous.FeelsLikeMax, t.day.FeelsLikeMin - t.previous.FeelsLikeMin, t.day.FeelsLike - t.previous.FeelsLike},
	}
}

func (t TempRiseFall) IsRise() bool {
	if t.day.Temp == t.previous.Temp {
		return t.day.TempMax > t.previous.TempMax
	}
	return t.day.Temp > t.previous.Temp
}

func (t TempRiseFall) IsExpired() bool {
	return t.day.IsExpired()
}

func (t TempRiseFall) DateInfo() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Date: %s %s", t.day.Date, t.day.Weekday())
	if until := t.day.Until(); until == 0 {
		fmt.Fprint(&b, " (today)")
	} else if until == 1 {
		fmt.Fprint(&b, " (tomorrow)")
	} else {
		fmt.Fprintf(&b, " (%s later)", fmtDuration(until))
	}
	return b.String()
}

func (t TempRiseFall) DiffInfo() string {
	var b strings.Builder
	diff := t.Difference()
	fmt.Fprintf(&b, "TempMaxDiff: %.1f°C, TempMinDiff: %.1f°C, TempDiff: %.1f°C\n", diff[0][0], diff[0][1], diff[0][2])
	fmt.Fprintf(&b, "FeelsLikeMaxDiff: %.1f°C, FeelsLikeMinDiff: %.1f°C, FeelsLikeDiff: %.1f°C", diff[1][0], diff[1][1], diff[1][2])
	return b.String()
}

func (t TempRiseFall) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, t.DateInfo())
	fmt.Fprintln(&b, t.DiffInfo())
	fmt.Fprintln(&b, "Forecast:")
	fmt.Fprintln(&b, "#0", t.previous.DateInfo(true))
	fmt.Fprintln(&b, t.previous.Temperature())
	fmt.Fprintln(&b, "#1", t.day.DateInfo(true))
	fmt.Fprint(&b, t.day.Temperature())
	return b.String()
}

func WillTempRiseFall(days []Day, difference float64) (res []TempRiseFall, err error) {
	var day, previous Day
	for _, i := range days {
		day = i
		if previous.Date != "" {
			if day.TempMax.AbsDiff(previous.TempMax) >= difference ||
				day.TempMin.AbsDiff(previous.TempMin) >= difference ||
				day.FeelsLikeMax.AbsDiff(previous.FeelsLikeMax) >= difference ||
				day.FeelsLikeMin.AbsDiff(previous.FeelsLikeMin) >= difference {
				res = append(res, NewTempRiseFall(day, previous))
			}
		}
		previous = i
	}
	return
}
