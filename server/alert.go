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
	if rainSnow != nil {
		if rainSnow.Start().IsExpired() && !rainSnow.End().IsExpired() {
			rainSnow.Start().Date = time.Now().Format("2006-01-02")
		}
	}

	if res, err := forecast.WillRainSnow(*query, *days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		var alert bool
		var title string
		var b strings.Builder
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
			fmt.Fprintln(&b, i.String())
			b.WriteRune('\n')
		}
		if alert {
			log.Print(title)
			go sendMail(title, b.String())
		}
	} else {
		rainSnow = nil
	}

	if res, err := forecast.WillTempRiseFall(*difference, *query, *days); err != nil {
		log.Print(err)
	} else if len(res) > 0 {
		var alert bool
		var title string
		var b strings.Builder
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
			fmt.Fprintln(&b, i.String())
			b.WriteRune('\n')
		}
		if alert {
			log.Print(title)
			go sendMail(title, b.String())
		}
	} else {
		tempRiseFall = nil
	}
}
