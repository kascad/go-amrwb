package encoder

// ROM tables used for interface functions (simplified)

var IF_ROM_modeBitrate = [9]int32{
	6600, 8850, 12650, 14250, 15850, 18250, 19850, 23050, 23850,
}

var IF_ROM_modeFrameSize = [9]int32{
	17, 23, 32, 36, 40, 46, 50, 58, 60,
}

var IF_ROM_frameType = [16]string{
	"FT_SPEECH_GOOD",
	"FT_SPEECH_DEGRADED",
	"FT_ONSET",
	"FT_SPEECH_BAD",
	"FT_SID_FIRST",
	"FT_SID_UPDATE",
	"FT_SID_BAD",
	"FT_NO_DATA",
	"FT_SPEECH_LOST",
	"FT_SID_LOST",
	"FT_EMPTY_FRAME",
	"FT_RESERVED",
	"FT_RESERVED",
	"FT_RESERVED",
	"FT_RESERVED",
	"FT_RESERVED",
}
