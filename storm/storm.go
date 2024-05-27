package storm

import (
	"cmp"
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
	ID     string
	Season string
	Name   string
	Title  string
	Active bool
	Place  string
	Track  []Track
	URL    string
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

func calcCoordinates(a, b Track, t time.Time) coordinates.Coordinates {
	if a == nil || b == nil {
		return nil
	}
	start, end, unix := a.Date().Unix(), b.Date().Unix(), t.Unix()
	rate := float64(unix-start) / float64(end-start)
	return coordinates.New(
		float64(a.Coordinates().Latitude())+float64(b.Coordinates().Latitude()-a.Coordinates().Latitude())*rate,
		float64(a.Coordinates().Longitude())+float64(b.Coordinates().Longitude()-a.Coordinates().Longitude())*rate,
	)
}

func (storm Data) Coordinates(t time.Time) coordinates.Coordinates {
	if len(storm.Track) == 0 {
		return nil
	}
	var track Track
	for _, i := range storm.Track {
		switch cmp.Compare(i.Date().Unix(), t.Unix()) {
		case 0:
			return i.Coordinates()
		case -1:
			track = i
		case 1:
			return calcCoordinates(track, i, t)
		}
	}
	return nil
}
