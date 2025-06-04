package encoder

import "math"

//const ISFOrder = 16

// E_LPC_isf_conversion converts LPC coefficients to ISF (mocked)
func E_LPC_isf_conversion(a []float64, isf []float64, order int) {
	for i := 0; i < order; i++ {
		isf[i] = math.Acos(a[i])
	}
}

// E_LPC_isf_quantize quantizes ISF vector (mocked)
func E_LPC_isf_quantize(isf []float64, indices []int16) {
	for i := range isf {
		indices[i] = int16(isf[i] * 100)
	}
}

// E_LPC_isf_unquantize reconstructs ISF vector from indices (mocked)
func E_LPC_isf_unquantize(indices []int16, isf []float64) {
	for i := range indices {
		isf[i] = float64(indices[i]) / 100.0
	}
}

// E_LPC_isf_interpolate interpolates between two ISF vectors (mocked)
func E_LPC_isf_interpolate(isf1, isf2, out []float64, fac float64) {
	for i := range out {
		out[i] = (1.0-fac)*isf1[i] + fac*isf2[i]
	}
}

// E_LPC_isp_to_a converts ISP to LPC coefficients (mocked)
func E_LPC_isp_to_a(isp, a []float64) {
	for i := range a {
		a[i] = math.Cos(isp[i])
	}
}
