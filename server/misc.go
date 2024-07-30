package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"

	"github.com/sunshineplan/utils/mail"
	"github.com/sunshineplan/utils/pool"
	"github.com/sunshineplan/utils/scheduler"
	"github.com/sunshineplan/weather/storm"
)

var (
	mailSchedule = scheduler.ClockSchedule(scheduler.ClockFromString(*start), scheduler.ClockFromString(*end), time.Second)
	msgPool      = pool.New[mail.Message]()
	bufPool      = pool.New[bytes.Buffer]()
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

func lastImage(path string) (img image.Image, err error) {
	files, err := filepath.Glob(path)
	if err != nil {
		return
	}
	if len(files) == 0 {
		return nil, errors.New("no files in this path: " + path)
	}
	f, err := os.Open(files[len(files)-1])
	if err != nil {
		return
	}
	defer f.Close()
	img, _, err = image.Decode(f)
	return
}

func attachLast() []*mail.Attachment {
	img, err := lastImage("daily/*")
	if err != nil {
		svc.Print(err)
		return nil
	}
	buf := bufPool.Get()
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: 90}); err != nil {
		svc.Print(err)
		return nil
	}
	return []*mail.Attachment{{Filename: "last.jpg", Bytes: buf.Bytes(), ContentID: "attachment"}}
}

func attach6hGIF() []*mail.Attachment {
	imgs, err := getImages("daily/*", 6*time.Hour, true)
	if err != nil {
		svc.Print(err)
		return nil
	}
	buf := bufPool.Get()
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	if err := encodeGIF(buf, imgs); err != nil {
		svc.Print(err)
		return nil
	}
	return []*mail.Attachment{{Filename: "6h.gif", Bytes: buf.Bytes(), ContentID: "attachment"}}
}

func attachStorm(i int, storm storm.Data) *mail.Attachment {
	imgs, err := getImages(fmt.Sprintf("%s/%s/%d-%s/*", *path, storm.Season, storm.No, storm.ID), 0, false)
	if err != nil {
		svc.Print(err)
		return nil
	}
	buf := bufPool.Get()
	defer func() {
		buf.Reset()
		bufPool.Put(buf)
	}()
	if err := encodeGIF(buf, imgs); err != nil {
		svc.Print(err)
		return nil
	}
	return &mail.Attachment{
		Filename:  fmt.Sprintf("image%d.gif", i),
		Bytes:     buf.Bytes(),
		ContentID: fmt.Sprintf("map%d", i),
	}
}

func timestamp() string {
	return time.Now().Format("(2006/01/02 15:04)")
}
