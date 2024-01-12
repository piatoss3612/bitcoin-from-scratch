package ecc

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestDER(t *testing.T) {
	tests := []struct {
		caseName string
		r        *big.Int
		s        *big.Int
	}{
		{
			caseName: "1",
			r:        big.NewInt(0x1),
			s:        big.NewInt(0x2),
		},
		{
			caseName: "2",
			r:        randInt1(),
			s:        randInt2(),
		},
		{
			caseName: "3",
			r:        randInt1(),
			s:        randInt2(),
		},
	}

	for _, test := range tests {
		sig := NewS256Signature(test.r, test.s)
		der := sig.DER()
		sig2, err := ParseSignature(der)
		if err != nil {
			t.Fatalf("failed to parse signature: %s", err)
		}

		if sig2.R().Cmp(test.r) != 0 {
			t.Fatalf("R mismatch: %s", sig2.R())
		}

		if sig2.S().Cmp(test.s) != 0 {
			t.Fatalf("S mismatch: %s", sig2.S())
		}
	}
}

func randInt1() *big.Int {
	r := rand.New(rand.NewSource(0))
	return big.NewInt(0).Rand(r, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil))
}

func randInt2() *big.Int {
	r := rand.New(rand.NewSource(1))
	return big.NewInt(0).Rand(r, big.NewInt(0).Exp(big.NewInt(2), big.NewInt(255), nil))
}
