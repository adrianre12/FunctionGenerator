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

// -----------------------------------------------------------
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

// -----------------------------------------------------------
type FieldInt32 struct {
	Label
	Value  int32
	Format string
}

func NewFieldInt32(font *tinyfont.Font, x int16, y int16, value int32) *FieldInt32 {
	return &FieldInt32{
		Label:  newLabel(font, x, y),
		Value:  value,
		Format: "%d",
	}
}

func (fi *FieldInt32) String() string {
	return fmt.Sprintf(fi.Format, fi.Value)
}

// -----------------------------------------------------------
type FieldFloat32 struct {
	Label
	Value  float32
	Format string
}

func NewFieldFloat32(font *tinyfont.Font, x int16, y int16, value float32) *FieldFloat32 {
	return &FieldFloat32{
		Label:  newLabel(font, x, y),
		Value:  value,
		Format: "%.2f",
	}
}

func (ff *FieldFloat32) String() string {
	return fmt.Sprintf(ff.Format, ff.Value)
}

// -----------------------------------------------------------

type FieldListItem struct {
	Text  string
	Value uint32
}

type FieldList struct {
	Label
	Selected int32
	Values   []FieldListItem
}

func NewFieldList(font *tinyfont.Font, x int16, y int16, values []FieldListItem) *FieldList {
	return &FieldList{
		Label:  newLabel(font, x, y),
		Values: values,
	}
}

func (fl *FieldList) String() string {
	return fl.Values[fl.Selected].Text
}

func (fl *FieldList) Value() uint32 {
	return fl.Values[fl.Selected].Value
}
