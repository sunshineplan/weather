package zoomearth

import (
	"errors"
	"image/jpeg"
	"log"
	"os"
	"testing"
	"time"

	"github.com/sunshineplan/chrome"
	"github.com/sunshineplan/weather/maps"
	"github.com/sunshineplan/weather/unit/coordinates"
)

func TestZoomEarth(t *testing.T) {
	if _, err := GetStorms(time.Now()); err != nil {
		t.Error(err)
	}
	c := chrome.Headless().NoSandbox()
	defer c.Close()
	for _, path := range []string{"satellite", "precipitation"} {
		log.Println("Test", path)
		_, img, err := MapWithContext(c, path, time.Time{}, coordinates.New(0, 0), nil)
		if err != nil && !errors.Is(err, maps.ErrInsufficientColor) {
			t.Fatal(err)
		}
		f, err := os.Create(path + ".jpg")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		if err := jpeg.Encode(f, img, nil); err != nil {
			t.Fatal(err)
		}
	}
}
