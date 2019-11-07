package codec

import (
	"time"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/consensus"
	ctypes "github.com/tendermint/tendermint/consensus/types"
	"github.com/tendermint/tendermint/evidence"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/blockchain/v0"
	"github.com/tendermint/tendermint/blockchain/v1"
	"github.com/tendermint/tendermint/p2p/pex"
	"github.com/tendermint/tendermint/p2p/conn"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/mempool"
	"github.com/tendermint/tendermint/version"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
)

type (
	PubKey  = crypto.PubKey
	PrivKey  = crypto.PrivKey

	ConsensusMessage = consensus.ConsensusMessage
	WALMessage = consensus.WALMessage

	TMEventData = tmtypes.TMEventData
	Evidence = tmtypes.Evidence

	EvidenceMessage = evidence.EvidenceMessage
	SignerMessage = privval.SignerMessage

	BlockchainMessageV1 = v1.BlockchainMessage
	BlockchainMessageV0 = v0.BlockchainMessage

	PexMessage = pex.PexMessage
	Packet = conn.Packet

	MempoolMessage = mempool.MempoolMessage
)

type (
	DuplicateVoteEvidence   = tmtypes.DuplicateVoteEvidence
	PrivKeyEd25519          = ed25519.PrivKeyEd25519
	PrivKeySecp256k1        = secp256k1.PrivKeySecp256k1
	PubKeyEd25519           = ed25519.PubKeyEd25519
	PubKeySecp256k1         = secp256k1.PubKeySecp256k1
	PubKeyMultisigThreshold = multisig.PubKeyMultisigThreshold
	SignedMsgType           = tmtypes.SignedMsgType
	Vote                    = tmtypes.Vote

	NewRoundStepMessage          = consensus.NewRoundStepMessage
	NewValidBlockMessage         = consensus.NewValidBlockMessage
	ProposalMessage              = consensus.ProposalMessage
	ProposalPOLMessage           = consensus.ProposalPOLMessage
	BlockPartMessage             = consensus.BlockPartMessage
	VoteMessage                  = consensus.VoteMessage
	HasVoteMessage               = consensus.HasVoteMessage
	VoteSetMaj23Message          = consensus.VoteSetMaj23Message
	VoteSetBitsMessage           = consensus.VoteSetBitsMessage
	EndHeightMessage             = consensus.EndHeightMessage
	EventDataNewBlock            = tmtypes.EventDataNewBlock
	EventDataNewBlockHeader      = tmtypes.EventDataNewBlockHeader
	EventDataTx                  = tmtypes.EventDataTx
	EventDataRoundState          = tmtypes.EventDataRoundState
	EventDataNewRound            = tmtypes.EventDataNewRound
	EventDataCompleteProposal    = tmtypes.EventDataCompleteProposal
	EventDataVote                = tmtypes.EventDataVote
	EventDataValidatorSetUpdates = tmtypes.EventDataValidatorSetUpdates
	EventDataString              = tmtypes.EventDataString
	MockGoodEvidence             = tmtypes.MockGoodEvidence
	MockRandomGoodEvidence       = tmtypes.MockRandomGoodEvidence
	MockBadEvidence              = tmtypes.MockBadEvidence
	EvidenceListMessage          = evidence.EvidenceListMessage
	PubKeyRequest                = privval.PubKeyRequest
	PubKeyResponse               = privval.PubKeyResponse
	SignVoteRequest              = privval.SignVoteRequest
	SignedVoteResponse           = privval.SignedVoteResponse
	SignProposalRequest          = privval.SignProposalRequest
	SignedProposalResponse       = privval.SignedProposalResponse
	PingRequest                  = privval.PingRequest
	PingResponse                 = privval.PingResponse
	PacketPing                   = conn.PacketPing
	PacketPong                   = conn.PacketPong
	PacketMsg                    = conn.PacketMsg
	TxMessage                    = mempool.TxMessage

	Bc1BlockRequestMessage    = v1.BcBlockRequestMessage
	Bc1BlockResponseMessage   = v1.BcBlockResponseMessage
	Bc1NoBlockResponseMessage = v1.BcNoBlockResponseMessage
	Bc1StatusResponseMessage  = v1.BcStatusResponseMessage
	Bc1StatusRequestMessage   = v1.BcStatusRequestMessage
	Bc0BlockRequestMessage    = v0.BcBlockRequestMessage
	Bc0BlockResponseMessage   = v0.BcBlockResponseMessage
	Bc0NoBlockResponseMessage = v0.BcNoBlockResponseMessage
	Bc0StatusResponseMessage  = v0.BcStatusResponseMessage
	Bc0StatusRequestMessage   = v0.BcStatusRequestMessage

	PexRequestMessage = pex.PexRequestMessage
	PexAddrsMessage   = pex.PexAddrsMessage

	MsgInfo = consensus.MsgInfo
	TimeoutInfo = consensus.TimeoutInfo
	RoundStepType = ctypes.RoundStepType
	BitArray = common.BitArray
	Proposal = tmtypes.Proposal
	Part = tmtypes.Part
	Block = tmtypes.Block
	Protocol = version.Protocol
	Tx = tmtypes.Tx
	Commit = tmtypes.Commit
	CommitSig = tmtypes.CommitSig
	Validator = tmtypes.Validator
	Event = abcitypes.Event
	ValidatorUpdate = abcitypes.ValidatorUpdate
	ConsensusParams = abcitypes.ConsensusParams
	BlockParams = abcitypes.BlockParams
	EvidenceParams = abcitypes.EvidenceParams
	ValidatorParams = abcitypes.ValidatorParams
	RemoteSignerError = privval.RemoteSignerError
	NetAddress = p2p.NetAddress
	P2PID = p2p.ID
	Duration = time.Duration
	KVPair = common.KVPair
)

