package airmatters

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/sunshineplan/weather/aqi"
	"github.com/sunshineplan/weather/unit"
	"github.com/sunshineplan/weather/unit/coordinates"
)

const baseURL = "https://api.air-matters.app"

var _ aqi.API = &AirMatters{}

type AirMatters struct{ key string }

func New(key string) *AirMatters { return &AirMatters{key} }

func (api *AirMatters) Request(endpoint string, query url.Values, data any) error {
	if !query.Has("lang") {
		query.Set("lang", "en")
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?%s", baseURL, endpoint, query.Encode()), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", api.key)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(data)
}

func (api *AirMatters) Standard(t aqi.Type) ([]aqi.AQI, error) {
	standard, ok := typeMap[t]
	if !ok {
		return nil, errors.New("aqi type not supported")
	}
	var res Standard
	if err := api.Request("standard", url.Values{"standard": {standard}}, &res); err != nil {
		return nil, err
	}
	if len(res.BreakPoint) == 0 {
		return nil, errors.New("no result")
	}
	if aqi := res.Standard(); aqi[0].Type() != t {
		return nil, errors.New("bad result")
	} else {
		return aqi, nil
	}
}

func (api *AirMatters) Places(query string) ([]Place, error) {
	var places struct{ Places []Place }
	if err := api.Request("place_search", url.Values{"content": {query}}, &places); err != nil {
		return nil, err
	}
	if len(places.Places) == 0 {
		return nil, errors.New("no result")
	}
	return places.Places, nil
}

var PlaceMap sync.Map

func (api *AirMatters) Coordinates(query string) (coordinates.Coordinates, error) {
	if v, ok := PlaceMap.Load(query); ok {
		return coordinates.New(v.(Place).Latitude, v.(Place).Longitude), nil
	}
	places, err := api.Places(query)
	if err != nil {
		return nil, err
	}
	var locality bool
	var place Place
	for _, i := range places {
		if i.Type == "locality" {
			place = i
			locality = true
		}
	}
	if !locality {
		place = places[0]
	}
	PlaceMap.Store(query, place)
	return coordinates.New(place.Latitude, place.Longitude), nil
}

func (api *AirMatters) Realtime(aqiType aqi.Type, query string) (aqi.Current, error) {
	standard, ok := typeMap[aqiType]
	if !ok {
		return nil, errors.New("aqi type not supported")
	}
	var id string
	if v, ok := PlaceMap.Load(query); ok {
		id = v.(Place).ID
	} else if _, err := api.Coordinates(query); err != nil {
		return nil, err
	} else {
		v, _ := PlaceMap.Load(query)
		id = v.(Place).ID
	}
	var res struct{ Latest Current }
	if err := api.Request("current_air_condition", url.Values{
		"place_id": {id},
		"standard": {standard},
	}, &res); err != nil {
		return nil, err
	}
	res.Latest.AQIType = aqiType
	for i := range res.Latest.Items {
		res.Latest.Items[i].AQIType = aqiType
	}
	if res.Latest.AQI() == nil {
		return nil, errors.New("no AQI index value found")
	}
	return res.Latest, nil
}

func (api *AirMatters) RealtimeNearby(aqiType aqi.Type, query coordinates.Coordinates) (Place, aqi.Current, error) {
	standard, ok := typeMap[aqiType]
	if !ok {
		return Place{}, nil, errors.New("aqi type not supported")
	}
	var res struct {
		Place  Place
		Latest Current
	}
	if err := api.Request("nearby_air_condition", url.Values{
		"lat":      {unit.FormatFloat64(query.Latitude(), 2)},
		"lon":      {unit.FormatFloat64(query.Longitude(), 2)},
		"standard": {standard},
	}, &res); err != nil {
		return Place{}, nil, err
	}
	res.Latest.AQIType = aqiType
	for i := range res.Latest.Items {
		res.Latest.Items[i].AQIType = aqiType
	}
	return res.Place, res.Latest, nil
}

func (api *AirMatters) Forecast(aqiType aqi.Type, query string, _ int) (days []aqi.Day, err error) {
	standard, ok := typeMap[aqiType]
	if !ok {
		return nil, errors.New("aqi type not supported")
	}
	var id string
	if v, ok := PlaceMap.Load(query); ok {
		id = v.(Place).ID
	} else if _, err = api.Coordinates(query); err != nil {
		return
	} else {
		v, _ := PlaceMap.Load(query)
		id = v.(Place).ID
	}
	var forecast struct{ Forecast []struct{ Data []Item } }
	if err = api.Request("aqi_forecast", url.Values{
		"place_id": {id},
		"standard": {standard},
	}, &forecast); err != nil {
		return
	}
	if len(forecast.Forecast) == 0 || len(forecast.Forecast[0].Data) == 0 {
		return nil, errors.New("no result")
	}
	for _, i := range forecast.Forecast[0].Data {
		i.AQIType = aqiType
		days = append(days, i)
	}
	return
}

func (api *AirMatters) History(aqiType aqi.Type, query string, t time.Time) (aqi.Day, error) {
	standard, ok := typeMap[aqiType]
	if !ok {
		return nil, errors.New("aqi type not supported")
	}
	var id string
	if v, ok := PlaceMap.Load(query); ok {
		id = v.(Place).ID
	} else if _, err := api.Coordinates(query); err != nil {
		return nil, err
	} else {
		v, _ := PlaceMap.Load(query)
		id = v.(Place).ID
	}
	date := t.Format("2006-01-02")
	var history struct{ History []struct{ Data []Item } }
	if err := api.Request("history_air_condition", url.Values{
		"place_id":    {id},
		"standard":    {standard},
		"daily_start": {date},
		"daily_end":   {date},
	}, &history); err != nil {
		return nil, err
	}
	if len(history.History) == 0 || len(history.History[0].Data) == 0 {
		return nil, errors.New("no result")
	}
	day := history.History[0].Data[0]
	day.AQIType = aqiType
	return day, nil
}
