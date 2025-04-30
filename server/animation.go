package main

import (
	"image"
	"os"

	"github.com/HugoSmits86/nativewebp"
	"github.com/sunshineplan/utils/pool"
)

const ext = ".webp"

var animationPool = pool.New[nativewebp.Animation]()

func encodeAnimation(file string, imgs []string) error {
	webp := animationPool.Get()
	defer func() {
		webp.Images = webp.Images[:0]
		webp.Disposals = webp.Disposals[:0]
		webp.Durations = webp.Durations[:0]
		animationPool.Put(webp)
	}()
	for i, img := range imgs {
		f, err := os.Open(img)
		if err != nil {
			return err
		}
		if img, _, err := image.Decode(f); err != nil {
			svc.Print(err)
		} else {
			webp.Images = append(webp.Images, img)
			webp.Disposals = append(webp.Disposals, 0)
			if i != len(imgs)-1 {
				webp.Durations = append(webp.Durations, 400)
			} else {
				webp.Durations = append(webp.Durations, 3000)
			}
		}
		f.Close()
	}
	f, err := os.Create(file + ext)
	if err != nil {
		return err
	}
	defer f.Close()
	return nativewebp.EncodeAll(f, webp, nil)
}
