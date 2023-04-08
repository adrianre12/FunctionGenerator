package main

import (
	"TinyGo/FunctionGenerator/ad9833"
	"TinyGo/FunctionGenerator/spix"
	"TinyGo/FunctionGenerator/ui"
	"fmt"
	"machine"
	"time"
)

var (
	spi0 = machine.SPI0

	fgen  *ad9833.Device
	pwm0  = machine.PWM0
	pwm0A uint8

	sweepStart    float32 //Sweep start frequency Hz
	sweepEnd      float32 //Sweep end frequency Hz
	sweepStep     float32 //Sweep step size Hz
	sweepStepTime uint16  //Target duration of each step mS
	sweepGate     bool    //If set the output will be off while not sweeping.
)

func SerialDelayStart(t int) {
	for i := t; i > 0; i-- {
		fmt.Printf("Starting... %d\n", i)
		time.Sleep(time.Second)
	}
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
		ui.Frequency = fgen.SetFrequency(f, ad9833.ADR_FREQ0)
		time.Sleep(time.Millisecond * time.Duration(sweepStepTime))
	}
	if sweepGate {
		ui.Frequency = fgen.SetFrequency(0, ad9833.ADR_FREQ0)
	}
}

func main() {
	SerialDelayStart(5)

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

	ConfigurePWM()

	ui.Configure(fgen) //ui should be last

	//sweepTest()

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
