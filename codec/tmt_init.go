package codec

import (
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/merkle"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/types"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/p2p/conn"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/state/txindex/kv"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/store"
	"github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/consensus"
	"github.com/tendermint/tendermint/blockchain/v0"
	"github.com/tendermint/tendermint/blockchain/v1"
	ctypes "github.com/tendermint/tendermint/consensus/types"

)

func Init() {
	amino.Stub = &CodonStub{}
	InitCdc()
}

func InitCdc() {
	cryptoAmino.InitCdc()
	conn.InitCdc()
	p2p.InitCdc()
	multisig.InitCdc()
	merkle.InitCdc()
	secp256k1.InitCdc()
	ed25519.InitCdc()
	types.InitCdc()
	mempool.InitCdc()
	state.InitCdc()
	node.InitCdc()
	store.InitCdc()
	client.InitCdc()
	privval.InitCdc()
	consensus.InitCdc()
	ctypes.InitCdc()
	v0.InitCdc()
	v1.InitCdc()
	kv.InitCdc()
}

