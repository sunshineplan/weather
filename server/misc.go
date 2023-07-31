package main

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	"github.com/sunshineplan/utils/mail"
)

var to mail.Receipts

func sendMail(subject, body string, attachments []*mail.Attachment) {
	msg := &mail.Message{
		Subject:     subject,
		Body:        body,
		Attachments: attachments,
		ContentType: mail.TextHTML,
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

func jpg2gif(jpgPath, output string) error {
	res, err := filepath.Glob(jpgPath)
	if err != nil {
		return err
	}
	var imgs []image.Image
	for _, i := range res {
		f, err := os.Open(i)
		if err != nil {
			return err
		}
		defer f.Close()
		img, err := jpeg.Decode(f)
		if err != nil {
			return err
		}
		imgs = append(imgs, img)
	}
	gifImg := new(gif.GIF)
	for _, img := range imgs {
		p := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(p, p.Rect, img, image.Point{}, draw.Over)
		gifImg.Image = append(gifImg.Image, p)
		gifImg.Delay = append(gifImg.Delay, 100)
	}
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, gifImg)
}
