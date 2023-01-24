package main

import (
	"bytes"
	"encoding/json"
	"log"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/weather"
)

func record(date time.Time) {
	resp, err := weather.HistoryWeather(*query, date)
	if err != nil {
		log.Print(err)
		return
	}

	for _, i := range resp.Forecast.Forecastday {
		if _, err := client.UpdateOne(
			mongodb.M{"date_epoch": i.DateEpoch, "date": i.Date},
			mongodb.M{"$set": mongodb.M{"day": i.Day}},
			&mongodb.UpdateOpt{Upsert: true},
		); err != nil {
			log.Print(err)
		} else {
			log.Printf("record %s %#v", i.Date, i.Day)
		}
	}
}

func export(month string, delete bool) (buf bytes.Buffer, err error) {
	var res []weather.ForecastForecastday
	if err = client.Find(
		mongodb.M{"date": mongodb.M{"$regex": month}},
		&mongodb.FindOpt{Sort: mongodb.M{"date": 1}},
		&res,
	); err != nil {
		return
	}

	buf.WriteRune('[')
	for index, i := range res {
		i.Day.Date = i.Date
		if i.Day.Condition != nil {
			i.Day.Weather = i.Day.Condition.Text
			i.Day.Condition = nil
		}
		b, err := json.Marshal(i.Day)
		if err != nil {
			log.Print(err)
			continue
		}
		buf.Write(b)
		if index < len(res)-1 {
			buf.WriteString(",\n")
		}
	}
	buf.WriteRune(']')

	if delete {
		go func() {
			if _, err := client.DeleteMany(mongodb.M{"date": mongodb.M{"$regex": month}}); err != nil {
				log.Print(err)
			}
		}()
	}

	return
}
