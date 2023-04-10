package ad9833

const (
	B28     uint16 = 0x2000 //D13 0x2000
	HLB     uint16 = 0x1000 //D12 0x1000
	FSELECT uint16 = 0x0800 //D11 0x0800
	PSELECT uint16 = 0x0400 //D10 0x0400
	//		= 0x0200 //D9 Reserved
	RESET   uint16 = 0x0100 //D8 0x0100
	SLEEP1  uint16 = 0x0080 //D7 0x0080
	SLEEP12 uint16 = 0x0040 //D6 0x0040
	OPBITEN uint16 = 0x0020 //D5 0x0020
	//		= 0x0010 //D4 Reserved
	DIV2 uint16 = 0x0008 //D3 0x0008
	//		= 0x0004 //D2 Reserved
	MODE uint16 = 0x0002 //D1 0x0002
	//		= 0x0001 //D0 Reserved

	ADR_PHASE1 uint16 = 0xE000 //D15/D14/D13 0xE000
	ADR_PHASE0 uint16 = 0xC000 //D15/D14/D13 0xC000
	ADR_FREQ1  uint16 = 0x8000 //D15/D14 0x8000
	ADR_FREQ0  uint16 = 0x4000 //D15/D14 0x4000

	SLEEPNONE   uint16 = 0x0000 //SleepNone No power down
	SLEEPDAC    uint16 = 0x0040 //SleepDac DAC off
	SLEEPINTCLK uint16 = 0x0080 //SleepIntClk Internal clock disabled

	BITS14H uint32 = 0x0FFFC000 //High 14 MSBs 0xFFFC000
	BITS14L uint32 = 0x00003FFF //Low 14 LSBs 0x3FFF
	BITS12  uint16 = 0x0FFF     //12 LSBs 0x0FFF
)

type Mode uint8

const (
	MODE_MASK uint16 = 0x002A //D5/D3/D1 Mask 0x002A
)
const (
	MODE_SINE = 0x0000 //D5/D3/D1 Sinewave 0x0000
	MODE_TRI  = 0x0002 //D5/D3/D1 Triangle 0x0008
	MODE_MSB2 = 0x0020 //D5/D3/D1 MSB/2 0x0020
	MODE_MSB  = 0x0022 //D5/D3/D1 MSB 0x0020
)

type register struct {
	value uint16
}

func (r *register) replaceBits(value uint16, mask uint16) {
	r.value = (r.value & ^mask) | (value & mask)
}
