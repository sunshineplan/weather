package main

import (
	"fmt"
	"log"
	"time"

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

func timestamp() string {
	now := time.Now()
	return fmt.Sprintf("(%02d:%02d)", now.Hour(), now.Minute())
}
