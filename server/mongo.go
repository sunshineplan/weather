package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/weather"
)

func record(day time.Time) {
	resp, err := weather.HistoryWeather(*query, day)
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

func export(month time.Time, delete bool) {
	if err := initWeather(); err != nil {
		log.Fatal(err)
	}

	var res []weather.ForecastForecastday
	if err := client.Find(
		mongodb.M{"date": mongodb.M{"$regex": month.Format("2006-01")}},
		&mongodb.FindOpt{Sort: mongodb.M{"date": 1}},
		&res,
	); err != nil {
		log.Fatal(err)
	}

	var output string
	for index, i := range res {
		i.DateEpoch = 0
		i.Day.Condition.Icon = ""
		b, err := json.Marshal(i)
		if err != nil {
			log.Print(err)
			continue
		}
		output += string(b)
		if index < len(res)-1 {
			output += ",\n"
		}
	}
	output = fmt.Sprintf("[%s]", output)
	if err := os.WriteFile(fmt.Sprintf("%s.json", month.Format("2006-01")), []byte(output), 0644); err != nil {
		log.Fatal(err)
	}

	if delete {
		if _, err := client.DeleteMany(mongodb.M{"date": mongodb.M{"$regex": month.Format("2006-01")}}); err != nil {
			log.Fatal(err)
		}
	}
}
