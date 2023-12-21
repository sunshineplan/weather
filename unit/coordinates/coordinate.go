package coordinates

import (
	"fmt"
	"math"
	"strings"

	"github.com/sunshineplan/weather/unit"
)

type Latitude float64

type Longitude float64

func (lat Latitude) String() string {
	return newDMS(lat).str(true)
}

func (long Longitude) String() string {
	return newDMS(long).str(false)
}

type dms struct {
	degrees  int
	minutes  int
	seconds  float64
	negative bool
}

func newDMS[T ~float64](f T) dms {
	var negative bool
	if f < 0 {
		negative = true
	}
	abs := math.Abs(float64(f))
	degrees := int(abs)
	decimalPart := abs - float64(degrees)
	minutes := int(decimalPart * 60)
	seconds := (decimalPart*60 - float64(minutes)) * 60
	return dms{degrees, minutes, seconds, negative}
}

// https://en.wikipedia.org/wiki/ISO_6709
func (dms dms) str(lat bool) string {
	var direction string
	if lat {
		if dms.degrees > 90 {
			panic(fmt.Sprintf("invalid latitude: %d", dms.degrees))
		}
		if dms.negative {
			direction = "S"
		} else {
			direction = "N"
		}
	} else {
		if dms.degrees > 180 {
			panic(fmt.Sprintf("invalid longitude: %d", dms.degrees))
		}
		if dms.negative {
			direction = "W"
		} else {
			direction = "E"
		}
	}
	var seconds, dec string
	s := strings.Split(unit.FormatFloat64(dms.seconds, 3), ".")
	seconds = s[0]
	if len(s) == 2 {
		dec = "." + s[1]
	}
	return fmt.Sprintf(`%d°%02d′%02s%s″%s`, dms.degrees, dms.minutes, seconds, dec, direction)
}
