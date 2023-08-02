package curve_test

import (
	"fmt"
	"gobox/src/curve"
	"testing"
)

func TestSort(t *testing.T) {
	crv := curve.NewCurve(2)
	crv.AddSample(0.1, 1)
	crv.AddSample(0.4, 4)
	crv.AddSample(0.2, 2)
	crv.AddSample(0.3, 3)
	crv.Fitting()
	fmt.Println(crv)
}
