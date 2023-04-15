package ui

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/lcdDisplay"

	"tinygo.org/x/tinyfont/proggy"
)

type ScreenManual struct {
	selectedField int32
	selected      bool
	label1        *lcdDisplay.FieldStr
	label2        *lcdDisplay.FieldStr
	modeList      *lcdDisplay.FieldList
	stepList      *lcdDisplay.FieldList
	frequency     *lcdDisplay.FieldFloat64
	actual        *lcdDisplay.FieldFloat64
}

func NewScreenManual() *ScreenManual {
	s := ScreenManual{selectedField: 1}
	font := &proggy.TinySZ8pt7b
	lcd.ClearBuffer()

	s.label1 = lcdDisplay.NewFieldStr(font, 0, 7, "Wave: ")
	s.modeList = lcdDisplay.NewFieldList(font, int16(lcd.LineWidth(s.label1)), 7, []lcdDisplay.FieldListItem{
		{Text: "Sine", Value: ad9833.MODE_SINE},
		{Text: "Tri", Value: ad9833.MODE_TRI},
		{Text: "Sqr", Value: ad9833.MODE_MSB2},
	})

	s.label2 = lcdDisplay.NewFieldStr(font, 0, 17, "Step: ")
	s.stepList = lcdDisplay.NewFieldList(font, int16(lcd.LineWidth(s.label2)), 17, []lcdDisplay.FieldListItem{
		{Text: "0.01 Hz", Value: 1},
		{Text: "Hz", Value: 100},
		{Text: "KHz", Value: 100000},
	})
	s.stepList.Selected = 1

	s.frequency = lcdDisplay.NewFieldFloat64(font, 0, 27, 0)
	s.frequency.Format = "%.2f Hz"

	s.actual = lcdDisplay.NewFieldFloat64(font, 0, 37, 0)
	s.actual.Format = "(%.3f)"
	return &s
}

func (s *ScreenManual) Update() {
	lcd.ClearBuffer()
	fgen.SetMode(uint16(s.modeList.Value()))
	s.actual.Value = fgen.SetFrequency(s.frequency.Value, ad9833.ADR_FREQ0)
	Changed = false

	lcd.WriteField(s.label1)
	lcd.WriteField(s.label2)

	s.modeList.Bold(s.selectedField == 1)
	s.stepList.Bold(s.selectedField == 2)
	s.frequency.Bold(s.selectedField == 3)

	lcd.WriteField(s.modeList)
	lcd.WriteField(s.stepList)
	lcd.WriteField(s.frequency)
	lcd.WriteField(s.actual)
}

func (s *ScreenManual) Push(result bool) {
	s.selected = !s.selected
	if result {
		ChangeScreen(NewScreenMenu())
		return
	}
}

func (s *ScreenManual) Rotate(increment int32) {
	Changed = true
	if s.selected {
		switch s.selectedField {
		case 1:
			{ //mode
				s.modeList.Selected = VaryInt32Between(s.modeList.Selected, increment, 0, 2)
			}
		case 2:
			{ //step
				s.stepList.Selected = VaryInt32Between(s.stepList.Selected, increment, 0, 2)
			}
		case 3:
			{ //frequency

				s.frequency.Value = VaryFloat64Between(s.frequency.Value, float64(increment)*float64(s.stepList.Value())/100, 0, 2e6)
			}
		default:
			{ // non selectable field
				s.selected = false
			}
		}
	} // this is not an if else
	if !s.selected { // not selected to scroll up an down
		s.selectedField = VaryInt32Between(s.selectedField, increment, 1, 3)
	}
}
