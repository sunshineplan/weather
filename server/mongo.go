package main

import (
	"bytes"
	"encoding/json"
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

	if _, err = client.UpdateOne(
		mongodb.M{"DateEpoch": day.DateEpoch, "Date": day.Date},
		mongodb.M{"$set": day},
		&mongodb.UpdateOpt{Upsert: true},
	); err == nil {
		log.Printf("record %#v", day)
	}
	return
}

func export(month string, delete bool) (buf bytes.Buffer, err error) {
	var res []weather.Day
	if err = client.Find(
		mongodb.M{"Date": mongodb.M{"$regex": month}},
		&mongodb.FindOpt{Sort: mongodb.M{"Date": 1}},
		&res,
	); err != nil {
		return
	}

	buf.WriteRune('[')
	for index, i := range res {
		i.DateEpoch = 0
		i.Hours = nil
		i.Icon = ""
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
			if _, err := client.DeleteMany(mongodb.M{"Date": mongodb.M{"$regex": month}}); err != nil {
				log.Print(err)
			}
		}()
	}

	return
}
