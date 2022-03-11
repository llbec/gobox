package benford

import "math"

//LogaN caculate Logarithm
func LogaN(a, N float64) float64 {
	return math.Log10(N) / math.Log10(a)
}

//IsEqual for float
func IsEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.000000000000001
}

//Benford return P(n) = log(b)((n+1)/n)
func Benford(n, b uint) float64 {
	n64 := math.Float64frombits(uint64(n))
	b64 := math.Float64frombits(uint64(b))
	return LogaN(b64, (n64+1)/n64)
}
