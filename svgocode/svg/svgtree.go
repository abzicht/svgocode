package svg

type SvgIdMap map[SvgId]SVGElement

func SvgToMap(s SVGElement) SvgIdMap {
	m := make(SvgIdMap)
	for el := range Seq(s) {
		m[el.ID()] = el
	}
	return m
}
