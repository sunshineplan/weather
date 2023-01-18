module weather

go 1.19

require (
	github.com/julienschmidt/httprouter v1.3.0
	github.com/sunshineplan/database/mongodb v1.0.5
	github.com/sunshineplan/metadata v1.1.1
	github.com/sunshineplan/service v1.0.6
	github.com/sunshineplan/utils v0.1.16
	github.com/sunshineplan/weather v0.0.0-00010101000000-000000000000
)

require (
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/pschlump/AesCCM v0.0.0-20160925022350-c5df73b5834e // indirect
	github.com/sunshineplan/cipher v1.0.4 // indirect
	golang.org/x/crypto v0.5.0 // indirect
	golang.org/x/sys v0.4.0 // indirect
)

replace github.com/sunshineplan/weather => ../
