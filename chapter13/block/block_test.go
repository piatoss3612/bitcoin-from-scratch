package block

import (
	"bytes"
	"encoding/hex"
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

	expected := 888171856257

	actual, _ := block.Difficulty().Int64()

	if actual != int64(expected) {
		t.Errorf("Difficulty() = %d, expected = %d", actual, expected)
	}
}

func TestCheckProofOfWork(t *testing.T) {
	blockRaw, _ := hex.DecodeString("04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec1")
	block, err := Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := block.CheckProofOfWork()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Errorf("CheckProofOfWork() = %t, expected = %t", ok, true)
	}

	blockRaw, _ = hex.DecodeString("04000000fbedbbf0cfdaf278c094f187f2eb987c86a199da22bbb20400000000000000007b7697b29129648fa08b4bcd13c9d5e60abb973a1efac9c8d573c71c807c56c3d6213557faa80518c3737ec0")
	block, err = Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	ok, err = block.CheckProofOfWork()
	if err != nil {
		t.Fatal(err)
	}

	if ok {
		t.Errorf("CheckProofOfWork() = %t, expected = %t", ok, false)
	}
}

func TestValidateMerkleRoot(t *testing.T) {
	hashesHex := []string{
		"f54cb69e5dc1bd38ee6901e4ec2007a5030e14bdd60afb4d2f3428c88eea17c1",
		"c57c2d678da0a7ee8cfa058f1cf49bfcb00ae21eda966640e312b464414731c1",
		"b027077c94668a84a5d0e72ac0020bae3838cb7f9ee3fa4e81d1eecf6eda91f3",
		"8131a1b8ec3a815b4800b43dff6c6963c75193c4190ec946b93245a9928a233d",
		"ae7d63ffcb3ae2bc0681eca0df10dda3ca36dedb9dbf49e33c5fbe33262f0910",
		"61a14b1bbdcdda8a22e61036839e8b110913832efd4b086948a6a64fd5b3377d",
		"fc7051c8b536ac87344c5497595d5d2ffdaba471c73fae15fe9228547ea71881",
		"77386a46e26f69b3cd435aa4faac932027f58d0b7252e62fb6c9c2489887f6df",
		"59cbc055ccd26a2c4c4df2770382c7fea135c56d9e75d3f758ac465f74c025b8",
		"7c2bf5687f19785a61be9f46e031ba041c7f93e2b7e9212799d84ba052395195",
		"08598eebd94c18b0d59ac921e9ba99e2b8ab7d9fccde7d44f2bd4d5e2e726d2e",
		"f0bb99ef46b029dd6f714e4b12a7d796258c48fee57324ebdc0bbc4700753ab1",
	}

	hashes := make([][]byte, len(hashesHex))
	for i, hashHex := range hashesHex {
		hashes[i], _ = hex.DecodeString(hashHex)
	}

	blockRaw, _ := hex.DecodeString("00000020fcb19f7895db08cadc9573e7915e3919fb76d59868a51d995201000000000000acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed691cfa85916ca061a00000000")
	block, err := Parse(blockRaw)
	if err != nil {
		t.Fatal(err)
	}

	block.TxHashes = hashes

	ok, err := block.ValidateMerkleRoot()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Errorf("ValidateMerkleRoot() = %t, expected = %t", ok, true)
	}
}
