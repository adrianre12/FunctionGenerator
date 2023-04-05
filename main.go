package main

import (
	ad9833 "TinyGo/FunctionGenerator/AD9833"
	text "TinyGo/FunctionGenerator/Text"
	"TinyGo/FunctionGenerator/spix"
	"TinyGo/FunctionGenerator/switches"
	"fmt"
	"machine"
	"time"

	"tinygo.org/x/drivers/pcd8544"
	"tinygo.org/x/tinyfont/proggy"
)

var (
	spi0      = machine.SPI0
	spi1      = machine.SPI1
	LCDdevice *pcd8544.Device
	fgen      *ad9833.Device
	pwm0      = machine.PWM0
	pwm0A     uint8

	mode           ad9833.Mode
	frequency      float32
	changed        bool
	rotaryLastTime int64

	sweepStart    float32 //Sweep start frequency Hz
	sweepEnd      float32 //Sweep end frequency Hz
	sweepStep     float32 //Sweep step size Hz
	sweepStepTime uint16  //Target duration of each step mS
	sweepGate     bool    //If set the output will be off while not sweeping.
)

func ConfigureLCD() {

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
	LCDdevice = pcd8544.New(spi1, dcPin, rstPin, scePin)
	LCDdevice.Configure(pcd8544.Config{})

	LCDdevice.SendCommand(pcd8544.FUNCTIONSET | pcd8544.EXTENDEDINSTRUCTION) // H = 1
	//LCDdevice.SendCommand(pcd8544.SETVOP | 0x3f)                             // 0x3f : Vop6 = 0, Vop5 to Vop0 = 1
	//LCDdevice.SendCommand(pcd8544.SETTEMP | 0x03)                            // Experimentally determined
	LCDdevice.SendCommand(pcd8544.SETBIAS | 0x04)

	LCDdevice.ClearBuffer()
	LCDdevice.ClearDisplay()
}

func SerialDelayStart(t int) {
	for i := t; i > 0; i-- {
		fmt.Printf("Starting... %d\n", i)
		time.Sleep(time.Second)
	}
}

func ConfigureRotary() {
	switches.SetupPush(machine.GP5, func(result bool) {
		if result { // ignore when pushed only use btn up
			return
		}
		fmt.Println("Released")
		mode = mode.Next()
		changed = true
		fmt.Printf("Mode %s value %v \n", mode.String(), mode.Uint16())
	})

	switches.SetupRotary(machine.GP6, machine.GP7, func(result bool) {
		delta := time.Now().UnixMilli() - rotaryLastTime
		var increment float32
		switch {
		case delta < 25:
			{
				increment = 50
			}
		case delta < 75:
			{
				increment = 10
			}
		case delta < 150:
			{
				increment = 5
			}
		default:
			{
				increment = 1
			}
		}
		if result {
			frequency += increment
		} else {
			frequency -= increment
		}
		if frequency < 0 {
			frequency = 0
		}
		changed = true
		rotaryLastTime = time.Now().UnixMilli()
	})
}

func ConfigureScreen() {
	ConfigureLCD()
	font := &proggy.TinySZ8pt7b

	label1 := text.NewLabel(LCDdevice, font, 0, 7, "Mode ")
	_, labelW := label1.LineWidth()
	text1 := text.NewLabel(LCDdevice, font, int16(labelW), 7, fmt.Sprintf("%s", mode))

	label2 := text.NewLabel(LCDdevice, font, 0, 18, "Freq ")
	_, labelW = label2.LineWidth()
	text2 := text.NewLabel(LCDdevice, font, int16(labelW), 18, fmt.Sprintf("%d", frequency))
	LCDdevice.Display()

	//Refresh screen periodically in a ro routine
	go func(t1 *text.Label, t2 *text.Label) {
		ticker := time.NewTicker(time.Millisecond * 250)

		for range ticker.C {
			if changed {
				fgen.SetMode(mode)
				frequency = fgen.SetFrequency(frequency, ad9833.ADR_FREQ0)
				changed = false
			}

			t1.Write(fmt.Sprintf("%s", mode))
			t2.Write(fmt.Sprintf("%.3f", frequency))
			LCDdevice.Display()
		}

	}(text1, text2)

}

func ConfigurePWM() {
	err := pwm0.Configure(machine.PWMConfig{
		Period: 100000,
	})
	if err != nil {
		println("failed to configure PWM")
		return
	}

	if err != nil {
		println("failed to configure channel A")
		return
	}

	pwm0A, err = pwm0.Channel(machine.GPIO0)
	pwm0.Set(pwm0A, 0)
}

func setPWM(ratio float32) {
	pwm0.Set(pwm0A, uint32(ratio*float32(pwm0.Top())))
}

// crude sweep
func StartSweep() {
	steps := (sweepEnd - sweepStart) / sweepStep
	var step float32 = 0
	for f := sweepStart; f <= sweepEnd; f += sweepStep {
		//fmt.Println(f)
		step++
		setPWM(step / steps)
		frequency = fgen.SetFrequency(f, ad9833.ADR_FREQ0)
		time.Sleep(time.Millisecond * time.Duration(sweepStepTime))
	}
	if sweepGate {
		frequency = fgen.SetFrequency(0, ad9833.ADR_FREQ0)
	}
}

func main() {
	SerialDelayStart(5)
	ConfigureRotary()

	spix := spix.NewSPIX(machine.SPI0)
	spix.Configure(machine.SPIConfig{
		Frequency: 100000,
		LSBFirst:  false,
		Mode:      2,
		DataBits:  8,
		SCK:       machine.GP18,
		SDO:       machine.GP19,
	})
	spix.SetCSn(machine.GP17)
	spix.SetDatabits(16)

	fgen = ad9833.NewAD9833(spix)
	fgen.Init()
	fgen.WriteErr = false

	ConfigureScreen()
	ConfigurePWM()

	sweepTest()
	//pwmTest()

	select {}
}

func pwmTest() {
	pwm := machine.PWM0

	err := pwm.Configure(machine.PWMConfig{
		Period: 100000,
	})
	if err != nil {
		println("failed to configure PWM")
		return
	}

	println("period=", pwm.Period())
	// The top value is the highest value that can be passed to PWMChannel.Set.
	// It is usually an even number.
	println("top:", pwm.Top())

	// Configure the two channels we'll use as outputs.
	channelA, err := pwm.Channel(machine.GPIO0)
	if err != nil {
		println("failed to configure channel A")
		return
	}

	pwm.Set(channelA, pwm.Top()/2)

	fmt.Printf("DIV = %x\n", pwm.DIV.Get())
	fmt.Printf("cpuFreq=%d\n", machine.CPUFrequency())
}

func sweepTest() {
	sweepStart = 100
	sweepEnd = 1000
	sweepStep = 1
	sweepStepTime = 10
	sweepGate = true

	StartSweep()
}
