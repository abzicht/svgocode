package math64

import (
	"math"
	"strconv"

	"github.com/abzicht/svgocode/llog"
)

// mm
type Float float64

func ParseFloat(s string) Float {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		llog.Panicf("Failed to parse string '%s' as float: %s", s, err.Error())
	}
	return Float(f)
}

func (f Float) Pow(f2 Float) Float {
	return Float(math.Pow(float64(f), float64(f2)))
}

func (f Float) Min(f2 Float) Float {
	if f < f2 {
		return f
	}
	return f2
}

func (f Float) Max(f2 Float) Float {
	if f > f2 {
		return f
	}
	return f2
}

func Min(floats ...Float) Float {
	if len(floats) == 0 {
		llog.Panic("Called Min with no arguments")
	}
	var f Float = floats[0]
	for _, f2 := range floats[1:] {
		f = f.Min(f2)
	}
	return f
}

func Max(floats ...Float) Float {
	if len(floats) == 0 {
		llog.Panic("Called Max with no arguments")
	}
	var f Float = floats[0]
	for _, f2 := range floats[1:] {
		f = f.Max(f2)
	}
	return f
}

// mm/s
type Speed Float
