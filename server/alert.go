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
	if subject, body, alert := alertRainSnow(); alert {
		log.Print(subject)
		go sendMail(subject, body.String())
	}

	if subject, body, alert := alertTempRiseFall(); alert {
		log.Print(subject)
		go sendMail(subject, body.String())
	}
}

func alertRainSnow() (subject string, body strings.Builder, alert bool) {
	if rainSnow != nil {
		if rainSnow.Start().IsExpired() && !rainSnow.End().IsExpired() {
			rainSnow.Start().Date = time.Now().Format("2006-01-02")
		}
	}

	if res, err := forecast.WillRainSnow(*query, *days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		for index, i := range res {
			if index == 1 {
				defer func() {
					rainSnow = &i
				}()
				if (rainSnow == nil || rainSnow.Start().Date != i.Start().Date) ||
					(rainSnow.End() == nil && i.End() != nil) || (rainSnow.End() != nil && i.End() == nil) ||
					(rainSnow.End() != nil && i.End() != nil && rainSnow.End().Date != i.End().Date) {
					alert = true
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

func alertTempRiseFall() (subject string, body strings.Builder, alert bool) {
	if res, err := forecast.WillTempRiseFall(*difference, *query, *days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		for index, i := range res {
			if index == 1 {
				defer func() {
					tempRiseFall = &i
				}()
				if tempRiseFall == nil || tempRiseFall.Day().Date != i.Day().Date {
					alert = true
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
