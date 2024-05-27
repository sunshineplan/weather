package zoomearth

import "time"

type ZoomEarthAPI struct{}

type StormID string

type MapOptions struct {
	width    int
	height   int
	zoom     float64
	overlays []string
	timezone *time.Location
}
