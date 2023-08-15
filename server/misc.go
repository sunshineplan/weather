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

func jpg2gif(jpgPath, output string, daily bool) error {
	res, err := filepath.Glob(jpgPath)
	if err != nil {
		return err
	}
	n := len(res)
	var step int
	if daily {
		step = 1
	} else if step = int(math.Round(math.Log(1+float64(n)))) - 2; step <= 0 {
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
	var delay int
	if daily {
		delay = 40
	} else if delay = 6000 / n; delay > 40 {
		delay = 40
	}
	for i, img := range imgs {
		p := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(p, p.Rect, img, image.Point{}, draw.Over)
		gifImg.Image = append(gifImg.Image, p)
		if i != n-1 {
			gifImg.Delay = append(gifImg.Delay, delay)
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
