package option

type Zoom interface {
	Zoom() float64
}

type Quality interface {
	Quality() int
}

type Overlays interface {
	Overlays() []string
}
