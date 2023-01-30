package weather

import (
	"fmt"
	"time"
)

type RainSnow struct {
	start, end *Day
}

func NewRainSnow(start, end *Day) *RainSnow {
	return &RainSnow{start, end}
}

func (rainsnow *RainSnow) Start() *Day {
	return rainsnow.start
}

func (rainsnow *RainSnow) End() *Day {
	return rainsnow.end
}

func (s *RainSnow) Duration() string {
	if s.start != nil && s.end != nil {
		if start, end := time.Unix(s.start.DateEpoch, 0), time.Unix(s.end.DateEpoch, 0); !start.IsZero() && !end.IsZero() {
			return end.Sub(start).String()
		}
	}
	return "unknown"
}

func (rainsnow RainSnow) String() string {
	if rainsnow.end == nil {
		return fmt.Sprintf(`
Begin at: %s
Duration: unknown
`, rainsnow.start.Date)
	}
	return fmt.Sprintf(`
Begin at: %s
End at: %s
Duration: %s
`, rainsnow.start.Date, rainsnow.end.Date, rainsnow.Duration())
}

func WillRainSnow(api API, query string, n int) (res []*RainSnow, err error) {
	_, days, err := api.Forecast(query, n)
	if err != nil {
		return
	}

	var start, last *Day
	for _, i := range days {
		switch i.Precip {
		case 0:
			if start != nil {
				res = append(res, NewRainSnow(start, last))
				start = nil
			}
		default:
			if start == nil {
				start = &i
			}
		}
		last = &i
	}
	if start != nil {
		res = append(res, NewRainSnow(start, nil))
	}
	return
}

func (api Weather) WillRainSnow(query string, n int) ([]*RainSnow, error) {
	return WillRainSnow(api, query, n)
}
