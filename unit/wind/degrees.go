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
	return unit.FormatFloat64(f.Degree(), 1) + "°"
}

func (f Degree) SVG() *html.Element {
	return html.Svg().
		Attribute("viewBox", "0 0 1000 1000").
		Style(fmt.Sprintf("transform:rotate(%fdeg)", f.Degree())).
		AppendChild(
			html.NewElement("g").AppendChild(
				html.NewElement("path").Attribute("d", "M510.5,749.6c-14.9-9.9-38.1-9.9-53.1,1.7l-262,207.3c-14.9,11.6-21.6,6.6-14.9-11.6L474,48.1c5-16.6,14.9-18.2,21.6,0l325,898.7c6.6,16.6-1.7,23.2-14.9,11.6L510.5,749.6z"),
				html.NewElement("path").Attribute("d", "M817.2,990c-8.3,0-16.6-3.3-26.5-9.9L497.2,769.5c-5-3.3-18.2-3.3-23.2,0L210.3,976.7c-19.9,16.6-41.5,14.9-51.4,0c-6.6-9.9-8.3-21.6-3.3-38.1L449.1,39.8C459,13.3,477.3,10,483.9,10c6.6,0,24.9,3.3,34.8,29.8l325,898.7c5,14.9,5,28.2-1.7,38.1C837.1,985,827.2,990,817.2,990z M485.6,716.4c14.9,0,28.2,5,39.8,11.6l255.4,182.4L485.6,92.9l-267,814.2l223.9-177.4C454.1,721.4,469,716.4,485.6,716.4z"),
			),
		)
}
