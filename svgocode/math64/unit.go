package math64

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/abzicht/svgocode/llog"
)

type UnitLength string

const (
	UnitMM = UnitLength("mm") // Millimeter
	UnitCM = UnitLength("cm") // Centimeter
	UnitIN = UnitLength("in") // Inches
	UnitPT = UnitLength("pt") // Points (Requires DPI)
	UnitPX = UnitLength("px") // Pixels (let's figure that one out later)
)

var UnitLengths []UnitLength = []UnitLength{UnitMM, UnitCM, UnitIN, UnitPT, UnitPX}

func UnitLengthFromString(s string) UnitLength {
	unitT := UnitLength(strings.ToLower(s))
	for _, t := range UnitLengths {
		if t == unitT {
			return t
		}
	}
	llog.Panicf("Unknown unit type: %s", unitT)
	return UnitLength("")
}

var unitMatcher *regexp.Regexp = regexp.MustCompile(`^([0-9]*\.?[0-9]+)([a-zA-Z%µ]+)$`)

// Given an input string such as "32mm", determine its value and unit
func NumberUnit(s string) (Float, UnitLength) {
	// Regex explanation:
	// ^([0-9]*\.?[0-9]+)   -> captures an integer or decimal number
	// ([a-zA-Z%µ]+)$       -> captures the unit (letters, %, µ, etc.)

	matches := unitMatcher.FindStringSubmatch(s)
	if len(matches) != 3 {
		llog.Panicf("Unknown number format: %s", s)
	}

	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		llog.Panicf("Failed to parse number as float: %s", err.Error())
	}

	unit := UnitLengthFromString(matches[2])
	return Float(value), unit
}

// Convert a length l from unit 'from' to unit 'to'
func LengthConvert(l Float, from, to UnitLength) Float {
	var tmp Float
	switch from {
	case UnitMM:
		tmp = l
	case UnitCM:
		tmp = l * 10
	case UnitIN:
		tmp = l * 25.4
	default:
		llog.Panicf("Conversion of type %s is not supported", from)
	}

	switch to {
	case UnitMM:
		return tmp
	case UnitCM:
		return tmp / 10.0
	case UnitIN:
		return tmp / 25.4
	default:
		llog.Panicf("Conversion of type %s is not supported", to)
	}
	llog.Panicf("NOT REACHED")
	return -1
}
