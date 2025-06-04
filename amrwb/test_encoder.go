package amrwb

import (
	"testing"
	"transcoder/amrwb/common"
	"transcoder/amrwb/encoder"
)

func TestE_MAIN_init_reset_exit(t *testing.T) {
	st := encoder.E_MAIN_init(7, false)
	if st == nil {
		t.Fatal("init returned nil")
	}
	if st.mode != 7 || st.allowDTX {
		t.Errorf("unexpected state %+v", st)
	}
	st.frameCount = 5
	encoder.E_MAIN_reset(st)
	if st.frameCount != 0 {
		t.Errorf("frameCount not reset: %d", st.frameCount)
	}
	encoder.E_MAIN_exit(&st)
	if st != nil {
		t.Error("exit did not nil pointer")
	}
}

func TestE_MAIN_encode(t *testing.T) {
	st := encoder.E_MAIN_init(7, false)
	defer encoder.E_MAIN_exit(&st)

	var input [common.L_FRAME16k]int16
	var serial [common.NB_SERIAL_MAX]int16
	encoder.E_MAIN_encode(st, input[:], serial[:])

	if st.frameCount != 1 {
		t.Errorf("frameCount=%d", st.frameCount)
	}
	for i, v := range serial {
		exp := int16(i * 2)
		if v != exp {
			t.Errorf("serial[%d]=%d want %d", i, v, exp)
			break
		}
	}
}
