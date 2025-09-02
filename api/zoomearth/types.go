package zoomearth

import "time"

type ZoomEarthAPI struct{}

type StormID string

type KeyEvent struct {
	Key                   string
	Code                  string
	WindowsVirtualKeyCode int64
	Listen                any
	MustListen            bool
}

type MapOptions struct {
	width      int
	height     int
	zoom       float64
	overlays   []string
	timezone   *time.Location
	listenList []any
	keyEvents  []KeyEvent
}
