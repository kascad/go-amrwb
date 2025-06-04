package decoder

import "math"

const (
	LSubfr    = 64
	LTPhist   = 5
	OnePer3   = 10923
	OnePerLTP = 6554
	Upsamp    = 4
	LInterp2  = 16
	MaxQ14    = 15565 // 0.95 in Q14
)

func D_GAIN_init(mem []int16) {
	mem[0] = -14336
	mem[1] = -14336
	mem[2] = -14336
	mem[3] = -14336
	for i := 4; i < 22; i++ {
		mem[i] = 0
	}
	mem[22] = 21845
}

func median5(x []int16) int16 {
	x1, x2, x3, x4, x5 := x[-2], x[-1], x[0], x[1], x[2]

	swap := func(a, b *int16) {
		if *a > *b {
			*a, *b = *b, *a
		}
	}

	swap(&x1, &x2)
	swap(&x1, &x3)
	swap(&x1, &x4)
	swap(&x1, &x5)
	swap(&x2, &x3)
	swap(&x2, &x4)
	swap(&x2, &x5)
	swap(&x3, &x4)
	swap(&x3, &x5)

	return x3
}

func D_GAIN_decode(index, nbits int16, code []int16, gainPit *int16, gainCode *int32, bfi, prevBfi, state, unusableFrame, vadHist int16, mem []int16) {
	pastQuaEn := mem[0:4]
	pastGainPit := &mem[4]
	pastGainCode := &mem[5]
	prevGC := &mem[6]
	pbuf := mem[7:12]
	gbuf := mem[12:17]
	pbuf2 := mem[17:22]

	// Placeholder energy computation
	gcodeInov := int16(16384) // Typically would be derived from code energy

	if bfi != 0 {
		tmp := median5(pbuf[2:])
		*pastGainPit = tmp
		if *pastGainPit > MaxQ14 {
			*pastGainPit = MaxQ14
		}
		*gainPit = *pastGainPit

		tmp = median5(gbuf[2:])
		if vadHist > 2 {
			*pastGainCode = tmp
		} else {
			*pastGainCode = tmp
		}

		Ltmp := int32(pastQuaEn[0]+pastQuaEn[1]+pastQuaEn[2]+pastQuaEn[3]) >> 2
		quaEner := Ltmp - 3072
		if quaEner < -14336 {
			quaEner = -14336
		}

		copy(pastQuaEn[1:], pastQuaEn[0:3])
		pastQuaEn[0] = int16(quaEner)

		copy(gbuf[0:4], gbuf[1:5])
		gbuf[4] = *pastGainCode
		copy(pbuf[0:4], pbuf[1:5])
		pbuf[4] = *pastGainPit

		*gainCode = int32(*pastGainCode) * int32(gcodeInov) * 2
		return
	}

	// Placeholder quantized gain tables
	gainTbl := []int16{
		16384, // pitch gain in Q14
		2048,  // code gain in Q11
	}

	*gainPit = gainTbl[0]
	gCode := gainTbl[1]
	gcode0 := int32(1 << 14) // typically computed using log predictors

	Ltmp := int32(gCode) * gcode0
	*gainCode = Ltmp >> 3

	*pastGainPit = *gainPit
	*pastGainCode = int16(*gainCode >> 13)
	if *pastGainCode > 32767 {
		*pastGainCode = 32767
	}
	*prevGC = *pastGainCode

	copy(gbuf[0:4], gbuf[1:5])
	gbuf[4] = *pastGainCode
	copy(pbuf[0:4], pbuf[1:5])
	pbuf[4] = *pastGainPit
	copy(pbuf2[0:4], pbuf2[1:5])
	pbuf2[4] = *pastGainPit

	Ltmp = int32(*gainCode)
	if Ltmp < 0x0FFFFFFF {
		*gainCode = Ltmp << 3
	} else {
		*gainCode = 0x7FFFFFFF
	}
}
