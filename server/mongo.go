package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/weather"
	"github.com/sunshineplan/weather/api/visualcrossing"
	"github.com/sunshineplan/weather/unit"
)

func record(date time.Time) (err error) {
	defer func() {
		if err != nil {
			svc.Print(err)
		}
	}()

	svc.Printf("Start recording %s's weather...", date.Format("2006-01-02"))
	day, err := history.History(*query, date)
	if err != nil {
		return
	}

	_, err = db.UpdateOne(
		mongodb.M{"dateEpoch": day.DateEpoch, "date": day.Date},
		mongodb.M{"$set": day},
		&mongodb.UpdateOpt{Upsert: true},
	)
	return
}

func historyRecord(t time.Time, round int) (weather.Day, weather.Day, error) {
	var lastYear, avg []visualcrossing.Day
	if err := db.Aggregate(
		[]mongodb.M{
			{"$match": mongodb.M{"date": t.AddDate(-1, 0, 0).Format("2006-01-02")}},
			{"$project": mongodb.M{
				"datetime":     "$date",
				"tempmax":      1,
				"tempmin":      1,
				"temp":         1,
				"feelslikemax": 1,
				"feelslikemin": 1,
				"feelslike":    1,
			}},
		},
		&lastYear,
	); err != nil {
		return weather.Day{}, weather.Day{}, err
	}
	if n := len(lastYear); n != 1 {
		return weather.Day{}, weather.Day{}, fmt.Errorf("incorrect quantity of last year results: %d", n)
	}
	if err := db.Aggregate(
		[]mongodb.M{
			{"$match": mongodb.M{"date": mongodb.M{"$regex": t.Format("01-02") + "$"}}},
			{"$group": mongodb.M{
				"_id":          mongodb.M{"$substr": []any{"$date", 5, -1}},
				"tempmax":      mongodb.M{"$avg": "$tempmax"},
				"tempmin":      mongodb.M{"$avg": "$tempmin"},
				"temp":         mongodb.M{"$avg": "$temp"},
				"feelslikemax": mongodb.M{"$avg": "$feelslikemax"},
				"feelslikemin": mongodb.M{"$avg": "$feelslikemin"},
				"feelslike":    mongodb.M{"$avg": "$feelslike"},
			}},
			{"$project": mongodb.M{
				"datetime":     "$_id",
				"tempmax":      mongodb.M{"$round": []any{"$tempmax", round}},
				"tempmin":      mongodb.M{"$round": []any{"$tempmin", round}},
				"temp":         mongodb.M{"$round": []any{"$temp", round}},
				"feelslikemax": mongodb.M{"$round": []any{"$feelslikemax", round}},
				"feelslikemin": mongodb.M{"$round": []any{"$feelslikemin", round}},
				"feelslike":    mongodb.M{"$round": []any{"$feelslike", round}},
			}},
		},
		&avg,
	); err != nil {
		return weather.Day{}, weather.Day{}, err
	}
	if n := len(avg); n != 1 {
		return weather.Day{}, weather.Day{}, fmt.Errorf("incorrect quantity of average results: %d", n)
	}
	return visualcrossing.ConvertDays(lastYear)[0], visualcrossing.ConvertDays(avg)[0], nil
}

func export(month string, delete bool) (string, error) {
	var res []struct {
		Date         string            `json:"date"`
		TempMax      any               `json:"tempmax"`
		TempMin      any               `json:"tempmin"`
		Temp         any               `json:"temp"`
		FeelsLikeMax any               `json:"feelslikemax"`
		FeelsLikeMin any               `json:"feelslikemin"`
		FeelsLike    any               `json:"feelslike"`
		Humidity     unit.Percent      `json:"humidity"`
		Dew          any               `json:"dew"`
		Precip       float64           `json:"precip"`
		PrecipCover  float64           `json:"precipcover"`
		WindSpeed    any               `json:"windspeed"`
		Pressure     float64           `json:"pressure"`
		Visibility   float64           `json:"visibility"`
		UVIndex      unit.UVIndex      `json:"uvindex"`
		Condition    weather.Condition `json:"condition"`
		Description  string            `json:"description"`
	}
	if err := db.Find(
		mongodb.M{"date": mongodb.M{"$regex": month}},
		&mongodb.FindOpt{Sort: mongodb.M{"date": 1}},
		&res,
	); err != nil {
		return "", err
	}

	b, err := json.Marshal(res)
	if err != nil {
		return "", err
	}
	if delete {
		go func() {
			if _, err := db.DeleteMany(mongodb.M{"date": mongodb.M{"$regex": month}}); err != nil {
				svc.Print(err)
			}
		}()
	}
	return strings.ReplaceAll(string(b), "},{", "},\n{"), nil
}
