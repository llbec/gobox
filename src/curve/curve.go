package curve

import (
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/gonum/mat"
)

type Curve struct {
	pow     int
	samples [][]float64
	thetas  []float64
}

func NewCurve(pow int) *Curve {
	return &Curve{
		pow: pow,
	}
}

func (crv *Curve) AddSample(x, y float64) {
	crv.samples = append(crv.samples, []float64{x, y})
	sort.Slice(crv.samples, func(i, j int) bool { return crv.samples[i][0] < crv.samples[j][0] })
}

func (crv *Curve) Fitting() error {
	m := len(crv.samples)
	Xv := []float64{}
	Y := []float64{}

	var matTheta, mat1, mat2 mat.Dense

	for i := 0; i < m; i++ {
		for j := 0; j < crv.pow+1; j++ {
			Xv = append(Xv, math.Pow(crv.samples[i][0], float64(j)))
		}
		Y = append(Y, crv.samples[i][1])
	}
	matXv := mat.NewDense(m, crv.pow+1, Xv)
	matY := mat.NewDense(m, 1, Y)

	mat1.Mul(matXv.T(), matXv)
	if err := mat1.Inverse(&mat1); err != nil {
		return err
	}
	mat2.Mul(&mat1, matXv.T())
	matTheta.Mul(&mat2, matY)
	row, _ := matTheta.Caps()
	for i := 0; i < row; i++ {
		crv.thetas = append(crv.thetas, matTheta.At(i, 0))
	}
	return nil
}

func (crv *Curve) Verify() error {
	for i := 0; i < len(crv.samples); i++ {
		var sum float64
		for j := 0; j < crv.pow+1; j++ {
			sum += math.Pow(crv.samples[i][0], float64(j)) * crv.thetas[j]
		}
		if math.Abs(sum-crv.samples[i][1]) > 0.000001 {
			return fmt.Errorf("sample f(%v) = %v, expect %v", crv.samples[i][0], sum, crv.samples[i][1])
		}
	}

	return nil
}
