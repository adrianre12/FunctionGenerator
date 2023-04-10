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
	Rotate(result bool)
}

var (
	lcd  *lcdDisplay.Device
	fgen *ad9833.Device

	Waveform  ad9833.Mode
	Frequency float32
	Changed   bool

	rotaryLastTime  int64
	nextScreen      Screen
	displayedScreen Screen
)

func Configure(frequencyGen *ad9833.Device) {
	println("UI Configure")
	fgen = frequencyGen
	configureScreen()
	configureKeyboard()
}

func configureScreen() {
	println("ConfigureScreen")
	lcd = lcdDisplay.NewDevice()

	nextScreen = nil
	displayedScreen = NewScreenMenu() // inital screen

	lcd.Display()

	//Refresh screen periodically in a go routine
	//go func(t1 *text.Label, t2 *text.Label) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 250)

		for range ticker.C {
			if nextScreen != nil {
				displayedScreen = nextScreen
				nextScreen = nil
			}
			lcd.ClearBuffer()
			displayedScreen.Update()
			lcd.Display()
		}

	}()

}

func configureKeyboard() {
	switches.SetupPush(machine.GP5, 1000, func(result bool) {
		displayedScreen.Push(result) //have to do it via this func to avoid runtime panic
	})

	switches.SetupRotary(machine.GP6, machine.GP7, func(result bool) {
		displayedScreen.Rotate(result) //have to do it via this func to avoid runtime panic
	})
}

func ChangeScreen(screen Screen) {
	nextScreen = screen
}
