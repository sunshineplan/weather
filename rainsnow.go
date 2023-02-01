package weather

import (
	"fmt"
	"strings"
	"time"
)

type RainSnow struct {
	start, end Day
}

func NewRainSnow(start, end Day) RainSnow {
	return RainSnow{start, end}
}

func (rainsnow *RainSnow) Start() *Day {
	return &rainsnow.start
}

func (rainsnow *RainSnow) End() *Day {
	return &rainsnow.end
}

func (s *RainSnow) Duration() string {
	if s.start.DateEpoch != 0 && s.end.DateEpoch != 0 {
		if start, end := time.Unix(s.start.DateEpoch, 0), time.Unix(s.end.DateEpoch, 0); !start.IsZero() && !end.IsZero() {
			return fmtDuration(end.Sub(start))
		}
	}
	return "unknown"
}

func (rainsnow RainSnow) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, "Begin:", rainsnow.start.Date)
	if rainsnow.end.Date != "" {
		fmt.Fprintln(&b, "End:", rainsnow.end.Date)
	}
	if duration := time.Until(time.Unix(rainsnow.start.DateEpoch, 0)); duration > 0 {
		fmt.Fprintln(&b, "Next:", fmtDuration(duration.Truncate(24*time.Hour)+24*time.Hour))
	}
	if duration := rainsnow.Duration(); duration != "0s" {
		fmt.Fprintln(&b, "Duration:", duration)
	}
	fmt.Fprintln(&b, "Detail:")
	fmt.Fprintln(&b, rainsnow.start)
	if rainsnow.end.Date != "" {
		fmt.Fprintln(&b, rainsnow.end)
	}
	return b.String()
}

func WillRainSnow(api API, query string, n int) (res []RainSnow, err error) {
	_, days, err := api.Forecast(query, n)
	if err != nil {
		return
	}

	var start, last Day
	for _, i := range days {
		switch i.Precip {
		case 0:
			if start.Date != "" {
				res = append(res, NewRainSnow(start, last))
				start = Day{}
			}
		default:
			if start.Date == "" {
				start = i
			}
		}
		last = i
	}
	if start.Date != "" {
		res = append(res, NewRainSnow(start, Day{}))
	}
	return
}

func (api Weather) WillRainSnow(query string, n int) ([]RainSnow, error) {
	return WillRainSnow(api, query, n)
}
