package main

import (
	"bytes"
	"time"

	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/utils/pool"
	"github.com/sunshineplan/utils/scheduler"
)

var (
	mailSchedule = scheduler.ClockSchedule(scheduler.ClockFromString(*start), scheduler.ClockFromString(*end), time.Second)
	msgPool      = pool.New[mail.Message]()
)

func sendMail[T ~string](subject string, body T, contentType mail.ContentType, attachments []*mail.Attachment, force bool) {
	if !force && !mailSchedule.IsMatched(time.Now()) {
		return
	}
	msg := msgPool.Get()
	defer msgPool.Put(msg)
	msg.Subject = subject
	msg.Body = string(body)
	msg.Attachments = attachments
	msg.ContentType = contentType
	for _, to := range to {
		msg.To = mail.Receipts{to}
		if err := dialer.Send(msg); err != nil {
			svc.Print(err)
		}
	}
}

func attachment6h() []*mail.Attachment {
	imgs, err := getImages("daily/*", 6*time.Hour, format, false)
	if err != nil {
		svc.Print(err)
		return nil
	}
	buf := new(bytes.Buffer)
	if err := encodeGIF(buf, imgs, 40); err != nil {
		svc.Print(err)
		return nil
	}
	return []*mail.Attachment{{Filename: "6h.gif", Bytes: buf.Bytes(), ContentID: "attachment"}}
}

func timestamp() string {
	return time.Now().Format("(2006/01/02 15:04)")
}
