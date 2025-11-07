package gcode

import (
	"math"
	"regexp"
	"slices"
	"strings"

	"github.com/abzicht/svgocode/llog"
	"github.com/abzicht/svgocode/svgocode/conf"
	"github.com/abzicht/svgocode/svgocode/math64"
)

// Match all lines that start with a comment
var reLineComment = regexp.MustCompile("^\\s*;(.*)$")

// Match all lines that contain a comment
var reComment = regexp.MustCompile("^.*;(.*)$")

// Match all lines that contain an instruction
var reInstruction = regexp.MustCompile("^\\s*([GMTgmt][\\d]+[^;]*)")

// Pure code
type Code struct {
	lines []string
}

func NewCode() *Code {
	return new(Code)
}

func (c *Code) Copy() *Code {
	c2 := NewCode()
	c2.lines = slices.Clone(c.lines)
	return c2
}

func (c *Code) NumLines() int {
	return len(c.lines)
}

func (c *Code) NumInstructions() int {
	counter := 0
	for _, line := range c.lines {
		if reInstruction.MatchString(line) {
			counter += 1
		}
	}
	return counter
}

func (c *Code) NumComments() int {
	counter := 0
	for _, line := range c.lines {
		if reComment.MatchString(line) {
			counter += 1
		}
	}
	return counter
}

func (c *Code) RemoveComments() *Code {
	c2 := NewCode()
	for _, line := range c.lines {
		if reLineComment.MatchString(line) {
			continue
		}
		if reComment.MatchString(line) {
			line = reInstruction.FindString(line)
		}
		c2.AppendLines(line)
	}
	return c2
}

func (c *Code) String() string {
	return strings.Join(c.lines, "\n") + "\n"
}

func (c *Code) AppendLines(lines ...string) {
	c.lines = append(c.lines, lines...)
}

func (c *Code) Append(c2 *Code) {
	c.AppendLines(c2.lines...)
}

// Gcode, also holds auxiliary information.
type Gcode struct {
	Code       *Code           // The actual gcode.
	StartCoord math64.VectorF3 // Start coordinates of the given Gcode (if known)
	EndCoord   math64.VectorF3 // End coordinates of the given Gcode (if known)
	BoundsMin  math64.VectorF3 // Minimum coordinates used in Gcode (if known)
	BoundsMax  math64.VectorF3 // Maximum coordinates used in Gcode (if known)
}

func NewGcode() *Gcode {
	g := new(Gcode)
	g.BoundsMin = math64.VectorF3{X: math.MaxFloat64, Y: math.MaxFloat64, Z: math.MaxFloat64}
	// better hope that BoundsMin will be overwritten
	g.Code = NewCode()
	return g
}

// Create a copy of the gcode, but only of its metadata without the actual code
func (g *Gcode) CopyMeta() *Gcode {
	g2 := NewGcode()
	g2.StartCoord = g.StartCoord
	g2.EndCoord = g.EndCoord
	g2.BoundsMin = g.BoundsMin
	g2.BoundsMax = g.BoundsMax
	return g2
}

// Full copy, including the code.
func (g *Gcode) Copy() *Gcode {
	g2 := g.CopyMeta()
	g2.Code = g.Code.Copy()
	return g2
}

func (g *Gcode) String() string {
	return g.Code.String()
}

// Add newline-separated code
func (g *Gcode) AppendCode(code string) {
	if len(code) == 0 {
		return
	}
	lines := strings.Split(code, "\n")
	g.Code.AppendLines(lines...)
}

// Join two gcodes, merging their boundaries, start/end coordinates, and code.
// Adds Retract command in-between both codes, if they end/start at different
// positions.
func (g *Gcode) Append(g2 *Gcode, plotterConf *conf.PlotterConfig) {
	ins := NewIns(plotterConf)
	g2StartRetracted := math64.VectorF3{X: g2.StartCoord.X, Y: g2.StartCoord.Y, Z: plotterConf.RetractHeight}
	gEndRetracted := math64.VectorF3{X: g.EndCoord.X, Y: g2.EndCoord.Y, Z: plotterConf.RetractHeight}
	if !g.EndCoord.Equal(g2.StartCoord) && !g.EndCoord.Equal(g2StartRetracted) && !gEndRetracted.Equal(g2.StartCoord) {
		ins.Retract(g)
		ins.Move(g, g2StartRetracted, plotterConf.RetractSpeed)
	}
	g.EndCoord = g2.EndCoord
	g.BoundsMin = g.BoundsMin.Min(g2.BoundsMin)
	g.BoundsMax = g.BoundsMax.Max(g2.BoundsMax)
	g.Code.Append(g2.Code)
}

func Join(gcodes []*Gcode, plotterConf *conf.PlotterConfig) *Gcode {
	if len(gcodes) == 0 {
		llog.Panic("Cannot join gcode segments, provided list is empty")
	}
	g := gcodes[0].Copy()
	for _, g2 := range gcodes[1:] {
		g.Append(g2, plotterConf)
	}
	return g
}

// Total euclidean distance between the end- and start coordinates of gcode
// segments. I.e., the distance travelled where nothing is being drawn.
func TotalDistanceInBetween(gcodes []*Gcode) math64.Float {
	var dist math64.Float = 0
	for i, _ := range gcodes[1:] {
		dist += gcodes[i].EndCoord.DistEuclid(gcodes[i+1].StartCoord)
	}
	return dist
}
