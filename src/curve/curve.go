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

	var matTheta, mat1, m1Inv, mat2 mat.Dense

	for i := 0; i < m; i++ {
		for j := 0; j < crv.pow+1; j++ {
			Xv = append(Xv, math.Pow(crv.samples[i][0], float64(j)))
		}
		Y = append(Y, crv.samples[i][1])
	}
	matXv := mat.NewDense(m, crv.pow+1, Xv)
	matY := mat.NewDense(m, 1, Y)

	mat1.Mul(matXv.T(), matXv)
	if err := m1Inv.Inverse(&mat1); err != nil {
		pinv, err := Pseudoinverse(&mat1)
		if err != nil {
			return err
		}
		m1Inv.CloneFrom(pinv)
	}
	mat2.Mul(&m1Inv, matXv.T())
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

func (crv *Curve) Theta(i int) float64 {
	if i < len(crv.thetas) {
		return crv.thetas[i]
	}
	return 0
}

func (crv *Curve) Thetas() []float64 {
	return crv.thetas
}

func Pseudoinverse(a mat.Matrix) (*mat.Dense, error) {
	var svd mat.SVD
	if ok := svd.Factorize(a, mat.SVDFull); !ok {
		return nil, fmt.Errorf("failed to factorize a")
	}
	r, c := a.Dims()
	svdS := svd.Values(nil)
	sData := []float64{}
	for i := 0; i < r*c; i++ {
		if i/c == i%c {
			sData = append(sData, 1/svdS[i/c])
			continue
		}
		sData = append(sData, 0)
	}
	Spinv := mat.NewDense(r, c, sData).T()
	U := mat.NewDense(r, r, nil)
	V := mat.NewDense(c, c, nil)

	var m1, m2 mat.Dense
	svd.UTo(U)
	svd.VTo(V)

	m1.Mul(V, Spinv)
	m2.Mul(&m1, U.T())

	return &m2, nil
}
