package encoder

import (
	"math"
)

func E_UTIL_saturate(inp int32) int16 {
	if inp > math.MaxInt16 {
		return math.MaxInt16
	} else if inp < math.MinInt16 {
		return math.MinInt16
	}
	return int16(inp)
}

func E_UTIL_saturate_31(inp int32) int32 {
	if inp > 0x3FFFFFFF {
		return 0x3FFFFFFF
	} else if inp < -0x40000000 {
		return -0x40000000
	}
	return inp
}

func E_UTIL_random(seed *int16) int16 {
	*seed = (*seed*31821 + 13849) & 0xFFFF
	return *seed
}

func E_UTIL_mpy_32_16(hi, lo, n int16) int32 {
	hiPart := int32(hi) * int32(n) << 1
	loPart := ((int32(lo) * int32(n)) >> 15) << 1
	return hiPart + loPart
}

func E_UTIL_pow2(exponent, fraction int16) int32 {
	return int32(1 << exponent)
}

func E_UTIL_log2_32(L_x int32, exponent, fraction *int16) {
	*exponent = 15
	*fraction = int16(L_x >> 16)
}

func E_UTIL_normalised_inverse_sqrt(frac *int32, exp *int16) {
	*frac = 0x40000000 // placeholder 1.0
	*exp = 0
}

func E_UTIL_dot_product12(x, y []int16, lg int32, exp *int32) int32 {
	var s int32
	for i := int32(0); i < lg; i++ {
		s += int32(x[i]) * int32(y[i])
	}
	*exp = 0
	return s
}

func E_UTIL_norm_s(var1 int16) int16 {
	if var1 == 0 {
		return 0
	}
	cnt := int16(0)
	val := int32(var1)
	if val < 0 {
		val = ^val
	}
	for val < 0x4000 {
		cnt++
		val <<= 1
	}
	return cnt
}

func E_UTIL_norm_l(L_var1 int32) int16 {
	if L_var1 == 0 {
		return 0
	}
	cnt := int16(0)
	val := L_var1
	if val < 0 {
		val = ^val
	}
	for val < 0x40000000 {
		cnt++
		val <<= 1
	}
	return cnt
}

func E_UTIL_l_extract(L_32 int32, hi, lo *int16) {
	*hi = int16(L_32 >> 16)
	*lo = int16((L_32 - int32(*hi)<<16) >> 1)
}
