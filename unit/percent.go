package unit

import (
	"fmt"
	"math"
)

type Percent float64

func (f Percent) Max(i Percent) Percent {
	return Percent(math.Max(float64(f), float64(i)))
}

func (f Percent) String() string {
	return fmt.Sprintf("%g%%", f)
}
