package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/sunshineplan/weather"
)

var (
	rainSnow     *weather.RainSnow
	tempRiseFall *weather.TempRiseFall
)

func alert() {
	days, err := forecast.Forecast(*query, *days)
	if err != nil {
		log.Print(err)
		return
	}

	go runAlert(days, alertRainSnow)
	go runAlert(days, alertTempRiseFall)
}

func runAlert(days []weather.Day, fn func([]weather.Day) (string, strings.Builder)) {
	if subject, body := fn(days); subject != "" {
		log.Print(subject)
		go sendMail(subject, body.String())
	}
}

func alertRainSnow(days []weather.Day) (subject string, body strings.Builder) {
	if rainSnow != nil {
		if rainSnow.IsExpired() {
			rainSnow = nil
		}
	}

	if res, err := weather.WillRainSnow(days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		var first weather.RainSnow
		for index, i := range res {
			if index == 0 {
				first = i
				if rainSnow == nil || rainSnow.Start().Date != i.Start().Date || rainSnow.Duration() != i.Duration() {
					if rainSnow != nil {
						log.Println(rainSnow.Start().Date, i.Start().Date) //test
					}
					subject = fmt.Sprintf("[Weather]Rain Snow Alert - %s", i.Start().Date)
				}
			}
			fmt.Fprintln(&body, i.String())
		}
		rainSnow = &first
	} else {
		rainSnow = nil
	}
	return
}

func alertTempRiseFall(days []weather.Day) (subject string, body strings.Builder) {
	if tempRiseFall != nil {
		if tempRiseFall.IsExpired() {
			tempRiseFall = nil
		}
	}

	if res, err := weather.WillTempRiseFall(days, *difference); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		var first weather.TempRiseFall
		for index, i := range res {
			if index == 0 {
				first = i
				if tempRiseFall == nil || tempRiseFall.Day().Date != i.Day().Date {
					if i.IsRise() {
						subject = fmt.Sprintf("[Weather]Temperature Rise Alert - %s", i.Day().Date)
					} else {
						subject = fmt.Sprintf("[Weather]Temperature Fall Alert - %s", i.Day().Date)
					}
				}
			}
			fmt.Fprintln(&body, i.String())
		}
		tempRiseFall = &first
	} else {
		tempRiseFall = nil
	}
	return
}
