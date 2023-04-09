package lcdDisplay

import (
	"tinygo.org/x/tinyfont"
)

type Label struct {
	font *tinyfont.Font
	bold bool
	x    int16
	y    int16
}

func newLabel(font *tinyfont.Font, x int16, y int16) Label {
	l := Label{
		font: font,
		x:    x,
		y:    y,
	}

	return l
}

func (l *Label) Bold(bold bool) {
	l.bold = bold
}

func (l *Label) Font() *tinyfont.Font {
	return l.font
}

func (l *Label) IsBold() bool {
	return l.bold
}

func (l *Label) X() int16 {
	return l.x
}

func (l *Label) Y() int16 {
	return l.y
}
