//go:build !webp

package main

import (
	"image"
	"os"

	"github.com/sunshineplan/apng"
	"github.com/sunshineplan/utils/pool"
)

const ext = ".png"

var (
	apngPool    = pool.New[apng.APNG]()
	apngEncoder = apng.Encoder{
		CompressionLevel: apng.BestCompression,
		BufferPool:       pool.New[apng.EncoderBuffer](),
	}
)

func encodeAnimation(file string, imgs []string) error {
	apngImg := apngPool.Get()
	defer func() {
		apngImg.Frames = apngImg.Frames[:0]
		apngPool.Put(apngImg)
	}()
	for i, img := range imgs {
		f, err := os.Open(img)
		if err != nil {
			return err
		}
		if img, _, err := image.Decode(f); err != nil {
			svc.Print(err)
		} else {
			if i != len(imgs)-1 {
				apngImg.Frames = append(apngImg.Frames, apng.Frame{Image: img, DelayNumerator: 40})
			} else {
				apngImg.Frames = append(apngImg.Frames, apng.Frame{Image: img, DelayNumerator: 300})
			}
		}
		f.Close()
	}
	f, err := os.Create(file + ext)
	if err != nil {
		return err
	}
	defer f.Close()
	return apngEncoder.Encode(f, *apngImg)
}
