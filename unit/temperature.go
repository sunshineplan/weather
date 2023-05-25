package unit

import (
	"encoding"
	"fmt"
	"strings"
)

var (
	_ encoding.TextMarshaler = Temperature(nil)

	_ Temperature = Celsius(0)
)

type Temperature interface {
	Float64() float64
	Difference(Temperature) Temperature
	String() string
	MarshalText() ([]byte, error)
	DiffHTML() string
}

type Celsius float64

func (f Celsius) Float64() float64 {
	return float64(f)
}

func (f Celsius) String() string {
	return fmt.Sprintf("%s°C", formatFloat64(f, 1))
}

func (f Celsius) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f Celsius) Difference(i Temperature) Temperature {
	return Celsius(f.Float64() - i.Float64())
}

func (f Celsius) DiffHTML() string {
	var b strings.Builder
	if f > 0 {
		fmt.Fprint(&b, `<span style="color:red">`, f, "↑")
	} else if f < 0 {
		fmt.Fprint(&b, `<span style="color:green">`, -f, "↓")
	} else {
		fmt.Fprint(&b, "<span>", f)
	}
	b.WriteString("</span>")
	return b.String()
}
