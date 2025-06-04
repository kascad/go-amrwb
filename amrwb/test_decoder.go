package amrwb

import (
	"testing"
	"transcoder/amrwb/common"
	"transcoder/amrwb/decoder"
)

func TestE_MAIN_init_exit(t *testing.T) {
	st := decoder.D_MAIN_init()
	if st == nil {
		t.Fatal("decoder init returned nil")
	}
	decoder.D_MAIN_close(&st)
	if st != nil {
		t.Error("decoder exit did not set pointer to nil")
	}
}

func TestE_MAIN_decode(t *testing.T) {
	st := decoder.D_MAIN_init()
	defer decoder.D_MAIN_close(&st)

	var serial [common.NB_SERIAL_MAX]int16
	var output [common.L_FRAME16k]int16
	serial[0] = 0x3c // mock SID / mode frame
	decoder.D_MAIN_decode(7, serial[:], output[:], st, 0)

	for i, s := range output {
		if s != 0 {
			t.Errorf("expected zero at %d got %d", i, s)
			break
		}
	}
}
