package main

import (
	"log"

	"github.com/sunshineplan/utils/mail"
)

var to []string

func sendMail(subject, body string) {
	for _, to := range to {
		if err := dialer.Send(
			&mail.Message{
				To:      []string{to},
				Subject: subject,
				Body:    body,
			},
		); err != nil {
			log.Print(err)
		}
	}
}
