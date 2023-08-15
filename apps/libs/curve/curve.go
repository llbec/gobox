package main

import "C"
import "gobox/src/curve"

//export FitCurve
func FitCurve(pow int, samples [][2]float64) []float64 {
	crv := curve.NewCurve(pow)

	for i := 0; i < len(samples); i++ {
		crv.AddSample(samples[i][0], samples[i][1])
	}

	if err := crv.Fitting(); err != nil {
		return []float64{}
	}

	return crv.Thetas()
}
