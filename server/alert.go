package main

import (
	"log"

	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/weather"
)

var (
	startRainSnow *weather.Hour
	stopRainSnow  *weather.Hour

	tempUp   *weather.Day
	tempDown *weather.Day
)

func alert() {
	if hour, start, err := history.WillRainSnow(*query, 3); err != nil {
		log.Print(err)
	} else if hour != nil {
		if start {
			defer func() {
				startRainSnow = hour
			}()
			if startRainSnow == nil || startRainSnow.Time != hour.Time {
				log.Printf("降雨警报 %#v", hour) //TODO
				//go sendMail("降雨警报")
			}
		} else {
			defer func() {
				stopRainSnow = hour
			}()
			if stopRainSnow == nil || stopRainSnow.Time != hour.Time {
				log.Printf("雨停预报 %#v", hour) //TODO
				//go sendMail("雨停预报")
			}
		}
	}

	if day, up, err := history.WillUpDown(*difference, *query, 3); err != nil {
		log.Print(err)
	} else if day != nil {
		if up {
			defer func() {
				tempUp = day
			}()
			if tempUp == nil || tempUp.Date != day.Date {
				log.Printf("升温警报 %#v", day) //TODO
				//go sendMail("升温警报")
			}
		} else {
			defer func() {
				tempDown = day
			}()
			if tempDown == nil || tempDown.Date != day.Date {
				log.Printf("降温警报 %#v", day) //TODO
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
