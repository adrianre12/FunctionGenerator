package spix

import (
	"device/rp"
	"machine"
	"time"
)

const _SPITimeout = 10000 //10ms

type SPIX struct {
	machine.SPI
}

func NewSPIX(spi *machine.SPI) *SPIX {
	spix := SPIX{}
	spix.Bus = spi.Bus
	return &spix
}

// Transfer upto 16 bits per frame, use SetDatabits to set the frame size after Configure.
func (spi *SPIX) Transfer16(w uint16) (uint16, error) {
	var deadline = time.Now().UnixMicro() + _SPITimeout
	for !spi.Bus.SSPSR.HasBits(rp.SPI0_SSPSR_TNF) {
		if time.Now().UnixMicro() > deadline {
			return 0, machine.ErrSPITimeout
		}
	}

	spi.Bus.SSPDR.Set(uint32(w))

	for !spi.Bus.SSPSR.HasBits(rp.SPI0_SSPSR_RNE) {
		if time.Now().UnixMicro() > deadline {
			return 0, machine.ErrSPITimeout
		}
	}
	return uint16(spi.Bus.SSPDR.Get()), nil
}

// Allow setting the databits in the rp2040
func (spi *SPIX) SetDatabits(databits uint32) {
	spi.Bus.SSPCR0.ReplaceBits(
		(uint32(databits-1) << rp.SPI0_SSPCR0_DSS_Pos), // Set databits (SPI word length). Valid inputs are 4-16.
		rp.SPI0_SSPCR0_DSS_Msk, 0)
}

func (SPIX) SetCSn(pin machine.Pin) {
	pin.Configure(machine.PinConfig{Mode: machine.PinSPI})
}
