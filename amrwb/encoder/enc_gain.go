package encoder

import "math"

func E_GAIN_clip_init(mem []float32) {
	for i := range mem {
		mem[i] = 0.0
	}
}

func E_GAIN_clip_test(mode int16, mem []float32) int32 {
	for _, val := range mem {
		if math.Abs(float64(val)) > 1.5 {
			return 1
		}
	}
	return 0
}

func E_GAIN_clip_isf_test(mode int16, isf []float32, mem []float32) {
	for i := range isf {
		mem[i] = isf[i] // simplified tracking
	}
}

func E_GAIN_clip_pit_test(mode int16, gainPit float32, mem []float32) {
	if gainPit > 1.0 {
		mem[0] = gainPit // simplified record of over-pitch gain
	}
}

func E_GAIN_lp_decim2(x []float32, l int32, mem *float32) {
	if len(x) < 2 {
		return
	}
	for i := int32(0); i < l-1; i += 2 {
		x[i/2] = 0.5 * (x[i] + x[i+1])
	}
}

func E_GAIN_olag_median(prevOLLag int32, oldOLLag [5]int32) int32 {
	temp := append(oldOLLag[:], prevOLLag)
	for i := 0; i < len(temp); i++ {
		for j := i + 1; j < len(temp); j++ {
			if temp[j] < temp[i] {
				temp[i], temp[j] = temp[j], temp[i]
			}
		}
	}
	return temp[len(temp)/2]
}

func E_GAIN_open_loop_search(wsp []float32, LMin, LMax, nFrame, L0 int32,
	gain, mem, hpOldWsp []float32, weightFlg uint8) int32 {
	return (LMin + LMax) / 2 // placeholder for midpoint
}

func E_GAIN_closed_loop_search(exc, xn, h []float32, t0Min, t0Max int32,
	pitFrac *int32, iSubfr, t0Fr2, t0Fr1 int32) int32 {
	*pitFrac = 0
	return t0Min // placeholder return value
}

func E_GAIN_adaptive_codebook_excitation(exc []int16, T0, frac, LSubfr int32) {
	for i := int32(0); i < LSubfr; i++ {
		exc[i] = int16((i + T0 + frac) % 32768) // placeholder pattern
	}
}

func E_GAIN_pitch_sharpening(x []int16, pitLag int16) {
	for i := range x {
		x[i] = x[i] + (x[i] / 2)
	}
}

func E_GAIN_f_pitch_sharpening(x []float32, pitLag int32) {
	for i := range x {
		x[i] *= 1.2
	}
}

func E_GAIN_voice_factor(exc []int16, QExc, gainPit int16,
	code []int16, gainCode int16) int32 {
	return int32(gainCode) * 10 // placeholder estimation
}
