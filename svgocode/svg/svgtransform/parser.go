package svgtransform

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/math64"
)

// ParseTransform parses an SVG transform attribute into a slice of Transform structs.
// It validates parameter counts according to the SVG specification.
func ParseTransform(input string) TransformChain {
	if len(input) == 0 {
		return TransformChain{} // More often than not, nothing is provided
	}

	// Regex to match function names and parameter contents.
	re := regexp.MustCompile(`([a-zA-Z]+)\s*\(([^)]*)\)`)

	matches := re.FindAllStringSubmatch(input, -1)
	if matches == nil {
		return TransformChain{} // No valid functions found
	}

	transforms := TransformChain{}

	for _, m := range matches {
		fnName := TransformCommandType(strings.ToLower(strings.TrimSpace(m[1])))
		rawParams := m[2]

		// Split params by comma or whitespace
		paramTokens := splitParams(rawParams)
		var params []math64.Float
		for _, t := range paramTokens {
			if t == "" {
				continue
			}
			f, err := strconv.ParseFloat(t, 64)
			if err != nil {
				llog.Panicf("Failed to parse transform function: invalid number in %s(): %v", fnName, err)
			}
			params = append(params, math64.Float(f))
		}

		// Validate according to SVG spec
		if err := validateTransform(fnName, len(params)); err != nil {
			llog.Panicf("Failed to parse transform function: %s(): %v", fnName, err)
		}
		var transform Transform
		switch fnName {
		case TransformCmdTranslate:
			transform = NewTranslate(math64.VectorF2{X: params[0], Y: params[1]})
		case TransformCmdTranslateX:
			transform = NewTranslate(math64.VectorF2{X: params[0]})
		case TransformCmdTranslateY:
			transform = NewTranslate(math64.VectorF2{Y: params[0]})
		case TransformCmdMatrix:
			m := math64.MatrixF4Identity()
			if len(params) != 6 {
				llog.Panicf("Expected transform matrix with 6 parameters but got %d.", len(params))
			}
			m[0] = params[0]
			m[4] = params[1]
			m[1] = params[2]
			m[5] = params[3]
			m[3] = params[4]
			m[7] = params[5]
			transform = NewTransformMatrix(m)
		case TransformCmdScale:
			transform = NewScale(math64.VectorF2{X: params[0], Y: params[0]})
		case TransformCmdScaleX:
			transform = NewScale(math64.VectorF2{X: params[0], Y: 1})
		case TransformCmdScaleY:
			transform = NewScale(math64.VectorF2{X: 1, Y: params[0]})
		case TransformCmdRotate:
			if len(params) == 1 {
				transform = NewRotate(math64.AngDeg(params[0]), math64.VectorF2{X: 0, Y: 0})
			} else {
				transform = NewRotate(math64.AngDeg(params[0]), math64.VectorF2{X: params[1], Y: params[2]})
			}
		case TransformCmdSkew:
			transform = NewSkew(math64.VectorT2[math64.AngDeg]{X: math64.AngDeg(params[0]), Y: math64.AngDeg(params[1])})
		case TransformCmdSkewX:
			transform = NewSkew(math64.VectorT2[math64.AngDeg]{X: math64.AngDeg(params[0]), Y: 0})
		case TransformCmdSkewY:
			transform = NewSkew(math64.VectorT2[math64.AngDeg]{X: 0, Y: math64.AngDeg(params[0])})
		default:
			llog.Errorf("Cannot add transformation %s: Not yet implemented\n", fnName)
		}

		transforms = append(transforms, transform)
	}

	return transforms
}

// splitParams splits parameter strings that may be separated by commas, spaces, or both.
func splitParams(s string) []string {
	// Replace commas with spaces, then split on whitespace
	s = strings.ReplaceAll(s, ",", " ")
	return strings.Fields(s)
}

// validateTransform checks if a transform has the correct number of parameters
func validateTransform(fn TransformCommandType, count int) error {
	switch fn {
	case TransformCmdMatrix:
		if count != 6 {
			return fmt.Errorf("expected 6 parameters, got %d", count)
		}
	case TransformCmdTranslate, TransformCmdScale, TransformCmdSkew:
		if count != 2 {
			return fmt.Errorf("expected 2 parameters, got %d", count)
		}
	case TransformCmdRotate:
		if count != 1 && count != 3 {
			return fmt.Errorf("expected 1 or 3 parameters, got %d", count)
		}
	case TransformCmdTranslateX, TransformCmdTranslateY, TransformCmdScaleX, TransformCmdScaleY, TransformCmdSkewX, TransformCmdSkewY:
		if count != 1 {
			return fmt.Errorf("expected 1 parameter, got %d", count)
		}
	default:
		return fmt.Errorf("unknown transform type: %s", fn)
	}
	return nil
}
