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

func (c *Curve) AddSample(x, y float64) {
	c.samples = append(c.samples, []float64{x, y})
	sort.Slice(c.samples, func(i, j int) bool { return c.samples[i][0] < c.samples[j][0] })
}

func (c *Curve) Fitting() {
	m := len(c.samples)
	Xv := []float64{}
	Y := []float64{}

	for i := 0; i < m; i++ {
		for j := 0; j < c.pow+1; j++ {
			Xv = append(Xv, math.Pow(c.samples[i][0], float64(j)))
		}
		Y = append(Y, c.samples[i][1])
	}
	matXv := mat.NewDense(m, c.pow+1, Xv)
	matY := mat.NewDense(m, 1, Y)
	var matTheta mat.Dense
	matTheta.Mul(matXv.T(), matXv)
	err := matTheta.Inverse(&matTheta)
	if err != nil {
		fmt.Println(err)
		return
	}
	matTheta.Mul(&matTheta, matXv.T())
	matTheta.Mul(&matTheta, matY)
	row, _ := matTheta.Caps()
	for i := 0; i < row; i++ {
		c.thetas = append(c.thetas, matTheta.At(i, 0))
	}
}
