package main

import (
	"log"
	"time"

	"github.com/sunshineplan/weather"
)

var (
	rainSnow     *weather.RainSnow
	tempRiseFall *weather.TempRiseFall
)

func alert() {
	if rainSnow != nil {
		if rainSnow.Start().IsExpired() && !rainSnow.End().IsExpired() {
			rainSnow.Start().Date = time.Now().Format("2006-01-02")
		}
	}

	if res, err := forecast.WillRainSnow(*query, *days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		var alert bool
		var output string
		for index, i := range res {
			if index == 1 {
				defer func() {
					rainSnow = i
				}()
				if (rainSnow == nil || rainSnow.Start().Date != i.Start().Date) ||
					(rainSnow.End() == nil && i.End() != nil) || (rainSnow.End() != nil && i.End() == nil) ||
					(rainSnow.End() != nil && i.End() != nil && rainSnow.End().Date != i.End().Date) {
					alert = true
				}
			}
			output += i.String()
		}
		if alert {
			log.Print(output) //TODO
			//go sendMail(output)
		}
	} else {
		rainSnow = nil
	}

	if res, err := forecast.WillTempRiseFall(*difference, *query, *days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		var alert bool
		var output string
		for index, i := range res {
			if index == 1 {
				defer func() {
					tempRiseFall = i
				}()
				if tempRiseFall == nil || tempRiseFall.Day().Date != i.Day().Date {
					alert = true
				}
			}
			output += i.String()
		}
		if alert {
			log.Print(output) //TODO
			//go sendMail(output)
		}
	} else {
		tempRiseFall = nil
	}
}
