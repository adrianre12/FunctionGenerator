package ui

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/lcdDisplay"

	"tinygo.org/x/tinyfont/proggy"
)

var ()

type ScreenManual struct {
	selectedField int32
	selected      bool
	label1        *lcdDisplay.FieldStr
	label2        *lcdDisplay.FieldStr
	modeList      *lcdDisplay.FieldList
	frequency     *lcdDisplay.FieldFloat32
}

func NewScreenManual() *ScreenManual {
	s := ScreenManual{selectedField: 1}
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Wave: ")
	_, labelW := lcd.LineWidth(s.label1)
	s.modeList = lcdDisplay.NewFieldList(font, int16(labelW), 7, []lcdDisplay.FieldListItem{
		{Text: "Sine", Value: ad9833.MODE_SINE},
		{Text: "Tri", Value: ad9833.MODE_TRI},
		{Text: "Sqr", Value: ad9833.MODE_MSB2},
	})

	s.label2 = lcdDisplay.NewFieldStr(font, 0, 17, "Freq: ")
	_, labelW = lcd.LineWidth(s.label2)
	s.frequency = lcdDisplay.NewFieldFloat32(font, int16(labelW), 17, 0)

	return &s
}

func (s *ScreenManual) Update() {
	if Changed {
		fgen.SetMode(uint16(s.modeList.Value()))
		s.frequency.Value = fgen.SetFrequency(s.frequency.Value, ad9833.ADR_FREQ0)
		Changed = false
	}

	lcd.WriteField(s.label1)
	lcd.WriteField(s.label2)

	s.modeList.Bold(s.selectedField == 1)
	s.frequency.Bold(s.selectedField == 2)

	lcd.WriteField(s.modeList)
	lcd.WriteField(s.frequency)
}

func (s *ScreenManual) Push(result bool) {
	s.selected = !s.selected
	if result {
		ChangeScreen(NewScreenMenu())
		return
	}

}

func (s *ScreenManual) Rotate(increment int32) {
	if s.selected {
		switch s.selectedField {
		case 1:
			{ //mode
				s.modeList.Selected = VaryBetween(s.modeList.Selected, increment, 0, 2)
				Changed = true
			}
		case 2:
			{ //frequency

				s.frequency.Value = float32(VaryBetween(int32(s.frequency.Value), increment, 0, 2e6))
				Changed = true
			}
		default:
			{ // non selectable field
				s.selected = false
			}
		}
	} // this is not an if else
	if !s.selected { // not selected to scroll up an down
		s.selectedField = VaryBetween(s.selectedField, increment, 1, 2)
	}
}
