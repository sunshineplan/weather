package wind

import (
	"encoding/json"
	"fmt"

	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/weather/unit"
)

var (
	_ json.Marshaler = Speed(nil)

	_ Speed = KPH(0)
)

type Speed interface {
	MPS() float64
	Force() int
	ForceColor() string
	String() string
	MarshalJSON() ([]byte, error)
	HTML() html.HTML
}

func ForceColor(force int) string {
	return map[int]string{
		0:  "#FFFFFF",
		1:  "#AEF1F9",
		2:  "#96F7DC",
		3:  "#96F7B4",
		4:  "#6FF46F",
		5:  "#73ED12",
		6:  "#A4ED12",
		7:  "#DAED12",
		8:  "#EDC212",
		9:  "#ED8F12",
		10: "#ED6312",
		11: "#ED2912",
		12: "#D5102D",
	}[force]
}

type KPH float64

func (f KPH) MPS() float64 {
	return float64(f * 1000 / 60 / 60)
}

func (f KPH) Force() int {
	if f < 2 {
		return 0
	} else if f < 6 {
		return 1
	} else if f < 12 {
		return 2
	} else if f < 20 {
		return 3
	} else if f < 29 {
		return 4
	} else if f < 39 {
		return 5
	} else if f < 50 {
		return 6
	} else if f < 62 {
		return 7
	} else if f < 75 {
		return 8
	} else if f < 89 {
		return 9
	} else if f < 103 {
		return 10
	} else if f < 118 {
		return 11
	}
	return 12
}

func (f KPH) ForceColor() string {
	return ForceColor(f.Force())
}

func (f KPH) String() string {
	return fmt.Sprintf("%sm/s(%d)", unit.FormatFloat64(f.MPS(), 1), f.Force())
}

func (f KPH) MarshalJSON() ([]byte, error) {
	return json.Marshal(float64(f))
}

func (f KPH) HTML() html.HTML {
	return html.Span().Style("color:"+f.ForceColor()).Contentf("%sm/s(%d)", unit.FormatFloat64(f.MPS(), 1), f.Force()).HTML()
}
