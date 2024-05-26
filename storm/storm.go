package storm

import (
	"time"

	"github.com/sunshineplan/weather/unit/coordinates"
)

type API interface {
	GetStorms(t time.Time) ([]Storm, error)
}

type Storm interface {
	Data() (Data, error)
}

type Data struct {
	ID          string
	Name        string
	Title       string
	Active      bool
	Place       string
	Track       []Track
	Coordinates coordinates.Coordinates
	URL         string
}

type Track interface {
	Date() time.Time
	Coordinates() coordinates.Coordinates
	Forecast() bool
}

func (storm Data) Affect(coords coordinates.Coordinates, radius float64) (affect, future bool) {
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
