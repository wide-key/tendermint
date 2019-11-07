package pex

import (
	amino "github.com/tendermint/go-amino"
)

var cdc *amino.Codec

func InitCdc() {
	cdc = amino.NewCodec()
	RegisterPexMessage(cdc)
}
