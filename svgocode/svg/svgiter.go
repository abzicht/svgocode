package svg

import (
	"iter"
	"slices"

	"github.com/abzicht/gogenericfunc/fun"
	"github.com/abzicht/svgocode/llog"
)

// Iterators for stepping through SVG trees

func recurse(yield func(SVGElement) bool, svgElements ...SVGElement) bool {
	for _, s := range svgElements {
		if !yield(s) {
			return false
		}
		childrenOpt := s.Children()
		switch childrenOpt.(type) {
		case fun.Some[[]SVGElement]:
			children := childrenOpt.GetValue()
			if !recurse(yield, children...) {
				return false
			}
		case fun.None[[]SVGElement]:
		default:
			llog.Panicf("Unknown Option type for SVG elements: %T\n", childrenOpt)
		}
	}
	return true
}

// Iterate over all children of the given element. Yields all nodes, including
// non-leafs and including the element itself.
func Seq(s SVGElement) iter.Seq[SVGElement] {
	return func(yield func(SVGElement) bool) {
		recurse(yield, s)
	}
}

func recursePath(yield func([]SVGElement) bool, currentPath []SVGElement, svgElements ...SVGElement) bool {
	for _, s := range svgElements {
		path := slices.Clone(currentPath)
		path = append(path, s)
		if !yield(path) {
			return false
		}
		childrenOpt := s.Children()
		switch childrenOpt.(type) {
		case fun.Some[[]SVGElement]:
			children := childrenOpt.GetValue()
			if !recursePath(yield, path, children...) {
				return false
			}
		case fun.None[[]SVGElement]:
		default:
			llog.Panicf("Unknown Option type for SVG elements: %T\n", childrenOpt)
		}
	}
	return true
}

// Iterate over all children of the given element. Yields all paths from root
// to sub-nodes (including non-leafs and a path that only contains the root
// node)
func PathSeq(s SVGElement) iter.Seq[[]SVGElement] {
	return func(yield func([]SVGElement) bool) {
		recursePath(yield, []SVGElement{}, s)
	}
}
