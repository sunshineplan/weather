package main

import (
	"bytes"
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
	"time"

	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/utils/scheduler"
)

var to mail.Receipts

func attachment6h() []*mail.Attachment {
	imgs, err := getImages("daily/*", 6*time.Hour, format, false)
	if err != nil {
		svc.Print(err)
		return nil
	}
	gifImg := new(gif.GIF)
	for i, img := range imgs {
		f, err := os.Open(img)
		if err != nil {
			svc.Print(err)
			return nil
		}
		if img, _, err := image.Decode(f); err != nil {
			svc.Print(err)
		} else {
			p := image.NewPaletted(img.Bounds(), palette.Plan9)
			draw.Draw(p, p.Rect, img, image.Point{}, draw.Over)
			gifImg.Image = append(gifImg.Image, p)
			if i != len(imgs)-1 {
				gifImg.Delay = append(gifImg.Delay, 40)
			} else {
				gifImg.Delay = append(gifImg.Delay, 300)
			}
		}
		f.Close()
	}
	var buf bytes.Buffer
	if err := gif.EncodeAll(&buf, gifImg); err != nil {
		svc.Print(err)
		return nil
	}
	return []*mail.Attachment{{Filename: "6h.gif", Bytes: buf.Bytes(), ContentID: "attachment"}}
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
