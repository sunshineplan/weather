package unit

import (
	"fmt"
	"strings"
)

var _ Temperature = Celsius(0)

type Temperature interface {
	Float64() float64
	Difference(Temperature) Temperature
	String() string
	DiffHTML() string
}

type Celsius float64

func (f Celsius) Float64() float64 {
	return float64(f)
}

func (f Celsius) String() string {
	return fmt.Sprintf("%.1f°C", f)
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
