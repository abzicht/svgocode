package gcode

import (
	"errors"
	"fmt"
	"io"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	e := new(Encoder)
	e.w = w
	return e
}

func (e *Encoder) Encode(g *Gcode) (err error) {
	defer func() {
		// Maybe g.String panics, but we want to return an error.
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprintf("failed to encode gcode: panicked with '%s'", r))
		}
		return
	}()
	gcodeBytes := []byte(g.String())
	n, err := e.w.Write(gcodeBytes)
	if err != nil {
		return err
	}
	if n != len(gcodeBytes) {
		return errors.New(fmt.Sprintf("failed to encode gcode: only %d of %d bytes written", n, len(gcodeBytes)))
	}
	return nil
}

func (e *Encoder) Close() error {
	return nil
}
