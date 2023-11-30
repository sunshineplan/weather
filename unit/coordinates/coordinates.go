package coordinates

import (
	"fmt"
	"math"
)

type GeoLocator interface {
	Coordinates(string) (Coordinates, error)
}

type Coordinates interface {
	Latitude() Latitude
	Longitude() Longitude
	String() string
}

func New(lat, long float64) Coordinates {
	return LatLong{lat, long}
}

func degreesToRadians[T ~float64](f T) float64 {
	return float64(f) * math.Pi / 180
}

func Distance(i, j Coordinates) float64 {
	lat1, long1 := degreesToRadians(i.Latitude()), degreesToRadians(i.Longitude())
	lat2, long2 := degreesToRadians(j.Latitude()), degreesToRadians(j.Longitude())

	dLat := lat2 - lat1
	dLong := long2 - long1

	a := math.Pow(math.Sin(dLat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dLong/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return c * 6371
}

var (
	_ Coordinates = LatLong{}
	_ Coordinates = LongLat{}
)

type LatLong [2]float64

type LongLat [2]float64

func (c LatLong) Latitude() Latitude {
	return Latitude(c[0])
}

func (c LatLong) Longitude() Longitude {
	return Longitude(c[1])
}

func (c LatLong) String() string {
	return fmt.Sprintf("%s, %s", c.Latitude(), c.Longitude())
}

func (c LongLat) Latitude() Latitude {
	return Latitude(c[1])
}

func (c LongLat) Longitude() Longitude {
	return Longitude(c[0])
}

func (c LongLat) String() string {
	return fmt.Sprintf("%s, %s", c.Latitude(), c.Longitude())
}
