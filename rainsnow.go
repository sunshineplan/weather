package weather

import (
	"fmt"
	"strings"
	"time"
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

func (rainsnow *RainSnow) Start() *Day {
	if len(rainsnow.days) > 0 {
		return &rainsnow.days[0]
	}
	return nil
}

func (rainsnow *RainSnow) End() *Day {
	if length := len(rainsnow.days); rainsnow.isEnd && length > 0 {
		return &rainsnow.days[length-1]
	}
	return nil
}

func (rainsnow *RainSnow) IsEnd() bool {
	return rainsnow.isEnd
}

func (rainsnow *RainSnow) Duration() int {
	if rainsnow.isEnd {
		return len(rainsnow.days)
	}
	return 0
}

func (rainsnow RainSnow) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, "Begin:", rainsnow.Start().Date)
	if rainsnow.isEnd {
		fmt.Fprintln(&b, "End:", rainsnow.End().Date)
	}
	if duration := time.Until(time.Unix(rainsnow.Start().DateEpoch, 0)); duration > 0 {
		fmt.Fprintln(&b, "Until:", fmtDuration(duration.Truncate(24*time.Hour)+24*time.Hour))
	}
	if duration := rainsnow.Duration(); duration != 0 {
		fmt.Fprintf(&b, "Duration: %dd\n", duration)
	} else {
		fmt.Fprintln(&b, "Duration: unknown")
	}
	fmt.Fprintln(&b, "Detail:")
	for index, i := range rainsnow.days {
		fmt.Fprintf(&b, "#%d %s\n", index+1, i)
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
