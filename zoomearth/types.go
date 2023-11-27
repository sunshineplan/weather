package zoomearth

import "time"

type ZoomEarthAPI struct{}

type StormID string

type StormData struct {
	ID     StormID
	Name   string
	Title  string
	Active bool
	Type   string
	Place  string
	Cone   [][2]float64
	Track  []Track
	JA     int

	Coordinates [2]float64
}

type Track struct {
	Date        time.Time
	Coordinates [2]float64
	Forecast    bool
}
