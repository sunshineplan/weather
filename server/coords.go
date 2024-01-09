package main

import (
	"sync"

	"github.com/sunshineplan/weather/unit/coordinates"
)

var coordsMap sync.Map

type coords struct{ coordinates.Coordinates }

func getCoords(query string, api coordinates.GeoLocator) (res *coords, err error) {
	if v, ok := coordsMap.Load(query); ok {
		res = v.(*coords)
		return
	}
	var c coordinates.Coordinates
	if api != nil {
		c, err = api.Coordinates(query)
	} else {
		c, err = forecast.Coordinates(query)
		if err != nil {
			c, err = realtime.Coordinates(query)
		}
	}
	if err != nil {
		return
	}
	res = &coords{c}
	coordsMap.Store(query, res)
	return
}

func (c *coords) offset(x, y float64) coords {
	return coords{coordinates.New(float64(c.Latitude())+x, float64(c.Longitude())+y)}
}
