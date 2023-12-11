package main

import (
	"sync"

	"github.com/sunshineplan/weather/unit/coordinates"
)

var coordsMap sync.Map

type coords struct{ coordinates.Coordinates }

func getCoords(query string) (*coords, error) {
	if v, ok := coordsMap.Load(query); ok {
		return v.(*coords), nil
	}
	coordinates, err := realtime.Coordinates(query)
	if err != nil {
		return nil, err
	}
	coords := &coords{coordinates}
	coordsMap.Store(query, coords)
	return coords, nil
}

func (c *coords) offset(x, y float64) coords {
	return coords{coordinates.New(float64(c.Latitude())+x, float64(c.Longitude())+y)}
}
