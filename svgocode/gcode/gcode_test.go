package gcode

import "testing"

func TestCodeNumComments(t *testing.T) {
	code := []string{
		";Comment 1",
		"; Comment 2 ; foo",
		"; Comment 3 ; G1 X1 Y1",
		"; G0 X1 Y2 ; Comment 4",
		"  ; Comment 5 ",
		"G0 X2 Y2; Comment 6",
		"G1 X1 Y2 ; Comment 7 ",
		"G1 X1 Y2 ;",
		"G1 X1 Y2",
		"G1",
	}
	numComments := []int{1, 2, 3, 4, 5, 6, 7, 8, 8, 8}
	numInstructions := []int{0, 0, 0, 0, 0, 1, 2, 3, 4, 5}
	for i, _ := range code {
		c := NewCode()
		c.AppendLines(code[0 : i+1]...)
		if numComments[i] != c.NumComments() {
			t.Errorf("Failed to detect line '%s' as comment", code[i])
		}
		if numInstructions[i] != c.NumInstructions() {
			t.Errorf("Failed to detect line '%s' as instruction (counting %d)", code[i], c.NumInstructions())
		}
	}
}
