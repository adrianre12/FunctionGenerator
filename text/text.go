package text

import (
	"image/color"

	"tinygo.org/x/drivers"
	"tinygo.org/x/tinyfont"
)

var (
	FG = color.RGBA{1, 1, 1, 255}
	BG = color.RGBA{0, 0, 0, 255}
)

type Label struct {
	displayer drivers.Displayer
	font      *tinyfont.Font
	text      string
	X         int16
	Y         int16
}

func NewLabel(displayer drivers.Displayer, font *tinyfont.Font, x int16, y int16, text string) *Label {
	t := Label{
		displayer: displayer,
		font:      font,
		X:         x,
		Y:         y,
	}
	t.Write(text)
	return &t
}

func (t *Label) Write(text string) {
	if len(t.text) != 0 {
		tinyfont.WriteLine(t.displayer, t.font, t.X, t.Y, t.text, BG)
	}
	t.text = text
	tinyfont.WriteLine(t.displayer, t.font, t.X, t.Y, t.text, FG)
}

func (t *Label) LineWidth() (innerWidth uint32, outboxWidth uint32) {
	return tinyfont.LineWidth(t.font, t.text)
}
