package ecc

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestS256Order(t *testing.T) {
	point, err := G.Mul(N)
	if err != nil {
		t.Error(err)
	}

	if point.X() != nil {
		t.Errorf("N*G is not infinity, got %s", point)
	}
}

func TestS256PubPoint(t *testing.T) {
	tests := []struct {
		caseName string
		secret   *big.Int
		x, y     string
	}{
		{
			caseName: "1",
			secret:   big.NewInt(7),
			x:        "0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc",
			y:        "0x6aebca40ba255960a3178d6d861a54dba813d0b813fde7b5a5082628087264da",
		},
		{
			caseName: "2",
			secret:   big.NewInt(1485),
			x:        "0xc982196a7466fbbbb0e27a940b6af926c1a74d5ad07128c82824a11b5398afda",
			y:        "0x7a91f9eae64438afb9ce6448a1c133db2d8fb9254e4546b6f001637d50901f55",
		},
		{
			caseName: "3",
			secret:   big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), nil),
			x:        "0x8f68b9d2f63b5f339239c1ad981f162ee88c5678723ea3351b7b444c9ec4c0da",
			y:        "0x662a9f2dba063986de1d90c2b6be215dbbea2cfe95510bfdf23cbf79501fff82",
		},
		{
			caseName: "4",
			secret:   big.NewInt(0).Add(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(240), nil), big.NewInt(0).Exp(big.NewInt(2), big.NewInt(31), nil)),
			x:        "0x9577ff57c8234558f293df502ca4f09cbc65a6572c842b39b366f21717945116",
			y:        "0x10b49c67fa9365ad7b90dab070be339a1daf9052373ec30ffae4f72d5e66d053",
		},
	}

	for _, test := range tests {
		secret := test.secret
		x, ok := big.NewInt(0).SetString(test.x, 0)
		if !ok {
			t.Fatalf("case %s: failed to parse x: %s", test.caseName, test.x)
		}
		y, ok := big.NewInt(0).SetString(test.y, 0)
		if !ok {
			t.Fatalf("case %s: failed to parse y: %s", test.caseName, test.y)
		}

		xel, err := NewS256FieldElement(x)
		if err != nil {
			t.Fatalf("case %s: failed to create x field element: %s", test.caseName, err)
		}
		yel, err := NewS256FieldElement(y)
		if err != nil {
			t.Fatalf("case %s: failed to create y field element: %s", test.caseName, err)
		}

		point, err := NewS256Point(xel, yel)
		if err != nil {
			t.Fatalf("case %s: failed to create point: %s", test.caseName, err)
		}

		point2, err := G.Mul(secret)
		if err != nil {
			t.Fatalf("case %s: failed to create point: %s", test.caseName, err)
		}

		if !point.Equal(point2) {
			t.Fatalf("case %s: expected %s, got %s", test.caseName, point, point2)
		}
	}
}

func TestS256Verify(t *testing.T) {
	x, _ := big.NewInt(0).SetString("0x887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c", 0)
	y, _ := big.NewInt(0).SetString("0x61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34", 0)

	xel, err := NewS256FieldElement(x)
	if err != nil {
		t.Fatalf("failed to create x field element: %s", err)
	}

	yel, err := NewS256FieldElement(y)
	if err != nil {
		t.Fatalf("failed to create y field element: %s", err)
	}

	point, err := NewS256Point(xel, yel)
	if err != nil {
		t.Fatalf("failed to create point: %s", err)
	}

	z, _ := big.NewInt(0).SetString("0xec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60", 0)
	r, _ := big.NewInt(0).SetString("0xac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395", 0)
	s, _ := big.NewInt(0).SetString("0x68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4", 0)

	sig := NewS256Signature(r, s)

	ok, err := point.Verify(z.Bytes(), sig)
	if err != nil {
		t.Fatalf("failed to verify: %s", err)
	}

	if !ok {
		t.Fatalf("failed to verify")
	}

	z, _ = big.NewInt(0).SetString("0x7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d", 0)
	r, _ = big.NewInt(0).SetString("0xeff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c", 0)
	s, _ = big.NewInt(0).SetString("0xc7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6", 0)

	sig = NewS256Signature(r, s)

	ok, err = point.Verify(z.Bytes(), sig)
	if err != nil {
		t.Fatalf("failed to verify: %s", err)
	}

	if !ok {
		t.Fatalf("failed to verify")
	}
}

func TestS256Sec(t *testing.T) {
	tests := []struct {
		caseName     string
		coefficient  *big.Int
		compressed   string
		uncompressed string
	}{
		{
			caseName:     "1",
			coefficient:  big.NewInt(0).Exp(big.NewInt(999), big.NewInt(3), nil),
			compressed:   "039d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d5",
			uncompressed: "049d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d56fa15cc7f3d38cda98dee2419f415b7513dde1301f8643cd9245aea7f3f911f9",
		},
		{
			caseName:     "2",
			coefficient:  big.NewInt(123),
			compressed:   "03a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5",
			uncompressed: "04a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5204b5d6f84822c307e4b4a7140737aec23fc63b65b35f86a10026dbd2d864e6b",
		},
		{
			caseName:     "3",
			coefficient:  big.NewInt(42424242),
			compressed:   "03aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e",
			uncompressed: "04aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e21ec53f40efac47ac1c5211b2123527e0e9b57ede790c4da1e72c91fb7da54a3",
		},
	}

	for _, test := range tests {
		point, err := G.Mul(test.coefficient)
		if err != nil {
			t.Fatalf("case %s: failed to create point: %s", test.caseName, err)
		}

		compressed := point.SEC(true)
		expected, _ := hex.DecodeString(test.compressed)

		if !bytes.EqualFold(compressed, expected) {
			t.Fatalf("case %s: expected %x, got %x", test.caseName, expected, compressed)
		}

		uncompressed := point.SEC(false)
		expected, _ = hex.DecodeString(test.uncompressed)

		if !bytes.EqualFold(uncompressed, expected) {
			t.Fatalf("case %s: expected %x, got %x", test.caseName, expected, uncompressed)
		}
	}
}

func TestS256Address(t *testing.T) {
	tests := []struct {
		caseName   string
		secret     *big.Int
		mainnet    string
		testnet    string
		compressed bool
	}{
		{
			caseName:   "1",
			secret:     big.NewInt(0).Exp(big.NewInt(888), big.NewInt(3), nil),
			mainnet:    "148dY81A9BmdpMhvYEVznrM45kWN32vSCN",
			testnet:    "mieaqB68xDCtbUBYFoUNcmZNwk74xcBfTP",
			compressed: true,
		},
		{
			caseName:   "2",
			secret:     big.NewInt(321),
			mainnet:    "1S6g2xBJSED7Qr9CYZib5f4PYVhHZiVfj",
			testnet:    "mfx3y63A7TfTtXKkv7Y6QzsPFY6QCBCXiP",
			compressed: false,
		},
		{
			caseName:   "3",
			secret:     big.NewInt(4242424242),
			mainnet:    "1226JSptcStqn4Yq9aAmNXdwdc2ixuH9nb",
			testnet:    "mgY3bVusRUL6ZB2Ss999CSrGVbdRwVpM8s",
			compressed: false,
		},
	}

	for _, test := range tests {
		point, err := G.Mul(test.secret)
		if err != nil {
			t.Fatalf("case %s: failed to create point: %s", test.caseName, err)
		}

		mainnet := point.Address(test.compressed, false)
		if mainnet != test.mainnet {
			t.Fatalf("case %s: expected %s, got %s", test.caseName, test.mainnet, mainnet)
		}

		testnet := point.Address(test.compressed, true)
		if testnet != test.testnet {
			t.Fatalf("case %s: expected %s, got %s", test.caseName, test.testnet, testnet)
		}
	}
}
