package zoomearth

import (
	"time"

	"github.com/sunshineplan/weather/unit/coordinates"
)

type ZoomEarthAPI struct{}

type StormID string

type StormData struct {
	ID     StormID
	Name   string
	Title  string
	Active bool
	Type   string
	Place  string
	Cone   []coordinates.LongLat
	Track  []Track
	JA     int

	Coordinates coordinates.Coordinates
}

type Track struct {
	Date        time.Time
	Coordinates coordinates.LongLat
	Forecast    bool
}

type MapOptions struct {
	zoom     float64
	quality  int
	overlays []string
}
