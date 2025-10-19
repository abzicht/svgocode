package math64

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/abzicht/svgocode/llog"
)

type UnitType string

const (
	UnitMM = UnitType("mm") // Millimeter
	UnitCM = UnitType("cm") // Centimeter
	UnitIN = UnitType("in") // Inches
	UnitPT = UnitType("pt") // Points (Requires DPI)
	UnitPX = UnitType("px") // Pixels (let's figure that one out later)
)

var UnitTypes []UnitType = []UnitType{UnitMM, UnitCM, UnitIN, UnitPT, UnitPX}

func UnitTypeFromString(s string) UnitType {
	unitT := UnitType(strings.ToLower(s))
	for _, t := range UnitTypes {
		if t == unitT {
			return t
		}
	}
	llog.Panicf("Unknown unit type: %s", unitT)
	return UnitType("")
}

var unitMatcher *regexp.Regexp = regexp.MustCompile(`^([0-9]*\.?[0-9]+)([a-zA-Z%µ]+)$`)

// Given an input string such as "32mm", determine its value and unit
func NumberUnit(s string) (Float, UnitType) {
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

	unit := UnitTypeFromString(matches[2])
	return Float(value), unit
}
