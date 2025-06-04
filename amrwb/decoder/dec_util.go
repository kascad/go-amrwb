package decoder

import "math"

// saturate limits a 32-bit integer to 16-bit range
func saturate(x int32) int16 {
	if x > math.MaxInt16 {
		return math.MaxInt16
	} else if x < math.MinInt16 {
		return math.MinInt16
	}
	return int16(x)
}

// shiftRight rounds and shifts with saturation
func shiftRight(x int32, shift int) int16 {
	if shift > 0 {
		x += int32(1 << (shift - 1))
		x >>= shift
	}
	return saturate(x)
}

// dotProduct computes scalar product of vectors with saturation
func dotProduct(x, y []int16, n int) int32 {
	var sum int32
	for i := 0; i < n; i++ {
		sum += int32(x[i]) * int32(y[i])
	}
	return sum
}

// copyBuffer copies one int16 slice to another
func copyBuffer(dst, src []int16) {
	copy(dst, src)
}

// setBuffer sets all elements of an int16 slice to a value
func setBuffer(dst []int16, val int16) {
	for i := range dst {
		dst[i] = val
	}
}
