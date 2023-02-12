package main

import (
	"bytes"
	"log"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/weather"
)

func record(date time.Time) (err error) {
	defer func() {
		if err != nil {
			log.Print(err)
		}
	}()

	day, err := history.History(*query, date)
	if err != nil {
		return
	}

	_, err = client.UpdateOne(
		mongodb.M{"dateEpoch": day.DateEpoch, "date": day.Date},
		mongodb.M{"$set": day},
		&mongodb.UpdateOpt{Upsert: true},
	)
	return
}

func export(month string, delete bool) (buf bytes.Buffer, err error) {
	var res []weather.Day
	if err = client.Find(
		mongodb.M{"date": mongodb.M{"$regex": month}},
		&mongodb.FindOpt{Sort: mongodb.M{"date": 1}},
		&res,
	); err != nil {
		return
	}

	buf.WriteRune('[')
	for index, i := range res {
		buf.WriteString(i.JSON())
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
