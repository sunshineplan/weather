package main

import (
	"log"

	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/weather"
)

var (
	startRainSnow *weather.ForecastHour
	stopRainSnow  *weather.ForecastHour

	tempUp   *weather.ForecastForecastday
	tempDown *weather.ForecastForecastday
)

func alert() {
	resp, err := weather.ForecastWeather(*query, 3)
	if err != nil {
		if *debug {
			log.Print(err)
		}
		return
	}

	if hour, start := resp.WillRainSnow(); hour != nil {
		if start {
			defer func() {
				startRainSnow = hour
			}()
			if startRainSnow == nil || startRainSnow.Time != hour.Time {
				log.Print("降雨警报", hour) //TODO
				//go sendMail("降雨警报")
			}
		} else {
			defer func() {
				stopRainSnow = hour
			}()
			if stopRainSnow == nil || stopRainSnow.Time != hour.Time {
				log.Print("雨停预报", hour) //TODO
				//go sendMail("雨停预报")
			}
		}
	}

	if day, up := resp.WillUpDown(*difference); day != nil {
		if up {
			defer func() {
				tempUp = day
			}()
			if tempUp == nil || tempUp.Date != day.Date {
				log.Print("升温警报", day) //TODO
				//go sendMail("升温警报")
			}
		} else {
			defer func() {
				tempDown = day
			}()
			if tempDown == nil || tempDown.Date != day.Date {
				log.Print("降温警报", day) //TODO
				//go sendMail("降温警报")
			}
		}
	}
}

var to []string

func sendMail(subject, body string) {
	for _, to := range to {
		if err := dialer.Send(
			&mail.Message{
				To:          []string{to},
				Subject:     subject,
				Body:        body,
				ContentType: mail.TextHTML,
			},
		); err != nil {
			log.Print(err)
		}
	}
}
