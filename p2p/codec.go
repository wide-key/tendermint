package p2p

import (
	amino "github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
)

var cdc *amino.Codec

func InitCdc() {
	cdc = amino.NewCodec()
	cryptoAmino.RegisterAmino(cdc)
}
