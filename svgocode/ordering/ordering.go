package ordering

import "github.com/abzicht/svgocode/svgocode/gcode"

type OrderingI interface {
	Order([]*gcode.Gcode) []*gcode.Gcode
}
