package decoder

import "math"

const (
	LPCOrder   = 16
	ISFGap     = 128
	ISFMeanVal = 200 // Example default mean offset for demonstration
)

var (
	D_ROM_mean_isf_noise = [LPCOrder]int16{
		200, 200, 200, 200, 200, 200, 200, 200,
		200, 200, 200, 200, 200, 200, 200, 200,
	}
)

func isfReorder(isf []int16, minDist int16, n int) {
	isfMin := int32(minDist)
	for i := 0; i < n-1; i++ {
		if int32(isf[i]) < isfMin {
			isf[i] = int16(isfMin)
		}
		isfMin = int32(isf[i]) + int32(minDist)
	}
}

func D_LPC_isf_noise_d(indice []int16, isfQ []int16) {
	// Example static decoding using fake codebooks (as real ones are extensive)
	for i := 0; i < 2; i++ {
		isfQ[i] = 50*indice[0] + int16(i)
	}
	for i := 0; i < 3; i++ {
		isfQ[i+2] = 60*indice[1] + int16(i)
	}
	for i := 0; i < 3; i++ {
		isfQ[i+5] = 70*indice[2] + int16(i)
	}
	for i := 0; i < 4; i++ {
		isfQ[i+8] = 80*indice[3] + int16(i)
	}
	for i := 0; i < 4; i++ {
		isfQ[i+12] = 90*indice[4] + int16(i)
	}
	for i := 0; i < LPCOrder; i++ {
		isfQ[i] += D_ROM_mean_isf_noise[i]
	}
	isfReorder(isfQ, ISFGap, LPCOrder)
}

// Converts ISF to ISP using linear interpolation with cosine lookup
func D_LPC_isf_isp_conversion(isf []int16, isp []int16, m int) {
	cosTable := make([]int16, 129) // simulate ROM cosine table [0..128]
	for i := range cosTable {
		cosTable[i] = int16(math.Cos(float64(i)*math.Pi/256) * 32767)
	}

	for i := 0; i < m-1; i++ {
		isp[i] = isf[i]
	}
	isp[m-1] = isf[m-1] << 1

	for i := 0; i < m; i++ {
		ind := isp[i] >> 7
		offset := isp[i] & 0x007F
		delta := int32(cosTable[ind+1]-cosTable[ind]) * int32(offset)
		isp[i] = int16(int32(cosTable[ind]) + (delta >> 7))
	}
}

// Convert ISP to LPC coefficients (simplified)
func D_LPC_isp_a_conversion(isp []int16, a []int16, adaptive bool, m int16) {
	nc := int(m / 2)
	f1 := make([]int32, nc+1)
	f2 := make([]int32, nc+1)

	ispPolGet(isp[:], f1, nc, false)
	ispPolGet(isp[1:], f2, nc-1, false)

	// Multiply F2 by (1 - z^-2)
	for i := nc - 1; i > 1; i-- {
		f2[i] -= f2[i-2]
	}

	// Apply symmetry/asymmetry combination
	a[0] = 4096 // Q12 (1.0)
	for i := 1; i < nc; i++ {
		t1 := (f1[i] + f2[i]) >> 1
		t2 := (f1[i] - f2[i]) >> 1
		a[i] = int16((t1 + (1 << 11)) >> 12)
		a[m-int16(i)] = int16((t2 + (1 << 11)) >> 12)
	}
	// Handle center tap
	a[nc] = int16(((f1[nc] + (1 << 11)) >> 12))
	a[m] = isp[m-1] >> 3 // Q15 to Q12
}

func ispPolGet(isp []int16, f []int32, n int, k16 bool) {
	s1 := int32(1 << 23)
	s2 := int32(1 << 9)
	if k16 {
		s1 >>= 2
		s2 >>= 2
	}
	f[0] = s1
	f[1] = int32(-isp[0]) * s2
	for i := 2; i <= n; i++ {
		f[i] = f[i-2]
		for j := i - 1; j >= 1; j-- {
			tmp := (f[j-1] * int32(isp[2*(i-1)])) >> 15
			f[j] += f[j-2] - tmp*2
		}
		f[1] -= int32(isp[2*(i-1)]) * s2
	}
}
