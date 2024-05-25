package option

import "time"

type Size interface {
	Size() (width int, height int)
}

type Zoom interface {
	Zoom() float64
}

type Overlays interface {
	Overlays() []string
}

type TimeZone interface {
	TimeZone() *time.Location
}
