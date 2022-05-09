package logarithm

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func Test_LogaN(t *testing.T) {
	a := rand.Float64()
	n := rand.Float64()
	m := LogaN(a, n)
	n1 := math.Pow(a, m)
	a1 := math.Pow(n, 1/m)
	if !IsEqual(a, a1) {
		t.Errorf("The base number is invalid, %v,%v\n", a, a1)
	}
	if !IsEqual(n, n1) {
		t.Errorf("The real number is invalid, %v %v\n", n, n1)
	}
	fmt.Printf("Log(%v)(%v) = %v\n%v ^ %v = %v\n%v ^ (1/%v) = %v\n", a, n, m, a, m, n1, n, m, a1)
}

func Test_benford(t *testing.T) {
	a := rand.Uint32()
	b := rand.Uint32()
	m := Benford(uint(a), uint(b))
	fmt.Printf("Log(%v)((%v+1)/%v) = %v\n", b, a, a, m)
}
