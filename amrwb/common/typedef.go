package common

// Go equivalents of fixed-width types used in AMR-WB C implementation
// Basic type definitions adapted from typedef.h in C

// Word8 and Word16 are fixed-size integer types

type Word16 = int16

type Word32 = int32

type UWord8 = uint8

type UWord16 = uint16

type UWord32 = uint32

type Word8 = uint8

// Boolean values
const (
	FALSE = 0
	TRUE  = 1
)

// Frame size constant (matching L_FRAME16k)
const L_FRAME16k = 320

// Max serial output size
const NB_SERIAL_MAX = 61
