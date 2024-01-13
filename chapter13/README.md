# Programming Bitcoin

## Install btcd

```bash
$ git clone https://github.com/btcsuite/btcd $GOPATH/src/github.com/btcsuite/btcd
$ cd $GOPATH/src/github.com/btcsuite/btcd
$ go install -v . ./cmd/...
```

## Run btcd simnet

```bash
$ btcd --configfile=btcd.conf
```

## Test

```bash
$ go test -cover ./...
?       chapter13       [no test files]
ok      chapter13/block 0.004s  coverage: 69.6% of statements
ok      chapter13/bloomfilter   0.004s  coverage: 82.4% of statements
ok      chapter13/ecc   0.060s  coverage: 62.8% of statements
ok      chapter13/merkleblock   0.003s  coverage: 79.8% of statements
ok      chapter13/network       0.479s  coverage: 50.2% of statements
ok      chapter13/script        0.030s  coverage: 14.6% of statements
ok      chapter13/tx    6.857s  coverage: 68.7% of statements
ok      chapter13/utils 0.003s  coverage: 49.3% of statements
```