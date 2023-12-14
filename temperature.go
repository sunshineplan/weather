package weather

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/weather/unit"
)

type TempRiseFall struct {
	day, previous Day
	standard      float64
}

func NewTempRiseFall(day, previous Day, standard float64) TempRiseFall {
	return TempRiseFall{day, previous, standard}
}

func (t TempRiseFall) Day() Day {
	return t.day
}

func (t TempRiseFall) Previous() Day {
	return t.previous
}

func (t TempRiseFall) Difference() [2][3]unit.Temperature {
	return [2][3]unit.Temperature{
		{
			t.day.TempMax.Difference(t.previous.TempMax),
			t.day.TempMin.Difference(t.previous.TempMin),
			t.day.Temp.Difference(t.previous.Temp),
		},
		{
			t.day.FeelsLikeMax.Difference(t.previous.FeelsLikeMax),
			t.day.FeelsLikeMin.Difference(t.previous.FeelsLikeMin),
			t.day.FeelsLike.Difference(t.previous.FeelsLike),
		},
	}
}

func (t TempRiseFall) IsRise() bool {
	if tempRiseFall, isRise := isRiseFall(t.day.TempMax, t.previous.TempMax, t.standard); tempRiseFall {
		return isRise
	} else if tempRiseFall, isRise = isRiseFall(t.day.TempMin, t.previous.TempMin, t.standard); tempRiseFall {
		return isRise
	} else if tempRiseFall, isRise = isRiseFall(t.day.FeelsLikeMax, t.previous.FeelsLikeMax, t.standard); tempRiseFall {
		return isRise
	} else if tempRiseFall, isRise = isRiseFall(t.day.FeelsLikeMin, t.previous.FeelsLikeMin, t.standard); tempRiseFall {
		return isRise
	}
	return t.day.Temp.Float64() > t.previous.Temp.Float64()
}

func (t TempRiseFall) IsExpired() bool {
	return t.day.DateEpoch.IsExpired()
}

func (t TempRiseFall) DateInfo() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s %s", t.day.Date, t.day.DateEpoch.Weekday())
	if until := t.day.DateEpoch.Until(); until == 0 {
		fmt.Fprint(&b, " (today)")
	} else if until == 24*time.Hour {
		fmt.Fprint(&b, " (tomorrow)")
	} else {
		fmt.Fprintf(&b, " (%dd later)", until/(24*time.Hour))
	}
	return b.String()
}

func (t TempRiseFall) DiffInfo() string {
	diff := t.Difference()
	var b strings.Builder
	fmt.Fprintf(&b, "TempMaxDiff: %s, TempMinDiff: %s, TempDiff: %s\n", diff[0][0], diff[0][1], diff[0][2])
	fmt.Fprintf(&b, "FeelsLikeMaxDiff: %s, FeelsLikeMinDiff: %s, FeelsLikeDiff: %s", diff[1][0], diff[1][1], diff[1][2])
	return b.String()
}

func (t TempRiseFall) DiffInfoHTML() html.HTML {
	diff := t.Difference()
	return html.Background().AppendChild(
		html.Tr(
			html.Td("TempMax:"), html.Td(diff[0][0].DiffHTML()),
			html.Td("TempMin:"), html.Td(diff[0][1].DiffHTML()),
			html.Td("Temp:"), html.Td(diff[0][2].DiffHTML()),
		),
		html.Tr(
			html.Td("FeelsLikeMax:"), html.Td(diff[1][0].DiffHTML()),
			html.Td("FeelsLikeMin:"), html.Td(diff[1][1].DiffHTML()),
			html.Td("FeelsLike:"), html.Td(diff[1][2].DiffHTML()),
		),
	).HTML()
}

func (t TempRiseFall) String() string {
	var b strings.Builder
	fmt.Fprintln(&b, "Date: ", t.DateInfo())
	fmt.Fprintln(&b, t.DiffInfo())
	fmt.Fprintln(&b, "Forecast:")
	fmt.Fprintln(&b, "#0", t.previous.DateInfo(true))
	fmt.Fprintln(&b, t.previous.Temperature())
	fmt.Fprintln(&b, "#1", t.day.DateInfo(true))
	fmt.Fprint(&b, t.day.Temperature())
	return b.String()
}

func (t TempRiseFall) HTML() html.HTML {
	return html.Div().AppendChild(
		html.Span().Style("display:list-item;margin-left:15px;list-style-type:disclosure-open").
			Content(t.DateInfo(), " ", t.day.Condition.Img(t.day.Icon)),
		html.Table().AppendChild(
			html.Tbody().Content(t.day.TemperatureHTML(), t.DiffInfoHTML())),
		html.Span().Style("display:list-item;margin-left:15px;list-style-type:circle").
			Content("Previous Day: ", t.previous.DateInfoHTML()),
		html.Table().AppendChild(
			html.Tbody().Content(t.previous.TemperatureHTML())),
	).HTML()
}

func WillTempRiseFall(days []Day, standard float64) (res []TempRiseFall) {
	var day, previous Day
	for _, i := range days {
		day = i
		if previous.Date != "" {
			var found bool
			if tempRiseFall, _ := isRiseFall(day.TempMax, previous.TempMax, standard); tempRiseFall {
				found = true
			} else if tempRiseFall, _ = isRiseFall(day.TempMin, previous.TempMin, standard); tempRiseFall {
				found = true
			} else if tempRiseFall, _ = isRiseFall(day.FeelsLikeMax, previous.FeelsLikeMax, standard); tempRiseFall {
				found = true
			} else if tempRiseFall, _ = isRiseFall(day.FeelsLikeMin, previous.FeelsLikeMin, standard); tempRiseFall {
				found = true
			}
			if found {
				res = append(res, NewTempRiseFall(day, previous, standard))
			}
		}
		previous = i
	}
	return
}

func isRiseFall(t, previous unit.Temperature, standard float64) (tempRiseFall bool, isRise bool) {
	diff := t.Difference(previous).Float64()
	return math.Abs(diff) >= math.Abs(standard), diff > 0
}
