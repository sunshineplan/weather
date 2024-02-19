package wind

import (
	"fmt"

	"github.com/sunshineplan/utils/html"
	"github.com/sunshineplan/weather/unit"
)

const (
	N Degree = iota * Degree(360) / 32
	NbE
	NNE
	NEbN
	NE
	NEbE
	ENE
	EbN
	E
	EbS
	ESE
	SEbE
	SE
	SEbS
	SSE
	SbE
	S
	SbW
	SSW
	SWbS
	SW
	SWbW
	WSW
	WbS
	W
	WbN
	WNW
	NWbW
	NW
	NWbN
	NNW
	NbW
)

type Degree float64

func (f Degree) Degree() Degree {
	for f < 0 {
		f += 360
	}
	for f >= 360 {
		f -= 360
	}
	return f
}

const difference = Degree(360) / 32 / 2

func (f Degree) Direction() Direction {
	for k, v := range chart {
		if f := f.Degree(); f >= v-difference && f < v+difference {
			return k
		}
	}
	return "N"
}

func (f Degree) String() string {
	return fmt.Sprintf("%s°(%s)", unit.FormatFloat64(f.Degree(), 1), f.Direction())
}

func (f Degree) HTML() html.HTML {
	return html.Span().Style(fmt.Sprintf("transform:rotate(%sdeg)", unit.FormatFloat64((f+180).Degree(), 1))).Content("⬆").HTML()
}
