package sweep

import (
	"TinyGo/FunctionGenerator/ad9833"
	"machine"
	"time"
)

var (
	fgen           *ad9833.Device
	sweepFreqeuncy float64
	sweepPercent   float32
	pwm0           = machine.PWM0
	pwm0A          uint8
	stopSweep      bool
)

func ConfigureSweep(frequencyGen *ad9833.Device) {
	fgen = frequencyGen
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

func SetPWM(ratio float32) {
	println(ratio)
	pwm0.Set(pwm0A, uint32(ratio*float32(pwm0.Top())))
}

func StopSweep() {
	stopSweep = true
}

// crude sweep
// mode     uint16 the waveform mode
// start    float32 Sweep start frequency Hz
// end      float32 Sweep end frequency Hz
// stepTime  int32  Target duration of each step mS
// steps     int32 Number of steps
// gate     bool    If set the output will be off while not sweeping.
// fn       func() called at the end of a sweep to update ui
func StartSweep(mode uint16, start float64, end float64, stepTime int32, steps int32, gate bool, fn func()) {
	stopSweep = false
	fgen.SetMode(mode)
	stepSize := (end - start) / float64(steps)
	var stepNum float32 = 0
	for f := start; f <= end; f += stepSize {
		if stopSweep {
			break
		}
		stepNum++
		sweepPercent = stepNum / float32(steps)
		SetPWM(sweepPercent)
		sweepFreqeuncy = fgen.SetFrequency(f, ad9833.ADR_FREQ0)
		time.Sleep(time.Millisecond * time.Duration(stepTime))
	}
	if gate {
		sweepFreqeuncy = fgen.SetFrequency(0, ad9833.ADR_FREQ0)
		SetPWM(0)
	}
	fn()
}
