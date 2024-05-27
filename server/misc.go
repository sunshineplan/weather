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
	"github.com/sunshineplan/utils/scheduler"
	"github.com/sunshineplan/weather/api/zoomearth"
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

func mapOptions(zoom float64) *zoomearth.MapOptions {
	return zoomearth.NewMapOptions().
		SetSize(600, 800).
		SetZoom(zoom).
		SetOverlays([]string{"radar", "wind"}).
		SetTimeZone(time.FixedZone("CST", 8*60*60))
}

func timestamp() string {
	return time.Now().Format("(2006/01/02 15:04)")
}

func jpg2gif(jpgPath, output string, count int) error {
	res, err := filepath.Glob(jpgPath)
	if err != nil {
		return err
	}
	n := len(res)
	if count != 0 && n > count {
		res = res[n-count:]
		n = count
	}
	var step int
	if count != 0 {
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
	if count != 0 {
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
