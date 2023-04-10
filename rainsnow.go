package weather

import (
	"fmt"
	"strings"
)

type RainSnow struct {
	days  []Day
	isEnd bool
}

func NewRainSnow(days []Day, isEnd bool) RainSnow {
	return RainSnow{days, isEnd}
}

func (rainsnow *RainSnow) Days() []Day {
	return rainsnow.days
}

func (rainsnow *RainSnow) Start() Day {
	if len(rainsnow.days) > 0 {
		return rainsnow.days[0]
	}
	return Day{}
}

func (rainsnow *RainSnow) End() Day {
	if length := len(rainsnow.days); rainsnow.isEnd && length > 0 {
		return rainsnow.days[length-1]
	}
	return Day{}
}

func (rainsnow RainSnow) IsEnd() bool {
	return rainsnow.isEnd
}

func (rainsnow RainSnow) Duration() int {
	if rainsnow.isEnd {
		return len(rainsnow.days)
	}
	return 0
}

func (rainsnow *RainSnow) IsExpired() bool {
	for day := rainsnow.Start(); day.Date != ""; day = rainsnow.Start() {
		if day.IsExpired() {
			rainsnow.days = rainsnow.days[1:]
		} else {
			return false
		}
	}
	return true
}

func (rainsnow RainSnow) DateInfo() string {
	var b strings.Builder
	fmt.Fprintf(&b, "Date: %s %s", rainsnow.Start().Date, rainsnow.Start().Weekday())
	if until := rainsnow.Start().Until(); until == 0 {
		fmt.Fprint(&b, " (today)")
	} else if until == 1 {
		fmt.Fprint(&b, " (tomorrow)")
	} else {
		fmt.Fprintf(&b, " (%s later)", fmtDuration(until))
	}
	if rainsnow.isEnd {
		if rainsnow.Duration() != 0 {
			fmt.Fprintf(&b, " ~ %s %s (last %dd)", rainsnow.End().Date, rainsnow.End().Weekday(), rainsnow.Duration())
		}
	} else {
		fmt.Fprint(&b, " ~ unknown")
	}
	return b.String()
}

func (rainsnow RainSnow) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, rainsnow.DateInfo())
	fmt.Fprint(&b, "Forecast:")
	for index, i := range rainsnow.days {
		fmt.Fprintf(&b, "\n#%d %s\n", index+1, i.DateInfo(true))
		fmt.Fprint(&b, i.Precipitation())
	}
	return b.String()
}

func WillRainSnow(days []Day) (res []RainSnow, err error) {
	var rainDays []Day
	for _, i := range days {
		switch i.Precip {
		case 0:
			if len(rainDays) > 0 {
				res = append(res, NewRainSnow(rainDays, true))
				rainDays = nil
			}
		default:
			rainDays = append(rainDays, i)
		}
	}
	if len(rainDays) > 0 {
		res = append(res, NewRainSnow(rainDays, false))
	}
	return
}
