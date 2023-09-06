package main

import (
	"chapter06/tx"
	"encoding/hex"
	"fmt"
)

func main() {

	txHex, err := hex.DecodeString("010000000456919960ac691763688d3d3bcea9ad6ecaf875df5339e148a1fc61c6ed7a069e010000006a47304402204585bcdef85e6b1c6af5c2669d4830ff86e42dd205c0e089bc2a821657e951c002201024a10366077f87d6bce1f7100ad8cfa8a064b39d4e8fe4ea13a7b71aa8180f012102f0da57e85eec2934a82a585ea337ce2f4998b50ae699dd79f5880e253dafafb7feffffffeb8f51f4038dc17e6313cf831d4f02281c2a468bde0fafd37f1bf882729e7fd3000000006a47304402207899531a52d59a6de200179928ca900254a36b8dff8bb75f5f5d71b1cdc26125022008b422690b8461cb52c3cc30330b23d574351872b7c361e9aae3649071c1a7160121035d5c93d9ac96881f19ba1f686f15f009ded7c62efe85a872e6a19b43c15a2937feffffff567bf40595119d1bb8a3037c356efd56170b64cbcc160fb028fa10704b45d775000000006a47304402204c7c7818424c7f7911da6cddc59655a70af1cb5eaf17c69dadbfc74ffa0b662f02207599e08bc8023693ad4e9527dc42c34210f7a7d1d1ddfc8492b654a11e7620a0012102158b46fbdff65d0172b7989aec8850aa0dae49abfb84c81ae6e5b251a58ace5cfeffffffd63a5e6c16e620f86f375925b21cabaf736c779f88fd04dcad51d26690f7f345010000006a47304402200633ea0d3314bea0d95b3cd8dadb2ef79ea8331ffe1e61f762c0f6daea0fabde022029f23b3e9c30f080446150b23852028751635dcee2be669c2a1686a4b5edf304012103ffd6f4a67e94aba353a00882e563ff2722eb4cff0ad6006e86ee20dfe7520d55feffffff0251430f00000000001976a914ab0c0b2e98b1ab6dbf67d4750b0a56244948a87988ac005a6202000000001976a9143c82d7df364eb6c75be8c80df2b3eda8db57397088ac46430600")
	if err != nil {
		panic(err)
	}

	tx, err := tx.ParseTx(txHex, false)
	if err != nil {
		panic(err)
	}

	fmt.Println("2nd input's ScriptSig:", tx.Inputs[1].ScriptSig)
	fmt.Println("1st output's ScriptPubKey:", tx.Outputs[0].ScriptPubKey)
	fmt.Println("2nd output's amount:", tx.Outputs[1].Amount)

	/*
		// 예제 코드에서는 잘 돌아가는데 최신 버전의 트랜잭션 아이디를 사용해 파싱을 하는 경우에는
		// 아이디 미스매치 오류가 발생하는 것을 봐서는 최신 버전의 트랜잭션은 파싱하는 방법이 다른 것 같다.

		tf := tx.NewTxFetcher()

		tx, err := tf.Fetch("c05fd9e2a85a716e2c2679052cc24839f91bb7df510169aa54ed02a6c67dae9e", false, true)
		if err != nil {
			panic(err)
		}

		fmt.Println("tx:", tx)
	*/
}
