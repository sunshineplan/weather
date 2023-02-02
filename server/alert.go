package main

import (
	"fmt"
	"log"
	"strings"
	"time"

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
		if rainSnow.Start().IsExpired() && !rainSnow.End().IsExpired() {
			rainSnow.Start().Date = time.Now().Format("2006-01-02")
		}
	}

	if res, err := weather.WillRainSnow(days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		for index, i := range res {
			if index == 1 {
				defer func() {
					rainSnow = &i
				}()
				if (rainSnow == nil || rainSnow.Start().Date != i.Start().Date) ||
					rainSnow.Duration() != i.Duration() {
					subject = fmt.Sprintf("[Weather]Rain Snow Alert - %s", i.Start().Date)
				}
			}
			fmt.Fprintln(&body, i.String())
		}
	} else {
		rainSnow = nil
	}
	return
}

func alertTempRiseFall(days []weather.Day) (subject string, body strings.Builder) {
	if res, err := weather.WillTempRiseFall(days, *difference); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		for index, i := range res {
			if index == 1 {
				defer func() {
					tempRiseFall = &i
				}()
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
	} else {
		tempRiseFall = nil
	}
	return
}
