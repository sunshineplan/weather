package main

import (
	"fmt"
	"log"

	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/weather"
)

var (
	rainSnow     *weather.RainSnow
	tempRiseFall *weather.TempRiseFall
)

func alert() {
	if rainSnow != nil {
		//TODO: check rainSnow and tempRiseFall expired
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
				if (rainSnow == nil || rainSnow.Start().Date != i.Start().Date) || rainSnow.End() != i.End() ||
					(rainSnow.End() != nil && i.End() != nil && rainSnow.End().Date != i.End().Date) {
					alert = true
				}
			}
			output += fmt.Sprintln(i)
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
			output += fmt.Sprintln(i)
		}
		if alert {
			log.Print(output) //TODO
			//go sendMail(output)
		}
	} else {
		tempRiseFall = nil
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
