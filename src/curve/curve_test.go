package curve_test

import (
	"fmt"
	"gobox/src/curve"
	"math"
	"math/rand"
	"testing"

	"gonum.org/v1/gonum/mat"
)

func TestCurve(t *testing.T) {
	crv := curve.NewCurve(2)
	crv.AddSample(29077, -15)
	crv.AddSample(30657, -11)
	crv.AddSample(35961, 30)
	crv.AddSample(39608, 65)
	crv.AddSample(41605, 85)
	crv.AddSample(42572, 95)
	fmt.Println(crv)
	if err := crv.Fitting(); err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(crv)

	if err := crv.Verify(); err != nil {
		t.Fatalf(err.Error())
	}
}

func TestDense(t *testing.T) {
	a := mat.NewDense(6, 2, nil)
	for i := 1; i < 7; i++ {
		a.SetRow(i-1, []float64{float64(i), rand.NormFloat64()})
	}
	t.Logf("\na^T =\n%v", mat.Formatted(a.T()))
	t.Logf("\na =\n%v", mat.Formatted(a))

	var b, c mat.Dense
	b.Mul(a.T(), a)
	t.Logf("\nb = a^T * a =\n%v", mat.Formatted(&b))
	err := b.Inverse(&b)
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("\nb^{-1} =\n%v", mat.Formatted(&b))

	c.Mul(&b, a.T())
	t.Logf("\nc = b^{-1} * a^T\n =\n%v", mat.Formatted(&b))
}

func TestSolve(t *testing.T) {
	x := []float64{29077, 30657, 35961, 39608, 41605, 42572}
	y := []float64{-15, -11, 30, 65, 85, 95}
	a := mat.NewDense(6, 3, nil)
	for i := 0; i < 6; i++ {
		a.SetRow(i, []float64{1, x[i], math.Pow(x[2], 2)})
	}
	b := mat.NewVecDense(6, y)
	var c mat.VecDense
	if err := c.SolveVec(a, b); err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("\n%v\n", c)
}

func TestSVD(t *testing.T) {
	s_x := []float64{29077, 30657, 35961, 39608, 41605, 42572}
	s_y := []float64{-15, -11, 30, 65, 85, 95}

	a := mat.NewDense(6, 3, nil)
	for i := 0; i < 6; i++ {
		a.SetRow(i, []float64{1, s_x[i], math.Pow(s_x[2], 2)})
	}

	var svd mat.SVD
	if ok := svd.Factorize(a, mat.SVDFull); !ok {
		t.Fatalf("failed to factorize a")
	}

	const rcond = 1e-15
	rank := svd.Rank(rcond)
	if rank == 0 {
		t.Fatalf("zero rank system")
	}

	b := mat.NewVecDense(6, s_y)

	var x mat.VecDense
	svd.SolveVecTo(&x, b, rank)

	t.Log(x)

	for i := 0; i < 6; i++ {
		var sum float64
		for j := 0; j < 3; j++ {
			sum += math.Pow(s_x[i], float64(j)) * x.At(j, 0)
		}
		if math.Abs(s_y[i]-sum) > 0.000001 {
			t.Fatalf("sample f(%v) = %v, expect %v", s_x[i], sum, s_y[i])
		}
	}
}
