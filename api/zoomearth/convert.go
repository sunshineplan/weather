package zoomearth

import (
	"time"

	"github.com/sunshineplan/weather/storm"
	"github.com/sunshineplan/weather/unit/coordinates"
)

func (data StormData) Convert() storm.Data {
	return storm.Data{
		ID:          string(data.ID),
		Name:        data.Name,
		Title:       data.Title,
		Active:      data.Active,
		Place:       data.Place,
		Track:       ConvertTracks(data.Track),
		Coordinates: data.Coordinates,
		URL:         data.ID.URL(),
	}
}

var _ storm.Track = track{}

type track struct {
	d time.Time
	c coordinates.Coordinates
	f bool
}

func (track track) Date() time.Time {
	return track.d
}

func (track track) Coordinates() coordinates.Coordinates {
	return track.c
}

func (track track) Forecast() bool {
	return track.f
}

func ConvertTrack(src Track) storm.Track {
	return track{src.Date, src.Coordinates, src.Forecast}
}

func ConvertTracks(src []Track) (tracks []storm.Track) {
	for _, i := range src {
		tracks = append(tracks, ConvertTrack(i))
	}
	return
}
