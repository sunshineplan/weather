package wind

import "github.com/sunshineplan/utils/html"

var chart = map[Direction]Degree{
	"N":    N,
	"NbE":  NbE,
	"NNE":  NNE,
	"NEbN": NEbN,
	"NE":   NE,
	"NEbE": NEbE,
	"ENE":  ENE,
	"EbN":  EbN,
	"E":    E,
	"EbS":  EbS,
	"ESE":  ESE,
	"SEbE": SEbE,
	"SE":   SE,
	"SEbS": SEbS,
	"SSE":  SSE,
	"SbE":  SbE,
	"S":    S,
	"SbW":  SbW,
	"SSW":  SSW,
	"SWbS": SWbS,
	"SW":   SW,
	"SWbW": SWbW,
	"WSW":  WSW,
	"WbS":  WbS,
	"W":    W,
	"WbN":  WbN,
	"WNW":  WNW,
	"NWbW": NWbW,
	"NW":   NW,
	"NWbN": NWbN,
	"NNW":  NNW,
	"NbW":  NbW,
}

type Direction string

func (s Direction) Degree() Degree {
	if f, ok := chart[s]; ok {
		return f
	}
	panic("unknown direction: " + s)
}

func (s Direction) Direction() Direction {
	return s.Degree().Direction()
}

func (s Direction) String() string {
	return string(s.Direction())
}

func (s Direction) SVG() *html.Element {
	return s.Degree().SVG()
}
