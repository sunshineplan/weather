package weather

import (
	"fmt"
	"math"
	"strings"
)

type Temperature float64

func (t Temperature) AbsDiff(i Temperature) float64 {
	return math.Abs(float64(t - i))
}

func (t Temperature) String() string {
	return fmt.Sprintf("%g°C", t)
}

type Percent float64

func (p Percent) Max(i Percent) Percent {
	return Percent(math.Max(float64(p), float64(i)))
}

func (p Percent) String() string {
	return fmt.Sprintf("%g%%", p)
}

type Condition string

func (c Condition) Short() string {
	s := strings.Split(string(c), ",")[0]
	switch s {
	case "Partially cloudy":
		return "Partly"
	default:
		return s
	}
}
