package store

import (
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/types"
)

var cdc *amino.Codec

func InitCdc() {
	cdc = amino.NewCodec()
	types.RegisterBlockAmino(cdc)
}
