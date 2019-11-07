package codec

import (
	"io"

	"github.com/coinexchain/codon"
)

func ShowInfo() {
	leafTypes := GetLeafTypes()

	//ShowInfo("",Account{})

	codon.ShowInfoForVar(leafTypes, DuplicateVoteEvidence{})
	codon.ShowInfoForVar(leafTypes, PrivKeyEd25519{})
	codon.ShowInfoForVar(leafTypes, PrivKeySecp256k1{})
	codon.ShowInfoForVar(leafTypes, PubKeyEd25519{})
	codon.ShowInfoForVar(leafTypes, PubKeySecp256k1{})
	codon.ShowInfoForVar(leafTypes, PubKeyMultisigThreshold{})
	codon.ShowInfoForVar(leafTypes, Vote{})
	codon.ShowInfoForVar(leafTypes, NewRoundStepMessage{})
	codon.ShowInfoForVar(leafTypes, NewValidBlockMessage{})
	codon.ShowInfoForVar(leafTypes, ProposalMessage{})
	codon.ShowInfoForVar(leafTypes, ProposalPOLMessage{})
	codon.ShowInfoForVar(leafTypes, BlockPartMessage{})
	codon.ShowInfoForVar(leafTypes, VoteMessage{})
	codon.ShowInfoForVar(leafTypes, HasVoteMessage{})
	codon.ShowInfoForVar(leafTypes, VoteSetMaj23Message{})
	codon.ShowInfoForVar(leafTypes, VoteSetBitsMessage{})
	codon.ShowInfoForVar(leafTypes, EventDataRoundState{})
	codon.ShowInfoForVar(leafTypes, EndHeightMessage{})
	codon.ShowInfoForVar(leafTypes, EventDataNewBlock{})
	codon.ShowInfoForVar(leafTypes, EventDataNewBlockHeader{})
	codon.ShowInfoForVar(leafTypes, EventDataTx{})
	codon.ShowInfoForVar(leafTypes, EventDataNewRound{})
	codon.ShowInfoForVar(leafTypes, EventDataCompleteProposal{})
	codon.ShowInfoForVar(leafTypes, EventDataVote{})
	codon.ShowInfoForVar(leafTypes, EventDataValidatorSetUpdates{})
	codon.ShowInfoForVar(leafTypes, EventDataString(""))
	codon.ShowInfoForVar(leafTypes, MockGoodEvidence{})
	codon.ShowInfoForVar(leafTypes, MockRandomGoodEvidence{})
	codon.ShowInfoForVar(leafTypes, MockBadEvidence{})
	codon.ShowInfoForVar(leafTypes, EvidenceListMessage{})
	codon.ShowInfoForVar(leafTypes, PubKeyRequest{})
	codon.ShowInfoForVar(leafTypes, PubKeyResponse{})
	codon.ShowInfoForVar(leafTypes, SignVoteRequest{})
	codon.ShowInfoForVar(leafTypes, SignedVoteResponse{})
	codon.ShowInfoForVar(leafTypes, SignProposalRequest{})
	codon.ShowInfoForVar(leafTypes, SignedProposalResponse{})
	codon.ShowInfoForVar(leafTypes, PingRequest{})
	codon.ShowInfoForVar(leafTypes, PingResponse{})
	codon.ShowInfoForVar(leafTypes, PacketPing{})
	codon.ShowInfoForVar(leafTypes, PacketPong{})
	codon.ShowInfoForVar(leafTypes, PacketMsg{})
	codon.ShowInfoForVar(leafTypes, TxMessage{})

	codon.ShowInfoForVar(leafTypes, Bc1BlockRequestMessage{})
	codon.ShowInfoForVar(leafTypes, Bc1BlockResponseMessage{})
	codon.ShowInfoForVar(leafTypes, Bc1NoBlockResponseMessage{})
	codon.ShowInfoForVar(leafTypes, Bc1StatusResponseMessage{})
	codon.ShowInfoForVar(leafTypes, Bc1StatusRequestMessage{})
	codon.ShowInfoForVar(leafTypes, Bc0BlockRequestMessage{})
	codon.ShowInfoForVar(leafTypes, Bc0BlockResponseMessage{})
	codon.ShowInfoForVar(leafTypes, Bc0NoBlockResponseMessage{})
	codon.ShowInfoForVar(leafTypes, Bc0StatusResponseMessage{})
	codon.ShowInfoForVar(leafTypes, Bc0StatusRequestMessage{})
	codon.ShowInfoForVar(leafTypes, PexRequestMessage{})
	codon.ShowInfoForVar(leafTypes, PexAddrsMessage{})
	codon.ShowInfoForVar(leafTypes, MsgInfo{})
	codon.ShowInfoForVar(leafTypes, TimeoutInfo{})
	codon.ShowInfoForVar(leafTypes, BitArray{})
	codon.ShowInfoForVar(leafTypes, Proposal{})
	codon.ShowInfoForVar(leafTypes, Part{})
	codon.ShowInfoForVar(leafTypes, Block{})
	codon.ShowInfoForVar(leafTypes, Tx{})
	codon.ShowInfoForVar(leafTypes, Commit{})
	codon.ShowInfoForVar(leafTypes, CommitSig{})
	codon.ShowInfoForVar(leafTypes, Validator{})
	codon.ShowInfoForVar(leafTypes, Event{})
	codon.ShowInfoForVar(leafTypes, ValidatorUpdate{})
	codon.ShowInfoForVar(leafTypes, ConsensusParams{})
	codon.ShowInfoForVar(leafTypes, BlockParams{})
	codon.ShowInfoForVar(leafTypes, EvidenceParams{})
	codon.ShowInfoForVar(leafTypes, ValidatorParams{})
	codon.ShowInfoForVar(leafTypes, RemoteSignerError{})
	codon.ShowInfoForVar(leafTypes, NetAddress{})
	codon.ShowInfoForVar(leafTypes, KVPair{})
	codon.ShowInfoForVar(leafTypes, P2PID(0))
	codon.ShowInfoForVar(leafTypes, Duration(0))
	codon.ShowInfoForVar(leafTypes, RoundStepType(0))
	codon.ShowInfoForVar(leafTypes, Protocol(0))
}

func GetLeafTypes() map[string]string {
	leafTypes := make(map[string]string, 20)
	leafTypes["time.Time"] = "time.Time"
	return leafTypes
}

var TypeEntryList = []codon.TypeEntry{
	{Alias: "PubKey", Value: (*PubKey)(nil)},
	{Alias: "PrivKey",              Value: (*PrivKey)(nil)},
	{Alias: "ConsensusMessage",     Value: (*ConsensusMessage)(nil)},
	{Alias: "WALMessage",           Value: (*WALMessage)(nil)},
	{Alias: "TMEventData",          Value: (*TMEventData)(nil)},
	{Alias: "Evidence",             Value: (*Evidence)(nil)},
	{Alias: "EvidenceMessage",      Value: (*EvidenceMessage)(nil)},
	{Alias: "SignerMessage",        Value: (*SignerMessage)(nil)},
	{Alias: "BlockchainMessageV1",  Value: (*BlockchainMessageV1)(nil)},
	{Alias: "BlockchainMessageV0",  Value: (*BlockchainMessageV0)(nil)},
	{Alias: "PexMessage",           Value: (*PexMessage)(nil)},
	{Alias: "Packet",               Value: (*Packet)(nil)},
	{Alias: "MempoolMessage",       Value: (*MempoolMessage)(nil)},

	{Alias: "DuplicateVoteEvidence", Value: DuplicateVoteEvidence{}},
	{Alias: "PrivKeyEd25519", Value: PrivKeyEd25519{}},
	{Alias: "PrivKeySecp256k1", Value: PrivKeySecp256k1{}},
	{Alias: "PubKeyEd25519", Value: PubKeyEd25519{}},
	{Alias: "PubKeySecp256k1", Value: PubKeySecp256k1{}},
	{Alias: "PubKeyMultisigThreshold", Value: PubKeyMultisigThreshold{}},
	{Alias: "SignedMsgType", Value: SignedMsgType(0)},
	{Alias: "Vote", Value: Vote{}},
	{Alias: "NewRoundStepMessage"         , Value: NewRoundStepMessage{}},
	{Alias: "NewValidBlockMessage"        , Value: NewValidBlockMessage{}},
	{Alias: "ProposalMessage"             , Value: ProposalMessage{}},
	{Alias: "ProposalPOLMessage"          , Value: ProposalPOLMessage{}},
	{Alias: "BlockPartMessage"            , Value: BlockPartMessage{}},
	{Alias: "VoteMessage"                 , Value: VoteMessage{}},
	{Alias: "HasVoteMessage"              , Value: HasVoteMessage{}},
	{Alias: "VoteSetMaj23Message"         , Value: VoteSetMaj23Message{}},
	{Alias: "VoteSetBitsMessage"          , Value: VoteSetBitsMessage{}},
	{Alias: "EventDataRoundState"         , Value: EventDataRoundState{}},
	{Alias: "EndHeightMessage"            , Value: EndHeightMessage{}},
	//{Alias: "EventDataNewBlock"           , Value: EventDataNewBlock{}},
	{Alias: "EventDataNewBlockHeader"     , Value: EventDataNewBlockHeader{}},
	{Alias: "EventDataTx"                 , Value: EventDataTx{}},
	{Alias: "EventDataNewRound"           , Value: EventDataNewRound{}},
	{Alias: "EventDataCompleteProposal"   , Value: EventDataCompleteProposal{}},
	{Alias: "EventDataVote"               , Value: EventDataVote{}},
	{Alias: "EventDataValidatorSetUpdates", Value: EventDataValidatorSetUpdates{}},
	{Alias: "EventDataString"             , Value: EventDataString("")},
	{Alias: "MockGoodEvidence"            , Value: MockGoodEvidence{}},
	//{Alias: "MockRandomGoodEvidence"      , Value: MockRandomGoodEvidence{}},
	{Alias: "MockBadEvidence"             , Value: MockBadEvidence{}},
	{Alias: "EvidenceListMessage"         , Value: EvidenceListMessage{}},
	//{Alias: "PubKeyRequest"               , Value: PubKeyRequest{}},
	{Alias: "PubKeyResponse"              , Value: PubKeyResponse{}},
	{Alias: "SignVoteRequest"             , Value: SignVoteRequest{}},
	{Alias: "SignedVoteResponse"          , Value: SignedVoteResponse{}},
	{Alias: "SignProposalRequest"         , Value: SignProposalRequest{}},
	{Alias: "SignedProposalResponse"      , Value: SignedProposalResponse{}},
	//{Alias: "PingRequest"                 , Value: PingRequest{}},
	//{Alias: "PingResponse"                , Value: PingResponse{}},
	//{Alias: "PacketPing"                  , Value: PacketPing{}},
	//{Alias: "PacketPong"                  , Value: PacketPong{}},
	{Alias: "PacketMsg"                   , Value: PacketMsg{}},
	{Alias: "TxMessage"                   , Value: TxMessage{}},

	{Alias: "Bc1BlockRequestMessage"    , Value: Bc1BlockRequestMessage{}},
	//{Alias: "Bc1BlockResponseMessage"   , Value: Bc1BlockResponseMessage{}},
	{Alias: "Bc1NoBlockResponseMessage" , Value: Bc1NoBlockResponseMessage{}},
	{Alias: "Bc1StatusResponseMessage"  , Value: Bc1StatusResponseMessage{}},
	{Alias: "Bc1StatusRequestMessage"   , Value: Bc1StatusRequestMessage{}},
	{Alias: "Bc0BlockRequestMessage"    , Value: Bc0BlockRequestMessage{}},
	//{Alias: "Bc0BlockResponseMessage"   , Value: Bc0BlockResponseMessage{}},
	{Alias: "Bc0NoBlockResponseMessage" , Value: Bc0NoBlockResponseMessage{}},
	{Alias: "Bc0StatusResponseMessage"  , Value: Bc0StatusResponseMessage{}},
	{Alias: "Bc0StatusRequestMessage"   , Value: Bc0StatusRequestMessage{}},
	//{Alias: "PexRequestMessage"         , Value: PexRequestMessage{}},
	//{Alias: "PexAddrsMessage"           , Value: PexAddrsMessage{}},
	{Alias: "MsgInfo"                   , Value: MsgInfo{}},
	{Alias: "TimeoutInfo"               , Value: TimeoutInfo{}},
	{Alias: "RoundStepType"             , Value: RoundStepType(0)},
	{Alias: "BitArray"             , Value: BitArray{}},
	{Alias: "Proposal"             , Value: Proposal{}},
	{Alias: "Part"             , Value: Part{}},
	//{Alias: "Block"             , Value: Block{}},
	{Alias: "Protocol"             , Value: Protocol(0)},
	{Alias: "Tx"             , Value: Tx{}},
	//{Alias: "Commit"             , Value: Commit{}},
	{Alias: "CommitSig"             , Value: CommitSig{}},
	{Alias: "Validator"             , Value: Validator{}},
	{Alias: "Event"             , Value: Event{}},
	{Alias: "ValidatorUpdate"             , Value: ValidatorUpdate{}},
	{Alias: "ConsensusParams"             , Value: ConsensusParams{}},
	{Alias: "BlockParams"             , Value: BlockParams{}},
	{Alias: "EvidenceParams"             , Value: EvidenceParams{}},
	{Alias: "ValidatorParams"             , Value: ValidatorParams{}},
	{Alias: "RemoteSignerError"             , Value: RemoteSignerError{}},
	//{Alias: "NetAddress"             , Value: NetAddress{}},
	{Alias: "P2PID"             , Value: P2PID(0)},
	{Alias: "Duration"             , Value: Duration(0)},
	{Alias: "KVPair"             , Value: KVPair{}},

}

func GenerateCodecFile(w io.Writer) {
	extraImports := []string{`"time"`}
	extraImports = append(extraImports, codon.ImportsForBridgeLogic...)
	extraLogics := extraLogicsForLeafTypes + codon.BridgeLogic
	ignoreImpl := make(map[string]string)
	ignoreImpl["PubKeyMultisigThreshold"] = "PubKey"
	codon.GenerateCodecFile(w, GetLeafTypes(), ignoreImpl, TypeEntryList, extraLogics, extraImports)
}


const MaxSliceLength = 10
const MaxStringLength = 100

var extraLogicsForLeafTypes = `

func EncodeTime(w *[]byte, t time.Time) {
	t = t.UTC()
	sec := t.Unix()
	var buf [10]byte
	n := binary.PutVarint(buf[:], sec)
	*w = append(*w, buf[0:n]...)

	nanosec := t.Nanosecond()
	n = binary.PutVarint(buf[:], int64(nanosec))
	*w = append(*w, buf[0:n]...)
}

func DecodeTime(bz []byte) (time.Time, int, error) {
	sec, n := binary.Varint(bz)
	var err error
	if n == 0 {
		// buf too small
		err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		err = errors.New("EOF decoding varint")
	}
	if err!=nil {
		return time.Unix(sec,0), n, err
	}

	nanosec, m := binary.Varint(bz[n:])
	if m == 0 {
		// buf too small
		err = errors.New("buffer too small")
	} else if m < 0 {
		// value larger than 64 bits (overflow)
		// and -m is the number of bytes read
		m = -m
		err = errors.New("EOF decoding varint")
	}
	if err!=nil {
		return time.Unix(sec,nanosec), n+m, err
	}

	return time.Unix(sec, nanosec).UTC(), n+m, nil
}

func T(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

var maxSec = T("9999-09-29T08:02:06.647266Z").Unix()

func RandTime(r RandSrc) time.Time {
	sec := r.GetInt64()
	nanosec := r.GetInt64()
	if sec < 0 {
		sec = -sec
	}
	if nanosec < 0 {
		nanosec = -nanosec
	}
	nanosec = nanosec%(1000*1000*1000)
	sec = sec%maxSec
	return time.Unix(sec, nanosec).UTC()
}


func DeepCopyTime(t time.Time) time.Time {
	return t.Add(time.Duration(0))
}
`

