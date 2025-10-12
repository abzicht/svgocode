package svg

import (
	"encoding/xml"
	"io"
)

type Decoder struct {
	xmlDecoder *xml.Decoder
}

func NewDecoder(r io.Reader) *Decoder {
	d := new(Decoder)
	d.xmlDecoder = xml.NewDecoder(r)
	return d
}

func (d *Decoder) Decode(svg *SVG) error {
	return d.xmlDecoder.Decode(svg)
}
