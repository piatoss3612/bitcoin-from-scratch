package main

import (
	"chapter11/utils"
	"encoding/hex"
)

func main() {
	practice4()
}

func practice1() {
	b1, _ := hex.DecodeString("c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5")
	b2, _ := hex.DecodeString("c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5")

	parent := utils.Hash256(append(b1, b2...))
	println(hex.EncodeToString(parent))

	// 8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd
}

func practice2() {
	hexHashes := []string{
		"c117ea8ec828342f4dfb0ad6bd140e03a50720ece40169ee38bdc15d9eb64cf5",
		"c131474164b412e3406696da1ee20ab0fc9bf41c8f05fa8ceea7a08d672d7cc5",
		"f391da6ecfeed1814efae39e7fcb3838ae0b02c02ae7d0a5848a66947c0727b0",
		"3d238a92a94532b946c90e19c49351c763696cff3db400485b813aecb8a13181",
		"10092f2633be5f3ce349bf9ddbde36caa3dd10dfa0ec8106bce23acbff637dae",
	}

	hashes := make([][]byte, len(hexHashes))
	for i, hexHash := range hexHashes {
		hashes[i], _ = hex.DecodeString(hexHash)
	}

	if len(hashes)%2 == 1 {
		hashes = append(hashes, hashes[len(hashes)-1])
	}

	parentLevel := make([][]byte, len(hashes)/2)
	for i := 0; i < len(hashes); i += 2 {
		parentLevel[i/2] = utils.Hash256(append(hashes[i], hashes[i+1]...))
	}

	for _, parent := range parentLevel {
		println(hex.EncodeToString(parent))
	}

	/*
		8b30c5ba100f6f2e5ad1e2a742e5020491240f8eb514fe97c713c31718ad7ecd
		7f4e6f9e224e20fda0ae4c44114237f97cd35aca38d83081c9bfd41feb907800
		3ecf6115380c77e8aae56660f5634982ee897351ba906a6837d15ebc3a225df0
	*/
}

func practice3() {
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

	hashes := make([][]byte, len(hexHashes))
	for i, hexHash := range hexHashes {
		hashes[i], _ = hex.DecodeString(hexHash)
	}

	currentLevel := hashes
	for len(currentLevel) > 1 {
		currentLevel = utils.MerkleParentLevel(currentLevel)
	}

	println(hex.EncodeToString(currentLevel[0]))

	// acbcab8bcc1af95d8d563b77d24c3d19b18f1486383d75a5085c4e86c86beed6
}

func practice4() {
	txHexHashes := []string{
		"42f6f52f17620653dcc909e58bb352e0bd4bd1381e2955d19c00959a22122b2e",
		"94c3af34b9667bf787e1c6a0a009201589755d01d02fe2877cc69b929d2418d4",
		"959428d7c48113cb9149d0566bde3d46e98cf028053c522b8fa8f735241aa953",
		"a9f27b99d5d108dede755710d4a1ffa2c74af70b4ca71726fa57d68454e609a2",
		"62af110031e29de1efcad103b3ad4bec7bdcf6cb9c9f4afdd586981795516577",
		"766900590ece194667e9da2984018057512887110bf54fe0aa800157aec796ba",
		"e8270fb475763bc8d855cfe45ed98060988c1bdcad2ffc8364f783c98999a208",
	}

	txHashes := make([][]byte, len(txHexHashes))
	for i, hexHash := range txHexHashes {
		b, _ := hex.DecodeString(hexHash)
		// reverse b
		for j := 0; j < len(b)/2; j++ {
			b[j], b[len(b)-j-1] = b[len(b)-j-1], b[j]
		}
		txHashes[i] = b
	}

	root := utils.MerkleRoot(txHashes)
	// reverse root
	for j := 0; j < len(root)/2; j++ {
		root[j], root[len(root)-j-1] = root[len(root)-j-1], root[j]
	}

	println(hex.EncodeToString(root))

	// 654d6181e18e4ac4368383fdc5eead11bf138f9b7ac1e15334e4411b3c4797d9
}
