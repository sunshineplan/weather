//go:build webp

package main

import (
	"image"
	"os"
	"time"

	"github.com/tidbyt/go-libwebp/webp"
)

const ext = ".webp"

func encodeAnimation(file string, imgs []string) error {
	webp, err := webp.NewAnimationEncoder(width, height, 0, 0)
	if err != nil {
		return err
	}
	defer webp.Close()
	for i, img := range imgs {
		f, err := os.Open(img)
		if err != nil {
			return err
		}
		if img, _, err := image.Decode(f); err != nil {
			svc.Print(err)
		} else {
			if i != len(imgs)-1 {
				webp.AddFrame(img, 400*time.Millisecond)
			} else {
				webp.AddFrame(img, 3*time.Second)
			}
		}
		f.Close()
	}
	b, err := webp.Assemble()
	if err != nil {
		return err
	}
	return os.WriteFile(file+ext, b, 0644)
}
