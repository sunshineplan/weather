package unit

import (
	"encoding/json"
	"fmt"

	"github.com/sunshineplan/utils/html"
)

var (
	_ json.Marshaler = Temperature(nil)

	_ Temperature = Celsius(0)
)

type Temperature interface {
	Float64() float64
	Difference(Temperature) Temperature
	String() string
	MarshalJSON() ([]byte, error)
	HTML() html.HTML
	DiffHTML() html.HTML
}

type Celsius float64

func (f Celsius) Float64() float64 {
	return float64(f)
}

func (f Celsius) String() string {
	return fmt.Sprintf("%s°C", FormatFloat64(f, 1))
}

func (f Celsius) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Float64())
}

func (f Celsius) Difference(i Temperature) Temperature {
	return Celsius(f.Float64() - i.Float64())
}

func (f Celsius) HTML() html.HTML {
	span := html.Span().Content(f.String())
	if f <= 0 {
		span.Style("color:blue")
	}
	return span.HTML()
}

func (f Celsius) DiffHTML() html.HTML {
	if f > 0 {
		return html.Span().Style("color:red").Contentf("%s↑", f).HTML()
	} else if f < 0 {
		return html.Span().Style("color:green").Contentf("%s↓", -f).HTML()
	} else {
		return html.Span().Content("0°C").HTML()
	}
}
