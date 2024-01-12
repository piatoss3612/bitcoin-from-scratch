package block

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
)

func TestSerialize(t *testing.T) {
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block, err := Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	expected := blockRaw

	b, err := block.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.EqualFold(b, expected) {
		t.Errorf("Serialize() = %x, expected = %x", b, expected)
	}
}

func TestHash(t *testing.T) {
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block, err := Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	expected, _ := hex.DecodeString("0000000000000000007e9e4c586439b0cdbe13b1370bdd9435d76a644d047523")

	actual, err := block.Hash()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.EqualFold(actual, expected) {
		t.Errorf("Hash() = %x, expected = %x", actual, expected)
	}
}

func TestBip9(t *testing.T) {
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block, err := Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	if !block.Bip9() {
		t.Errorf("Bip9() = %t, expected = %t", block.Bip9(), true)
	}

	blockRaw, _ = hex.DecodeString("0400000039fa821848781f027a2e6dfabbf6bda920d9ae61b63400030000000000000000ecae536a304042e3154be0e3e9a8220e5568c3433a9ab49ac4cbb74f8df8e8b0cc2acf569fb9061806652c27")
	block, err = Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	if block.Bip9() {
		t.Errorf("Bip9() = %t, expected = %t", block.Bip9(), false)
	}
}

func TestBip91(t *testing.T) {
	blockRaw, _ := hex.DecodeString("1200002028856ec5bca29cf76980d368b0a163a0bb81fc192951270100000000000000003288f32a2831833c31a25401c52093eb545d28157e200a64b21b3ae8f21c507401877b5935470118144dbfd1")
	block, err := Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	if !block.Bip91() {
		t.Errorf("Bip91() = %t, expected = %t", block.Bip91(), true)
	}

	blockRaw, _ = hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block, err = Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	if block.Bip91() {
		t.Errorf("Bip91() = %t, expected = %t", block.Bip91(), false)
	}
}

func TestBip141(t *testing.T) {
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block, err := Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	if !block.Bip141() {
		t.Errorf("Bip141() = %t, expected = %t", block.Bip141(), true)
	}

	blockRaw, _ = hex.DecodeString("0000002066f09203c1cf5ef1531f24ed21b1915ae9abeb691f0d2e0100000000000000003de0976428ce56125351bae62c5b8b8c79d8297c702ea05d60feabb4ed188b59c36fa759e93c0118b74b2618")
	block, err = Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	if block.Bip141() {
		t.Errorf("Bip141() = %t, expected = %t", block.Bip141(), false)
	}
}

func TestTarget(t *testing.T) {
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block, err := Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	expected := big.NewInt(0).SetBytes([]byte{0x01, 0x3c, 0xe9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	if block.Target().Cmp(expected) != 0 {
		t.Errorf("Target() = %x, expected = %x", block.Target(), expected)
	}
}

// def test_difficulty(self):
// block_raw = bytes.fromhex('020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d')
// stream = BytesIO(block_raw)
// block = Block.parse(stream)
// self.assertEqual(int(block.difficulty()), 888171856257)

func TestDifficulty(t *testing.T) {
	blockRaw, _ := hex.DecodeString("020000208ec39428b17323fa0ddec8e887b4a7c53b8c0a0a220cfd0000000000000000005b0750fce0a889502d40508d39576821155e9c9e3f5c3157f961db38fd8b25be1e77a759e93c0118a4ffd71d")
	block, err := Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	expected := big.NewInt(888171856257)

	actual, _ := block.Difficulty().Int64()

	fmt.Println(actual)
	fmt.Println(expected)
}
