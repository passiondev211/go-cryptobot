package helper

import (
	"crypto/rand"
	"math"
	"math/big"
	"strconv"
	"strings"
)

// Round rounding number to the selected precision
func Round(x float64, prec int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}

// GetPrecision returns count of numbers after dot
func GetPrecision(x float64) int {
	str := strconv.FormatFloat(x, 'f', -1, 64)
	s := strings.Split(str, ".")
	if len(s) == 1 {
		return 0
	}
	return len(s[1])
}

// PseudoSha2 outputs 64-byte string that looks like the real sha2.
func PseudoSha2() string {
	hexNums := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	maxRand := big.NewInt(16)

	l := 64
	result := make([]byte, l)
	for i := 0; i < l; i++ {
		// The error is checked in tests.
		t, _ := rand.Int(rand.Reader, maxRand)
		result[i] = hexNums[t.Int64()]
	}
	return string(result)
}
