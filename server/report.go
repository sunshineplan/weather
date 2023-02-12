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

func daily() {
	days, err := forecast.Forecast(*query, 1)
	if err != nil {
		log.Print(err)
		return
	}

	var body strings.Builder
	fmt.Fprintln(&body, days[0].String())
	if rainSnow != nil {
		fmt.Fprintln(&body, "Recent Rain Snow Alert:")
		fmt.Fprintln(&body, rainSnow.String())
	}
	if tempRiseFall != nil {
		if tempRiseFall.IsRise() {
			fmt.Fprintln(&body, "Recent Temperature Rise Alert:")
		} else {
			fmt.Fprintln(&body, "Recent Temperature Fall Alert:")
		}
		fmt.Fprintln(&body, tempRiseFall.String())
	}
	go sendMail("[Weather]Daily Report"+timestamp(), body.String())
}

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
					subject = "[Weather]Rain Snow Alert - " + i.Start().Date + timestamp()
				}
			}
			fmt.Fprintln(&body, i.String())
		}
		rainSnow = &first
	} else if rainSnow != nil {
		subject = "[Weather]Rain Snow Alert - Canceled" + timestamp()
		body.WriteString("No more rain snow")
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
						subject = "[Weather]Temperature Rise Alert - " + i.Day().Date + timestamp()
					} else {
						subject = "[Weather]Temperature Fall Alert - " + i.Day().Date + timestamp()
					}
				}
			}
			fmt.Fprintln(&body, i.String())
		}
		tempRiseFall = &first
	} else if tempRiseFall != nil {
		if tempRiseFall.IsRise() {
			subject = "[Weather]Temperature Rise Alert - Canceled" + timestamp()
			body.WriteString("No more temperature rise")
		} else {
			subject = "[Weather]Temperature Fall Alert- Canceled" + timestamp()
			body.WriteString("No more temperature fall")
		}
		tempRiseFall = nil
	}
	return
}
