package main

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"math"
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
	n := len(res)
	step := int(math.Round(float64(n) / 24))
	if step == 0 {
		step = 1
	}
	var imgs []image.Image
	for i, name := range res {
		if i%step == 0 || i == n-1 {
			f, err := os.Open(name)
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
	}
	gifImg, n := new(gif.GIF), len(imgs)
	for i, img := range imgs {
		p := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(p, p.Rect, img, image.Point{}, draw.Over)
		gifImg.Image = append(gifImg.Image, p)
		if i != n-1 {
			gifImg.Delay = append(gifImg.Delay, 40)
		} else {
			gifImg.Delay = append(gifImg.Delay, 300)
		}
	}
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	return gif.EncodeAll(f, gifImg)
}
