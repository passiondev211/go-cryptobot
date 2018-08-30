package helper

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPrecision(t *testing.T) {
	// Arrange.
	expected := make(map[float64]int)

	expected[2.00000001] = 8
	expected[2.0000001] = 7
	expected[2.000001] = 6
	expected[2.000000] = 0
	expected[0.000000] = 0
	expected[0.002000] = 3
	expected[0.002000] = 3
	expected[1] = 0
	expected[10] = 0

	for x, expPrec := range expected {
		testName := fmt.Sprintf("test get prec for %v", x)
		t.Run(testName, func(t *testing.T) {
			// Act.
			actual := getPrecision(x)

			// Assert.
			require.Equal(t, expPrec, actual)
		})
	}
}

func TestRound(t *testing.T) {
	// Arrange.
	expected := []struct {
		arrange  float64
		prec     int
		expected float64
	}{
		{1.2385, 3, 1.239},
		{1.2385, 2, 1.24},
		{1.2385, 1, 1.2},
		{1.7385, 0, 2},
		{237845, -1, 237850},
		{237845, -2, 237800},
		{237845, -3, 238000},
		{200000.001, 10, 200000.001},
	}

	for _, e := range expected {
		testName := fmt.Sprintf("test rounding for %f", e.arrange)
		t.Run(testName, func(t *testing.T) {
			// Act.
			actual := round(e.arrange, e.prec)
			// Assert.
			require.Equal(t, e.expected, actual)
		})
	}
}
