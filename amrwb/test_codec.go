package amrwb

import (
	"testing"
	"transcoder/amrwb/common"
)

func TestAMRWBCodecEncodeDecode(t *testing.T) {
	codec := NewAMRWBCodec(7, false)
	defer codec.Close()

	input := make([]int16, common.L_FRAME16k)
	for i := range input {
		input[i] = int16(i)
	}
	serial, err := codec.Encode(input)
	if err != nil {
		t.Fatalf("encode error: %v", err)
	}
	if len(serial) != common.NB_SERIAL_MAX {
		t.Errorf("unexpected serial length %d", len(serial))
	}

	output, err := codec.Decode(serial)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(output) != common.L_FRAME16k {
		t.Errorf("unexpected output length %d", len(output))
	}
	for _, v := range output {
		if v != 0 {
			t.Errorf("expected zero output, got %d", v)
			break
		}
	}
}
