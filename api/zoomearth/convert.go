package zoomearth

import "github.com/sunshineplan/weather/storm"

func (data stormData) Convert() storm.Data {
	return storm.Data{
		ID:     string(data.ID),
		Season: data.Season,
		Name:   data.Name,
		Title:  data.Title,
		Active: data.Active,
		Place:  data.Place,
		Track:  ConvertTracks(data.Track),
		URL:    data.ID.URL(),
	}
}

func ConvertTracks(src []track) (tracks []storm.Track) {
	for _, i := range src {
		tracks = append(tracks, i)
	}
	return
}
