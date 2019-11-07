package merkle

import (
	amino "github.com/tendermint/go-amino"
)

var cdc *amino.Codec

func InitCdc() {
	cdc = amino.NewCodec()
	cdc.Seal()
}
