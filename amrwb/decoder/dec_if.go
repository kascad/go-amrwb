package decoder

import "bytes"

const (
	LFrame16k   = 320
	NBSerialMax = 61
	EHFMask     = 0x0008

	Mode7k    = 0
	Mode9k    = 1
	Mode12k   = 2
	Mode14k   = 3
	Mode16k   = 4
	Mode18k   = 5
	Mode20k   = 6
	Mode23k   = 7
	Mode24k   = 8
	MRDTX     = 9
	MRNOData  = 15
	LOSTFrame = 14
)

type IFState struct {
	ResetFlagOld int16
	PrevFT       int16
	PrevMode     int16
	Decoder      any
}

func IFInit() *IFState {
	return &IFState{
		ResetFlagOld: 0,
		PrevFT:       0,
		PrevMode:     0,
		Decoder:      nil, // Replace with actual decoder state initialization
	}
}

func IFExit(st *IFState) {
	st.Decoder = nil
}

func IsHomingFrame(input []int16, mode int) bool {
	dummyRef := make([]int16, len(input))
	// Normally compare with decoder homing frame pattern
	return bytes.Equal(toBytes(input), toBytes(dummyRef))
}

func toBytes(input []int16) []byte {
	buf := new(bytes.Buffer)
	for _, v := range input {
		buf.WriteByte(byte(v >> 8))
		buf.WriteByte(byte(v))
	}
	return buf.Bytes()
}

func IFDecode(state *IFState, bits []uint8, synth []int16, bfi int32) {
	// Placeholder decoding logic: fill synth buffer with dummy samples
	for i := 0; i < LFrame16k; i++ {
		synth[i] = int16(i % 32768)
	}
	// Real decoder should invoke frame unpacking, bit parsing, and synthesis here
}
