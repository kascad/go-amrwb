package amrwb

import (
	"fmt"
	"transcoder/amrwb/common"
	"transcoder/amrwb/decoder"
	"transcoder/amrwb/encoder"
)

// AMRWBCodec is a unified structure for both encoding and decoding AMR-WB

type AMRWBCodec struct {
	EncoderState *encoder.EncoderMainState
	DecoderState *decoder.DecoderMainState
	Mode         int16
	DTX          bool
}

func NewAMRWBCodec(mode int16, dtx bool) *AMRWBCodec {
	enc := encoder.E_MAIN_init(mode, dtx)
	dec := decoder.D_MAIN_init()
	return &AMRWBCodec{
		EncoderState: enc,
		DecoderState: dec,
		Mode:         mode,
		DTX:          dtx,
	}
}

func (c *AMRWBCodec) Close() {
	encoder.E_MAIN_exit(&c.EncoderState)
	decoder.D_MAIN_exit(&c.DecoderState)
}

func (c *AMRWBCodec) Encode(frame []int16) ([]byte, error) {
	if len(frame) != common.L_FRAME16k {
		return nil, fmt.Errorf("invalid frame size: %d", len(frame))
	}
	var serial [common.NB_SERIAL_MAX]int16
	encoder.E_MAIN_encode(c.EncoderState, frame, serial[:])
	out := make([]byte, len(serial))
	for i, val := range serial {
		out[i] = byte(val & 0xFF)
	}
	return out, nil
}

func (c *AMRWBCodec) Decode(serial []byte) ([]int16, error) {
	if len(serial) > common.NB_SERIAL_MAX {
		return nil, fmt.Errorf("serial frame too long: %d", len(serial))
	}
	var output [common.L_FRAME16k]int16
	var input [common.NB_SERIAL_MAX]int16
	for i := range serial {
		input[i] = int16(serial[i])
	}
	decoder.D_MAIN_decode(c.DecoderState, input[:], output[:])
	return output[:], nil
}
