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
		if day.DateEpoch.IsExpired() {
			rainsnow.days = rainsnow.days[1:]
		} else {
			return false
		}
	}
	return true
}

func (rainsnow RainSnow) DateInfo() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s", rainsnow.Start().Date, rainsnow.Start().DateEpoch.Weekday())
	if until := rainsnow.Start().DateEpoch.Until(); until == 0 {
		fmt.Fprint(&b, " (today)")
	} else if until == 24*time.Hour {
		fmt.Fprint(&b, " (tomorrow)")
	} else {
		fmt.Fprintf(&b, " (%dd later)", until/(24*time.Hour))
	}
	if rainsnow.isEnd {
		if rainsnow.Duration() != 0 {
			fmt.Fprintf(&b, " ~ %s %s (last %dd)", rainsnow.End().Date, rainsnow.End().DateEpoch.Weekday(), rainsnow.Duration())
		}
	} else {
		fmt.Fprint(&b, " ~ unknown")
	}
	return b.String()
}

func (rainsnow RainSnow) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, "Date:", rainsnow.DateInfo())
	for index, i := range rainsnow.days {
		fmt.Fprintf(&b, "\n#%d %s\n", index+1, i.DateInfo(true))
		fmt.Fprint(&b, i.Precipitation())
	}
	return b.String()
}

func (rainsnow RainSnow) HTML(t time.Time, highlight ...int) string {
	var b strings.Builder
	fmt.Fprint(&b, "<div>")
	fmt.Fprintf(&b, `<span style="display:list-item;margin-left:15px;list-style-type:disclosure-open">%s</span>`, rainsnow.DateInfo())
	for index, i := range rainsnow.days {
		fmt.Fprint(&b, `<div style="display:grid">`)
		fmt.Fprintf(&b, "<span>%d.  %s</span>", index+1, i.DateInfoHTML())
		if i.Date == t.Format("2006-01-02") {
			fmt.Fprint(&b, i.PrecipitationHTML(highlight...))
		} else {
			fmt.Fprint(&b, i.PrecipitationHTML())
		}
		fmt.Fprint(&b, "</div>")
	}
	fmt.Fprint(&b, "</div>")
	return b.String()
}

func WillRainSnow(days []Day) (res []RainSnow) {
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
