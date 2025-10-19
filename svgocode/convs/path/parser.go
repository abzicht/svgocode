package path

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// God dammit, this code was produced by passing the SVG path ENBF grammar to
// AI.
// Not fully revised yet, but will do.

// -------- AST types --------

type CommandType int

const (
	CmdMoveTo CommandType = iota
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

func (c CommandType) String() string {
	return [...]string{
		"MoveTo", "ClosePath", "LineTo", "HLineTo", "VLineTo",
		"CurveTo", "SmoothCurveTo", "QuadraticBezierTo",
		"SmoothQuadraticBezierTo", "EllipticalArc",
	}[c]
}

type Point struct {
	X float64
	Y float64
}

type EllipticalArcArg struct {
	Rx    float64
	Ry    float64
	XAxis float64 // x-axis rotation
	Large int     // flag 0/1
	Sweep int     // flag 0/1
	To    Point
}

type Command struct {
	Type        CommandType
	Relative    bool
	Points      []Point   // used for moveto/lineto/curveto combos
	Coordinates []float64 // used for H/V coordinates (list of floats)
	ArcArgs     []EllipticalArcArg
}

func (c Command) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s (relative=%v): ", c.Type, c.Relative)
	switch c.Type {
	case CmdMoveTo, CmdLineTo, CmdCurveTo, CmdSmoothCurveTo, CmdQuadraticBezierTo, CmdSmoothQuadraticBezierTo:
		fmt.Fprintf(&b, "Points=%v", c.Points)
	case CmdHLineTo, CmdVLineTo:
		fmt.Fprintf(&b, "Coords=%v", c.Coordinates)
	case CmdEllipticalArc:
		fmt.Fprintf(&b, "Arcs=%v", c.ArcArgs)
	case CmdClosePath:
		// nothing more
	default:
		fmt.Fprintf(&b, "Points=%v", c.Points)
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
func (p *parser) parseFloat() (float64, error) {
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
	return val, nil
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
func (p *parser) parseCoordinate() (float64, error) {
	p.skipWsp()
	sign := 1.0
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
func (p *parser) parseCoordinatePair() (Point, error) {
	x, err := p.parseCoordinate()
	if err != nil {
		return Point{}, err
	}
	p.consumeCommaWspOptional()
	y, err := p.parseCoordinate()
	if err != nil {
		return Point{}, err
	}
	return Point{X: x, Y: y}, nil
}

// coordinate_pair_sequence ::= coordinate_pair | (coordinate_pair comma_wsp? coordinate_pair_sequence)
func (p *parser) parseCoordinatePairSequence() ([]Point, error) {
	var pts []Point
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
func (p *parser) parseCoordinateSequence() ([]float64, error) {
	var coords []float64
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
func (p *parser) parseCoordinatePairTriplet() ([]Point, error) {
	var pts []Point
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

func (p *parser) parseCurvetoCoordinateSequence() ([][]Point, error) {
	// one or more triplets
	first, err := p.parseCoordinatePairTriplet()
	if err != nil {
		return nil, err
	}
	result := [][]Point{first}
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
func (p *parser) parseCoordinatePairDouble() ([]Point, error) {
	var pts []Point
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

func (p *parser) parseSmoothCurvetoCoordinateSequence() ([][]Point, error) {
	first, err := p.parseCoordinatePairDouble()
	if err != nil {
		return nil, err
	}
	res := [][]Point{first}
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

func (p *parser) parseQuadraticBezierCoordinateSequence() ([][]Point, error) {
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
	large, err := p.parseFlag()
	if err != nil {
		return EllipticalArcArg{}, fmt.Errorf("elliptical arc large-arc-flag: %w", err)
	}
	p.consumeCommaWspOptional()
	sweep, err := p.parseFlag()
	if err != nil {
		return EllipticalArcArg{}, fmt.Errorf("elliptical arc sweep-flag: %w", err)
	}
	p.consumeCommaWspOptional()
	pt, err := p.parseCoordinatePair()
	if err != nil {
		return EllipticalArcArg{}, fmt.Errorf("elliptical arc end point: %w", err)
	}
	return EllipticalArcArg{
		Rx:    rx,
		Ry:    ry,
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

// -------- Top-level parsing according to grammar --------

func ParseSVGPath(s string) ([]Command, error) {
	p := newParser(s)
	cmds := []Command{}
	// // wsp*
	// p.skipWsp()

	// // optional moveto?
	// if p.eof() {
	// 	return cmds, nil
	// }
	// // if next is moveto (M/m)
	// if ch := p.peek(); ch == 'M' || ch == 'm' {
	// 	moveCmds, err := p.parseMoveto()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	cmds = append(cmds, moveCmds...)
	// }

	// (moveto drawto_command*)?  The EBNF allows an optional outer moveto followed by drawto commands.
	// But drawto_command itself can also include moveto. We'll parse remaining commands until EOF.
	for {
		p.skipWsp()
		if p.eof() {
			break
		}
		ch := p.peek()
		// Commands are letters among: Mm Zz L l H h V v C c S s Q q T t A a
		if !isCommandLetter(ch) {
			// If unexpected char, try to skip whitespace/comma_wsp; otherwise error
			return nil, fmt.Errorf("unexpected char %q at pos %d", ch, p.i)
		}
		// parse the next drawto_command (which may be moveto again)
		subCmds, err := p.parseDrawtoCommand()
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, subCmds...)
	}
	return cmds, nil
}

func isCommandLetter(r rune) bool {
	switch r {
	case 'M', 'm', 'Z', 'z', 'L', 'l', 'H', 'h', 'V', 'v', 'C', 'c', 'S', 's', 'Q', 'q', 'T', 't', 'A', 'a':
		return true
	default:
		return false
	}
}

func (p *parser) parseDrawtoCommand() ([]Command, error) {
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
		return []Command{{Type: CmdClosePath, Relative: unicode.IsLower(ch)}}, nil
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

func (p *parser) parseMoveto() ([]Command, error) {
	// ("M" | "m") wsp* coordinate_pair_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	points, err := p.parseCoordinatePairSequence()
	if err != nil {
		return nil, fmt.Errorf("moveto: %w", err)
	}
	// According to SVG: moveto with multiple coordinate pairs: first is MoveTo, the rest are implicit LineTo
	cmds := []Command{}
	if len(points) >= 1 {
		cmds = append(cmds, Command{
			Type:     CmdMoveTo,
			Relative: rel,
			Points:   []Point{points[0]},
		})
		for j := 1; j < len(points); j++ {
			cmds = append(cmds, Command{
				Type:     CmdLineTo,
				Relative: rel,
				Points:   []Point{points[j]},
			})
		}
	}
	return cmds, nil
}

func (p *parser) parseLineto() ([]Command, error) {
	// ("L"|"l") wsp* coordinate_pair_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	points, err := p.parseCoordinatePairSequence()
	if err != nil {
		return nil, fmt.Errorf("lineto: %w", err)
	}
	cmds := []Command{}
	for _, pt := range points {
		cmds = append(cmds, Command{
			Type:     CmdLineTo,
			Relative: rel,
			Points:   []Point{pt},
		})
	}
	return cmds, nil
}

func (p *parser) parseHorizontalLineto() ([]Command, error) {
	// ("H"|"h") wsp* coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	coords, err := p.parseCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("horizontal_lineto: %w", err)
	}
	return []Command{{Type: CmdHLineTo, Relative: rel, Coordinates: coords}}, nil
}

func (p *parser) parseVerticalLineto() ([]Command, error) {
	// ("V"|"v") wsp* coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	coords, err := p.parseCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("vertical_lineto: %w", err)
	}
	return []Command{{Type: CmdVLineTo, Relative: rel, Coordinates: coords}}, nil
}

func (p *parser) parseCurveto() ([]Command, error) {
	// ("C"|"c") wsp* curveto_coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	trips, err := p.parseCurvetoCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("curveto: %w", err)
	}
	cmds := []Command{}
	for _, trip := range trips {
		// trip is 3 points (control1, control2, end)
		cmds = append(cmds, Command{Type: CmdCurveTo, Relative: rel, Points: trip})
	}
	return cmds, nil
}

func (p *parser) parseSmoothCurveto() ([]Command, error) {
	// ("S"|"s") wsp* smooth_curveto_coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	seq, err := p.parseSmoothCurvetoCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("smooth_curveto: %w", err)
	}
	cmds := []Command{}
	for _, dbl := range seq {
		// each dbl is 2 points (control2, end)
		cmds = append(cmds, Command{Type: CmdSmoothCurveTo, Relative: rel, Points: dbl})
	}
	return cmds, nil
}

func (p *parser) parseQuadraticBezierCurveto() ([]Command, error) {
	// ("Q"|"q") wsp* quadratic_bezier_curveto_coordinate_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	seq, err := p.parseQuadraticBezierCoordinateSequence()
	if err != nil {
		return nil, fmt.Errorf("quadratic_bezier_curveto: %w", err)
	}
	cmds := []Command{}
	for _, dbl := range seq {
		// dbl: 2 points (control, end)
		cmds = append(cmds, Command{Type: CmdQuadraticBezierTo, Relative: rel, Points: dbl})
	}
	return cmds, nil
}

func (p *parser) parseSmoothQuadraticBezierCurveto() ([]Command, error) {
	// ("T"|"t") wsp* coordinate_pair_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	points, err := p.parseCoordinatePairSequence()
	if err != nil {
		return nil, fmt.Errorf("smooth_quadratic_bezier_curveto: %w", err)
	}
	cmds := []Command{}
	for _, pt := range points {
		cmds = append(cmds, Command{Type: CmdSmoothQuadraticBezierTo, Relative: rel, Points: []Point{pt}})
	}
	return cmds, nil
}

func (p *parser) parseEllipticalArc() ([]Command, error) {
	// ( "A" | "a" ) wsp* elliptical_arc_argument_sequence
	ch := p.next()
	rel := unicode.IsLower(ch)
	p.skipWsp()
	args, err := p.parseEllipticalArcArgumentSequence()
	if err != nil {
		return nil, fmt.Errorf("elliptical_arc: %w", err)
	}
	cmds := []Command{}
	for _, a := range args {
		cmds = append(cmds, Command{Type: CmdEllipticalArc, Relative: rel, ArcArgs: []EllipticalArcArg{a}})
	}
	return cmds, nil
}
