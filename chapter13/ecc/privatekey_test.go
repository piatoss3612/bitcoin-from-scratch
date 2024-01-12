package ecc

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
)

func TestSign(t *testing.T) {
	b := []byte("hello world")

	pvk, err := NewS256PrivateKey(b)
	if err != nil {
		t.Fatal(err)
	}

	z := randInt1()
	sig, err := pvk.Sign(z.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	ok, err := pvk.Point().Verify(z.Bytes(), sig)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("verify failed")
	}
}

func TestWIF(t *testing.T) {
	tests := []struct {
		caseName   string
		secret     []byte
		compressed bool
		testnet    bool
		expect     string
	}{
		{
			caseName:   "1",
			secret:     big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(0).Exp(big.NewInt(2), big.NewInt(199), nil)).Bytes(),
			compressed: true,
			testnet:    false,
			expect:     "L5oLkpV3aqBJ4BgssVAsax1iRa77G5CVYnv9adQ6Z87te7TyUdSC",
		},
		{
			caseName:   "2",
			secret:     big.NewInt(0).Sub(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(0).Exp(big.NewInt(2), big.NewInt(201), nil)).Bytes(),
			compressed: false,
			testnet:    true,
			expect:     "93XfLeifX7Jx7n7ELGMAf1SUR6f9kgQs8Xke8WStMwUtrDucMzn",
		},
		{
			caseName:   "3",
			secret:     decodeHex("0dba685b4511dbd3d368e5c4358a1277de9486447af7b3604a69b8d9d8b7889d"),
			compressed: false,
			testnet:    false,
			expect:     "5HvLFPDVgFZRK9cd4C5jcWki5Skz6fmKqi1GQJf5ZoMofid2Dty",
		},
		{
			caseName:   "4",
			secret:     decodeHex("1cca23de92fd1862fb5b76e5f4f50eb082165e5191e116c18ed1a6b24be6a53f"),
			compressed: true,
			testnet:    true,
			expect:     "cNYfWuhDpbNM1JWc3c6JTrtrFVxU4AGhUKgw5f93NP2QaBqmxKkg",
		},
	}

	for _, test := range tests {
		pvk, err := NewS256PrivateKey(test.secret)
		if err != nil {
			t.Fatalf("%s: %v", test.caseName, err)
		}

		wif := pvk.WIF(test.compressed, test.testnet)

		if !strings.EqualFold(wif, test.expect) {
			t.Fatalf("%s: expect %s, got %s", test.caseName, test.expect, wif)
		}
	}
}

func decodeHex(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}
