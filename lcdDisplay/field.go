package lcdDisplay

import (
	"fmt"

	"tinygo.org/x/tinyfont"
)

type Field interface {
	Bold(bold bool) //Set state of bold flag
	Font() *tinyfont.Font
	IsBold() bool
	X() int16
	Y() int16
	String() string
}

type FieldStr struct {
	Label
	Value string
}

func NewFieldStr(font *tinyfont.Font, x int16, y int16, value string) *FieldStr {
	return &FieldStr{
		Label: newLabel(font, x, y),
		Value: value,
	}
}

func (fs *FieldStr) String() string {
	return fs.Value
}

type FieldInt struct {
	Label
	Value int
}

func NewFieldInt(font *tinyfont.Font, x int16, y int16, value int) *FieldInt {
	return &FieldInt{
		Label: newLabel(font, x, y),
		Value: value,
	}
}

func (fi *FieldInt) String() string {
	return fmt.Sprintf("%d", fi.Value)
}
