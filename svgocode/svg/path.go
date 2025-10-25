package svg

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/abzicht/svgocode/svgocode/math64"
)

// SVG Path command parser.

type PathCommandType int

const (
	CmdMoveTo PathCommandType = iota
	CmdClosePath
	CmdLineTo
	CmdHLineTo
	CmdVLineTo
	CmdCurveTo
	CmdSmoothCurveTo
	CmdQuadraticBezierTo
	CmdSmoothQuadraticBezierTo
	CmdEllipticalArc
)

func (c PathCommandType) String() string {
	return [...]string{
		"MoveTo", "ClosePath", "LineTo", "HLineTo", "VLineTo",
		"CurveTo", "SmoothCurveTo", "QuadraticBezierTo",
		"SmoothQuadraticBezierTo", "EllipticalArc",
	}[c]
}

type EllipticalArcArg struct {
	R     math64.VectorF2
	XAxis math64.Float // x-axis rotation
	Large bool         // flag 0/1
	Sweep bool         // flag 0/1
	To    math64.VectorF2
}

type PathCommand struct {
	Type        PathCommandType
	Relative    bool
	PathPoints  []math64.VectorF2 // used for moveto/lineto/curveto combos
	Coordinates []math64.Float    // used for H/V coordinates (list of floats)
	ArcArgs     []EllipticalArcArg
}

func (c PathCommand) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s (relative=%v): ", c.Type, c.Relative)
	switch c.Type {
	case CmdMoveTo, CmdLineTo, CmdCurveTo, CmdSmoothCurveTo, CmdQuadraticBezierTo, CmdSmoothQuadraticBezierTo:
		fmt.Fprintf(&b, "PathPoints=%v", c.PathPoints)
	case CmdHLineTo, CmdVLineTo:
		fmt.Fprintf(&b, "Coords=%v", c.Coordinates)
	case CmdEllipticalArc:
		fmt.Fprintf(&b, "Arcs=%v", c.ArcArgs)
	case CmdClosePath:
		// nothing more
	default:
		fmt.Fprintf(&b, "PathPoints=%v", c.PathPoints)
	}
	return b.String()
}

// -------- Parser --------

type parser struct {
	s string
	i int
	n int
}

func newParser(s string) *parser {
	return &parser{s: s, i: 0, n: len(s)}
}

func (p *parser) peek() rune {
	if p.i >= p.n {
		return 0
	}
	return rune(p.s[p.i])
}

func (p *parser) next() rune {
	if p.i >= p.n {
		return 0
	}
	ch := rune(p.s[p.i])
	p.i++
	return ch
}

func (p *parser) eof() bool {
	return p.i >= p.n
}

func (p *parser) skipWsp() {
	for !p.eof() {
		r := p.peek()
		// wsp ::= #x9 | #x20 | #xA | #xC | #xD
		if r == '\t' || r == ' ' || r == '\n' || r == '\f' || r == '\r' {
			p.i++
			continue
		}
		break
	}
}

// comma_wsp::=(wsp+ ","? wsp*) | ("," wsp*)
// We'll implement a permissive helper that consumes any mix of whitespace and at most one comma, but
// it should allow zero or more whitespace and optional comma according to contexts invoked.
func (p *parser) consumeCommaWspOptional() {
	// consume leading whitespace
	start := p.i
	hasWsp := false
	for !p.eof() {
		r := p.peek()
		if r == '\t' || r == ' ' || r == '\n' || r == '\f' || r == '\r' {
			hasWsp = true
			p.i++
			continue
		}
		break
	}
	if p.eof() {
		return
	}
	if p.peek() == ',' {
		// comma then optional whitespace
		p.i++
		for !p.eof() {
			r := p.peek()
			if r == '\t' || r == ' ' || r == '\n' || r == '\f' || r == '\r' {
				p.i++
				continue
			}
			break
		}
		return
	}
	// if had whitespace and optional comma after it
	if hasWsp {
		if !p.eof() && p.peek() == ',' {
			p.i++
			for !p.eof() {
				r := p.peek()
				if r == '\t' || r == ' ' || r == '\n' || r == '\f' || r == '\r' {
					p.i++
					continue
				}
				break
			}
		}
		return
	}
	// else nothing consumed
	_ = start
}

// Some contexts require "*", i.e., not throwing if absent. We'll use the optional variant above.
// For strict consumption where at least one wsp required, we provide:
func (p *parser) consumeMandatoryWsp() bool {
	start := p.i
	count := 0
	for !p.eof() {
		r := p.peek()
		if r == '\t' || r == ' ' || r == '\n' || r == '\f' || r == '\r' {
			p.i++
			count++
			continue
		}
		break
	}
	return count > 0 || p.i != start
}

// number ::= ([0-9])+
func (p *parser) parseNumber() (int, error) {
	p.skipWsp()
	start := p.i
	if p.eof() {
		return 0, fmt.Errorf("expected number at pos %d", p.i)
	}
	if !isDigit(p.peek()) {
		return 0, fmt.Errorf("expected digit at pos %d", p.i)
	}
	val := 0
	for !p.eof() && isDigit(p.peek()) {
		val = val*10 + int(p.next()-'0')
	}
	if p.i == start {
		return 0, fmt.Errorf("expected number at pos %d", start)
	}
	return val, nil
}

// parseFloat reads a floating-point number from the stream.
func (p *parser) parseFloat() (math64.Float, error) {
	var b strings.Builder

	p.skipWsp()
	if p.eof() {
		return 0, fmt.Errorf("expected float at pos %d", p.i)
	}
	// optional sign
	if ch := p.peek(); ch == '+' || ch == '-' {
		b.WriteRune(p.next())
	}

	// digits before decimal
	for isDigit(p.peek()) {
		b.WriteRune(p.next())
	}

	// optional decimal part
	if p.peek() == '.' {
		b.WriteRune(p.next())
		for isDigit(p.peek()) {
			b.WriteRune(p.next())
		}
	}

	// optional exponent
	if ch := p.peek(); ch == 'e' || ch == 'E' {
		b.WriteRune(p.next())
		if ch2 := p.peek(); ch2 == '+' || ch2 == '-' {
			b.WriteRune(p.next())
		}
		for isDigit(p.peek()) {
			b.WriteRune(p.next())
		}
	}

	numStr := b.String()
	if numStr == "" || numStr == "+" || numStr == "-" {
		return 0, fmt.Errorf("no valid float")
	}

	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid float %q: %v", numStr, err)
	}
	return math64.Float(val), nil
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9' || r == '.'
}

func (p *parser) parseFlag() (int, error) {
	p.skipWsp()
	if p.eof() {
		return 0, fmt.Errorf("expected flag at %d", p.i)
	}
	ch := p.peek()
	if ch != '0' && ch != '1' {
		return 0, fmt.Errorf("expected flag 0 or 1 at %d, got %q", p.i, ch)
	}
	p.i++
	return int(ch - '0'), nil
}

// coordinate ::= sign? number
// sign ::= "+" | "-"
func (p *parser) parseCoordinate() (math64.Float, error) {
	p.skipWsp()
	var sign math64.Float = 1.0
	if !p.eof() {
		if p.peek() == '+' {
			p.i++
		} else if p.peek() == '-' {
			sign = -1
			p.i++
		}
	}
	num, err := p.parseFloat()
	if err != nil {
		return 0, err
	}
	return sign * num, nil
}

// coordinate_pair ::= coordinate comma_wsp? coordinate
func (p *parser) parseCoordinatePair() (math64.VectorF2, error) {
	x, err := p.parseCoordinate()
	if err != nil {
		return math64.VectorF2{}, err
	}
	p.consumeCommaWspOptional()
	y, err := p.parseCoordinate()
	if err != nil {
		return math64.VectorF2{}, err
	}
	return math64.VectorF2{X: x, Y: y}, nil
}

// coordinate_pair_sequence ::= coordinate_pair | (coordinate_pair comma_wsp? coordinate_pair_sequence)
func (p *parser) parseCoordinatePairSequence() ([]math64.VectorF2, error) {
	var pts []math64.VectorF2
	pt, err := p.parseCoordinatePair()
	if err != nil {
		return nil, err
	}
	pts = append(pts, pt)
	// loop for additional pairs
	for {
		// try to see if next is comma_wsp or whitespace or number (coordinate)
		// save index
		before := p.i
		p.consumeCommaWspOptional()
		// peek: could be sign or digit for more coordinate pair
		if p.eof() {
			break
		}
		ch := p.peek()
		if ch == '+' || ch == '-' || isDigit(ch) {
			pt2, err := p.parseCoordinatePair()
			if err != nil {
				// rollback if parse fails
				p.i = before
				break
			}
			pts = append(pts, pt2)
			continue
		}
		// else rollback and break
		p.i = before
		break
	}
	return pts, nil
}

// coordinate_sequence ::= coordinate | (coordinate comma_wsp? coordinate_sequence)
func (p *parser) parseCoordinateSequence() ([]math64.Float, error) {
	var coords []math64.Float
	c, err := p.parseCoordinate()
	if err != nil {
		return nil, err
	}
	coords = append(coords, c)
	for {
		before := p.i
		p.consumeCommaWspOptional()
		if p.eof() {
			break
		}
		ch := p.peek()
		if ch == '+' || ch == '-' || isDigit(ch) {
			c2, err := p.parseCoordinate()
			if err != nil {
				p.i = before
				break
			}
			coords = append(coords, c2)
			continue
		}
		p.i = before
		break
	}
	return coords, nil
}

// curveto_coordinate_sequence: sequence of coordinate_pair_triplet
// coordinate_pair_triplet ::= coordinate_pair comma_wsp? coordinate_pair comma_wsp? coordinate_pair
func (p *parser) parseCoordinatePairTriplet() ([]math64.VectorF2, error) {
	var pts []math64.VectorF2
	for j := 0; j < 3; j++ {
		pt, err := p.parseCoordinatePair()
		if err != nil {
			return nil, err
		}
		pts = append(pts, pt)
		if j < 2 {
			p.consumeCommaWspOptional()
		}
	}
	return pts, nil
}

func (p *parser) parseCurvetoCoordinateSequence() ([][]math64.VectorF2, error) {
	// one or more triplets
	first, err := p.parseCoordinatePairTriplet()
	if err != nil {
		return nil, err
	}
	result := [][]math64.VectorF2{first}
	for {
		before := p.i
		p.consumeCommaWspOptional()
		// check if next looks like a coordinate_pair (sign/digit)
		if p.eof() {
			break
		}
		ch := p.peek()
		if ch == '+' || ch == '-' || isDigit(ch) {
			trip, err := p.parseCoordinatePairTriplet()
			if err != nil {
				p.i = before
				break
			}
			result = append(result, trip)
			continue
		}
		p.i = before
		break
	}
	return result, nil
}

// coordinate_pair_double ::= coordinate_pair comma_wsp? coordinate_pair
func (p *parser) parseCoordinatePairDouble() ([]math64.VectorF2, error) {
	var pts []math64.VectorF2
	for j := 0; j < 2; j++ {
		pt, err := p.parseCoordinatePair()
		if err != nil {
			return nil, err
		}
		pts = append(pts, pt)
		if j < 1 {
			p.consumeCommaWspOptional()
		}
	}
	return pts, nil
}

func (p *parser) parseSmoothCurvetoCoordinateSequence() ([][]math64.VectorF2, error) {
	first, err := p.parseCoordinatePairDouble()
	if err != nil {
		return nil, err
	}
	res := [][]math64.VectorF2{first}
	for {
		before := p.i
		p.consumeCommaWspOptional()
		if p.eof() {
			break
		}
		ch := p.peek()
		if ch == '+' || ch == '-' || isDigit(ch) {
			dbl, err := p.parseCoordinatePairDouble()
			if err != nil {
				p.i = before
				break
			}
			res = append(res, dbl)
			continue
		}
		p.i = before
		break
	}
	return res, nil
}

func (p *parser) parseQuadraticBezierCoordinateSequence() ([][]math64.VectorF2, error) {
	// similar to smooth curveto: two pairs per element
	return p.parseSmoothCurvetoCoordinateSequence()
}

// elliptical_arc_argument ::= number comma_wsp? number comma_wsp? number comma_wsp flag comma_wsp? flag comma_wsp? coordinate_pair
func (p *parser) parseEllipticalArcArgument() (EllipticalArcArg, error) {
	rx, err := p.parseFloat()
	if err != nil {
		return EllipticalArcArg{}, fmt.Errorf("elliptical arc rx: %w", err)
	}
	p.consumeCommaWspOptional()
	ry, err := p.parseFloat()
	if err != nil {
		return EllipticalArcArg{}, fmt.Errorf("elliptical arc ry: %w", err)
	}
	p.consumeCommaWspOptional()
	xAxis, err := p.parseFloat()
	if err != nil {
		return EllipticalArcArg{}, fmt.Errorf("elliptical arc x-axis-rotation: %w", err)
	}
	p.consumeCommaWspOptional()
	// flag
	largeInt, err := p.parseFlag()
	if err != nil {
		return EllipticalArcArg{}, fmt.Errorf("elliptical arc large-arc-flag: %w", err)
	}
	large := false
	if largeInt == 1 {
		large = true
	}
	p.consumeCommaWspOptional()
	sweepInt, err := p.parseFlag()
	if err != nil {
		return EllipticalArcArg{}, fmt.Errorf("elliptical arc sweep-flag: %w", err)
	}
	sweep := false
	if sweepInt == 1 {
		sweep = true
	}
	p.consumeCommaWspOptional()
	pt, err := p.parseCoordinatePair()
	if err != nil {
		return EllipticalArcArg{}, fmt.Errorf("elliptical arc end point: %w", err)
	}
	return EllipticalArcArg{
		R:     math64.VectorF2{X: rx, Y: ry},
		XAxis: xAxis,
		Large: large,
		Sweep: sweep,
		To:    pt,
	}, nil
}

func (p *parser) parseEllipticalArcArgumentSequence() ([]EllipticalArcArg, error) {
	first, err := p.parseEllipticalArcArgument()
	if err != nil {
		return nil, err
	}
	args := []EllipticalArcArg{first}
	for {
		before := p.i
		p.consumeCommaWspOptional()
		if p.eof() {
			break
		}
		ch := p.peek()
		// next should be digit for number (rx)
		if isDigit(ch) {
			a, err := p.parseEllipticalArcArgument()
			if err != nil {
				p.i = before
				break
			}
			args = append(args, a)
			continue
		}
		p.i = before
		break
	}
	return args, nil
}

func isPathCommandLetter(r rune) bool {
	switch r {
	case 'M', 'm', 'Z', 'z', 'L', 'l', 'H', 'h', 'V', 'v', 'C', 'c', 'S', 's', 'Q', 'q', 'T', 't', 'A', 'a':
		return true
	default:
		return false
	}
}

func (p *parser) parseDrawtoPathCommand() ([]PathCommand, error) {
	// drawto_command can be moveto | closepath | lineto | hline | vline | curveto | smoothcurveto | quadratic | smoothquadratic | elliptical_arc
	// Just detect letter and dispatch.
	if p.eof() {
		return nil, errors.New("unexpected EOF parsing drawto_command")
	}
	ch := p.peek()
	switch ch {
	case 'M', 'm':
		return p.parseMoveto()
	case 'Z', 'z':
		p.next()
		return []PathCommand{{Type: CmdClosePath, Relative: unicode.IsLower(ch)}}, nil
	case 'L', 'l':
		return p.parseLineto()
	case 'H', 'h':
		return p.parseHorizontalLineto()
	case 'V', 'v':
		return p.parseVerticalLineto()
	case 'C', 'c':
		return p.parseCurveto()
	case 'S', 's':
		return p.parseSmoothCurveto()
	case 'Q', 'q':
		return p.parseQuadraticBezierCurveto()
	case 'T', 't':
		return p.parseSmoothQuadraticBezierCurveto()
	case 'A', 'a':
		return p.parseEllipticalArc()
	default:
		return nil, fmt.Errorf("unknown drawto command %q at %d", ch, p.i)
	}
}

func (p *parser) parseMoveto() ([]PathCommand, error) {
	// ("M" | "m") wsp* coordinate_pair_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	points, err := p.parseCoordinatePairSequence()
	if err != nil {
		return nil, fmt.Errorf("moveto: %w", err)
	}
	// According to SVG: moveto with multiple coordinate pairs: first is MoveTo, the rest are implicit LineTo
	cmds := []PathCommand{}
	if len(points) >= 1 {
		cmds = append(cmds, PathCommand{
			Type:       CmdMoveTo,
			Relative:   rel,
			PathPoints: []math64.VectorF2{points[0]},
		})
		for j := 1; j < len(points); j++ {
			cmds = append(cmds, PathCommand{
				Type:       CmdLineTo,
				Relative:   rel,
				PathPoints: []math64.VectorF2{points[j]},
			})
		}
	}
	return cmds, nil
}

func (p *parser) parseLineto() ([]PathCommand, error) {
	// ("L"|"l") wsp* coordinate_pair_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	points, err := p.parseCoordinatePairSequence()
	if err != nil {
		return nil, fmt.Errorf("lineto: %w", err)
	}
	cmds := []PathCommand{}
	for _, pt := range points {
		cmds = append(cmds, PathCommand{
			Type:       CmdLineTo,
			Relative:   rel,
			PathPoints: []math64.VectorF2{pt},
		})
	}
	return cmds, nil
}

func (p *parser) parseHorizontalLineto() ([]PathCommand, error) {
	// ("H"|"h") wsp* coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	coords, err := p.parseCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("horizontal_lineto: %w", err)
	}
	return []PathCommand{{Type: CmdHLineTo, Relative: rel, Coordinates: coords}}, nil
}

func (p *parser) parseVerticalLineto() ([]PathCommand, error) {
	// ("V"|"v") wsp* coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	coords, err := p.parseCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("vertical_lineto: %w", err)
	}
	return []PathCommand{{Type: CmdVLineTo, Relative: rel, Coordinates: coords}}, nil
}

func (p *parser) parseCurveto() ([]PathCommand, error) {
	// ("C"|"c") wsp* curveto_coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	trips, err := p.parseCurvetoCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("curveto: %w", err)
	}
	cmds := []PathCommand{}
	for _, trip := range trips {
		// trip is 3 points (control1, control2, end)
		cmds = append(cmds, PathCommand{Type: CmdCurveTo, Relative: rel, PathPoints: trip})
	}
	return cmds, nil
}

func (p *parser) parseSmoothCurveto() ([]PathCommand, error) {
	// ("S"|"s") wsp* smooth_curveto_coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	seq, err := p.parseSmoothCurvetoCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("smooth_curveto: %w", err)
	}
	cmds := []PathCommand{}
	for _, dbl := range seq {
		// each dbl is 2 points (control2, end)
		cmds = append(cmds, PathCommand{Type: CmdSmoothCurveTo, Relative: rel, PathPoints: dbl})
	}
	return cmds, nil
}

func (p *parser) parseQuadraticBezierCurveto() ([]PathCommand, error) {
	// ("Q"|"q") wsp* quadratic_bezier_curveto_coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	seq, err := p.parseQuadraticBezierCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("quadratic_bezier_curveto: %w", err)
	}
	cmds := []PathCommand{}
	for _, dbl := range seq {
		// dbl: 2 points (control, end)
		cmds = append(cmds, PathCommand{Type: CmdQuadraticBezierTo, Relative: rel, PathPoints: dbl})
	}
	return cmds, nil
}

func (p *parser) parseSmoothQuadraticBezierCurveto() ([]PathCommand, error) {
	// ("T"|"t") wsp* coordinate_pair_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	points, err := p.parseCoordinatePairSequence()
	if err != nil {
		return nil, fmt.Errorf("smooth_quadratic_bezier_curveto: %w", err)
	}
	cmds := []PathCommand{}
	for _, pt := range points {
		cmds = append(cmds, PathCommand{Type: CmdSmoothQuadraticBezierTo, Relative: rel, PathPoints: []math64.VectorF2{pt}})
	}
	return cmds, nil
}

func (p *parser) parseEllipticalArc() ([]PathCommand, error) {
	// ( "A" | "a" ) wsp* elliptical_arc_argument_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	args, err := p.parseEllipticalArcArgumentSequence()
	if err != nil {
		return nil, fmt.Errorf("elliptical_arc: %w", err)
	}
	cmds := []PathCommand{}
	for _, a := range args {
		cmds = append(cmds, PathCommand{Type: CmdEllipticalArc, Relative: rel, ArcArgs: []EllipticalArcArg{a}})
	}
	return cmds, nil
}

func ParseSVGPath(s string) ([]PathCommand, error) {
	p := newParser(s)
	cmds := []PathCommand{}

	for {
		p.skipWsp()
		if p.eof() {
			break
		}
		ch := p.peek()
		// PathCommands are letters among: Mm Zz L l H h V v C c S s Q q T t A a
		if !isPathCommandLetter(ch) {
			// If unexpected char, try to skip whitespace/comma_wsp; otherwise error
			return nil, fmt.Errorf("unexpected char %q at pos %d", ch, p.i)
		}
		// parse the next drawto_command (which may be moveto again)
		subCmds, err := p.parseDrawtoPathCommand()
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, subCmds...)
	}
	return cmds, nil
}
