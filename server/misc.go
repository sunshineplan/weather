package main

import (
	"time"

	"github.com/sunshineplan/utils/mail"
)

var to mail.Receipts

func sendMail(subject, body string) {
	for _, to := range to {
		if err := dialer.Send(
			&mail.Message{
				To:          mail.Receipts{to},
				Subject:     subject,
				Body:        body,
				ContentType: mail.TextHTML,
			},
		); err != nil {
			svc.Print(err)
		}
	}
}

func timestamp() string {
	return time.Now().Format("(2006/01/02 15:04)")
}
