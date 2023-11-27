package storm

import "time"

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
	Coordinates [2]float64
	URL         string
}

type Track interface {
	Date() time.Time
	Coordinates() [2]float64
	Forecast() bool
}
