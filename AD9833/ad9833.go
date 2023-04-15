package ad9833

import (
	"TinyGo/FunctionGenerator/spix"
	"fmt"
)

type Device struct {
	spi        *spix.SPIX
	controlReg register
	WriteErr   bool //Enable writing of errors to STDOUT
}

// New creates a new AD9833 connection. The SPI bus must already be configured.
func NewAD9833(spi *spix.SPIX) *Device {
	return &Device{
		spi:        spi,
		controlReg: register{value: 0},
	}
}

func (d *Device) Init() {
	d.controlReg.value = uint16(B28 | RESET)
	d.Write(d.controlReg.value)

	//set freq and phase to 0
	d.Write(ADR_FREQ0) //LSB
	d.Write(ADR_FREQ0) //MSB

	d.Write(ADR_FREQ1) //LSB
	d.Write(ADR_FREQ1) //MSB

	d.Write(ADR_PHASE0)
	d.Write(ADR_PHASE1)

	d.controlReg.replaceBits(0, RESET)
	d.Write(d.controlReg.value)
}

// Write directly to the AD9833
// If WriteErr is set, writes errors to STDOUT
func (d *Device) Write(tx uint16) {
	_, err := d.spi.Transfer16(tx)
	if d.WriteErr && err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// Set the sleep mode
// e.g. Sleep(SLEEPDAC) to disable the DAC
func (d *Device) Sleep(mode uint16) {
	d.controlReg.replaceBits(mode, SLEEP1|SLEEP12)
}

// Select which frequency register is used.
// FREQ0 enable = false
// FREQ1 enable = true
func (d *Device) EnableFREQ1(enable bool) {
	if enable {
		d.controlReg.replaceBits(FSELECT, FSELECT)
	} else {
		d.controlReg.replaceBits(0, FSELECT)
	}
	d.Write(d.controlReg.value)
}

// Select which phase register is used.
// PHASE0 enable = false
// PHASE1 enable = true
func (d *Device) EnablePHASE1(enable bool) {
	if enable {
		d.controlReg.replaceBits(PSELECT, PSELECT)
	} else {
		d.controlReg.replaceBits(0, PSELECT)
	}
	d.Write(d.controlReg.value)
}

// Set the selected frequency register in Hz
// returns the actual frequency set.
// FREQ register MSB and LSB are set using B28
// e.g. SetFrequency(1000.0,ADR_FREQ0)
func (d *Device) SetFrequency(freq float64, freqReg uint16) float64 {
	d.Write(d.controlReg.value | B28)
	freqReg = freqReg & (ADR_FREQ0 | ADR_FREQ1)
	//freqValue := uint32(freq * math.Pow(2, 28) / 25e6)
	//freqValue := uint32(freq * 0x10000000 / 25000000)
	freqValue := uint32(freq*10.73741824 + 0.5) //round up to make more acurate

	d.Write(freqReg | uint16(freqValue&BITS14L))
	d.Write(freqReg | uint16((freqValue&BITS14H)>>14))
	return float64(freqValue) / 10.73741824
}

// Set the selected phase register in degrees or radians
// e.g. SetPhase(90.0, ADR_PHASE0, false)
func (d *Device) SetPhase(phase float32, phaseReg uint16, radians bool) {
	phaseReg = phaseReg & (ADR_PHASE0 | ADR_PHASE1)
	if !radians {
		phase = phase * 2 * 3.1415926 / 360
	}
	phaseValue := uint16(phase / 4096)
	d.Write(phaseReg | (phaseValue & BITS12))
}

// Set the output waveform mode using the Mode constants
// SetMode(Mode_TRI)
func (d *Device) SetMode(mode uint16) {
	d.controlReg.replaceBits(mode, MODE_MASK)
	d.Write(d.controlReg.value)
}
