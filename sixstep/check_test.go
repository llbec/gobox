package sixstep

import (
	"testing"
)

func TestCalcCode(t *testing.T) {
	if _, err := CalcCode(""); err == nil {
		t.Errorf("length check error!")
	}
	s := "34052419800101001X"
	v, err := CalcCode(s[:17])
	if err != nil {
		t.Errorf("lenght error!")
	}
	if v != "X" {
		t.Errorf("Expect X, but is %s", v)
	}
}
