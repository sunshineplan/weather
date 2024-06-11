package maps

import (
	"errors"
	"fmt"
	"image"
	"time"

	"github.com/sunshineplan/weather/unit/coordinates"
)

type MapType int

const (
	Satellite MapType = iota + 1
	Radar
	Precipitation
	Wind
	Temperature
	Humidity
	DewPoint
	Pressure
)

type MapAPI interface {
	URL(mt MapType, t time.Time, coords coordinates.Coordinates, opt any) string
	Map(mt MapType, t time.Time, coords coordinates.Coordinates, opt any) (time.Time, image.Image, error)
	Realtime(mt MapType, coords coordinates.Coordinates, opt any) (time.Time, image.Image, error)
}

var ErrInsufficientColor = errors.New("image has insufficient color depth")

var _ error = InsufficientColor(0)

type InsufficientColor int

func (i InsufficientColor) Error() string {
	return fmt.Sprintf("image has insufficient color depth: %d", i)
}
func (i InsufficientColor) Is(target error) bool { return target == ErrInsufficientColor }
func (i InsufficientColor) Unwrap() error        { return ErrInsufficientColor }