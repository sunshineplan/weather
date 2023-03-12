package main

import (
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
			svc.Print(err)
		}
	}
}

func timestamp() string {
	return time.Now().Format("(2006/01/02 15:04)")
}
