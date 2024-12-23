package helper

import "math"

// Amount2Decimal return amount only 2 decima places.
func Amount2Decimal(s float64) float64 {
	return math.Floor(s*100) / 100
}
