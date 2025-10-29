package svg

import (
	"fmt"
	"iter"
	"slices"

	"github.com/abzicht/svgocode/llog"
)

// Iterators for stepping through SVG trees

func recurse(yield func(SVGElement) bool, svgElements ...SVGElement) bool {
	for _, s := range svgElements {
		if !yield(s) {
			return false
		}
		children := s.Children()
		if len(children) != 0 {
			if !recurse(yield, children...) {
				return false
			}
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

func recursePath(yield func([]SVGElement) bool, resolveUses bool, currentPath []SVGElement, svgElements ...SVGElement) bool {
	for _, s := range svgElements {
		if resolveUses {
			switch s.(type) {
			case *Defs:
				// If we resolve uses, we don't want to step through defs
				return true
			case *Use:
				if len(currentPath) == 0 {
					llog.Panicf("Cannot resolve 'use' tag's reference due to missing root node")
				}
				idMap := SvgToMap(currentPath[0])
				refElement := s.(*Use).GetRefElement(idMap).CloneSVGElement()
				refElement.AppendTransform(fmt.Sprintf("translate(%f, %f)", s.(*Use).X, s.(*Use).Y), true)
				if !recursePath(yield, resolveUses, currentPath, refElement) {
					return false
				}
				continue
			}
		}
		path := slices.Clone(currentPath)
		path = append(path, s)
		if !yield(path) {
			return false
		}
		children := s.Children()
		if len(children) != 0 {
			if !recursePath(yield, resolveUses, path, children...) {
				return false
			}
		}
	}
	return true
}

// Iterate over all children of the given element. If resolveUses is true,
// "use" tags are resolved to the referenced path and path that contain "defs"
// tags are skipped. Yields all paths from root to sub-nodes (including
// non-leafs and a path that only contains the root node)
func PathSeq(s SVGElement, resolveUses bool) iter.Seq[[]SVGElement] {
	return func(yield func([]SVGElement) bool) {
		recursePath(yield, resolveUses, []SVGElement{}, s)
	}
}
