package maps

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrTimeParse         = errors.New("failed to parse time")
	ErrInsufficientColor = errors.New("image has insufficient color depth")
)

type parseTimeError struct{ *time.ParseError }

func ParseTimeError(err error) (error, bool) {
	if err, ok := err.(*time.ParseError); ok && err != nil {
		return &parseTimeError{err}, true
	}
	return err, false
}
func (parseTimeError) Is(target error) bool { return target == ErrTimeParse }
func (parseTimeError) Unwrap() error        { return ErrTimeParse }

var _ error = InsufficientColor(0)

type InsufficientColor int

func (i InsufficientColor) Error() string {
	return fmt.Sprintf("image has insufficient color depth: %d", i)
}
func (InsufficientColor) Is(target error) bool { return target == ErrInsufficientColor }
func (InsufficientColor) Unwrap() error        { return ErrInsufficientColor }
