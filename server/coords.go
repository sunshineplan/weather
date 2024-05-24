package main

import (
	"sync"

	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit/coordinates"
)

var coordsMap sync.Map

func getCoords(query string, api coordinates.GeoLocator) (res coordinates.Coordinates, err error) {
	if v, ok := coordsMap.Load(query); ok {
		res = v.(coordinates.Coordinates)
		return
	}
	if api != nil {
		res, err = api.Coordinates(query)
	} else {
		res, err = forecast.Coordinates(query)
		if err != nil {
			res, err = realtime.Coordinates(query)
		}
	}
	if err != nil {
		return
	}
	coordsMap.Store(query, res)
	return
}

func willAffect(storm storm.Data, coords coordinates.Coordinates, radius float64) (affect, future bool) {
	if !storm.Active {
		return
	}
	for _, i := range storm.Track {
		if coordinates.Distance(i.Coordinates(), coords) <= radius {
			affect = true
			if i.Forecast() {
				future = true
				break
			}
		}
	}
	return
}
