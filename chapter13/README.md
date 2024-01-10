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