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
	Invert    bool
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
	fg := FG
	bg := BG
	if t.Invert {
		fg = BG
		bg = FG
	}
	if len(t.text) != 0 {
		//tinyfont.WriteLine(t.displayer, t.font, t.X, t.Y, t.text, bg)
		t.Clear(bg)
	}
	t.text = text
	tinyfont.WriteLine(t.displayer, t.font, t.X, t.Y, t.text, fg)
}

func (t *Label) LineWidth() (innerWidth uint32, outboxWidth uint32) {
	return tinyfont.LineWidth(t.font, t.text)
}

func (t *Label) Clear(colour color.RGBA) {
	_, outboxWidth := t.LineWidth()
	bbox := t.font.BBox
	var x int16
	var y int16
	for i := int16(0); i < int16(outboxWidth); i++ {
		x = t.X + i
		for j := int16(0); j < int16(bbox[1]); j++ {
			y = t.Y + j + int16(bbox[3])
			t.displayer.SetPixel(x, y, colour)
		}
	}
}
