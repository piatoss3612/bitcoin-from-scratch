package tx

import (
	"bytes"
	"chapter13/ecc"
	"encoding/hex"
	"math/big"
	"testing"
)

func TestParseTx(t *testing.T) {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")
	tx, err := ParseTx(rawTx)
	if err != nil {
		t.Fatal(err)
	}

	if tx.Version != 1 {
		t.Errorf("expected version 1, got %d", tx.Version)
	}

	if len(tx.Inputs) != 1 {
		t.Errorf("expected 1 input, got %d", len(tx.Inputs))
	}

	expected1 := "d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81"

	if tx.Inputs[0].PrevTx != expected1 {
		t.Errorf("expected %s, got %s", expected1, tx.Inputs[0].PrevTx)
	}

	if tx.Inputs[0].PrevIndex != 0 {
		t.Errorf("expected 0, got %d", tx.Inputs[0].PrevIndex)
	}

	expected2, _ := hex.DecodeString("6b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278a")

	b, err := tx.Inputs[0].ScriptSig.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, expected2) {
		t.Errorf("expected %x, got %x", expected2, b)
	}

	if tx.Inputs[0].SeqNo != 0xfffffffe {
		t.Errorf("expected 0xfffffffe, got %d", tx.Inputs[0].SeqNo)
	}

	if len(tx.Outputs) != 2 {
		t.Errorf("expected 2 outputs, got %d", len(tx.Outputs))
	}

	expected3 := 32454049

	if tx.Outputs[0].Amount != expected3 {
		t.Errorf("expected %d, got %d", expected3, tx.Outputs[0].Amount)
	}

	expected4, _ := hex.DecodeString("1976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac")

	b, err = tx.Outputs[0].ScriptPubKey.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, expected4) {
		t.Errorf("expected %x, got %x", expected4, b)
	}

	expected5 := 10011545

	if tx.Outputs[1].Amount != expected5 {
		t.Errorf("expected %d, got %d", expected5, tx.Outputs[1].Amount)
	}

	expected6, _ := hex.DecodeString("1976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac")

	b, err = tx.Outputs[1].ScriptPubKey.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, expected6) {
		t.Errorf("expected %x, got %x", expected6, b)
	}

	if tx.Locktime != 410393 {
		t.Errorf("expected 410393, got %d", tx.Locktime)
	}
}

func TestSerialize(t *testing.T) {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")

	tx, err := ParseTx(rawTx)
	if err != nil {
		t.Fatal(err)
	}

	b, err := tx.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, rawTx) {
		t.Errorf("expected %x, got %x", rawTx, b)
	}
}

func TestInput(t *testing.T) {
	txHash := "d1c789a9c60383bf715f3f6ad9d14b91fe55f3deb369fe5d9280cb1a01793f81"

	index := 0
	expected1 := 42505594
	txin := TxIn{
		PrevTx:    txHash,
		PrevIndex: index,
	}
	if txin.PrevIndex != index {
		t.Errorf("expected %d, got %d", index, txin.PrevIndex)
	}

	f := NewTxFetcher()

	actual, err := txin.Value(f, false)
	if err != nil {
		t.Fatal(err)
	}

	if actual != expected1 {
		t.Errorf("expected %d, got %d", expected1, actual)
	}

	expected2, _ := hex.DecodeString("1976a914a802fc56c704ce87c42d7c92eb75e7896bdc41ae88ac")

	pubkey, err := txin.ScriptPubKey(f, false)
	if err != nil {
		t.Fatal(err)
	}

	b, err := pubkey.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, expected2) {
		t.Errorf("expected %x, got %x", expected2, b)
	}
}

func TestFee(t *testing.T) {
	rawTx, _ := hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")

	tx, err := ParseTx(rawTx)
	if err != nil {
		t.Fatal(err)
	}

	f := NewTxFetcher()

	fee1, err := tx.Fee(f)
	if err != nil {
		t.Fatal(err)
	}

	if fee1 != 40000 {
		t.Errorf("expected 40000, got %d", fee1)
	}

	rawTx, _ = hex.DecodeString("010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600")

	tx, err = ParseTx(rawTx)
	if err != nil {
		t.Fatal(err)
	}

	fee2, err := tx.Fee(f)
	if err != nil {
		t.Fatal(err)
	}

	if fee2 != 140500 {
		t.Errorf("expected 140500, got %d", fee2)
	}
}

func TestSigHash(t *testing.T) {
	f := NewTxFetcher()

	txId := "452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03"

	tx, err := f.Fetch(txId, false, false)
	if err != nil {
		t.Fatal(err)
	}

	sigHash, err := tx.SigHash(0)
	if err != nil {
		t.Fatal(err)
	}

	expected, _ := hex.DecodeString("27e0c5994dec7824e56dec6b2fcb342eb7cdb0d0957c2fce9882f715e85d81a6")

	if !bytes.Equal(sigHash, expected) {
		t.Errorf("expected %x, got %x", expected, sigHash)
	}
}

func TestVerifyP2pkh(t *testing.T) {
	f := NewTxFetcher()

	txId := "452c629d67e41baec3ac6f04fe744b4b9617f8f859c63b3002f8684e7a4fee03"

	tx, err := f.Fetch(txId, false, false)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := tx.Verify()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("tx verification failed")
	}

	txId2 := "5418099cc755cb9dd3ebc6cf1a7888ad53a1a3beb5a025bce89eb1bf7f1650a2"

	tx, err = f.Fetch(txId2, true, false) // testnet
	if err != nil {
		t.Fatal(err)
	}

	ok, err = tx.Verify()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("tx verification failed")
	}
}

func TestVerifyP2sh(t *testing.T) {
	f := NewTxFetcher()

	txId := "46df1a9484d0a81d03ce0ee543ab6e1a23ed06175c104a178268fad381216c2b"

	tx, err := f.Fetch(txId, false, false)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := tx.Verify()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("tx verification failed")
	}
}

// 통과
func TestVerifyP2wpkh(t *testing.T) {
	f := NewTxFetcher()

	txId := "d869f854e1f8788bcff294cc83b280942a8c728de71eb709a2c29d10bfe21b7c"

	tx, err := f.Fetch(txId, true, false)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := tx.Verify()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("tx verification failed")
	}
}

func TestVerifyP2shP2wpkh(t *testing.T) {
	f := NewTxFetcher()

	txId := "c586389e5e4b3acb9d6c8be1c19ae8ab2795397633176f5a6442a261bbdefc3a"

	tx, err := f.Fetch(txId, false, false)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := tx.Verify()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("tx verification failed")
	}
}

func TestVerifyP2wsh(t *testing.T) {
	f := NewTxFetcher()

	txId := "78457666f82c28aa37b74b506745a7c7684dc7842a52a457b09f09446721e11c"

	tx, err := f.Fetch(txId, true, false)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := tx.Verify()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("tx verification failed")
	}
}

func TestVerifyP2shP2wsh(t *testing.T) {
	f := NewTxFetcher()

	txId := "954f43dbb30ad8024981c07d1f5eb6c9fd461e2cf1760dd1283f052af746fc88"

	tx, err := f.Fetch(txId, true, false)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := tx.Verify()
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("tx verification failed")
	}
}

func TestSignInput(t *testing.T) {
	secret := big.NewInt(8675309)
	pvk, _ := ecc.NewS256PrivateKey(secret.Bytes())

	rawTx, _ := hex.DecodeString("010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d00000000ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000")

	tx, err := ParseTx(rawTx, true)
	if err != nil {
		t.Fatal(err)
	}

	ok, err := tx.SignInput(0, pvk, true)
	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("signing failed")
	}

	expected, _ := hex.DecodeString("010000000199a24308080ab26e6fb65c4eccfadf76749bb5bfa8cb08f291320b3c21e56f0d0d0000006b4830450221008ed46aa2cf12d6d81065bfabe903670165b538f65ee9a3385e6327d80c66d3b502203124f804410527497329ec4715e18558082d489b218677bd029e7fa306a72236012103935581e52c354cd2f484fe8ed83af7a3097005b2f9c60bff71d35bd795f54b67ffffffff02408af701000000001976a914d52ad7ca9b3d096a38e752c2018e6fbc40cdf26f88ac80969800000000001976a914507b27411ccf7f16f10297de6cef3f291623eddf88ac00000000")

	b, err := tx.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(b, expected) {
		t.Errorf("expected %x, got %x", expected, b)
	}
}

func TestIsCoinbase(t *testing.T) {
	rawTx, _ := hex.DecodeString("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00ffffffff01faf20b58000000001976a914338c84849423992471bffb1a54a8d9b1d69dc28a88ac00000000")

	tx, err := ParseTx(rawTx)
	if err != nil {
		t.Fatal(err)
	}

	if !tx.IsCoinbase() {
		t.Error("expected coinbase tx")
	}
}

func TestCoinbaseHeight(t *testing.T) {
	rawTx, _ := hex.DecodeString("01000000010000000000000000000000000000000000000000000000000000000000000000ffffffff5e03d71b07254d696e656420627920416e74506f6f6c20626a31312f4542312f4144362f43205914293101fabe6d6d678e2c8c34afc36896e7d9402824ed38e856676ee94bfdb0c6c4bcd8b2e5666a0400000000000000c7270000a5e00e00ffffffff01faf20b58000000001976a914338c84849423992471bffb1a54a8d9b1d69dc28a88ac00000000")

	tx, err := ParseTx(rawTx)
	if err != nil {
		t.Fatal(err)
	}

	expected := 465879

	ok, actual := tx.CoinbaseHeight()
	if !ok {
		t.Error("expected coinbase tx")
	}

	if actual != expected {
		t.Errorf("expected %d, got %d", expected, actual)
	}

	rawTx, _ = hex.DecodeString("0100000001813f79011acb80925dfe69b3def355fe914bd1d96a3f5f71bf8303c6a989c7d1000000006b483045022100ed81ff192e75a3fd2304004dcadb746fa5e24c5031ccfcf21320b0277457c98f02207a986d955c6e0cb35d446a89d3f56100f4d7f67801c31967743a9c8e10615bed01210349fc4e631e3624a545de3f89f5d8684c7b8138bd94bdd531d2e213bf016b278afeffffff02a135ef01000000001976a914bc3b654dca7e56b04dca18f2566cdaf02e8d9ada88ac99c39800000000001976a9141c4bc762dd5423e332166702cb75f40df79fea1288ac19430600")

	tx, err = ParseTx(rawTx)
	if err != nil {
		t.Fatal(err)
	}

	ok, actual = tx.CoinbaseHeight()
	if ok {
		t.Error("expected non-coinbase tx")
	}

	if actual != 0 {
		t.Errorf("expected 0, got %d", actual)
	}
}
