package encoder

// DTX encoder state

const (
	DTXHistorySize = 8
	ISFOrder       = 16
)

type DTXEncState struct {
	isfHist    [DTXHistorySize][ISFOrder]int16
	energyHist [DTXHistorySize]int16
	histPtr    int
}

func E_DTX_init() *DTXEncState {
	return &DTXEncState{}
}

func E_DTX_exit(state **DTXEncState) {
	*state = nil
}

func E_DTX_reset(st *DTXEncState) {
	for i := range st.isfHist {
		for j := range st.isfHist[i] {
			st.isfHist[i][j] = 0
		}
		st.energyHist[i] = 0
	}
	st.histPtr = 0
}

// E_DTX_activity_update updates comfort noise estimation
func E_DTX_activity_update(st *DTXEncState, isf []int16, energy int16) {
	copy(st.isfHist[st.histPtr][:], isf)
	st.energyHist[st.histPtr] = energy
	st.histPtr = (st.histPtr + 1) % DTXHistorySize
}

// E_DTX_encode encodes comfort noise frame (placeholder logic)
func E_DTX_encode(st *DTXEncState, prms []int16) {
	// Use last ISF and average energy to produce SID frame
	avgEnergy := int16(0)
	for _, e := range st.energyHist {
		avgEnergy += e / DTXHistorySize
	}
	for i := range prms {
		prms[i] = avgEnergy // fake parameter encoding
	}
}
