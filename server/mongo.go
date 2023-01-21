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
		if *debug {
			log.Print(err)
		}
		return
	}

	for _, i := range resp.Forecast.Forecastday {
		res, err := client.UpdateOne(
			mongodb.M{"date_epoch": i.DateEpoch, "date": i.Date},
			mongodb.M{"$set": mongodb.M{"day": i.Day}},
			&mongodb.UpdateOpt{Upsert: true},
		)
		if err != nil {
			if *debug {
				log.Print(err)
			}
			return
		}

		if n := res.MatchedCount; n != 0 && *debug {
			log.Printf("Updated %d record", n)
		}
		if n := res.UpsertedCount; n != 0 && *debug {
			log.Printf("Upserted %d record", n)
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
		i.DateEpoch = 0
		if i.Day.Condition != nil {
			i.Day.Condition.Icon = ""
			i.Day.Condition.Code = 0
		}
		b, err := json.Marshal(i)
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
