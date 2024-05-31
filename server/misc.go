package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/utils/scheduler"
)

var to mail.Receipts

func attachment(file string) (attachments []*mail.Attachment) {
	zoomMutex.Lock()
	defer zoomMutex.Unlock()
	if b, err := os.ReadFile(file); err != nil {
		svc.Print(err)
	} else {
		attachments = append(attachments, &mail.Attachment{Filename: filepath.Base(file), Bytes: b, ContentID: "attachment"})
	}
	return
}

func sendMail[T ~string](subject string, body T, contentType mail.ContentType, attachments []*mail.Attachment, force bool) {
	if !force && !scheduler.ClockSchedule(
		scheduler.ClockFromString(*start),
		scheduler.ClockFromString(*end),
		time.Second,
	).IsMatched(time.Now()) {
		return
	}
	msg := &mail.Message{
		Subject:     subject,
		Body:        string(body),
		Attachments: attachments,
		ContentType: contentType,
	}
	for _, to := range to {
		msg.To = mail.Receipts{to}
		if err := dialer.Send(msg); err != nil {
			svc.Print(err)
		}
	}
}

func timestamp() string {
	return time.Now().Format("(2006/01/02 15:04)")
}
