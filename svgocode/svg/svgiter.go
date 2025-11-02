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

// Sub function of PathSeq. Iterates over all svg elements recursively.
// Resolves 'use' references, if resolveUses is true and sMap holds the
// referenced element.
func recursePath(yield func([]SVGElement) bool, resolveUses bool, sMap SvgIdMap, root SVGElement, currentPath []SVGElement, svgElements ...SVGElement) bool {
	for _, s := range svgElements {
		s.SetRoot(root)
		if resolveUses {
			switch s.(type) {
			case *Defs:
				// If we resolve uses, we don't want to step through defs
				return true
			case *Use:
				if len(currentPath) == 0 {
					llog.Panicf("Cannot resolve 'use' tag's reference due to missing root node")
				}
				refElement := s.(*Use).GetRefElement(sMap).CloneSVGElement()
				refElement.AppendTransform(fmt.Sprintf("translate(%f, %f)", s.(*Use).X, s.(*Use).Y), true)
				if !recursePath(yield, resolveUses, sMap, root, currentPath, refElement) {
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
			if !recursePath(yield, resolveUses, sMap, root, path, children...) {
				return false
			}
		}
	}
	return true
}

// Iterate over all paths traversable from the given element. If resolveUses is true,
// "use" tags are resolved to the referenced path (referenced path must be
// present in s); paths that contain "defs" tags are skipped. Yields all paths
// from root to sub-nodes (including non-leafs and a path that only contains
// the root node). For all iterated elements, the referenced root element is
// set to the provided "root".
func PathSeq_(s SVGElement, root SVGElement, resolveUses bool) iter.Seq[[]SVGElement] {
	var sMap SvgIdMap
	if resolveUses {
		sMap = SvgToMap(s)
	}
	return func(yield func([]SVGElement) bool) {
		recursePath(yield, resolveUses, sMap, root, []SVGElement{}, s)
	}
}

// Iterate over all paths that can be reached from the SVG element. 'use' tags
// are translated to their referenced element.
func PathSeq(s SVGElement) iter.Seq[[]SVGElement] {
	return PathSeq_(s, s, true)
}
