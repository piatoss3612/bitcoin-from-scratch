package utils

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestP2pkhAddress(t *testing.T) {
	h160, _ := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
	expected := "1BenRpVUFK65JFWcQSuHnJKzc4M8ZP8Eqa"
	actual := H160ToP2pkhAddress(h160, false)
	if actual != expected {
		t.Fatalf("H160ToP2pkhAddress(%x, false) = %s, expected = %s", h160, actual, expected)
	}

	expected = "mrAjisaT4LXL5MzE81sfcDYKU3wqWSvf9q"
	actual = H160ToP2pkhAddress(h160, true)
	if actual != expected {
		t.Fatalf("H160ToP2pkhAddress(%x, true) = %s, expected = %s", h160, actual, expected)
	}
}

func TestP2shAddress(t *testing.T) {
	h160, _ := hex.DecodeString("74d691da1574e6b3c192ecfb52cc8984ee7b6c56")
	expected := "3CLoMMyuoDQTPRD3XYZtCvgvkadrAdvdXh"
	actual := H160ToP2shAddress(h160, false)
	if actual != expected {
		t.Fatalf("H160ToP2shAddress(%x, false) = %s, expected = %s", h160, actual, expected)
	}

	expected = "2N3u1R6uwQfuobCqbCgBkpsgBxvr1tZpe7B"
	actual = H160ToP2shAddress(h160, true)
	if actual != expected {
		t.Fatalf("H160ToP2shAddress(%x, true) = %s, expected = %s", h160, actual, expected)
	}
}

func TestMerkleParent(t *testing.T) {
	txHash0, _ := hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	txHash1, _ := hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")
	want, _ := hex.DecodeString("8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd")
	actual := MerkleParent(txHash0, txHash1)
	if !bytes.EqualFold(actual, want) {
		t.Fatalf("MerkleParent(%x, %x) = %x, expected = %x", txHash0, txHash1, actual, want)
	}
}

func TestMerkleParentLevel(t *testing.T) {
	hexHashes := []string{
		"c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5",
		"c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5",
		"f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0",
		"3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181",
		"10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae",
		"7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161",
		"8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc",
		"dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877",
		"b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59",
		"95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c",
		"2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908",
	}
	var txHashes [][]byte
	for _, hexHash := range hexHashes {
		txHash, _ := hex.DecodeString(hexHash)
		txHashes = append(txHashes, txHash)
	}
	wantHexHashes := []string{
		"8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd",
		"7f4e6f9e224e20fda0ae4c44114237f97cd35aca38d83081c9bfd41feb907800",
		"ade48f2bbb57318cc79f3a8678febaa827599c509dce5940602e54c7733332e7",
		"68b3e2ab8182dfd646f13fdf01c335cf32476482d963f5cd94e934e6b3401069",
		"43e7274e77fbe8e5a42a8fb58f7decdb04d521f319f332d88e6b06f8e6c09e27",
		"1796cd3ca4fef00236e07b723d3ed88e1ac433acaaa21da64c4b33c946cf3d10",
	}

	var wantTxHashes [][]byte
	for _, wantHexHash := range wantHexHashes {
		wantTxHash, _ := hex.DecodeString(wantHexHash)
		wantTxHashes = append(wantTxHashes, wantTxHash)
	}

	actualTxHashes := MerkleParentLevel(txHashes)

	for i, actualTxHash := range actualTxHashes {
		if !bytes.EqualFold(actualTxHash, wantTxHashes[i]) {
			t.Fatalf("MerkleParentLevel(txHashes) = %x, expected = %x", actualTxHash, wantTxHashes[i])
		}
	}
}

func TestMerkleRoot(t *testing.T) {
	hexHashes := []string{
		"c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5",
		"c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5",
		"f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0",
		"3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181",
		"10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae",
		"7d37b3d54fa6a64869084bfd2e831309118b9e833610e6228adacdbd1b4ba161",
		"8118a77e542892fe15ae3fc771a4abfd2f5d5d5997544c3487ac36b5c85170fc",
		"dff6879848c2c9b62fe652720b8df5272093acfaa45a43cdb3696fe2466a3877",
		"b825c0745f46ac58f7d3759e6dc535a1fec7820377f24d4c2c6ad2cc55c0cb59",
		"95513952a04bd8992721e9b7e2937f1c04ba31e0469fbe615a78197f68f52b7c",
		"2e6d722e5e4dbdf2447ddecc9f7dabb8e299bae921c99ad5b0184cd9eb8e5908",
		"b13a750047bc0bdceb2473e5fe488c2596d7a7124b4e716fdd29b046ef99bbf0",
	}
	var txHashes [][]byte
	for _, hexHash := range hexHashes {
		txHash, _ := hex.DecodeString(hexHash)
		txHashes = append(txHashes, txHash)
	}

	wantHexHash := "acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed6"
	wantHash, _ := hex.DecodeString(wantHexHash)

	actualHash := MerkleRoot(txHashes)

	if !bytes.EqualFold(actualHash, wantHash) {
		t.Fatalf("MerkleRoot(txHashes) = %x, expected = %x", actualHash, wantHash)
	}
}
