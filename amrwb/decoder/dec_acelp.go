// Package amrwb provides a Go translation of selected C files from the 3GPP AMR-WB floating-point codec.
// The goal is to faithfully reimplement C logic in idiomatic Go, while retaining signal processing semantics.
package decoder

const (
	L_SUBFR    = 64 // Subframe size
	PRED_ORDER = 4
	MEAN_ENER  = 30
)

// D_ACELP_add_pulse adds pulse positions into the fixed codebook
func D_ACELP_add_pulse(pos []int32, nb_pulse, track int32, code []int16) {
	for k := int32(0); k < nb_pulse; k++ {
		i := ((pos[k] & 15) << 2) + track
		if (pos[k] & 16) == 0 {
			code[i] += 512
		} else {
			code[i] -= 512
		}
	}
}

func D_ACELP_decode_1p_N1(index, N, offset int32, pos []int32) {
	mask := (1 << N) - 1
	pos1 := (index & int32(mask)) + offset
	if (index>>N)&1 == 1 {
		pos1 += 16
	}
	pos[0] = pos1
}

func D_ACELP_decode_2p_2N1(index, N, offset int32, pos []int32) {
	mask := (1 << N) - 1
	pos1 := ((index >> N) & int32(mask)) + offset
	pos2 := (index & int32(mask)) + offset
	sign := (index >> (2 * N)) & 1
	if pos2 < pos1 {
		if sign == 1 {
			pos1 += 16
		} else {
			pos2 += 16
		}
	} else {
		if sign == 1 {
			pos1 += 16
			pos2 += 16
		}
	}
	pos[0], pos[1] = pos1, pos2
}

func D_ACELP_decode_3p_3N1(index, N, offset int32, pos []int32) {
	mask := (1 << ((2 * N) - 1)) - 1
	idx := index & int32(mask)
	j := offset
	if (index>>((2*N)-1))&1 == 1 {
		j += 1 << (N - 1)
	}
	D_ACELP_decode_2p_2N1(idx, N-1, j, pos)
	mask = (1 << (N + 1)) - 1
	idx = (index >> (2 * N)) & int32(mask)
	D_ACELP_decode_1p_N1(idx, N, offset, pos[2:])
}

func D_ACELP_decode_4p_4N1(index, N, offset int32, pos []int32) {
	mask := (1 << ((2 * N) - 1)) - 1
	idx := index & int32(mask)
	j := offset
	if (index>>((2*N)-1))&1 == 1 {
		j += 1 << (N - 1)
	}
	D_ACELP_decode_2p_2N1(idx, N-1, j, pos)
	mask = (1 << ((2 * N) + 1)) - 1
	idx = (index >> (2 * N)) & int32(mask)
	D_ACELP_decode_2p_2N1(idx, N, offset, pos[2:])
}

func D_ACELP_decode_4p_4N(index, N, offset int32, pos []int32) {
	n_1 := N - 1
	j := offset + (1 << n_1)
	switch (index >> ((4 * N) - 2)) & 3 {
	case 0:
		if ((index >> ((4 * n_1) + 1)) & 1) == 0 {
			D_ACELP_decode_4p_4N1(index, n_1, offset, pos)
		} else {
			D_ACELP_decode_4p_4N1(index, n_1, j, pos)
		}
	case 1:
		D_ACELP_decode_1p_N1(index>>(3*n_1+1), n_1, offset, pos)
		D_ACELP_decode_3p_3N1(index, n_1, j, pos[1:])
	case 2:
		D_ACELP_decode_2p_2N1(index>>(2*n_1+1), n_1, offset, pos)
		D_ACELP_decode_2p_2N1(index, n_1, j, pos[2:])
	case 3:
		D_ACELP_decode_3p_3N1(index>>(n_1+1), n_1, offset, pos)
		D_ACELP_decode_1p_N1(index, n_1, j, pos[3:])
	}
}

func D_ACELP_decode_5p_5N(index, N, offset int32, pos []int32) {
	n_1 := N - 1
	j := offset + (1 << n_1)
	idx := index >> ((2 * N) + 1)
	if (index>>((5*N)-1))&1 == 0 {
		D_ACELP_decode_3p_3N1(idx, n_1, offset, pos)
		D_ACELP_decode_2p_2N1(index, N, offset, pos[3:])
	} else {
		D_ACELP_decode_3p_3N1(idx, n_1, j, pos)
		D_ACELP_decode_2p_2N1(index, N, offset, pos[3:])
	}
}

func D_ACELP_decode_6p_6N_2(index, N, offset int32, pos []int32) {
	n_1 := N - 1
	j := offset + (1 << n_1)
	offA, offB := j, j
	if (index>>((6*N)-5))&1 == 0 {
		offA = offset
	} else {
		offB = offset
	}
	switch (index >> ((6 * N) - 4)) & 3 {
	case 0:
		D_ACELP_decode_5p_5N(index>>N, n_1, offA, pos)
		D_ACELP_decode_1p_N1(index, n_1, offA, pos[5:])
	case 1:
		D_ACELP_decode_1p_N1(index>>((5*N)+1), n_1, offA, pos)
		D_ACELP_decode_5p_5N(index, n_1, offB, pos[1:])
	case 2:
		D_ACELP_decode_2p_2N1(index>>((4*N)+1), n_1, offA, pos)
		D_ACELP_decode_4p_4N1(index, n_1, offB, pos[2:])
	case 3:
		D_ACELP_decode_4p_4N1(index>>((2*N)+1), n_1, offA, pos)
		D_ACELP_decode_2p_2N1(index, n_1, offB, pos[4:])
	}
}
