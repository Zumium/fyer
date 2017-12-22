package alg

import "testing"

func TestFragCutter(t *testing.T) {
	fc := FragCutter{Size: 5, FragSize: 2}
	r := fc.Cut()
	if len(r) != 3 {
		t.Fatal("number of frags mismatching")
	}
	if r[0].Index != 0 || r[0].Start != 0 || r[0].Size != 2 {
		t.Fatalf("frag 0 error: %v", r[0])
	}
	if r[1].Index != 1 || r[1].Start != 2 || r[1].Size != 2 {
		t.Fatalf("frag 1 error: %v", r[1])
	}
	if r[2].Index != 2 || r[2].Start != 4 || r[2].Size != 1 {
		t.Fatalf("frag 2 error: %v", r[2])
	}
}
