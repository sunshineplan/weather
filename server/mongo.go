package main

import (
	"bytes"
	"fmt"
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

func average(date string, round int) (weather.Day, error) {
	var res []weather.Day
	if err := client.Aggregate(
		[]mongodb.M{
			{"$match": mongodb.M{"date": mongodb.M{"$regex": date}}},
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
				"tempmax":      mongodb.M{"$round": []any{"$tempmax", round}},
				"tempmin":      mongodb.M{"$round": []any{"$tempmin", round}},
				"temp":         mongodb.M{"$round": []any{"$temp", round}},
				"feelslikemax": mongodb.M{"$round": []any{"$feelslikemax", round}},
				"feelslikemin": mongodb.M{"$round": []any{"$feelslikemin", round}},
				"feelslike":    mongodb.M{"$round": []any{"$feelslike", round}},
			}},
		},
		&res,
	); err != nil {
		return weather.Day{}, err
	}
	if n := len(res); n != 1 {
		return weather.Day{}, fmt.Errorf("incorrect quantity of average results: %d", n)
	}
	return res[0], nil
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
