package main

import (
	"log"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/weather"
)

func record() {
	resp, err := weather.HistoryWeather(*query, time.Now().AddDate(0, 0, -1))
	if err != nil {
		if *debug {
			log.Print(err)
		}
		return
	}

	for _, i := range resp.Forecast.Forecastday {
		res, err := client.UpdateOne(
			mongodb.M{"date_epoch": i.DateEpoch, "date": i.Date},
			mongodb.M{"$set": i.Day},
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
