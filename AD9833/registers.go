package ad9833

const (
	PHASE1  = 0xE000 //D15/D14/D13 0xE000
	PHASE0  = 0xC000 //D15/D14/D13 0xC000
	FREQ1   = 0x8000 //D15/D14 0x8000
	FREQ0   = 0x4000 //D15/D14 0x4000
	B28     = 0x2000 //D13 0x2000
	HLB     = 0x1000 //D12 0x1000
	FSELECT = 0x0800 //D11 0x0800
	PSELECT = 0x0400 //D10 0x0400
	//		= 0x0200 //D9 Reserved
	RESET   = 0x0100 //D8 0x0100
	SLEEP1  = 0x0080 //D7 0x0080
	SLEEP12 = 0x0040 //D6 0x0040
	OPBITEN = 0x0020 //D5 0x0020
	//		= 0x0010 //D4 Reserved
	DIV2 = 0x0008 //D3 0x0008
	//		= 0x0004 //D2 Reserved
	MODE = 0x0002 //D1 0x0002
	//		= 0x0001 //D0 Reserved

	BITS14H = 0xFFFC000 //High 14 MSBs 0xFFFC000
	BITS14L = 0x3FFF    //Low 14 LSBs 0x3FFF
	BITS12  = 0x0FFF    //12 LSBs 0x0FFF
)
