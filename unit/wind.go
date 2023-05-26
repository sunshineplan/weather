package unit

import (
	"encoding"
	"fmt"
)

var (
	_ encoding.TextMarshaler = WindSpeed(nil)

	_ WindSpeed = WindKPH(0)
)

type WindSpeed interface {
	MPS() float64
	Force() int
	ForceColor() string
	String() string
	MarshalText() ([]byte, error)
	HTML() string
}

func WindForceColor(force int) string {
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

type WindKPH float64

func (f WindKPH) MPS() float64 {
	return float64(f * 1000 / 60 / 60)
}

func (f WindKPH) Force() int {
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

func (f WindKPH) ForceColor() string {
	return WindForceColor(f.Force())
}

func (f WindKPH) String() string {
	return fmt.Sprintf("%sm/s(%d)", formatFloat64(f.MPS(), 1), f.Force())
}

func (f WindKPH) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f WindKPH) HTML() string {
	return fmt.Sprintf(`<span style="color:%s">%sm/s(%d)</span>`, f.ForceColor(), formatFloat64(f.MPS(), 1), f.Force())
}
