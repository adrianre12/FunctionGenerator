package ui

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/switches"
	"machine"
	"time"

	"tinygo.org/x/drivers"
	"tinygo.org/x/drivers/pcd8544"
)

type Display interface {
	drivers.Displayer
	ClearBuffer()
	ClearDisplay()
}

type Screen interface {
	Setup()
	Update()
	Push(result bool)
	Rotate(result bool)
}

var (
	spi1 = machine.SPI1
	lcd  Display
	fgen *ad9833.Device

	Waveform  ad9833.Mode
	Frequency float32
	Changed   bool

	rotaryLastTime  int64
	nextScreen      Screen
	displayedScreen Screen
	Menu            *ScreenMenu
	Manual          *ScreenManual
)

func Configure(frequencyGen *ad9833.Device) {
	fgen = frequencyGen
	configureScreen()
	configureKeyboard()
}

func configureScreen() {
	Menu = &ScreenMenu{}
	Manual = &ScreenManual{}
	nextScreen = nil
	displayedScreen = Menu // inital screen
	setupLCD()

	displayedScreen.Setup()
	lcd.Display()

	//Refresh screen periodically in a go routine
	//go func(t1 *text.Label, t2 *text.Label) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 250)

		for range ticker.C {
			if nextScreen != nil {
				displayedScreen = nextScreen
				nextScreen = nil
				displayedScreen.Setup()
			}
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

func setupLCD() {
	spi1.Configure(machine.SPIConfig{
		Frequency: 4000000,
		LSBFirst:  false,
		Mode:      0,
		DataBits:  8,
		SCK:       machine.GP10,
		SDO:       machine.GP11,
	})

	dcPin := machine.GP14
	dcPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin := machine.GP15
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	scePin := machine.GP9
	scePin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	lcdDevice := pcd8544.New(spi1, dcPin, rstPin, scePin)
	lcdDevice.Configure(pcd8544.Config{})

	lcdDevice.SendCommand(pcd8544.FUNCTIONSET | pcd8544.EXTENDEDINSTRUCTION)
	lcdDevice.SendCommand(pcd8544.SETBIAS | 0x04)
	lcd = lcdDevice //use Device interface for the rest so the lcd device can be changed.
	lcd.ClearDisplay()
}

func ChangeScreen(screen Screen) {
	nextScreen = screen
}
