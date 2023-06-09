package ui

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/lcdDisplay"
	"TinyGo/FunctionGenerator/switches"
	"machine"
	"time"
)

type Screen interface {
	Update()
	Push(result bool)
	Rotate(increment int32)
}

var (
	lcd  *lcdDisplay.Device
	fgen *ad9833.Device

	Changed bool // indicates the screen needs updating

	rotaryLastTime  int64
	nextScreen      Screen
	displayedScreen Screen
)

func Configure(frequencyGen *ad9833.Device) {
	fgen = frequencyGen
	configureScreen()
	configureKeyboard()
}

func configureScreen() {
	lcd = lcdDisplay.NewDevice()

	nextScreen = nil
	displayedScreen = NewScreenMenu() // inital screen
	Changed = true

	//Refresh screen periodically in a go routine
	//go func(t1 *text.Label, t2 *text.Label) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 250)

		for range ticker.C {
			if nextScreen != nil {
				displayedScreen = nextScreen
				nextScreen = nil
				Changed = true
			}
			if Changed {
				displayedScreen.Update()
				lcd.Display()
				Changed = false
			}
		}

	}()

}

func configureKeyboard() {
	switches.SetupPush(machine.GP5, 1000, func(result bool) {
		displayedScreen.Push(result) //have to do it via this func to avoid runtime panic
	})

	switches.NewRotary(machine.GP6, machine.GP7, func(increment int32) {
		displayedScreen.Rotate(increment) //have to do it via this func to avoid runtime panic
	})
}

func ChangeScreen(screen Screen) {
	fgen.SetFrequency(0, ad9833.ADR_FREQ0) // turn off output
	nextScreen = screen
}

func VaryInt32Between(value int32, increment int32, min int32, max int32) int32 {
	value += increment
	if value < min { //down
		value = min
	}
	if value > max { //up
		value = max
	}
	return value
}

func VaryFloat64Between(value float64, increment float64, min float64, max float64) float64 {
	value += increment
	if value < min { //down
		value = min
	}
	if value > max { //up
		value = max
	}
	return value
}
