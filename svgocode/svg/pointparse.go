package svg

import (
	"fmt"
	"strings"

	"github.com/abzicht/svgocode/svgocode/math64"
)

func PointsToPathStr(pts []math64.VectorF2, close bool) string {
	if len(pts) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("M %g %g", pts[0].X, pts[0].Y))
	for _, p := range pts[1:] {
		sb.WriteString(fmt.Sprintf(" L %g %g", p.X, p.Y))
	}
	if close {
		sb.WriteString(" Z")
	}
	return sb.String()
}

// Parse "x1,y1 x2,y2 ..." or "x1 y1 x2 y2 ..." formats
func ParsePointString(s string) []math64.VectorF2 {
	var pts []math64.VectorF2
	s = strings.TrimSpace(s)
	if s == "" {
		return pts
	}

	// Normalize commas and spaces
	replacer := strings.NewReplacer(",", " ", "\n", " ")
	fields := strings.Fields(replacer.Replace(s))

	if len(fields)%2 != 0 {
		return pts // malformed
	}

	for i := 0; i < len(fields); i += 2 {
		x := math64.ParseFloat(fields[i])
		y := math64.ParseFloat(fields[i+1])
		pts = append(pts, math64.VectorF2{X: x, Y: y})
	}
	return pts
}
