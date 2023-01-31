package main

import (
	"fmt"
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
		var title, output string
		for index, i := range res {
			if index == 1 {
				defer func() {
					rainSnow = i
				}()
				if (rainSnow == nil || rainSnow.Start().Date != i.Start().Date) ||
					(rainSnow.End() == nil && i.End() != nil) || (rainSnow.End() != nil && i.End() == nil) ||
					(rainSnow.End() != nil && i.End() != nil && rainSnow.End().Date != i.End().Date) {
					alert = true
					title = fmt.Sprintf("[Weather]Rain Snow Alert - %s", i.Start().Date)
				}
			}
			output += i.String()
		}
		if alert {
			log.Print(output)
			go sendMail(title, output)
		}
	} else {
		rainSnow = nil
	}

	if res, err := forecast.WillTempRiseFall(*difference, *query, *days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		var alert bool
		var title, output string
		for index, i := range res {
			if index == 1 {
				defer func() {
					tempRiseFall = i
				}()
				if tempRiseFall == nil || tempRiseFall.Day().Date != i.Day().Date {
					alert = true
					if i.IsRise() {
						title = fmt.Sprintf("[Weather]Temperature Rise Alert - %s", i.Day().Date)
					} else {
						title = fmt.Sprintf("[Weather]Temperature Fall Alert - %s", i.Day().Date)
					}
				}
			}
			output += i.String()
		}
		if alert {
			log.Print(output)
			go sendMail(title, output)
		}
	} else {
		tempRiseFall = nil
	}
}
