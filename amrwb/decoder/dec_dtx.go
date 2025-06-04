package decoder

const (
	LFrame             = 256
	M                  = 16
	DDTXHistSize       = 8
	RandomInitSeed     = 21845
	DDTXHangConst      = 7
	DDTXMaxEmptyThresh = 50
	ISFGap             = 128
	GainFactor         = 75
	ISFFactorLow       = 256
	ISFFactorStep      = 2
	ISFDithGap         = 448
)

type DTXState struct {
	ISFHist       [M * DDTXHistSize]int16
	ISF           [M]int16
	ISFPrev       [M]int16
	LogEnergyHist [DDTXHistSize]int16
	TrueSIDPeriod int16
	LogEnergy     int16
	LogEnergyPrev int16
	CNGSeed       int16
	HistPtr       int16
	DitherSeed    int16
	CNDither      int16
	SinceLastSID  int16
	ElapsedCount  uint8
	GlobalState   uint8
	DataUpdated   uint8
	HangoverCount uint8
	SIDFrame      uint8
	ValidData     uint8
	HangoverAdded uint8
	VADHist       int16
}

func NewDTXState(initISF []int16) *DTXState {
	st := &DTXState{
		TrueSIDPeriod: 8192, // 0.25 in Q15
		LogEnergy:     3500,
		LogEnergyPrev: 3500,
		CNGSeed:       RandomInitSeed,
		HistPtr:       0,
		DitherSeed:    RandomInitSeed,
		GlobalState:   0,
		ElapsedCount:  127,
		HangoverCount: DDTXHangConst,
	}

	copy(st.ISF[:], initISF)
	copy(st.ISFPrev[:], initISF)

	for i := 0; i < DDTXHistSize; i++ {
		copy(st.ISFHist[i*M:], initISF)
		st.LogEnergyHist[i] = 3500 / 8
	}

	return st
}

func cnDithering(isf []int16, logEnergy *int32, seed *int16) {
	rand1 := random(seed) >> 1
	rand2 := random(seed) >> 1
	rand := rand1 + rand2
	*logEnergy += int32(rand * GainFactor * 2)
	if *logEnergy < 0 {
		*logEnergy = 0
	}

	ditherFac := ISFFactorLow
	rand = (random(seed) >> 1) + (random(seed) >> 1)
	temp := int32(isf[0]) + ((int32(rand)*int32(ditherFac) + 0x4000) >> 15)
	if temp < ISFGap {
		isf[0] = ISFGap
	} else {
		isf[0] = int16(temp)
	}

	for i := 1; i < M-1; i++ {
		ditherFac += ISFFactorStep
		rand = (random(seed) >> 1) + (random(seed) >> 1)
		temp = int32(isf[i]) + ((int32(rand)*int32(ditherFac) + 0x4000) >> 15)
		temp1 := temp - int32(isf[i-1])
		if temp1 < ISFDithGap {
			isf[i] = isf[i-1] + ISFDithGap
		} else {
			isf[i] = int16(temp)
		}
	}
}

func random(seed *int16) int16 {
	x := int32(*seed)
	x = (x*31821 + 13849) & 0xFFFF
	*seed = int16(x)
	return int16(x)
}
