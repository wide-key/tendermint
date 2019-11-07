//nolint
package codec
import (
"time"
"io"
"fmt"
"reflect"
amino "github.com/tendermint/go-amino"
"encoding/binary"
"math"
"errors"
)

type RandSrc interface {
	GetBool() bool
	GetInt() int
	GetInt8() int8
	GetInt16() int16
	GetInt32() int32
	GetInt64() int64
	GetUint() uint
	GetUint8() uint8
	GetUint16() uint16
	GetUint32() uint32
	GetUint64() uint64
	GetFloat32() float32
	GetFloat64() float64
	GetString(n int) string
	GetBytes(n int) []byte
}

func codonEncodeBool(w *[]byte, v bool) {
	if v {
		*w = append(*w, byte(1))
	} else {
		*w = append(*w, byte(0))
	}
}
func codonEncodeVarint(w *[]byte, v int64) {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], v)
	*w = append(*w, buf[0:n]...)
}
func codonEncodeInt8(w *[]byte, v int8) {
	*w = append(*w, byte(v))
}
func codonEncodeInt16(w *[]byte, v int16) {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], uint16(v))
	*w = append(*w, buf[:]...)
}
func codonEncodeUvarint(w *[]byte, v uint64) {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], v)
	*w = append(*w, buf[0:n]...)
}
func codonEncodeUint8(w *[]byte, v uint8) {
	*w = append(*w, byte(v))
}
func codonEncodeUint16(w *[]byte, v uint16) {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], v)
	*w = append(*w, buf[:]...)
}
func codonEncodeFloat32(w *[]byte, v float32) {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], math.Float32bits(v))
	*w = append(*w, buf[:]...)
}
func codonEncodeFloat64(w *[]byte, v float64) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], math.Float64bits(v))
	*w = append(*w, buf[:]...)
}
func codonEncodeByteSlice(w *[]byte, v []byte) {
	codonEncodeVarint(w, int64(len(v)))
	*w = append(*w, v...)
}
func codonEncodeString(w *[]byte, v string) {
	codonEncodeByteSlice(w, []byte(v))
}
func codonDecodeBool(bz []byte, n *int, err *error) bool {
	if len(bz) < 1 {
		*err = errors.New("Not enough bytes to read")
		return false
	}
	*n = 1
	*err = nil
	return bz[0]!=0
}
func codonDecodeInt(bz []byte, m *int, err *error) int {
	i, n := binary.Varint(bz)
	if n == 0 {
		// buf too small
		*err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		*err = errors.New("EOF decoding varint")
	}
	*m = n
	return int(i)
}
func codonDecodeInt8(bz []byte, n *int, err *error) int8 {
	if len(bz) < 1 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*err = nil
	*n = 1
	return int8(bz[0])
}
func codonDecodeInt16(bz []byte, n *int, err *error) int16 {
	if len(bz) < 2 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 2
	*err = nil
	return int16(binary.LittleEndian.Uint16(bz[:2]))
}
func codonDecodeInt32(bz []byte, n *int, err *error) int32 {
	i := codonDecodeInt64(bz, n, err)
	return int32(i)
}
func codonDecodeInt64(bz []byte, m *int, err *error) int64 {
	i, n := binary.Varint(bz)
	if n == 0 {
		// buf too small
		*err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		*err = errors.New("EOF decoding varint")
	}
	*m = n
	*err = nil
	return int64(i)
}
func codonDecodeUint(bz []byte, n *int, err *error) uint {
	i := codonDecodeUint64(bz, n, err)
	return uint(i)
}
func codonDecodeUint8(bz []byte, n *int, err *error) uint8 {
	if len(bz) < 1 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 1
	*err = nil
	return uint8(bz[0])
}
func codonDecodeUint16(bz []byte, n *int, err *error) uint16 {
	if len(bz) < 2 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 2
	*err = nil
	return uint16(binary.LittleEndian.Uint16(bz[:2]))
}
func codonDecodeUint32(bz []byte, n *int, err *error) uint32 {
	i := codonDecodeUint64(bz, n, err)
	return uint32(i)
}
func codonDecodeUint64(bz []byte, m *int, err *error) uint64 {
	i, n := binary.Uvarint(bz)
	if n == 0 {
		// buf too small
		*err = errors.New("buffer too small")
	} else if n < 0 {
		// value larger than 64 bits (overflow)
		// and -n is the number of bytes read
		n = -n
		*err = errors.New("EOF decoding varint")
	}
	*m = n
	*err = nil
	return uint64(i)
}
func codonDecodeFloat64(bz []byte, n *int, err *error) float64 {
	if len(bz) < 8 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 8
	*err = nil
	i := binary.LittleEndian.Uint64(bz[:8])
	return math.Float64frombits(i)
}
func codonDecodeFloat32(bz []byte, n *int, err *error) float32 {
	if len(bz) < 4 {
		*err = errors.New("Not enough bytes to read")
		return 0
	}
	*n = 4
	*err = nil
	i := binary.LittleEndian.Uint32(bz[:4])
	return math.Float32frombits(i)
}
func codonGetByteSlice(bz []byte, length int) ([]byte, int, error) {
	if len(bz) < length {
		return nil, 0, errors.New("Not enough bytes to read")
	}
	return bz[:length], length, nil
}
func codonDecodeString(bz []byte, n *int, err *error) string {
	var m int
	length := codonDecodeInt64(bz, &m, err)
	if *err != nil {
		return ""
	}
	var bs []byte
	var l int
	bs, l, *err = codonGetByteSlice(bz[m:], int(length))
	*n = m + l
	return string(bs)
}


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

// ========= BridgeBegin ============
type CodecImp struct {
	sealed          bool
	structPath2Name map[string]string
}

func (cdc *CodecImp) MarshalBinaryBare(o interface{}) ([]byte, error) {
	s := CodonStub{}
	return s.MarshalBinaryBare(o)
}
func (cdc *CodecImp) MarshalBinaryLengthPrefixed(o interface{}) ([]byte, error) {
	s := CodonStub{}
	return s.MarshalBinaryLengthPrefixed(o)
}
func (cdc *CodecImp) MarshalBinaryLengthPrefixedWriter(w io.Writer, o interface{}) (n int64, err error) {
	bz, err := cdc.MarshalBinaryLengthPrefixed(o)
	m, err := w.Write(bz)
	return int64(m), err
}
func (cdc *CodecImp) UnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	s := CodonStub{}
	return s.UnmarshalBinaryBare(bz, ptr)
}
func (cdc *CodecImp) UnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error {
	s := CodonStub{}
	return s.UnmarshalBinaryLengthPrefixed(bz, ptr)
}
func (cdc *CodecImp) UnmarshalBinaryLengthPrefixedReader(r io.Reader, ptr interface{}, maxSize int64) (n int64, err error) {
	if maxSize < 0 {
		panic("maxSize cannot be negative.")
	}

	// Read byte-length prefix.
	var l int64
	var buf [binary.MaxVarintLen64]byte
	for i := 0; i < len(buf); i++ {
		_, err = r.Read(buf[i : i+1])
		if err != nil {
			return
		}
		n += 1
		if buf[i]&0x80 == 0 {
			break
		}
		if n >= maxSize {
			err = fmt.Errorf("Read overflow, maxSize is %v but uvarint(length-prefix) is itself greater than maxSize.", maxSize)
		}
	}
	u64, _ := binary.Uvarint(buf[:])
	if err != nil {
		return
	}
	if maxSize > 0 {
		if uint64(maxSize) < u64 {
			err = fmt.Errorf("Read overflow, maxSize is %v but this amino binary object is %v bytes.", maxSize, u64)
			return
		}
		if (maxSize - n) < int64(u64) {
			err = fmt.Errorf("Read overflow, maxSize is %v but this length-prefixed amino binary object is %v+%v bytes.", maxSize, n, u64)
			return
		}
	}
	l = int64(u64)
	if l < 0 {
		err = fmt.Errorf("Read overflow, this implementation can't read this because, why would anyone have this much data?")
	}

	// Read that many bytes.
	var bz = make([]byte, l, l)
	_, err = io.ReadFull(r, bz)
	if err != nil {
		return
	}
	n += l

	// Decode.
	err = cdc.UnmarshalBinaryBare(bz, ptr)
	return
}

//------

func (cdc *CodecImp) MustMarshalBinaryBare(o interface{}) []byte {
	bz, err := cdc.MarshalBinaryBare(o)
	if err!=nil {
		panic(err)
	}
	return bz
}
func (cdc *CodecImp) MustMarshalBinaryLengthPrefixed(o interface{}) []byte {
	bz, err := cdc.MarshalBinaryLengthPrefixed(o)
	if err!=nil {
		panic(err)
	}
	return bz
}
func (cdc *CodecImp) MustUnmarshalBinaryBare(bz []byte, ptr interface{}) {
	err := cdc.UnmarshalBinaryBare(bz, ptr)
	if err!=nil {
		panic(err)
	}
}
func (cdc *CodecImp) MustUnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) {
	err := cdc.UnmarshalBinaryLengthPrefixed(bz, ptr)
	if err!=nil {
		panic(err)
	}
}

// ====================
func derefPtr(v interface{}) reflect.Type {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func (cdc *CodecImp) PrintTypes(out io.Writer) error {
	for _, entry := range GetSupportList() {
		_, err := out.Write([]byte(entry))
		if err != nil {
			return err
		}
		_, err = out.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}
	return nil
}
func (cdc *CodecImp) RegisterConcrete(o interface{}, name string, copts *amino.ConcreteOptions) {
	if cdc.sealed {
		panic("Codec is already sealed")
	}
	t := derefPtr(o)
	path := t.PkgPath() + "." + t.Name()
	cdc.structPath2Name[path] = name
}
func (cdc *CodecImp) RegisterInterface(_ interface{}, _ *amino.InterfaceOptions) {
	if cdc.sealed {
		panic("Codec is already sealed")
	}
	//Nop
}
func (cdc *CodecImp) SealImp() amino.CodecIfc {
	if cdc.sealed {
		panic("Codec is already sealed")
	}
	cdc.sealed = true
	return cdc
}

// ========================================

type CodonStub struct {
}

func (_ *CodonStub) NewCodecImp() amino.CodecIfc {
	return &CodecImp{
		structPath2Name: make(map[string]string),
	}
}
func (_ *CodonStub) DeepCopy(o interface{}) (r interface{}) {
	r = DeepCopyAny(o)
	return
}

func (_ *CodonStub) MarshalBinaryBare(o interface{}) ([]byte, error) {
	if _, err := getMagicBytesOfVar(o); err!=nil {
		return nil, err
	}
	buf := make([]byte, 0, 1024)
	EncodeAny(&buf, o)
	return buf, nil
}
func (s *CodonStub) MarshalBinaryLengthPrefixed(o interface{}) ([]byte, error) {
	if _, err := getMagicBytesOfVar(o); err!=nil {
		return nil, err
	}
	bz, err := s.MarshalBinaryBare(o)
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], uint64(len(bz)))
	return append(buf[:n], bz...), err
}
func (_ *CodonStub) UnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	if _, err := getMagicBytesOfVar(ptr); err!=nil {
		return err
	}
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		panic("Unmarshal expects a pointer")
	}
	if len(bz) <= 4 {
		return fmt.Errorf("Byte slice is too short: %d", len(bz))
	}
	magicBytes, err := getMagicBytesOfVar(ptr)
	if err!=nil {
		return err
	}
	if bz[0]!=magicBytes[0] || bz[1]!=magicBytes[1] || bz[2]!=magicBytes[2] || bz[3]!=magicBytes[3] {
		return fmt.Errorf("MagicBytes Missmatch %v vs %v", bz[0:4], magicBytes[:])
	}
	o, _, err := DecodeAny(bz)
	rv.Elem().Set(reflect.ValueOf(o))
	return err
}
func (s *CodonStub) UnmarshalBinaryLengthPrefixed(bz []byte, ptr interface{}) error {
	if _, err := getMagicBytesOfVar(ptr); err!=nil {
		return err
	}
	if len(bz) == 0 {
		return errors.New("UnmarshalBinaryLengthPrefixed cannot decode empty bytes")
	}
	// Read byte-length prefix.
	u64, n := binary.Uvarint(bz)
	if n < 0 {
		return fmt.Errorf("Error reading msg byte-length prefix: got code %v", n)
	}
	if u64 > uint64(len(bz)-n) {
		return fmt.Errorf("Not enough bytes to read in UnmarshalBinaryLengthPrefixed, want %v more bytes but only have %v",
			u64, len(bz)-n)
	} else if u64 < uint64(len(bz)-n) {
		return fmt.Errorf("Bytes left over in UnmarshalBinaryLengthPrefixed, should read %v more bytes but have %v",
			u64, len(bz)-n)
	}
	bz = bz[n:]
	return s.UnmarshalBinaryBare(bz, ptr)
}
func (s *CodonStub) MustMarshalBinaryLengthPrefixed(o interface{}) []byte {
	bz, err := s.MarshalBinaryLengthPrefixed(o)
	if err!=nil {
		panic(err)
	}
	return bz
}

// ========================================
func (_ *CodonStub) UvarintSize(u uint64) int {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutUvarint(buf[:], u)
	return n
}
func (_ *CodonStub) EncodeByteSlice(w io.Writer, bz []byte) error {
	buf := make([]byte, 0, binary.MaxVarintLen64+len(bz))
	codonEncodeByteSlice(&buf, bz)
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeUvarint(w io.Writer, u uint64) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeUvarint(&buf, u)
	_, err := w.Write(buf)
	return err
}
func (s *CodonStub) ByteSliceSize(bz []byte) int {
	return s.UvarintSize(uint64(len(bz))) + len(bz)
}
func (_ *CodonStub) EncodeInt8(w io.Writer, i int8) error {
	_, err := w.Write([]byte{byte(i)})
	return err
}
func (_ *CodonStub) EncodeInt16(w io.Writer, i int16) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeInt16(&buf, i)
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeInt32(w io.Writer, i int32) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeVarint(&buf, int64(i))
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeInt64(w io.Writer, i int64) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeVarint(&buf, i)
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeVarint(w io.Writer, i int64) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeVarint(&buf, i)
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeByte(w io.Writer, b byte) error {
	_, err := w.Write([]byte{b})
	return err
}
func (_ *CodonStub) EncodeUint8(w io.Writer, u uint8) error {
	_, err := w.Write([]byte{u})
	return err
}
func (_ *CodonStub) EncodeUint16(w io.Writer, u uint16) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeUint16(&buf, u)
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeUint32(w io.Writer, u uint32) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeUvarint(&buf, uint64(u))
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeUint64(w io.Writer, u uint64) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeUvarint(&buf, u)
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeBool(w io.Writer, b bool) error {
	u := byte(0)
	if b {
		u = byte(1)
	}
	_, err := w.Write([]byte{u})
	return err
}
func (_ *CodonStub) EncodeFloat32(w io.Writer, f float32) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeFloat32(&buf, f)
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeFloat64(w io.Writer, f float64) error {
	buf := make([]byte, 0, binary.MaxVarintLen64)
	codonEncodeFloat64(&buf, f)
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) EncodeString(w io.Writer, s string) error {
	buf := make([]byte, 0, binary.MaxVarintLen64+len(s))
	codonEncodeString(&buf, s)
	_, err := w.Write(buf)
	return err
}
func (_ *CodonStub) DecodeInt8(bz []byte) (i int8, n int, err error) {
	i = codonDecodeInt8(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeInt16(bz []byte) (i int16, n int, err error) {
	i = codonDecodeInt16(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeInt32(bz []byte) (i int32, n int, err error) {
	i = codonDecodeInt32(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeInt64(bz []byte) (i int64, n int, err error) {
	i = codonDecodeInt64(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeVarint(bz []byte) (i int64, n int, err error) {
	i = codonDecodeInt64(bz, &n, &err)
	return
}
func (s *CodonStub) DecodeByte(bz []byte) (b byte, n int, err error) {
	b = codonDecodeUint8(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUint8(bz []byte) (u uint8, n int, err error) {
	u = codonDecodeUint8(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUint16(bz []byte) (u uint16, n int, err error) {
	u = codonDecodeUint16(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUint32(bz []byte) (u uint32, n int, err error) {
	u = codonDecodeUint32(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUint64(bz []byte) (u uint64, n int, err error) {
	u = codonDecodeUint64(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeUvarint(bz []byte) (u uint64, n int, err error) {
	u = codonDecodeUint64(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeBool(bz []byte) (b bool, n int, err error) {
	b = codonDecodeBool(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeFloat32(bz []byte) (f float32, n int, err error) {
	f = codonDecodeFloat32(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeFloat64(bz []byte) (f float64, n int, err error) {
	f = codonDecodeFloat64(bz, &n, &err)
	return
}
func (_ *CodonStub) DecodeByteSlice(bz []byte) (bz2 []byte, n int, err error) {
	length := codonDecodeInt(bz, &n, &err)
	if err != nil {
		return
	}
	bz = bz[n:]
	n += length
	bz2, m, err := codonGetByteSlice(bz, length)
	n += m
	return
}
func (_ *CodonStub) DecodeString(bz []byte) (s string, n int, err error) {
	s = codonDecodeString(bz, &n, &err)
	return
}
func (_ *CodonStub) VarintSize(i int64) int {
	var buf [binary.MaxVarintLen64]byte
	n := binary.PutVarint(buf[:], i)
	return n
}
// ========= BridgeEnd ============
// Non-Interface
func EncodeDuplicateVoteEvidence(w *[]byte, v DuplicateVoteEvidence) {
EncodePubKey(w, v.PubKey) // interface_encode
codonEncodeUint8(w, uint8(v.VoteA.Type))
codonEncodeVarint(w, int64(v.VoteA.Height))
codonEncodeVarint(w, int64(v.VoteA.Round))
codonEncodeByteSlice(w, v.VoteA.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.VoteA.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.VoteA.BlockID.PartsHeader.Hash[:])
// end of v.VoteA.BlockID.PartsHeader
// end of v.VoteA.BlockID
EncodeTime(w, v.VoteA.Timestamp)
codonEncodeByteSlice(w, v.VoteA.ValidatorAddress[:])
codonEncodeVarint(w, int64(v.VoteA.ValidatorIndex))
codonEncodeByteSlice(w, v.VoteA.Signature[:])
// end of v.VoteA
codonEncodeUint8(w, uint8(v.VoteB.Type))
codonEncodeVarint(w, int64(v.VoteB.Height))
codonEncodeVarint(w, int64(v.VoteB.Round))
codonEncodeByteSlice(w, v.VoteB.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.VoteB.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.VoteB.BlockID.PartsHeader.Hash[:])
// end of v.VoteB.BlockID.PartsHeader
// end of v.VoteB.BlockID
EncodeTime(w, v.VoteB.Timestamp)
codonEncodeByteSlice(w, v.VoteB.ValidatorAddress[:])
codonEncodeVarint(w, int64(v.VoteB.ValidatorIndex))
codonEncodeByteSlice(w, v.VoteB.Signature[:])
// end of v.VoteB
} //End of EncodeDuplicateVoteEvidence

func DecodeDuplicateVoteEvidence(bz []byte) (DuplicateVoteEvidence, int, error) {
var err error
var length int
var v DuplicateVoteEvidence
var n int
var total int
v.PubKey, n, err = DecodePubKey(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n // interface_decode
v.VoteA = &Vote{}
v.VoteA.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteA.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteA.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteA.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteA.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteA.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.VoteA.BlockID.PartsHeader
// end of v.VoteA.BlockID
v.VoteA.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteA.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteA.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteA.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.VoteA
v.VoteB = &Vote{}
v.VoteB.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteB.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteB.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteB.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteB.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteB.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.VoteB.BlockID.PartsHeader
// end of v.VoteB.BlockID
v.VoteB.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteB.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteB.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.VoteB.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.VoteB
return v, total, nil
} //End of DecodeDuplicateVoteEvidence

func RandDuplicateVoteEvidence(r RandSrc) DuplicateVoteEvidence {
var length int
var v DuplicateVoteEvidence
v.PubKey = RandPubKey(r) // interface_decode
v.VoteA = &Vote{}
v.VoteA.Type = SignedMsgType(r.GetUint8())
v.VoteA.Height = r.GetInt64()
v.VoteA.Round = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.VoteA.BlockID.Hash = r.GetBytes(length)
v.VoteA.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.VoteA.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.VoteA.BlockID.PartsHeader
// end of v.VoteA.BlockID
v.VoteA.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.VoteA.ValidatorAddress = r.GetBytes(length)
v.VoteA.ValidatorIndex = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.VoteA.Signature = r.GetBytes(length)
// end of v.VoteA
v.VoteB = &Vote{}
v.VoteB.Type = SignedMsgType(r.GetUint8())
v.VoteB.Height = r.GetInt64()
v.VoteB.Round = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.VoteB.BlockID.Hash = r.GetBytes(length)
v.VoteB.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.VoteB.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.VoteB.BlockID.PartsHeader
// end of v.VoteB.BlockID
v.VoteB.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.VoteB.ValidatorAddress = r.GetBytes(length)
v.VoteB.ValidatorIndex = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.VoteB.Signature = r.GetBytes(length)
// end of v.VoteB
return v
} //End of RandDuplicateVoteEvidence

func DeepCopyDuplicateVoteEvidence(in DuplicateVoteEvidence) (out DuplicateVoteEvidence) {
var length int
out.PubKey = DeepCopyPubKey(in.PubKey)
out.VoteA = &Vote{}
out.VoteA.Type = in.VoteA.Type
out.VoteA.Height = in.VoteA.Height
out.VoteA.Round = in.VoteA.Round
length = len(in.VoteA.BlockID.Hash)
out.VoteA.BlockID.Hash = make([]uint8, length)
copy(out.VoteA.BlockID.Hash[:], in.VoteA.BlockID.Hash[:])
out.VoteA.BlockID.PartsHeader.Total = in.VoteA.BlockID.PartsHeader.Total
length = len(in.VoteA.BlockID.PartsHeader.Hash)
out.VoteA.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.VoteA.BlockID.PartsHeader.Hash[:], in.VoteA.BlockID.PartsHeader.Hash[:])
// end of .VoteA.BlockID.PartsHeader
// end of .VoteA.BlockID
out.VoteA.Timestamp = DeepCopyTime(in.VoteA.Timestamp)
length = len(in.VoteA.ValidatorAddress)
out.VoteA.ValidatorAddress = make([]uint8, length)
copy(out.VoteA.ValidatorAddress[:], in.VoteA.ValidatorAddress[:])
out.VoteA.ValidatorIndex = in.VoteA.ValidatorIndex
length = len(in.VoteA.Signature)
out.VoteA.Signature = make([]uint8, length)
copy(out.VoteA.Signature[:], in.VoteA.Signature[:])
// end of .VoteA
out.VoteB = &Vote{}
out.VoteB.Type = in.VoteB.Type
out.VoteB.Height = in.VoteB.Height
out.VoteB.Round = in.VoteB.Round
length = len(in.VoteB.BlockID.Hash)
out.VoteB.BlockID.Hash = make([]uint8, length)
copy(out.VoteB.BlockID.Hash[:], in.VoteB.BlockID.Hash[:])
out.VoteB.BlockID.PartsHeader.Total = in.VoteB.BlockID.PartsHeader.Total
length = len(in.VoteB.BlockID.PartsHeader.Hash)
out.VoteB.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.VoteB.BlockID.PartsHeader.Hash[:], in.VoteB.BlockID.PartsHeader.Hash[:])
// end of .VoteB.BlockID.PartsHeader
// end of .VoteB.BlockID
out.VoteB.Timestamp = DeepCopyTime(in.VoteB.Timestamp)
length = len(in.VoteB.ValidatorAddress)
out.VoteB.ValidatorAddress = make([]uint8, length)
copy(out.VoteB.ValidatorAddress[:], in.VoteB.ValidatorAddress[:])
out.VoteB.ValidatorIndex = in.VoteB.ValidatorIndex
length = len(in.VoteB.Signature)
out.VoteB.Signature = make([]uint8, length)
copy(out.VoteB.Signature[:], in.VoteB.Signature[:])
// end of .VoteB
return
} //End of DeepCopyDuplicateVoteEvidence

// Non-Interface
func EncodePrivKeyEd25519(w *[]byte, v PrivKeyEd25519) {
codonEncodeByteSlice(w, v[:])
} //End of EncodePrivKeyEd25519

func DecodePrivKeyEd25519(bz []byte) (PrivKeyEd25519, int, error) {
var err error
var length int
var v PrivKeyEd25519
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
v[_0] = uint8(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
return v, total, nil
} //End of DecodePrivKeyEd25519

func RandPrivKeyEd25519(r RandSrc) PrivKeyEd25519 {
var length int
var v PrivKeyEd25519
length = 64
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
v[_0] = r.GetUint8()
}
return v
} //End of RandPrivKeyEd25519

func DeepCopyPrivKeyEd25519(in PrivKeyEd25519) (out PrivKeyEd25519) {
var length int
length = len(in)
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
out[_0] = in[_0]
}
return
} //End of DeepCopyPrivKeyEd25519

// Non-Interface
func EncodePrivKeySecp256k1(w *[]byte, v PrivKeySecp256k1) {
codonEncodeByteSlice(w, v[:])
} //End of EncodePrivKeySecp256k1

func DecodePrivKeySecp256k1(bz []byte) (PrivKeySecp256k1, int, error) {
var err error
var length int
var v PrivKeySecp256k1
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
v[_0] = uint8(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
return v, total, nil
} //End of DecodePrivKeySecp256k1

func RandPrivKeySecp256k1(r RandSrc) PrivKeySecp256k1 {
var length int
var v PrivKeySecp256k1
length = 32
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
v[_0] = r.GetUint8()
}
return v
} //End of RandPrivKeySecp256k1

func DeepCopyPrivKeySecp256k1(in PrivKeySecp256k1) (out PrivKeySecp256k1) {
var length int
length = len(in)
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
out[_0] = in[_0]
}
return
} //End of DeepCopyPrivKeySecp256k1

// Non-Interface
func EncodePubKeyEd25519(w *[]byte, v PubKeyEd25519) {
codonEncodeByteSlice(w, v[:])
} //End of EncodePubKeyEd25519

func DecodePubKeyEd25519(bz []byte) (PubKeyEd25519, int, error) {
var err error
var length int
var v PubKeyEd25519
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
v[_0] = uint8(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
return v, total, nil
} //End of DecodePubKeyEd25519

func RandPubKeyEd25519(r RandSrc) PubKeyEd25519 {
var length int
var v PubKeyEd25519
length = 32
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
v[_0] = r.GetUint8()
}
return v
} //End of RandPubKeyEd25519

func DeepCopyPubKeyEd25519(in PubKeyEd25519) (out PubKeyEd25519) {
var length int
length = len(in)
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
out[_0] = in[_0]
}
return
} //End of DeepCopyPubKeyEd25519

// Non-Interface
func EncodePubKeySecp256k1(w *[]byte, v PubKeySecp256k1) {
codonEncodeByteSlice(w, v[:])
} //End of EncodePubKeySecp256k1

func DecodePubKeySecp256k1(bz []byte) (PubKeySecp256k1, int, error) {
var err error
var length int
var v PubKeySecp256k1
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
v[_0] = uint8(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
return v, total, nil
} //End of DecodePubKeySecp256k1

func RandPubKeySecp256k1(r RandSrc) PubKeySecp256k1 {
var length int
var v PubKeySecp256k1
length = 33
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
v[_0] = r.GetUint8()
}
return v
} //End of RandPubKeySecp256k1

func DeepCopyPubKeySecp256k1(in PubKeySecp256k1) (out PubKeySecp256k1) {
var length int
length = len(in)
for _0, length_0 := 0, length; _0<length_0; _0++ { //array of uint8
out[_0] = in[_0]
}
return
} //End of DeepCopyPubKeySecp256k1

// Non-Interface
func EncodePubKeyMultisigThreshold(w *[]byte, v PubKeyMultisigThreshold) {
codonEncodeUvarint(w, uint64(v.K))
codonEncodeVarint(w, int64(len(v.PubKeys)))
for _0:=0; _0<len(v.PubKeys); _0++ {
EncodePubKey(w, v.PubKeys[_0]) // interface_encode
}
} //End of EncodePubKeyMultisigThreshold

func DecodePubKeyMultisigThreshold(bz []byte) (PubKeyMultisigThreshold, int, error) {
var err error
var length int
var v PubKeyMultisigThreshold
var n int
var total int
v.K = uint(codonDecodeUint(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.PubKeys = make([]PubKey, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of interface
v.PubKeys[_0], n, err = DecodePubKey(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
return v, total, nil
} //End of DecodePubKeyMultisigThreshold

func RandPubKeyMultisigThreshold(r RandSrc) PubKeyMultisigThreshold {
var length int
var v PubKeyMultisigThreshold
v.K = r.GetUint()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.PubKeys = make([]PubKey, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of interface
v.PubKeys[_0] = RandPubKey(r)
}
return v
} //End of RandPubKeyMultisigThreshold

func DeepCopyPubKeyMultisigThreshold(in PubKeyMultisigThreshold) (out PubKeyMultisigThreshold) {
var length int
out.K = in.K
length = len(in.PubKeys)
out.PubKeys = make([]PubKey, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of interface
out.PubKeys[_0] = DeepCopyPubKey(in.PubKeys[_0])
}
return
} //End of DeepCopyPubKeyMultisigThreshold

// Non-Interface
func EncodeSignedMsgType(w *[]byte, v SignedMsgType) {
codonEncodeUint8(w, uint8(v))
} //End of EncodeSignedMsgType

func DecodeSignedMsgType(bz []byte) (SignedMsgType, int, error) {
var err error
var v SignedMsgType
var n int
var total int
v = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeSignedMsgType

func RandSignedMsgType(r RandSrc) SignedMsgType {
var v SignedMsgType
v = SignedMsgType(r.GetUint8())
return v
} //End of RandSignedMsgType

func DeepCopySignedMsgType(in SignedMsgType) (out SignedMsgType) {
out = in
return
} //End of DeepCopySignedMsgType

// Non-Interface
func EncodeVote(w *[]byte, v Vote) {
codonEncodeUint8(w, uint8(v.Type))
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeByteSlice(w, v.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.BlockID.PartsHeader.Hash[:])
// end of v.BlockID.PartsHeader
// end of v.BlockID
EncodeTime(w, v.Timestamp)
codonEncodeByteSlice(w, v.ValidatorAddress[:])
codonEncodeVarint(w, int64(v.ValidatorIndex))
codonEncodeByteSlice(w, v.Signature[:])
} //End of EncodeVote

func DecodeVote(bz []byte) (Vote, int, error) {
var err error
var length int
var v Vote
var n int
var total int
v.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.BlockID.PartsHeader
// end of v.BlockID
v.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeVote

func RandVote(r RandSrc) Vote {
var length int
var v Vote
v.Type = SignedMsgType(r.GetUint8())
v.Height = r.GetInt64()
v.Round = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.Hash = r.GetBytes(length)
v.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.BlockID.PartsHeader
// end of v.BlockID
v.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ValidatorAddress = r.GetBytes(length)
v.ValidatorIndex = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Signature = r.GetBytes(length)
return v
} //End of RandVote

func DeepCopyVote(in Vote) (out Vote) {
var length int
out.Type = in.Type
out.Height = in.Height
out.Round = in.Round
length = len(in.BlockID.Hash)
out.BlockID.Hash = make([]uint8, length)
copy(out.BlockID.Hash[:], in.BlockID.Hash[:])
out.BlockID.PartsHeader.Total = in.BlockID.PartsHeader.Total
length = len(in.BlockID.PartsHeader.Hash)
out.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.BlockID.PartsHeader.Hash[:], in.BlockID.PartsHeader.Hash[:])
// end of .BlockID.PartsHeader
// end of .BlockID
out.Timestamp = DeepCopyTime(in.Timestamp)
length = len(in.ValidatorAddress)
out.ValidatorAddress = make([]uint8, length)
copy(out.ValidatorAddress[:], in.ValidatorAddress[:])
out.ValidatorIndex = in.ValidatorIndex
length = len(in.Signature)
out.Signature = make([]uint8, length)
copy(out.Signature[:], in.Signature[:])
return
} //End of DeepCopyVote

// Non-Interface
func EncodeNewRoundStepMessage(w *[]byte, v NewRoundStepMessage) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeUint8(w, uint8(v.Step))
codonEncodeVarint(w, int64(v.SecondsSinceStartTime))
codonEncodeVarint(w, int64(v.LastCommitRound))
} //End of EncodeNewRoundStepMessage

func DecodeNewRoundStepMessage(bz []byte) (NewRoundStepMessage, int, error) {
var err error
var v NewRoundStepMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Step = RoundStepType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.SecondsSinceStartTime = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.LastCommitRound = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeNewRoundStepMessage

func RandNewRoundStepMessage(r RandSrc) NewRoundStepMessage {
var v NewRoundStepMessage
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.Step = RoundStepType(r.GetUint8())
v.SecondsSinceStartTime = r.GetInt()
v.LastCommitRound = r.GetInt()
return v
} //End of RandNewRoundStepMessage

func DeepCopyNewRoundStepMessage(in NewRoundStepMessage) (out NewRoundStepMessage) {
out.Height = in.Height
out.Round = in.Round
out.Step = in.Step
out.SecondsSinceStartTime = in.SecondsSinceStartTime
out.LastCommitRound = in.LastCommitRound
return
} //End of DeepCopyNewRoundStepMessage

// Non-Interface
func EncodeNewValidBlockMessage(w *[]byte, v NewValidBlockMessage) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeVarint(w, int64(v.BlockPartsHeader.Total))
codonEncodeByteSlice(w, v.BlockPartsHeader.Hash[:])
// end of v.BlockPartsHeader
codonEncodeVarint(w, int64(v.BlockParts.Bits))
codonEncodeVarint(w, int64(len(v.BlockParts.Elems)))
for _0:=0; _0<len(v.BlockParts.Elems); _0++ {
codonEncodeUvarint(w, uint64(v.BlockParts.Elems[_0]))
}
// end of v.BlockParts
codonEncodeBool(w, v.IsCommit)
} //End of EncodeNewValidBlockMessage

func DecodeNewValidBlockMessage(bz []byte) (NewValidBlockMessage, int, error) {
var err error
var length int
var v NewValidBlockMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockPartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockPartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.BlockPartsHeader
v.BlockParts = &BitArray{}
v.BlockParts.Bits = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockParts.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
v.BlockParts.Elems[_0] = uint64(codonDecodeUint64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.BlockParts
v.IsCommit = bool(codonDecodeBool(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeNewValidBlockMessage

func RandNewValidBlockMessage(r RandSrc) NewValidBlockMessage {
var length int
var v NewValidBlockMessage
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.BlockPartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockPartsHeader.Hash = r.GetBytes(length)
// end of v.BlockPartsHeader
v.BlockParts = &BitArray{}
v.BlockParts.Bits = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockParts.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
v.BlockParts.Elems[_0] = r.GetUint64()
}
// end of v.BlockParts
v.IsCommit = r.GetBool()
return v
} //End of RandNewValidBlockMessage

func DeepCopyNewValidBlockMessage(in NewValidBlockMessage) (out NewValidBlockMessage) {
var length int
out.Height = in.Height
out.Round = in.Round
out.BlockPartsHeader.Total = in.BlockPartsHeader.Total
length = len(in.BlockPartsHeader.Hash)
out.BlockPartsHeader.Hash = make([]uint8, length)
copy(out.BlockPartsHeader.Hash[:], in.BlockPartsHeader.Hash[:])
// end of .BlockPartsHeader
out.BlockParts = &BitArray{}
out.BlockParts.Bits = in.BlockParts.Bits
length = len(in.BlockParts.Elems)
out.BlockParts.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
out.BlockParts.Elems[_0] = in.BlockParts.Elems[_0]
}
// end of .BlockParts
out.IsCommit = in.IsCommit
return
} //End of DeepCopyNewValidBlockMessage

// Non-Interface
func EncodeProposalMessage(w *[]byte, v ProposalMessage) {
codonEncodeUint8(w, uint8(v.Proposal.Type))
codonEncodeVarint(w, int64(v.Proposal.Height))
codonEncodeVarint(w, int64(v.Proposal.Round))
codonEncodeVarint(w, int64(v.Proposal.POLRound))
codonEncodeByteSlice(w, v.Proposal.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.Proposal.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.Proposal.BlockID.PartsHeader.Hash[:])
// end of v.Proposal.BlockID.PartsHeader
// end of v.Proposal.BlockID
EncodeTime(w, v.Proposal.Timestamp)
codonEncodeByteSlice(w, v.Proposal.Signature[:])
// end of v.Proposal
} //End of EncodeProposalMessage

func DecodeProposalMessage(bz []byte) (ProposalMessage, int, error) {
var err error
var length int
var v ProposalMessage
var n int
var total int
v.Proposal = &Proposal{}
v.Proposal.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.POLRound = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Proposal.BlockID.PartsHeader
// end of v.Proposal.BlockID
v.Proposal.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Proposal
return v, total, nil
} //End of DecodeProposalMessage

func RandProposalMessage(r RandSrc) ProposalMessage {
var length int
var v ProposalMessage
v.Proposal = &Proposal{}
v.Proposal.Type = SignedMsgType(r.GetUint8())
v.Proposal.Height = r.GetInt64()
v.Proposal.Round = r.GetInt()
v.Proposal.POLRound = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposal.BlockID.Hash = r.GetBytes(length)
v.Proposal.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposal.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.Proposal.BlockID.PartsHeader
// end of v.Proposal.BlockID
v.Proposal.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposal.Signature = r.GetBytes(length)
// end of v.Proposal
return v
} //End of RandProposalMessage

func DeepCopyProposalMessage(in ProposalMessage) (out ProposalMessage) {
var length int
out.Proposal = &Proposal{}
out.Proposal.Type = in.Proposal.Type
out.Proposal.Height = in.Proposal.Height
out.Proposal.Round = in.Proposal.Round
out.Proposal.POLRound = in.Proposal.POLRound
length = len(in.Proposal.BlockID.Hash)
out.Proposal.BlockID.Hash = make([]uint8, length)
copy(out.Proposal.BlockID.Hash[:], in.Proposal.BlockID.Hash[:])
out.Proposal.BlockID.PartsHeader.Total = in.Proposal.BlockID.PartsHeader.Total
length = len(in.Proposal.BlockID.PartsHeader.Hash)
out.Proposal.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.Proposal.BlockID.PartsHeader.Hash[:], in.Proposal.BlockID.PartsHeader.Hash[:])
// end of .Proposal.BlockID.PartsHeader
// end of .Proposal.BlockID
out.Proposal.Timestamp = DeepCopyTime(in.Proposal.Timestamp)
length = len(in.Proposal.Signature)
out.Proposal.Signature = make([]uint8, length)
copy(out.Proposal.Signature[:], in.Proposal.Signature[:])
// end of .Proposal
return
} //End of DeepCopyProposalMessage

// Non-Interface
func EncodeProposalPOLMessage(w *[]byte, v ProposalPOLMessage) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.ProposalPOLRound))
codonEncodeVarint(w, int64(v.ProposalPOL.Bits))
codonEncodeVarint(w, int64(len(v.ProposalPOL.Elems)))
for _0:=0; _0<len(v.ProposalPOL.Elems); _0++ {
codonEncodeUvarint(w, uint64(v.ProposalPOL.Elems[_0]))
}
// end of v.ProposalPOL
} //End of EncodeProposalPOLMessage

func DecodeProposalPOLMessage(bz []byte) (ProposalPOLMessage, int, error) {
var err error
var length int
var v ProposalPOLMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ProposalPOLRound = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ProposalPOL = &BitArray{}
v.ProposalPOL.Bits = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ProposalPOL.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
v.ProposalPOL.Elems[_0] = uint64(codonDecodeUint64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.ProposalPOL
return v, total, nil
} //End of DecodeProposalPOLMessage

func RandProposalPOLMessage(r RandSrc) ProposalPOLMessage {
var length int
var v ProposalPOLMessage
v.Height = r.GetInt64()
v.ProposalPOLRound = r.GetInt()
v.ProposalPOL = &BitArray{}
v.ProposalPOL.Bits = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ProposalPOL.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
v.ProposalPOL.Elems[_0] = r.GetUint64()
}
// end of v.ProposalPOL
return v
} //End of RandProposalPOLMessage

func DeepCopyProposalPOLMessage(in ProposalPOLMessage) (out ProposalPOLMessage) {
var length int
out.Height = in.Height
out.ProposalPOLRound = in.ProposalPOLRound
out.ProposalPOL = &BitArray{}
out.ProposalPOL.Bits = in.ProposalPOL.Bits
length = len(in.ProposalPOL.Elems)
out.ProposalPOL.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
out.ProposalPOL.Elems[_0] = in.ProposalPOL.Elems[_0]
}
// end of .ProposalPOL
return
} //End of DeepCopyProposalPOLMessage

// Non-Interface
func EncodeBlockPartMessage(w *[]byte, v BlockPartMessage) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeVarint(w, int64(v.Part.Index))
codonEncodeByteSlice(w, v.Part.Bytes[:])
codonEncodeVarint(w, int64(v.Part.Proof.Total))
codonEncodeVarint(w, int64(v.Part.Proof.Index))
codonEncodeByteSlice(w, v.Part.Proof.LeafHash[:])
codonEncodeVarint(w, int64(len(v.Part.Proof.Aunts)))
for _0:=0; _0<len(v.Part.Proof.Aunts); _0++ {
codonEncodeByteSlice(w, v.Part.Proof.Aunts[_0][:])
}
// end of v.Part.Proof
// end of v.Part
} //End of EncodeBlockPartMessage

func DecodeBlockPartMessage(bz []byte) (BlockPartMessage, int, error) {
var err error
var length int
var v BlockPartMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Part = &Part{}
v.Part.Index = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Part.Bytes, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Part.Proof.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Part.Proof.Index = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Part.Proof.LeafHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Part.Proof.Aunts = make([][]byte, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of slice
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Part.Proof.Aunts[_0], n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.Part.Proof
// end of v.Part
return v, total, nil
} //End of DecodeBlockPartMessage

func RandBlockPartMessage(r RandSrc) BlockPartMessage {
var length int
var v BlockPartMessage
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.Part = &Part{}
v.Part.Index = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Part.Bytes = r.GetBytes(length)
v.Part.Proof.Total = r.GetInt()
v.Part.Proof.Index = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Part.Proof.LeafHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Part.Proof.Aunts = make([][]byte, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of slice
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Part.Proof.Aunts[_0] = r.GetBytes(length)
}
// end of v.Part.Proof
// end of v.Part
return v
} //End of RandBlockPartMessage

func DeepCopyBlockPartMessage(in BlockPartMessage) (out BlockPartMessage) {
var length int
out.Height = in.Height
out.Round = in.Round
out.Part = &Part{}
out.Part.Index = in.Part.Index
length = len(in.Part.Bytes)
out.Part.Bytes = make([]uint8, length)
copy(out.Part.Bytes[:], in.Part.Bytes[:])
out.Part.Proof.Total = in.Part.Proof.Total
out.Part.Proof.Index = in.Part.Proof.Index
length = len(in.Part.Proof.LeafHash)
out.Part.Proof.LeafHash = make([]uint8, length)
copy(out.Part.Proof.LeafHash[:], in.Part.Proof.LeafHash[:])
length = len(in.Part.Proof.Aunts)
out.Part.Proof.Aunts = make([][]byte, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of slice
length = len(in.Part.Proof.Aunts[_0])
out.Part.Proof.Aunts[_0] = make([]uint8, length)
copy(out.Part.Proof.Aunts[_0][:], in.Part.Proof.Aunts[_0][:])
}
// end of .Part.Proof
// end of .Part
return
} //End of DeepCopyBlockPartMessage

// Non-Interface
func EncodeVoteMessage(w *[]byte, v VoteMessage) {
codonEncodeUint8(w, uint8(v.Vote.Type))
codonEncodeVarint(w, int64(v.Vote.Height))
codonEncodeVarint(w, int64(v.Vote.Round))
codonEncodeByteSlice(w, v.Vote.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.Vote.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.Vote.BlockID.PartsHeader.Hash[:])
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
EncodeTime(w, v.Vote.Timestamp)
codonEncodeByteSlice(w, v.Vote.ValidatorAddress[:])
codonEncodeVarint(w, int64(v.Vote.ValidatorIndex))
codonEncodeByteSlice(w, v.Vote.Signature[:])
// end of v.Vote
} //End of EncodeVoteMessage

func DecodeVoteMessage(bz []byte) (VoteMessage, int, error) {
var err error
var length int
var v VoteMessage
var n int
var total int
v.Vote = &Vote{}
v.Vote.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
v.Vote.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Vote
return v, total, nil
} //End of DecodeVoteMessage

func RandVoteMessage(r RandSrc) VoteMessage {
var length int
var v VoteMessage
v.Vote = &Vote{}
v.Vote.Type = SignedMsgType(r.GetUint8())
v.Vote.Height = r.GetInt64()
v.Vote.Round = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.BlockID.Hash = r.GetBytes(length)
v.Vote.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
v.Vote.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.ValidatorAddress = r.GetBytes(length)
v.Vote.ValidatorIndex = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.Signature = r.GetBytes(length)
// end of v.Vote
return v
} //End of RandVoteMessage

func DeepCopyVoteMessage(in VoteMessage) (out VoteMessage) {
var length int
out.Vote = &Vote{}
out.Vote.Type = in.Vote.Type
out.Vote.Height = in.Vote.Height
out.Vote.Round = in.Vote.Round
length = len(in.Vote.BlockID.Hash)
out.Vote.BlockID.Hash = make([]uint8, length)
copy(out.Vote.BlockID.Hash[:], in.Vote.BlockID.Hash[:])
out.Vote.BlockID.PartsHeader.Total = in.Vote.BlockID.PartsHeader.Total
length = len(in.Vote.BlockID.PartsHeader.Hash)
out.Vote.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.Vote.BlockID.PartsHeader.Hash[:], in.Vote.BlockID.PartsHeader.Hash[:])
// end of .Vote.BlockID.PartsHeader
// end of .Vote.BlockID
out.Vote.Timestamp = DeepCopyTime(in.Vote.Timestamp)
length = len(in.Vote.ValidatorAddress)
out.Vote.ValidatorAddress = make([]uint8, length)
copy(out.Vote.ValidatorAddress[:], in.Vote.ValidatorAddress[:])
out.Vote.ValidatorIndex = in.Vote.ValidatorIndex
length = len(in.Vote.Signature)
out.Vote.Signature = make([]uint8, length)
copy(out.Vote.Signature[:], in.Vote.Signature[:])
// end of .Vote
return
} //End of DeepCopyVoteMessage

// Non-Interface
func EncodeHasVoteMessage(w *[]byte, v HasVoteMessage) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeUint8(w, uint8(v.Type))
codonEncodeVarint(w, int64(v.Index))
} //End of EncodeHasVoteMessage

func DecodeHasVoteMessage(bz []byte) (HasVoteMessage, int, error) {
var err error
var v HasVoteMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Index = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeHasVoteMessage

func RandHasVoteMessage(r RandSrc) HasVoteMessage {
var v HasVoteMessage
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.Type = SignedMsgType(r.GetUint8())
v.Index = r.GetInt()
return v
} //End of RandHasVoteMessage

func DeepCopyHasVoteMessage(in HasVoteMessage) (out HasVoteMessage) {
out.Height = in.Height
out.Round = in.Round
out.Type = in.Type
out.Index = in.Index
return
} //End of DeepCopyHasVoteMessage

// Non-Interface
func EncodeVoteSetMaj23Message(w *[]byte, v VoteSetMaj23Message) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeUint8(w, uint8(v.Type))
codonEncodeByteSlice(w, v.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.BlockID.PartsHeader.Hash[:])
// end of v.BlockID.PartsHeader
// end of v.BlockID
} //End of EncodeVoteSetMaj23Message

func DecodeVoteSetMaj23Message(bz []byte) (VoteSetMaj23Message, int, error) {
var err error
var length int
var v VoteSetMaj23Message
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.BlockID.PartsHeader
// end of v.BlockID
return v, total, nil
} //End of DecodeVoteSetMaj23Message

func RandVoteSetMaj23Message(r RandSrc) VoteSetMaj23Message {
var length int
var v VoteSetMaj23Message
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.Type = SignedMsgType(r.GetUint8())
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.Hash = r.GetBytes(length)
v.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.BlockID.PartsHeader
// end of v.BlockID
return v
} //End of RandVoteSetMaj23Message

func DeepCopyVoteSetMaj23Message(in VoteSetMaj23Message) (out VoteSetMaj23Message) {
var length int
out.Height = in.Height
out.Round = in.Round
out.Type = in.Type
length = len(in.BlockID.Hash)
out.BlockID.Hash = make([]uint8, length)
copy(out.BlockID.Hash[:], in.BlockID.Hash[:])
out.BlockID.PartsHeader.Total = in.BlockID.PartsHeader.Total
length = len(in.BlockID.PartsHeader.Hash)
out.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.BlockID.PartsHeader.Hash[:], in.BlockID.PartsHeader.Hash[:])
// end of .BlockID.PartsHeader
// end of .BlockID
return
} //End of DeepCopyVoteSetMaj23Message

// Non-Interface
func EncodeVoteSetBitsMessage(w *[]byte, v VoteSetBitsMessage) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeUint8(w, uint8(v.Type))
codonEncodeByteSlice(w, v.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.BlockID.PartsHeader.Hash[:])
// end of v.BlockID.PartsHeader
// end of v.BlockID
codonEncodeVarint(w, int64(v.Votes.Bits))
codonEncodeVarint(w, int64(len(v.Votes.Elems)))
for _0:=0; _0<len(v.Votes.Elems); _0++ {
codonEncodeUvarint(w, uint64(v.Votes.Elems[_0]))
}
// end of v.Votes
} //End of EncodeVoteSetBitsMessage

func DecodeVoteSetBitsMessage(bz []byte) (VoteSetBitsMessage, int, error) {
var err error
var length int
var v VoteSetBitsMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.BlockID.PartsHeader
// end of v.BlockID
v.Votes = &BitArray{}
v.Votes.Bits = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Votes.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
v.Votes.Elems[_0] = uint64(codonDecodeUint64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.Votes
return v, total, nil
} //End of DecodeVoteSetBitsMessage

func RandVoteSetBitsMessage(r RandSrc) VoteSetBitsMessage {
var length int
var v VoteSetBitsMessage
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.Type = SignedMsgType(r.GetUint8())
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.Hash = r.GetBytes(length)
v.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.BlockID.PartsHeader
// end of v.BlockID
v.Votes = &BitArray{}
v.Votes.Bits = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Votes.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
v.Votes.Elems[_0] = r.GetUint64()
}
// end of v.Votes
return v
} //End of RandVoteSetBitsMessage

func DeepCopyVoteSetBitsMessage(in VoteSetBitsMessage) (out VoteSetBitsMessage) {
var length int
out.Height = in.Height
out.Round = in.Round
out.Type = in.Type
length = len(in.BlockID.Hash)
out.BlockID.Hash = make([]uint8, length)
copy(out.BlockID.Hash[:], in.BlockID.Hash[:])
out.BlockID.PartsHeader.Total = in.BlockID.PartsHeader.Total
length = len(in.BlockID.PartsHeader.Hash)
out.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.BlockID.PartsHeader.Hash[:], in.BlockID.PartsHeader.Hash[:])
// end of .BlockID.PartsHeader
// end of .BlockID
out.Votes = &BitArray{}
out.Votes.Bits = in.Votes.Bits
length = len(in.Votes.Elems)
out.Votes.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
out.Votes.Elems[_0] = in.Votes.Elems[_0]
}
// end of .Votes
return
} //End of DeepCopyVoteSetBitsMessage

// Non-Interface
func EncodeEventDataRoundState(w *[]byte, v EventDataRoundState) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeString(w, v.Step)
} //End of EncodeEventDataRoundState

func DecodeEventDataRoundState(bz []byte) (EventDataRoundState, int, error) {
var err error
var v EventDataRoundState
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Step = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeEventDataRoundState

func RandEventDataRoundState(r RandSrc) EventDataRoundState {
var v EventDataRoundState
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.Step = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
return v
} //End of RandEventDataRoundState

func DeepCopyEventDataRoundState(in EventDataRoundState) (out EventDataRoundState) {
out.Height = in.Height
out.Round = in.Round
out.Step = in.Step
return
} //End of DeepCopyEventDataRoundState

// Non-Interface
func EncodeEndHeightMessage(w *[]byte, v EndHeightMessage) {
codonEncodeVarint(w, int64(v.Height))
} //End of EncodeEndHeightMessage

func DecodeEndHeightMessage(bz []byte) (EndHeightMessage, int, error) {
var err error
var v EndHeightMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeEndHeightMessage

func RandEndHeightMessage(r RandSrc) EndHeightMessage {
var v EndHeightMessage
v.Height = r.GetInt64()
return v
} //End of RandEndHeightMessage

func DeepCopyEndHeightMessage(in EndHeightMessage) (out EndHeightMessage) {
out.Height = in.Height
return
} //End of DeepCopyEndHeightMessage

// Non-Interface
func EncodeEventDataNewBlockHeader(w *[]byte, v EventDataNewBlockHeader) {
codonEncodeUvarint(w, uint64(v.Header.Version.Block))
codonEncodeUvarint(w, uint64(v.Header.Version.App))
// end of v.Header.Version
codonEncodeString(w, v.Header.ChainID)
codonEncodeVarint(w, int64(v.Header.Height))
EncodeTime(w, v.Header.Time)
codonEncodeVarint(w, int64(v.Header.NumTxs))
codonEncodeVarint(w, int64(v.Header.TotalTxs))
codonEncodeByteSlice(w, v.Header.LastBlockID.Hash[:])
codonEncodeVarint(w, int64(v.Header.LastBlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.Header.LastBlockID.PartsHeader.Hash[:])
// end of v.Header.LastBlockID.PartsHeader
// end of v.Header.LastBlockID
codonEncodeByteSlice(w, v.Header.LastCommitHash[:])
codonEncodeByteSlice(w, v.Header.DataHash[:])
codonEncodeByteSlice(w, v.Header.ValidatorsHash[:])
codonEncodeByteSlice(w, v.Header.NextValidatorsHash[:])
codonEncodeByteSlice(w, v.Header.ConsensusHash[:])
codonEncodeByteSlice(w, v.Header.AppHash[:])
codonEncodeByteSlice(w, v.Header.LastResultsHash[:])
codonEncodeByteSlice(w, v.Header.EvidenceHash[:])
codonEncodeByteSlice(w, v.Header.ProposerAddress[:])
// end of v.Header
codonEncodeVarint(w, int64(len(v.ResultBeginBlock.Events)))
for _0:=0; _0<len(v.ResultBeginBlock.Events); _0++ {
codonEncodeString(w, v.ResultBeginBlock.Events[_0].Type)
codonEncodeVarint(w, int64(len(v.ResultBeginBlock.Events[_0].Attributes)))
for _1:=0; _1<len(v.ResultBeginBlock.Events[_0].Attributes); _1++ {
codonEncodeByteSlice(w, v.ResultBeginBlock.Events[_0].Attributes[_1].Key[:])
codonEncodeByteSlice(w, v.ResultBeginBlock.Events[_0].Attributes[_1].Value[:])
// end of v.ResultBeginBlock.Events[_0].Attributes[_1].XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultBeginBlock.Events[_0].Attributes[_1].XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultBeginBlock.Events[_0].Attributes[_1].XXX_sizecache))
// end of v.ResultBeginBlock.Events[_0].Attributes[_1]
}
// end of v.ResultBeginBlock.Events[_0].XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultBeginBlock.Events[_0].XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultBeginBlock.Events[_0].XXX_sizecache))
// end of v.ResultBeginBlock.Events[_0]
}
// end of v.ResultBeginBlock.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultBeginBlock.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultBeginBlock.XXX_sizecache))
// end of v.ResultBeginBlock
codonEncodeVarint(w, int64(len(v.ResultEndBlock.ValidatorUpdates)))
for _0:=0; _0<len(v.ResultEndBlock.ValidatorUpdates); _0++ {
codonEncodeString(w, v.ResultEndBlock.ValidatorUpdates[_0].PubKey.Type)
codonEncodeByteSlice(w, v.ResultEndBlock.ValidatorUpdates[_0].PubKey.Data[:])
// end of v.ResultEndBlock.ValidatorUpdates[_0].PubKey.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultEndBlock.ValidatorUpdates[_0].PubKey.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultEndBlock.ValidatorUpdates[_0].PubKey.XXX_sizecache))
// end of v.ResultEndBlock.ValidatorUpdates[_0].PubKey
codonEncodeVarint(w, int64(v.ResultEndBlock.ValidatorUpdates[_0].Power))
// end of v.ResultEndBlock.ValidatorUpdates[_0].XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultEndBlock.ValidatorUpdates[_0].XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultEndBlock.ValidatorUpdates[_0].XXX_sizecache))
// end of v.ResultEndBlock.ValidatorUpdates[_0]
}
codonEncodeVarint(w, int64(v.ResultEndBlock.ConsensusParamUpdates.Block.MaxBytes))
codonEncodeVarint(w, int64(v.ResultEndBlock.ConsensusParamUpdates.Block.MaxGas))
// end of v.ResultEndBlock.ConsensusParamUpdates.Block.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultEndBlock.ConsensusParamUpdates.Block.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultEndBlock.ConsensusParamUpdates.Block.XXX_sizecache))
// end of v.ResultEndBlock.ConsensusParamUpdates.Block
codonEncodeVarint(w, int64(v.ResultEndBlock.ConsensusParamUpdates.Evidence.MaxAge))
// end of v.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_sizecache))
// end of v.ResultEndBlock.ConsensusParamUpdates.Evidence
codonEncodeVarint(w, int64(len(v.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes)))
for _0:=0; _0<len(v.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes); _0++ {
codonEncodeString(w, v.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes[_0])
}
// end of v.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_sizecache))
// end of v.ResultEndBlock.ConsensusParamUpdates.Validator
// end of v.ResultEndBlock.ConsensusParamUpdates.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultEndBlock.ConsensusParamUpdates.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultEndBlock.ConsensusParamUpdates.XXX_sizecache))
// end of v.ResultEndBlock.ConsensusParamUpdates
codonEncodeVarint(w, int64(len(v.ResultEndBlock.Events)))
for _0:=0; _0<len(v.ResultEndBlock.Events); _0++ {
codonEncodeString(w, v.ResultEndBlock.Events[_0].Type)
codonEncodeVarint(w, int64(len(v.ResultEndBlock.Events[_0].Attributes)))
for _1:=0; _1<len(v.ResultEndBlock.Events[_0].Attributes); _1++ {
codonEncodeByteSlice(w, v.ResultEndBlock.Events[_0].Attributes[_1].Key[:])
codonEncodeByteSlice(w, v.ResultEndBlock.Events[_0].Attributes[_1].Value[:])
// end of v.ResultEndBlock.Events[_0].Attributes[_1].XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultEndBlock.Events[_0].Attributes[_1].XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultEndBlock.Events[_0].Attributes[_1].XXX_sizecache))
// end of v.ResultEndBlock.Events[_0].Attributes[_1]
}
// end of v.ResultEndBlock.Events[_0].XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultEndBlock.Events[_0].XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultEndBlock.Events[_0].XXX_sizecache))
// end of v.ResultEndBlock.Events[_0]
}
// end of v.ResultEndBlock.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.ResultEndBlock.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.ResultEndBlock.XXX_sizecache))
// end of v.ResultEndBlock
} //End of EncodeEventDataNewBlockHeader

func DecodeEventDataNewBlockHeader(bz []byte) (EventDataNewBlockHeader, int, error) {
var err error
var length int
var v EventDataNewBlockHeader
var n int
var total int
v.Header.Version.Block = Protocol(codonDecodeUint64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.Version.App = Protocol(codonDecodeUint64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Header.Version
v.Header.ChainID = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.Time, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.NumTxs = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.TotalTxs = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.LastBlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.LastBlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.LastBlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Header.LastBlockID.PartsHeader
// end of v.Header.LastBlockID
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.LastCommitHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.DataHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.ValidatorsHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.NextValidatorsHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.ConsensusHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.AppHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.LastResultsHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.EvidenceHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Header.ProposerAddress, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Header
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultBeginBlock.Events = make([]Event, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.ResultBeginBlock.Events[_0], n, err = DecodeEvent(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.ResultBeginBlock.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultBeginBlock.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultBeginBlock.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.ResultBeginBlock
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ValidatorUpdates = make([]ValidatorUpdate, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.ResultEndBlock.ValidatorUpdates[_0], n, err = DecodeValidatorUpdate(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
v.ResultEndBlock.ConsensusParamUpdates = &ConsensusParams{}
v.ResultEndBlock.ConsensusParamUpdates.Block = &BlockParams{}
v.ResultEndBlock.ConsensusParamUpdates.Block.MaxBytes = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.Block.MaxGas = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.ResultEndBlock.ConsensusParamUpdates.Block.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.Block.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.Block.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.ResultEndBlock.ConsensusParamUpdates.Block
v.ResultEndBlock.ConsensusParamUpdates.Evidence = &EvidenceParams{}
v.ResultEndBlock.ConsensusParamUpdates.Evidence.MaxAge = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.ResultEndBlock.ConsensusParamUpdates.Evidence
v.ResultEndBlock.ConsensusParamUpdates.Validator = &ValidatorParams{}
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes = make([]string, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of string
v.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes[_0] = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.ResultEndBlock.ConsensusParamUpdates.Validator
// end of v.ResultEndBlock.ConsensusParamUpdates.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.ConsensusParamUpdates.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.ResultEndBlock.ConsensusParamUpdates
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.Events = make([]Event, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.ResultEndBlock.Events[_0], n, err = DecodeEvent(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.ResultEndBlock.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ResultEndBlock.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.ResultEndBlock
return v, total, nil
} //End of DecodeEventDataNewBlockHeader

func RandEventDataNewBlockHeader(r RandSrc) EventDataNewBlockHeader {
var length int
var v EventDataNewBlockHeader
v.Header.Version.Block = Protocol(r.GetUint64())
v.Header.Version.App = Protocol(r.GetUint64())
// end of v.Header.Version
v.Header.ChainID = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
v.Header.Height = r.GetInt64()
v.Header.Time = RandTime(r)
v.Header.NumTxs = r.GetInt64()
v.Header.TotalTxs = r.GetInt64()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.LastBlockID.Hash = r.GetBytes(length)
v.Header.LastBlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.LastBlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.Header.LastBlockID.PartsHeader
// end of v.Header.LastBlockID
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.LastCommitHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.DataHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.ValidatorsHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.NextValidatorsHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.ConsensusHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.AppHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.LastResultsHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.EvidenceHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Header.ProposerAddress = r.GetBytes(length)
// end of v.Header
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultBeginBlock.Events = make([]Event, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.ResultBeginBlock.Events[_0] = RandEvent(r)
}
// end of v.ResultBeginBlock.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultBeginBlock.XXX_unrecognized = r.GetBytes(length)
v.ResultBeginBlock.XXX_sizecache = r.GetInt32()
// end of v.ResultBeginBlock
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultEndBlock.ValidatorUpdates = make([]ValidatorUpdate, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.ResultEndBlock.ValidatorUpdates[_0] = RandValidatorUpdate(r)
}
v.ResultEndBlock.ConsensusParamUpdates = &ConsensusParams{}
v.ResultEndBlock.ConsensusParamUpdates.Block = &BlockParams{}
v.ResultEndBlock.ConsensusParamUpdates.Block.MaxBytes = r.GetInt64()
v.ResultEndBlock.ConsensusParamUpdates.Block.MaxGas = r.GetInt64()
// end of v.ResultEndBlock.ConsensusParamUpdates.Block.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultEndBlock.ConsensusParamUpdates.Block.XXX_unrecognized = r.GetBytes(length)
v.ResultEndBlock.ConsensusParamUpdates.Block.XXX_sizecache = r.GetInt32()
// end of v.ResultEndBlock.ConsensusParamUpdates.Block
v.ResultEndBlock.ConsensusParamUpdates.Evidence = &EvidenceParams{}
v.ResultEndBlock.ConsensusParamUpdates.Evidence.MaxAge = r.GetInt64()
// end of v.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_unrecognized = r.GetBytes(length)
v.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_sizecache = r.GetInt32()
// end of v.ResultEndBlock.ConsensusParamUpdates.Evidence
v.ResultEndBlock.ConsensusParamUpdates.Validator = &ValidatorParams{}
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes = make([]string, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of string
v.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes[_0] = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
}
// end of v.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_unrecognized = r.GetBytes(length)
v.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_sizecache = r.GetInt32()
// end of v.ResultEndBlock.ConsensusParamUpdates.Validator
// end of v.ResultEndBlock.ConsensusParamUpdates.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultEndBlock.ConsensusParamUpdates.XXX_unrecognized = r.GetBytes(length)
v.ResultEndBlock.ConsensusParamUpdates.XXX_sizecache = r.GetInt32()
// end of v.ResultEndBlock.ConsensusParamUpdates
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultEndBlock.Events = make([]Event, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.ResultEndBlock.Events[_0] = RandEvent(r)
}
// end of v.ResultEndBlock.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ResultEndBlock.XXX_unrecognized = r.GetBytes(length)
v.ResultEndBlock.XXX_sizecache = r.GetInt32()
// end of v.ResultEndBlock
return v
} //End of RandEventDataNewBlockHeader

func DeepCopyEventDataNewBlockHeader(in EventDataNewBlockHeader) (out EventDataNewBlockHeader) {
var length int
out.Header.Version.Block = in.Header.Version.Block
out.Header.Version.App = in.Header.Version.App
// end of .Header.Version
out.Header.ChainID = in.Header.ChainID
out.Header.Height = in.Header.Height
out.Header.Time = DeepCopyTime(in.Header.Time)
out.Header.NumTxs = in.Header.NumTxs
out.Header.TotalTxs = in.Header.TotalTxs
length = len(in.Header.LastBlockID.Hash)
out.Header.LastBlockID.Hash = make([]uint8, length)
copy(out.Header.LastBlockID.Hash[:], in.Header.LastBlockID.Hash[:])
out.Header.LastBlockID.PartsHeader.Total = in.Header.LastBlockID.PartsHeader.Total
length = len(in.Header.LastBlockID.PartsHeader.Hash)
out.Header.LastBlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.Header.LastBlockID.PartsHeader.Hash[:], in.Header.LastBlockID.PartsHeader.Hash[:])
// end of .Header.LastBlockID.PartsHeader
// end of .Header.LastBlockID
length = len(in.Header.LastCommitHash)
out.Header.LastCommitHash = make([]uint8, length)
copy(out.Header.LastCommitHash[:], in.Header.LastCommitHash[:])
length = len(in.Header.DataHash)
out.Header.DataHash = make([]uint8, length)
copy(out.Header.DataHash[:], in.Header.DataHash[:])
length = len(in.Header.ValidatorsHash)
out.Header.ValidatorsHash = make([]uint8, length)
copy(out.Header.ValidatorsHash[:], in.Header.ValidatorsHash[:])
length = len(in.Header.NextValidatorsHash)
out.Header.NextValidatorsHash = make([]uint8, length)
copy(out.Header.NextValidatorsHash[:], in.Header.NextValidatorsHash[:])
length = len(in.Header.ConsensusHash)
out.Header.ConsensusHash = make([]uint8, length)
copy(out.Header.ConsensusHash[:], in.Header.ConsensusHash[:])
length = len(in.Header.AppHash)
out.Header.AppHash = make([]uint8, length)
copy(out.Header.AppHash[:], in.Header.AppHash[:])
length = len(in.Header.LastResultsHash)
out.Header.LastResultsHash = make([]uint8, length)
copy(out.Header.LastResultsHash[:], in.Header.LastResultsHash[:])
length = len(in.Header.EvidenceHash)
out.Header.EvidenceHash = make([]uint8, length)
copy(out.Header.EvidenceHash[:], in.Header.EvidenceHash[:])
length = len(in.Header.ProposerAddress)
out.Header.ProposerAddress = make([]uint8, length)
copy(out.Header.ProposerAddress[:], in.Header.ProposerAddress[:])
// end of .Header
length = len(in.ResultBeginBlock.Events)
out.ResultBeginBlock.Events = make([]Event, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
out.ResultBeginBlock.Events[_0] = DeepCopyEvent(in.ResultBeginBlock.Events[_0])
}
// end of .ResultBeginBlock.XXX_NoUnkeyedLiteral
length = len(in.ResultBeginBlock.XXX_unrecognized)
out.ResultBeginBlock.XXX_unrecognized = make([]uint8, length)
copy(out.ResultBeginBlock.XXX_unrecognized[:], in.ResultBeginBlock.XXX_unrecognized[:])
out.ResultBeginBlock.XXX_sizecache = in.ResultBeginBlock.XXX_sizecache
// end of .ResultBeginBlock
length = len(in.ResultEndBlock.ValidatorUpdates)
out.ResultEndBlock.ValidatorUpdates = make([]ValidatorUpdate, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
out.ResultEndBlock.ValidatorUpdates[_0] = DeepCopyValidatorUpdate(in.ResultEndBlock.ValidatorUpdates[_0])
}
out.ResultEndBlock.ConsensusParamUpdates = &ConsensusParams{}
out.ResultEndBlock.ConsensusParamUpdates.Block = &BlockParams{}
out.ResultEndBlock.ConsensusParamUpdates.Block.MaxBytes = in.ResultEndBlock.ConsensusParamUpdates.Block.MaxBytes
out.ResultEndBlock.ConsensusParamUpdates.Block.MaxGas = in.ResultEndBlock.ConsensusParamUpdates.Block.MaxGas
// end of .ResultEndBlock.ConsensusParamUpdates.Block.XXX_NoUnkeyedLiteral
length = len(in.ResultEndBlock.ConsensusParamUpdates.Block.XXX_unrecognized)
out.ResultEndBlock.ConsensusParamUpdates.Block.XXX_unrecognized = make([]uint8, length)
copy(out.ResultEndBlock.ConsensusParamUpdates.Block.XXX_unrecognized[:], in.ResultEndBlock.ConsensusParamUpdates.Block.XXX_unrecognized[:])
out.ResultEndBlock.ConsensusParamUpdates.Block.XXX_sizecache = in.ResultEndBlock.ConsensusParamUpdates.Block.XXX_sizecache
// end of .ResultEndBlock.ConsensusParamUpdates.Block
out.ResultEndBlock.ConsensusParamUpdates.Evidence = &EvidenceParams{}
out.ResultEndBlock.ConsensusParamUpdates.Evidence.MaxAge = in.ResultEndBlock.ConsensusParamUpdates.Evidence.MaxAge
// end of .ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_NoUnkeyedLiteral
length = len(in.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_unrecognized)
out.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_unrecognized = make([]uint8, length)
copy(out.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_unrecognized[:], in.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_unrecognized[:])
out.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_sizecache = in.ResultEndBlock.ConsensusParamUpdates.Evidence.XXX_sizecache
// end of .ResultEndBlock.ConsensusParamUpdates.Evidence
out.ResultEndBlock.ConsensusParamUpdates.Validator = &ValidatorParams{}
length = len(in.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes)
out.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes = make([]string, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of string
out.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes[_0] = in.ResultEndBlock.ConsensusParamUpdates.Validator.PubKeyTypes[_0]
}
// end of .ResultEndBlock.ConsensusParamUpdates.Validator.XXX_NoUnkeyedLiteral
length = len(in.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_unrecognized)
out.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_unrecognized = make([]uint8, length)
copy(out.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_unrecognized[:], in.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_unrecognized[:])
out.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_sizecache = in.ResultEndBlock.ConsensusParamUpdates.Validator.XXX_sizecache
// end of .ResultEndBlock.ConsensusParamUpdates.Validator
// end of .ResultEndBlock.ConsensusParamUpdates.XXX_NoUnkeyedLiteral
length = len(in.ResultEndBlock.ConsensusParamUpdates.XXX_unrecognized)
out.ResultEndBlock.ConsensusParamUpdates.XXX_unrecognized = make([]uint8, length)
copy(out.ResultEndBlock.ConsensusParamUpdates.XXX_unrecognized[:], in.ResultEndBlock.ConsensusParamUpdates.XXX_unrecognized[:])
out.ResultEndBlock.ConsensusParamUpdates.XXX_sizecache = in.ResultEndBlock.ConsensusParamUpdates.XXX_sizecache
// end of .ResultEndBlock.ConsensusParamUpdates
length = len(in.ResultEndBlock.Events)
out.ResultEndBlock.Events = make([]Event, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
out.ResultEndBlock.Events[_0] = DeepCopyEvent(in.ResultEndBlock.Events[_0])
}
// end of .ResultEndBlock.XXX_NoUnkeyedLiteral
length = len(in.ResultEndBlock.XXX_unrecognized)
out.ResultEndBlock.XXX_unrecognized = make([]uint8, length)
copy(out.ResultEndBlock.XXX_unrecognized[:], in.ResultEndBlock.XXX_unrecognized[:])
out.ResultEndBlock.XXX_sizecache = in.ResultEndBlock.XXX_sizecache
// end of .ResultEndBlock
return
} //End of DeepCopyEventDataNewBlockHeader

// Non-Interface
func EncodeEventDataTx(w *[]byte, v EventDataTx) {
codonEncodeVarint(w, int64(v.TxResult.Height))
codonEncodeUvarint(w, uint64(v.TxResult.Index))
codonEncodeByteSlice(w, v.TxResult.Tx[:])
codonEncodeUvarint(w, uint64(v.TxResult.Result.Code))
codonEncodeByteSlice(w, v.TxResult.Result.Data[:])
codonEncodeString(w, v.TxResult.Result.Log)
codonEncodeString(w, v.TxResult.Result.Info)
codonEncodeVarint(w, int64(v.TxResult.Result.GasWanted))
codonEncodeVarint(w, int64(v.TxResult.Result.GasUsed))
codonEncodeVarint(w, int64(len(v.TxResult.Result.Events)))
for _0:=0; _0<len(v.TxResult.Result.Events); _0++ {
codonEncodeString(w, v.TxResult.Result.Events[_0].Type)
codonEncodeVarint(w, int64(len(v.TxResult.Result.Events[_0].Attributes)))
for _1:=0; _1<len(v.TxResult.Result.Events[_0].Attributes); _1++ {
codonEncodeByteSlice(w, v.TxResult.Result.Events[_0].Attributes[_1].Key[:])
codonEncodeByteSlice(w, v.TxResult.Result.Events[_0].Attributes[_1].Value[:])
// end of v.TxResult.Result.Events[_0].Attributes[_1].XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.TxResult.Result.Events[_0].Attributes[_1].XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.TxResult.Result.Events[_0].Attributes[_1].XXX_sizecache))
// end of v.TxResult.Result.Events[_0].Attributes[_1]
}
// end of v.TxResult.Result.Events[_0].XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.TxResult.Result.Events[_0].XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.TxResult.Result.Events[_0].XXX_sizecache))
// end of v.TxResult.Result.Events[_0]
}
codonEncodeString(w, v.TxResult.Result.Codespace)
// end of v.TxResult.Result.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.TxResult.Result.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.TxResult.Result.XXX_sizecache))
// end of v.TxResult.Result
// end of v.TxResult
} //End of EncodeEventDataTx

func DecodeEventDataTx(bz []byte) (EventDataTx, int, error) {
var err error
var length int
var v EventDataTx
var n int
var total int
v.TxResult.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Index = uint32(codonDecodeUint32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Tx, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Result.Code = uint32(codonDecodeUint32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Result.Data, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Result.Log = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Result.Info = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Result.GasWanted = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Result.GasUsed = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Result.Events = make([]Event, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.TxResult.Result.Events[_0], n, err = DecodeEvent(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
v.TxResult.Result.Codespace = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.TxResult.Result.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Result.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.TxResult.Result.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.TxResult.Result
// end of v.TxResult
return v, total, nil
} //End of DecodeEventDataTx

func RandEventDataTx(r RandSrc) EventDataTx {
var length int
var v EventDataTx
v.TxResult.Height = r.GetInt64()
v.TxResult.Index = r.GetUint32()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.TxResult.Tx = r.GetBytes(length)
v.TxResult.Result.Code = r.GetUint32()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.TxResult.Result.Data = r.GetBytes(length)
v.TxResult.Result.Log = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
v.TxResult.Result.Info = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
v.TxResult.Result.GasWanted = r.GetInt64()
v.TxResult.Result.GasUsed = r.GetInt64()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.TxResult.Result.Events = make([]Event, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.TxResult.Result.Events[_0] = RandEvent(r)
}
v.TxResult.Result.Codespace = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
// end of v.TxResult.Result.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.TxResult.Result.XXX_unrecognized = r.GetBytes(length)
v.TxResult.Result.XXX_sizecache = r.GetInt32()
// end of v.TxResult.Result
// end of v.TxResult
return v
} //End of RandEventDataTx

func DeepCopyEventDataTx(in EventDataTx) (out EventDataTx) {
var length int
out.TxResult.Height = in.TxResult.Height
out.TxResult.Index = in.TxResult.Index
length = len(in.TxResult.Tx)
out.TxResult.Tx = make([]uint8, length)
copy(out.TxResult.Tx[:], in.TxResult.Tx[:])
out.TxResult.Result.Code = in.TxResult.Result.Code
length = len(in.TxResult.Result.Data)
out.TxResult.Result.Data = make([]uint8, length)
copy(out.TxResult.Result.Data[:], in.TxResult.Result.Data[:])
out.TxResult.Result.Log = in.TxResult.Result.Log
out.TxResult.Result.Info = in.TxResult.Result.Info
out.TxResult.Result.GasWanted = in.TxResult.Result.GasWanted
out.TxResult.Result.GasUsed = in.TxResult.Result.GasUsed
length = len(in.TxResult.Result.Events)
out.TxResult.Result.Events = make([]Event, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
out.TxResult.Result.Events[_0] = DeepCopyEvent(in.TxResult.Result.Events[_0])
}
out.TxResult.Result.Codespace = in.TxResult.Result.Codespace
// end of .TxResult.Result.XXX_NoUnkeyedLiteral
length = len(in.TxResult.Result.XXX_unrecognized)
out.TxResult.Result.XXX_unrecognized = make([]uint8, length)
copy(out.TxResult.Result.XXX_unrecognized[:], in.TxResult.Result.XXX_unrecognized[:])
out.TxResult.Result.XXX_sizecache = in.TxResult.Result.XXX_sizecache
// end of .TxResult.Result
// end of .TxResult
return
} //End of DeepCopyEventDataTx

// Non-Interface
func EncodeEventDataNewRound(w *[]byte, v EventDataNewRound) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeString(w, v.Step)
codonEncodeByteSlice(w, v.Proposer.Address[:])
codonEncodeVarint(w, int64(v.Proposer.Index))
// end of v.Proposer
} //End of EncodeEventDataNewRound

func DecodeEventDataNewRound(bz []byte) (EventDataNewRound, int, error) {
var err error
var length int
var v EventDataNewRound
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Step = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposer.Address, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposer.Index = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Proposer
return v, total, nil
} //End of DecodeEventDataNewRound

func RandEventDataNewRound(r RandSrc) EventDataNewRound {
var length int
var v EventDataNewRound
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.Step = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposer.Address = r.GetBytes(length)
v.Proposer.Index = r.GetInt()
// end of v.Proposer
return v
} //End of RandEventDataNewRound

func DeepCopyEventDataNewRound(in EventDataNewRound) (out EventDataNewRound) {
var length int
out.Height = in.Height
out.Round = in.Round
out.Step = in.Step
length = len(in.Proposer.Address)
out.Proposer.Address = make([]uint8, length)
copy(out.Proposer.Address[:], in.Proposer.Address[:])
out.Proposer.Index = in.Proposer.Index
// end of .Proposer
return
} //End of DeepCopyEventDataNewRound

// Non-Interface
func EncodeEventDataCompleteProposal(w *[]byte, v EventDataCompleteProposal) {
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeString(w, v.Step)
codonEncodeByteSlice(w, v.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.BlockID.PartsHeader.Hash[:])
// end of v.BlockID.PartsHeader
// end of v.BlockID
} //End of EncodeEventDataCompleteProposal

func DecodeEventDataCompleteProposal(bz []byte) (EventDataCompleteProposal, int, error) {
var err error
var length int
var v EventDataCompleteProposal
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Step = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.BlockID.PartsHeader
// end of v.BlockID
return v, total, nil
} //End of DecodeEventDataCompleteProposal

func RandEventDataCompleteProposal(r RandSrc) EventDataCompleteProposal {
var length int
var v EventDataCompleteProposal
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.Step = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.Hash = r.GetBytes(length)
v.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.BlockID.PartsHeader
// end of v.BlockID
return v
} //End of RandEventDataCompleteProposal

func DeepCopyEventDataCompleteProposal(in EventDataCompleteProposal) (out EventDataCompleteProposal) {
var length int
out.Height = in.Height
out.Round = in.Round
out.Step = in.Step
length = len(in.BlockID.Hash)
out.BlockID.Hash = make([]uint8, length)
copy(out.BlockID.Hash[:], in.BlockID.Hash[:])
out.BlockID.PartsHeader.Total = in.BlockID.PartsHeader.Total
length = len(in.BlockID.PartsHeader.Hash)
out.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.BlockID.PartsHeader.Hash[:], in.BlockID.PartsHeader.Hash[:])
// end of .BlockID.PartsHeader
// end of .BlockID
return
} //End of DeepCopyEventDataCompleteProposal

// Non-Interface
func EncodeEventDataVote(w *[]byte, v EventDataVote) {
codonEncodeUint8(w, uint8(v.Vote.Type))
codonEncodeVarint(w, int64(v.Vote.Height))
codonEncodeVarint(w, int64(v.Vote.Round))
codonEncodeByteSlice(w, v.Vote.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.Vote.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.Vote.BlockID.PartsHeader.Hash[:])
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
EncodeTime(w, v.Vote.Timestamp)
codonEncodeByteSlice(w, v.Vote.ValidatorAddress[:])
codonEncodeVarint(w, int64(v.Vote.ValidatorIndex))
codonEncodeByteSlice(w, v.Vote.Signature[:])
// end of v.Vote
} //End of EncodeEventDataVote

func DecodeEventDataVote(bz []byte) (EventDataVote, int, error) {
var err error
var length int
var v EventDataVote
var n int
var total int
v.Vote = &Vote{}
v.Vote.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
v.Vote.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Vote
return v, total, nil
} //End of DecodeEventDataVote

func RandEventDataVote(r RandSrc) EventDataVote {
var length int
var v EventDataVote
v.Vote = &Vote{}
v.Vote.Type = SignedMsgType(r.GetUint8())
v.Vote.Height = r.GetInt64()
v.Vote.Round = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.BlockID.Hash = r.GetBytes(length)
v.Vote.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
v.Vote.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.ValidatorAddress = r.GetBytes(length)
v.Vote.ValidatorIndex = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.Signature = r.GetBytes(length)
// end of v.Vote
return v
} //End of RandEventDataVote

func DeepCopyEventDataVote(in EventDataVote) (out EventDataVote) {
var length int
out.Vote = &Vote{}
out.Vote.Type = in.Vote.Type
out.Vote.Height = in.Vote.Height
out.Vote.Round = in.Vote.Round
length = len(in.Vote.BlockID.Hash)
out.Vote.BlockID.Hash = make([]uint8, length)
copy(out.Vote.BlockID.Hash[:], in.Vote.BlockID.Hash[:])
out.Vote.BlockID.PartsHeader.Total = in.Vote.BlockID.PartsHeader.Total
length = len(in.Vote.BlockID.PartsHeader.Hash)
out.Vote.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.Vote.BlockID.PartsHeader.Hash[:], in.Vote.BlockID.PartsHeader.Hash[:])
// end of .Vote.BlockID.PartsHeader
// end of .Vote.BlockID
out.Vote.Timestamp = DeepCopyTime(in.Vote.Timestamp)
length = len(in.Vote.ValidatorAddress)
out.Vote.ValidatorAddress = make([]uint8, length)
copy(out.Vote.ValidatorAddress[:], in.Vote.ValidatorAddress[:])
out.Vote.ValidatorIndex = in.Vote.ValidatorIndex
length = len(in.Vote.Signature)
out.Vote.Signature = make([]uint8, length)
copy(out.Vote.Signature[:], in.Vote.Signature[:])
// end of .Vote
return
} //End of DeepCopyEventDataVote

// Non-Interface
func EncodeEventDataValidatorSetUpdates(w *[]byte, v EventDataValidatorSetUpdates) {
codonEncodeVarint(w, int64(len(v.ValidatorUpdates)))
for _0:=0; _0<len(v.ValidatorUpdates); _0++ {
codonEncodeByteSlice(w, v.ValidatorUpdates[_0].Address[:])
EncodePubKey(w, v.ValidatorUpdates[_0].PubKey) // interface_encode
codonEncodeVarint(w, int64(v.ValidatorUpdates[_0].VotingPower))
codonEncodeVarint(w, int64(v.ValidatorUpdates[_0].ProposerPriority))
// end of v.ValidatorUpdates[_0]
}
} //End of EncodeEventDataValidatorSetUpdates

func DecodeEventDataValidatorSetUpdates(bz []byte) (EventDataValidatorSetUpdates, int, error) {
var err error
var length int
var v EventDataValidatorSetUpdates
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ValidatorUpdates = make([]*Validator, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of ptr
v.ValidatorUpdates[_0] = &Validator{}
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ValidatorUpdates[_0].Address, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ValidatorUpdates[_0].PubKey, n, err = DecodePubKey(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n // interface_decode
v.ValidatorUpdates[_0].VotingPower = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ValidatorUpdates[_0].ProposerPriority = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.ValidatorUpdates[_0]
}
return v, total, nil
} //End of DecodeEventDataValidatorSetUpdates

func RandEventDataValidatorSetUpdates(r RandSrc) EventDataValidatorSetUpdates {
var length int
var v EventDataValidatorSetUpdates
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ValidatorUpdates = make([]*Validator, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of ptr
v.ValidatorUpdates[_0] = &Validator{}
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ValidatorUpdates[_0].Address = r.GetBytes(length)
v.ValidatorUpdates[_0].PubKey = RandPubKey(r) // interface_decode
v.ValidatorUpdates[_0].VotingPower = r.GetInt64()
v.ValidatorUpdates[_0].ProposerPriority = r.GetInt64()
// end of v.ValidatorUpdates[_0]
}
return v
} //End of RandEventDataValidatorSetUpdates

func DeepCopyEventDataValidatorSetUpdates(in EventDataValidatorSetUpdates) (out EventDataValidatorSetUpdates) {
var length int
length = len(in.ValidatorUpdates)
out.ValidatorUpdates = make([]*Validator, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of ptr
out.ValidatorUpdates[_0] = &Validator{}
length = len(in.ValidatorUpdates[_0].Address)
out.ValidatorUpdates[_0].Address = make([]uint8, length)
copy(out.ValidatorUpdates[_0].Address[:], in.ValidatorUpdates[_0].Address[:])
out.ValidatorUpdates[_0].PubKey = DeepCopyPubKey(in.ValidatorUpdates[_0].PubKey)
out.ValidatorUpdates[_0].VotingPower = in.ValidatorUpdates[_0].VotingPower
out.ValidatorUpdates[_0].ProposerPriority = in.ValidatorUpdates[_0].ProposerPriority
// end of .ValidatorUpdates[_0]
}
return
} //End of DeepCopyEventDataValidatorSetUpdates

// Non-Interface
func EncodeEventDataString(w *[]byte, v EventDataString) {
codonEncodeString(w, string(v))
} //End of EncodeEventDataString

func DecodeEventDataString(bz []byte) (EventDataString, int, error) {
var err error
var v EventDataString
var n int
var total int
v = EventDataString(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeEventDataString

func RandEventDataString(r RandSrc) EventDataString {
var v EventDataString
v = EventDataString(r.GetString(1+int(r.GetUint()%(MaxStringLength-1))))
return v
} //End of RandEventDataString

func DeepCopyEventDataString(in EventDataString) (out EventDataString) {
out = in
return
} //End of DeepCopyEventDataString

// Non-Interface
func EncodeMockGoodEvidence(w *[]byte, v MockGoodEvidence) {
codonEncodeVarint(w, int64(v.Height_))
codonEncodeByteSlice(w, v.Address_[:])
} //End of EncodeMockGoodEvidence

func DecodeMockGoodEvidence(bz []byte) (MockGoodEvidence, int, error) {
var err error
var length int
var v MockGoodEvidence
var n int
var total int
v.Height_ = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Address_, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeMockGoodEvidence

func RandMockGoodEvidence(r RandSrc) MockGoodEvidence {
var length int
var v MockGoodEvidence
v.Height_ = r.GetInt64()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Address_ = r.GetBytes(length)
return v
} //End of RandMockGoodEvidence

func DeepCopyMockGoodEvidence(in MockGoodEvidence) (out MockGoodEvidence) {
var length int
out.Height_ = in.Height_
length = len(in.Address_)
out.Address_ = make([]uint8, length)
copy(out.Address_[:], in.Address_[:])
return
} //End of DeepCopyMockGoodEvidence

// Non-Interface
func EncodeMockBadEvidence(w *[]byte, v MockBadEvidence) {
codonEncodeVarint(w, int64(v.MockGoodEvidence.Height_))
codonEncodeByteSlice(w, v.MockGoodEvidence.Address_[:])
// end of v.MockGoodEvidence
} //End of EncodeMockBadEvidence

func DecodeMockBadEvidence(bz []byte) (MockBadEvidence, int, error) {
var err error
var length int
var v MockBadEvidence
var n int
var total int
v.MockGoodEvidence.Height_ = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.MockGoodEvidence.Address_, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.MockGoodEvidence
return v, total, nil
} //End of DecodeMockBadEvidence

func RandMockBadEvidence(r RandSrc) MockBadEvidence {
var length int
var v MockBadEvidence
v.MockGoodEvidence.Height_ = r.GetInt64()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.MockGoodEvidence.Address_ = r.GetBytes(length)
// end of v.MockGoodEvidence
return v
} //End of RandMockBadEvidence

func DeepCopyMockBadEvidence(in MockBadEvidence) (out MockBadEvidence) {
var length int
out.MockGoodEvidence.Height_ = in.MockGoodEvidence.Height_
length = len(in.MockGoodEvidence.Address_)
out.MockGoodEvidence.Address_ = make([]uint8, length)
copy(out.MockGoodEvidence.Address_[:], in.MockGoodEvidence.Address_[:])
// end of .MockGoodEvidence
return
} //End of DeepCopyMockBadEvidence

// Non-Interface
func EncodeEvidenceListMessage(w *[]byte, v EvidenceListMessage) {
codonEncodeVarint(w, int64(len(v.Evidence)))
for _0:=0; _0<len(v.Evidence); _0++ {
EncodeEvidence(w, v.Evidence[_0]) // interface_encode
}
} //End of EncodeEvidenceListMessage

func DecodeEvidenceListMessage(bz []byte) (EvidenceListMessage, int, error) {
var err error
var length int
var v EvidenceListMessage
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Evidence = make([]Evidence, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of interface
v.Evidence[_0], n, err = DecodeEvidence(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
return v, total, nil
} //End of DecodeEvidenceListMessage

func RandEvidenceListMessage(r RandSrc) EvidenceListMessage {
var length int
var v EvidenceListMessage
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Evidence = make([]Evidence, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of interface
v.Evidence[_0] = RandEvidence(r)
}
return v
} //End of RandEvidenceListMessage

func DeepCopyEvidenceListMessage(in EvidenceListMessage) (out EvidenceListMessage) {
var length int
length = len(in.Evidence)
out.Evidence = make([]Evidence, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of interface
out.Evidence[_0] = DeepCopyEvidence(in.Evidence[_0])
}
return
} //End of DeepCopyEvidenceListMessage

// Non-Interface
func EncodePubKeyResponse(w *[]byte, v PubKeyResponse) {
EncodePubKey(w, v.PubKey) // interface_encode
codonEncodeVarint(w, int64(v.Error.Code))
codonEncodeString(w, v.Error.Description)
// end of v.Error
} //End of EncodePubKeyResponse

func DecodePubKeyResponse(bz []byte) (PubKeyResponse, int, error) {
var err error
var v PubKeyResponse
var n int
var total int
v.PubKey, n, err = DecodePubKey(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n // interface_decode
v.Error = &RemoteSignerError{}
v.Error.Code = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Error.Description = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Error
return v, total, nil
} //End of DecodePubKeyResponse

func RandPubKeyResponse(r RandSrc) PubKeyResponse {
var v PubKeyResponse
v.PubKey = RandPubKey(r) // interface_decode
v.Error = &RemoteSignerError{}
v.Error.Code = r.GetInt()
v.Error.Description = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
// end of v.Error
return v
} //End of RandPubKeyResponse

func DeepCopyPubKeyResponse(in PubKeyResponse) (out PubKeyResponse) {
out.PubKey = DeepCopyPubKey(in.PubKey)
out.Error = &RemoteSignerError{}
out.Error.Code = in.Error.Code
out.Error.Description = in.Error.Description
// end of .Error
return
} //End of DeepCopyPubKeyResponse

// Non-Interface
func EncodeSignVoteRequest(w *[]byte, v SignVoteRequest) {
codonEncodeUint8(w, uint8(v.Vote.Type))
codonEncodeVarint(w, int64(v.Vote.Height))
codonEncodeVarint(w, int64(v.Vote.Round))
codonEncodeByteSlice(w, v.Vote.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.Vote.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.Vote.BlockID.PartsHeader.Hash[:])
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
EncodeTime(w, v.Vote.Timestamp)
codonEncodeByteSlice(w, v.Vote.ValidatorAddress[:])
codonEncodeVarint(w, int64(v.Vote.ValidatorIndex))
codonEncodeByteSlice(w, v.Vote.Signature[:])
// end of v.Vote
} //End of EncodeSignVoteRequest

func DecodeSignVoteRequest(bz []byte) (SignVoteRequest, int, error) {
var err error
var length int
var v SignVoteRequest
var n int
var total int
v.Vote = &Vote{}
v.Vote.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
v.Vote.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Vote
return v, total, nil
} //End of DecodeSignVoteRequest

func RandSignVoteRequest(r RandSrc) SignVoteRequest {
var length int
var v SignVoteRequest
v.Vote = &Vote{}
v.Vote.Type = SignedMsgType(r.GetUint8())
v.Vote.Height = r.GetInt64()
v.Vote.Round = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.BlockID.Hash = r.GetBytes(length)
v.Vote.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
v.Vote.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.ValidatorAddress = r.GetBytes(length)
v.Vote.ValidatorIndex = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.Signature = r.GetBytes(length)
// end of v.Vote
return v
} //End of RandSignVoteRequest

func DeepCopySignVoteRequest(in SignVoteRequest) (out SignVoteRequest) {
var length int
out.Vote = &Vote{}
out.Vote.Type = in.Vote.Type
out.Vote.Height = in.Vote.Height
out.Vote.Round = in.Vote.Round
length = len(in.Vote.BlockID.Hash)
out.Vote.BlockID.Hash = make([]uint8, length)
copy(out.Vote.BlockID.Hash[:], in.Vote.BlockID.Hash[:])
out.Vote.BlockID.PartsHeader.Total = in.Vote.BlockID.PartsHeader.Total
length = len(in.Vote.BlockID.PartsHeader.Hash)
out.Vote.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.Vote.BlockID.PartsHeader.Hash[:], in.Vote.BlockID.PartsHeader.Hash[:])
// end of .Vote.BlockID.PartsHeader
// end of .Vote.BlockID
out.Vote.Timestamp = DeepCopyTime(in.Vote.Timestamp)
length = len(in.Vote.ValidatorAddress)
out.Vote.ValidatorAddress = make([]uint8, length)
copy(out.Vote.ValidatorAddress[:], in.Vote.ValidatorAddress[:])
out.Vote.ValidatorIndex = in.Vote.ValidatorIndex
length = len(in.Vote.Signature)
out.Vote.Signature = make([]uint8, length)
copy(out.Vote.Signature[:], in.Vote.Signature[:])
// end of .Vote
return
} //End of DeepCopySignVoteRequest

// Non-Interface
func EncodeSignedVoteResponse(w *[]byte, v SignedVoteResponse) {
codonEncodeUint8(w, uint8(v.Vote.Type))
codonEncodeVarint(w, int64(v.Vote.Height))
codonEncodeVarint(w, int64(v.Vote.Round))
codonEncodeByteSlice(w, v.Vote.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.Vote.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.Vote.BlockID.PartsHeader.Hash[:])
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
EncodeTime(w, v.Vote.Timestamp)
codonEncodeByteSlice(w, v.Vote.ValidatorAddress[:])
codonEncodeVarint(w, int64(v.Vote.ValidatorIndex))
codonEncodeByteSlice(w, v.Vote.Signature[:])
// end of v.Vote
codonEncodeVarint(w, int64(v.Error.Code))
codonEncodeString(w, v.Error.Description)
// end of v.Error
} //End of EncodeSignedVoteResponse

func DecodeSignedVoteResponse(bz []byte) (SignedVoteResponse, int, error) {
var err error
var length int
var v SignedVoteResponse
var n int
var total int
v.Vote = &Vote{}
v.Vote.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
v.Vote.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Vote.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Vote
v.Error = &RemoteSignerError{}
v.Error.Code = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Error.Description = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Error
return v, total, nil
} //End of DecodeSignedVoteResponse

func RandSignedVoteResponse(r RandSrc) SignedVoteResponse {
var length int
var v SignedVoteResponse
v.Vote = &Vote{}
v.Vote.Type = SignedMsgType(r.GetUint8())
v.Vote.Height = r.GetInt64()
v.Vote.Round = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.BlockID.Hash = r.GetBytes(length)
v.Vote.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.Vote.BlockID.PartsHeader
// end of v.Vote.BlockID
v.Vote.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.ValidatorAddress = r.GetBytes(length)
v.Vote.ValidatorIndex = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Vote.Signature = r.GetBytes(length)
// end of v.Vote
v.Error = &RemoteSignerError{}
v.Error.Code = r.GetInt()
v.Error.Description = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
// end of v.Error
return v
} //End of RandSignedVoteResponse

func DeepCopySignedVoteResponse(in SignedVoteResponse) (out SignedVoteResponse) {
var length int
out.Vote = &Vote{}
out.Vote.Type = in.Vote.Type
out.Vote.Height = in.Vote.Height
out.Vote.Round = in.Vote.Round
length = len(in.Vote.BlockID.Hash)
out.Vote.BlockID.Hash = make([]uint8, length)
copy(out.Vote.BlockID.Hash[:], in.Vote.BlockID.Hash[:])
out.Vote.BlockID.PartsHeader.Total = in.Vote.BlockID.PartsHeader.Total
length = len(in.Vote.BlockID.PartsHeader.Hash)
out.Vote.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.Vote.BlockID.PartsHeader.Hash[:], in.Vote.BlockID.PartsHeader.Hash[:])
// end of .Vote.BlockID.PartsHeader
// end of .Vote.BlockID
out.Vote.Timestamp = DeepCopyTime(in.Vote.Timestamp)
length = len(in.Vote.ValidatorAddress)
out.Vote.ValidatorAddress = make([]uint8, length)
copy(out.Vote.ValidatorAddress[:], in.Vote.ValidatorAddress[:])
out.Vote.ValidatorIndex = in.Vote.ValidatorIndex
length = len(in.Vote.Signature)
out.Vote.Signature = make([]uint8, length)
copy(out.Vote.Signature[:], in.Vote.Signature[:])
// end of .Vote
out.Error = &RemoteSignerError{}
out.Error.Code = in.Error.Code
out.Error.Description = in.Error.Description
// end of .Error
return
} //End of DeepCopySignedVoteResponse

// Non-Interface
func EncodeSignProposalRequest(w *[]byte, v SignProposalRequest) {
codonEncodeUint8(w, uint8(v.Proposal.Type))
codonEncodeVarint(w, int64(v.Proposal.Height))
codonEncodeVarint(w, int64(v.Proposal.Round))
codonEncodeVarint(w, int64(v.Proposal.POLRound))
codonEncodeByteSlice(w, v.Proposal.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.Proposal.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.Proposal.BlockID.PartsHeader.Hash[:])
// end of v.Proposal.BlockID.PartsHeader
// end of v.Proposal.BlockID
EncodeTime(w, v.Proposal.Timestamp)
codonEncodeByteSlice(w, v.Proposal.Signature[:])
// end of v.Proposal
} //End of EncodeSignProposalRequest

func DecodeSignProposalRequest(bz []byte) (SignProposalRequest, int, error) {
var err error
var length int
var v SignProposalRequest
var n int
var total int
v.Proposal = &Proposal{}
v.Proposal.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.POLRound = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Proposal.BlockID.PartsHeader
// end of v.Proposal.BlockID
v.Proposal.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Proposal
return v, total, nil
} //End of DecodeSignProposalRequest

func RandSignProposalRequest(r RandSrc) SignProposalRequest {
var length int
var v SignProposalRequest
v.Proposal = &Proposal{}
v.Proposal.Type = SignedMsgType(r.GetUint8())
v.Proposal.Height = r.GetInt64()
v.Proposal.Round = r.GetInt()
v.Proposal.POLRound = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposal.BlockID.Hash = r.GetBytes(length)
v.Proposal.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposal.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.Proposal.BlockID.PartsHeader
// end of v.Proposal.BlockID
v.Proposal.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposal.Signature = r.GetBytes(length)
// end of v.Proposal
return v
} //End of RandSignProposalRequest

func DeepCopySignProposalRequest(in SignProposalRequest) (out SignProposalRequest) {
var length int
out.Proposal = &Proposal{}
out.Proposal.Type = in.Proposal.Type
out.Proposal.Height = in.Proposal.Height
out.Proposal.Round = in.Proposal.Round
out.Proposal.POLRound = in.Proposal.POLRound
length = len(in.Proposal.BlockID.Hash)
out.Proposal.BlockID.Hash = make([]uint8, length)
copy(out.Proposal.BlockID.Hash[:], in.Proposal.BlockID.Hash[:])
out.Proposal.BlockID.PartsHeader.Total = in.Proposal.BlockID.PartsHeader.Total
length = len(in.Proposal.BlockID.PartsHeader.Hash)
out.Proposal.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.Proposal.BlockID.PartsHeader.Hash[:], in.Proposal.BlockID.PartsHeader.Hash[:])
// end of .Proposal.BlockID.PartsHeader
// end of .Proposal.BlockID
out.Proposal.Timestamp = DeepCopyTime(in.Proposal.Timestamp)
length = len(in.Proposal.Signature)
out.Proposal.Signature = make([]uint8, length)
copy(out.Proposal.Signature[:], in.Proposal.Signature[:])
// end of .Proposal
return
} //End of DeepCopySignProposalRequest

// Non-Interface
func EncodeSignedProposalResponse(w *[]byte, v SignedProposalResponse) {
codonEncodeUint8(w, uint8(v.Proposal.Type))
codonEncodeVarint(w, int64(v.Proposal.Height))
codonEncodeVarint(w, int64(v.Proposal.Round))
codonEncodeVarint(w, int64(v.Proposal.POLRound))
codonEncodeByteSlice(w, v.Proposal.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.Proposal.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.Proposal.BlockID.PartsHeader.Hash[:])
// end of v.Proposal.BlockID.PartsHeader
// end of v.Proposal.BlockID
EncodeTime(w, v.Proposal.Timestamp)
codonEncodeByteSlice(w, v.Proposal.Signature[:])
// end of v.Proposal
codonEncodeVarint(w, int64(v.Error.Code))
codonEncodeString(w, v.Error.Description)
// end of v.Error
} //End of EncodeSignedProposalResponse

func DecodeSignedProposalResponse(bz []byte) (SignedProposalResponse, int, error) {
var err error
var length int
var v SignedProposalResponse
var n int
var total int
v.Proposal = &Proposal{}
v.Proposal.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.POLRound = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Proposal.BlockID.PartsHeader
// end of v.Proposal.BlockID
v.Proposal.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proposal.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Proposal
v.Error = &RemoteSignerError{}
v.Error.Code = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Error.Description = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Error
return v, total, nil
} //End of DecodeSignedProposalResponse

func RandSignedProposalResponse(r RandSrc) SignedProposalResponse {
var length int
var v SignedProposalResponse
v.Proposal = &Proposal{}
v.Proposal.Type = SignedMsgType(r.GetUint8())
v.Proposal.Height = r.GetInt64()
v.Proposal.Round = r.GetInt()
v.Proposal.POLRound = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposal.BlockID.Hash = r.GetBytes(length)
v.Proposal.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposal.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.Proposal.BlockID.PartsHeader
// end of v.Proposal.BlockID
v.Proposal.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proposal.Signature = r.GetBytes(length)
// end of v.Proposal
v.Error = &RemoteSignerError{}
v.Error.Code = r.GetInt()
v.Error.Description = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
// end of v.Error
return v
} //End of RandSignedProposalResponse

func DeepCopySignedProposalResponse(in SignedProposalResponse) (out SignedProposalResponse) {
var length int
out.Proposal = &Proposal{}
out.Proposal.Type = in.Proposal.Type
out.Proposal.Height = in.Proposal.Height
out.Proposal.Round = in.Proposal.Round
out.Proposal.POLRound = in.Proposal.POLRound
length = len(in.Proposal.BlockID.Hash)
out.Proposal.BlockID.Hash = make([]uint8, length)
copy(out.Proposal.BlockID.Hash[:], in.Proposal.BlockID.Hash[:])
out.Proposal.BlockID.PartsHeader.Total = in.Proposal.BlockID.PartsHeader.Total
length = len(in.Proposal.BlockID.PartsHeader.Hash)
out.Proposal.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.Proposal.BlockID.PartsHeader.Hash[:], in.Proposal.BlockID.PartsHeader.Hash[:])
// end of .Proposal.BlockID.PartsHeader
// end of .Proposal.BlockID
out.Proposal.Timestamp = DeepCopyTime(in.Proposal.Timestamp)
length = len(in.Proposal.Signature)
out.Proposal.Signature = make([]uint8, length)
copy(out.Proposal.Signature[:], in.Proposal.Signature[:])
// end of .Proposal
out.Error = &RemoteSignerError{}
out.Error.Code = in.Error.Code
out.Error.Description = in.Error.Description
// end of .Error
return
} //End of DeepCopySignedProposalResponse

// Non-Interface
func EncodePacketMsg(w *[]byte, v PacketMsg) {
codonEncodeUint8(w, v.ChannelID)
codonEncodeUint8(w, v.EOF)
codonEncodeByteSlice(w, v.Bytes[:])
} //End of EncodePacketMsg

func DecodePacketMsg(bz []byte) (PacketMsg, int, error) {
var err error
var length int
var v PacketMsg
var n int
var total int
v.ChannelID = uint8(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.EOF = uint8(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Bytes, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodePacketMsg

func RandPacketMsg(r RandSrc) PacketMsg {
var length int
var v PacketMsg
v.ChannelID = r.GetUint8()
v.EOF = r.GetUint8()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Bytes = r.GetBytes(length)
return v
} //End of RandPacketMsg

func DeepCopyPacketMsg(in PacketMsg) (out PacketMsg) {
var length int
out.ChannelID = in.ChannelID
out.EOF = in.EOF
length = len(in.Bytes)
out.Bytes = make([]uint8, length)
copy(out.Bytes[:], in.Bytes[:])
return
} //End of DeepCopyPacketMsg

// Non-Interface
func EncodeTxMessage(w *[]byte, v TxMessage) {
codonEncodeByteSlice(w, v.Tx[:])
} //End of EncodeTxMessage

func DecodeTxMessage(bz []byte) (TxMessage, int, error) {
var err error
var length int
var v TxMessage
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Tx, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeTxMessage

func RandTxMessage(r RandSrc) TxMessage {
var length int
var v TxMessage
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Tx = r.GetBytes(length)
return v
} //End of RandTxMessage

func DeepCopyTxMessage(in TxMessage) (out TxMessage) {
var length int
length = len(in.Tx)
out.Tx = make([]uint8, length)
copy(out.Tx[:], in.Tx[:])
return
} //End of DeepCopyTxMessage

// Non-Interface
func EncodeBc1BlockRequestMessage(w *[]byte, v Bc1BlockRequestMessage) {
codonEncodeVarint(w, int64(v.Height))
} //End of EncodeBc1BlockRequestMessage

func DecodeBc1BlockRequestMessage(bz []byte) (Bc1BlockRequestMessage, int, error) {
var err error
var v Bc1BlockRequestMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeBc1BlockRequestMessage

func RandBc1BlockRequestMessage(r RandSrc) Bc1BlockRequestMessage {
var v Bc1BlockRequestMessage
v.Height = r.GetInt64()
return v
} //End of RandBc1BlockRequestMessage

func DeepCopyBc1BlockRequestMessage(in Bc1BlockRequestMessage) (out Bc1BlockRequestMessage) {
out.Height = in.Height
return
} //End of DeepCopyBc1BlockRequestMessage

// Non-Interface
func EncodeBc1NoBlockResponseMessage(w *[]byte, v Bc1NoBlockResponseMessage) {
codonEncodeVarint(w, int64(v.Height))
} //End of EncodeBc1NoBlockResponseMessage

func DecodeBc1NoBlockResponseMessage(bz []byte) (Bc1NoBlockResponseMessage, int, error) {
var err error
var v Bc1NoBlockResponseMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeBc1NoBlockResponseMessage

func RandBc1NoBlockResponseMessage(r RandSrc) Bc1NoBlockResponseMessage {
var v Bc1NoBlockResponseMessage
v.Height = r.GetInt64()
return v
} //End of RandBc1NoBlockResponseMessage

func DeepCopyBc1NoBlockResponseMessage(in Bc1NoBlockResponseMessage) (out Bc1NoBlockResponseMessage) {
out.Height = in.Height
return
} //End of DeepCopyBc1NoBlockResponseMessage

// Non-Interface
func EncodeBc1StatusResponseMessage(w *[]byte, v Bc1StatusResponseMessage) {
codonEncodeVarint(w, int64(v.Height))
} //End of EncodeBc1StatusResponseMessage

func DecodeBc1StatusResponseMessage(bz []byte) (Bc1StatusResponseMessage, int, error) {
var err error
var v Bc1StatusResponseMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeBc1StatusResponseMessage

func RandBc1StatusResponseMessage(r RandSrc) Bc1StatusResponseMessage {
var v Bc1StatusResponseMessage
v.Height = r.GetInt64()
return v
} //End of RandBc1StatusResponseMessage

func DeepCopyBc1StatusResponseMessage(in Bc1StatusResponseMessage) (out Bc1StatusResponseMessage) {
out.Height = in.Height
return
} //End of DeepCopyBc1StatusResponseMessage

// Non-Interface
func EncodeBc1StatusRequestMessage(w *[]byte, v Bc1StatusRequestMessage) {
codonEncodeVarint(w, int64(v.Height))
} //End of EncodeBc1StatusRequestMessage

func DecodeBc1StatusRequestMessage(bz []byte) (Bc1StatusRequestMessage, int, error) {
var err error
var v Bc1StatusRequestMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeBc1StatusRequestMessage

func RandBc1StatusRequestMessage(r RandSrc) Bc1StatusRequestMessage {
var v Bc1StatusRequestMessage
v.Height = r.GetInt64()
return v
} //End of RandBc1StatusRequestMessage

func DeepCopyBc1StatusRequestMessage(in Bc1StatusRequestMessage) (out Bc1StatusRequestMessage) {
out.Height = in.Height
return
} //End of DeepCopyBc1StatusRequestMessage

// Non-Interface
func EncodeBc0BlockRequestMessage(w *[]byte, v Bc0BlockRequestMessage) {
codonEncodeVarint(w, int64(v.Height))
} //End of EncodeBc0BlockRequestMessage

func DecodeBc0BlockRequestMessage(bz []byte) (Bc0BlockRequestMessage, int, error) {
var err error
var v Bc0BlockRequestMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeBc0BlockRequestMessage

func RandBc0BlockRequestMessage(r RandSrc) Bc0BlockRequestMessage {
var v Bc0BlockRequestMessage
v.Height = r.GetInt64()
return v
} //End of RandBc0BlockRequestMessage

func DeepCopyBc0BlockRequestMessage(in Bc0BlockRequestMessage) (out Bc0BlockRequestMessage) {
out.Height = in.Height
return
} //End of DeepCopyBc0BlockRequestMessage

// Non-Interface
func EncodeBc0NoBlockResponseMessage(w *[]byte, v Bc0NoBlockResponseMessage) {
codonEncodeVarint(w, int64(v.Height))
} //End of EncodeBc0NoBlockResponseMessage

func DecodeBc0NoBlockResponseMessage(bz []byte) (Bc0NoBlockResponseMessage, int, error) {
var err error
var v Bc0NoBlockResponseMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeBc0NoBlockResponseMessage

func RandBc0NoBlockResponseMessage(r RandSrc) Bc0NoBlockResponseMessage {
var v Bc0NoBlockResponseMessage
v.Height = r.GetInt64()
return v
} //End of RandBc0NoBlockResponseMessage

func DeepCopyBc0NoBlockResponseMessage(in Bc0NoBlockResponseMessage) (out Bc0NoBlockResponseMessage) {
out.Height = in.Height
return
} //End of DeepCopyBc0NoBlockResponseMessage

// Non-Interface
func EncodeBc0StatusResponseMessage(w *[]byte, v Bc0StatusResponseMessage) {
codonEncodeVarint(w, int64(v.Height))
} //End of EncodeBc0StatusResponseMessage

func DecodeBc0StatusResponseMessage(bz []byte) (Bc0StatusResponseMessage, int, error) {
var err error
var v Bc0StatusResponseMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeBc0StatusResponseMessage

func RandBc0StatusResponseMessage(r RandSrc) Bc0StatusResponseMessage {
var v Bc0StatusResponseMessage
v.Height = r.GetInt64()
return v
} //End of RandBc0StatusResponseMessage

func DeepCopyBc0StatusResponseMessage(in Bc0StatusResponseMessage) (out Bc0StatusResponseMessage) {
out.Height = in.Height
return
} //End of DeepCopyBc0StatusResponseMessage

// Non-Interface
func EncodeBc0StatusRequestMessage(w *[]byte, v Bc0StatusRequestMessage) {
codonEncodeVarint(w, int64(v.Height))
} //End of EncodeBc0StatusRequestMessage

func DecodeBc0StatusRequestMessage(bz []byte) (Bc0StatusRequestMessage, int, error) {
var err error
var v Bc0StatusRequestMessage
var n int
var total int
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeBc0StatusRequestMessage

func RandBc0StatusRequestMessage(r RandSrc) Bc0StatusRequestMessage {
var v Bc0StatusRequestMessage
v.Height = r.GetInt64()
return v
} //End of RandBc0StatusRequestMessage

func DeepCopyBc0StatusRequestMessage(in Bc0StatusRequestMessage) (out Bc0StatusRequestMessage) {
out.Height = in.Height
return
} //End of DeepCopyBc0StatusRequestMessage

// Non-Interface
func EncodeMsgInfo(w *[]byte, v MsgInfo) {
EncodeConsensusMessage(w, v.Msg) // interface_encode
codonEncodeString(w, string(v.PeerID))
} //End of EncodeMsgInfo

func DecodeMsgInfo(bz []byte) (MsgInfo, int, error) {
var err error
var v MsgInfo
var n int
var total int
v.Msg, n, err = DecodeConsensusMessage(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n // interface_decode
v.PeerID = P2PID(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeMsgInfo

func RandMsgInfo(r RandSrc) MsgInfo {
var v MsgInfo
v.Msg = RandConsensusMessage(r) // interface_decode
v.PeerID = P2PID(r.GetString(1+int(r.GetUint()%(MaxStringLength-1))))
return v
} //End of RandMsgInfo

func DeepCopyMsgInfo(in MsgInfo) (out MsgInfo) {
out.Msg = DeepCopyConsensusMessage(in.Msg)
out.PeerID = in.PeerID
return
} //End of DeepCopyMsgInfo

// Non-Interface
func EncodeTimeoutInfo(w *[]byte, v TimeoutInfo) {
codonEncodeVarint(w, int64(v.Duration))
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeUint8(w, uint8(v.Step))
} //End of EncodeTimeoutInfo

func DecodeTimeoutInfo(bz []byte) (TimeoutInfo, int, error) {
var err error
var v TimeoutInfo
var n int
var total int
v.Duration = Duration(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Step = RoundStepType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeTimeoutInfo

func RandTimeoutInfo(r RandSrc) TimeoutInfo {
var v TimeoutInfo
v.Duration = Duration(r.GetInt64())
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.Step = RoundStepType(r.GetUint8())
return v
} //End of RandTimeoutInfo

func DeepCopyTimeoutInfo(in TimeoutInfo) (out TimeoutInfo) {
out.Duration = in.Duration
out.Height = in.Height
out.Round = in.Round
out.Step = in.Step
return
} //End of DeepCopyTimeoutInfo

// Non-Interface
func EncodeRoundStepType(w *[]byte, v RoundStepType) {
codonEncodeUint8(w, uint8(v))
} //End of EncodeRoundStepType

func DecodeRoundStepType(bz []byte) (RoundStepType, int, error) {
var err error
var v RoundStepType
var n int
var total int
v = RoundStepType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeRoundStepType

func RandRoundStepType(r RandSrc) RoundStepType {
var v RoundStepType
v = RoundStepType(r.GetUint8())
return v
} //End of RandRoundStepType

func DeepCopyRoundStepType(in RoundStepType) (out RoundStepType) {
out = in
return
} //End of DeepCopyRoundStepType

// Non-Interface
func EncodeBitArray(w *[]byte, v BitArray) {
codonEncodeVarint(w, int64(v.Bits))
codonEncodeVarint(w, int64(len(v.Elems)))
for _0:=0; _0<len(v.Elems); _0++ {
codonEncodeUvarint(w, uint64(v.Elems[_0]))
}
} //End of EncodeBitArray

func DecodeBitArray(bz []byte) (BitArray, int, error) {
var err error
var length int
var v BitArray
var n int
var total int
v.Bits = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
v.Elems[_0] = uint64(codonDecodeUint64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
return v, total, nil
} //End of DecodeBitArray

func RandBitArray(r RandSrc) BitArray {
var length int
var v BitArray
v.Bits = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
v.Elems[_0] = r.GetUint64()
}
return v
} //End of RandBitArray

func DeepCopyBitArray(in BitArray) (out BitArray) {
var length int
out.Bits = in.Bits
length = len(in.Elems)
out.Elems = make([]uint64, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of uint64
out.Elems[_0] = in.Elems[_0]
}
return
} //End of DeepCopyBitArray

// Non-Interface
func EncodeProposal(w *[]byte, v Proposal) {
codonEncodeUint8(w, uint8(v.Type))
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeVarint(w, int64(v.POLRound))
codonEncodeByteSlice(w, v.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.BlockID.PartsHeader.Hash[:])
// end of v.BlockID.PartsHeader
// end of v.BlockID
EncodeTime(w, v.Timestamp)
codonEncodeByteSlice(w, v.Signature[:])
} //End of EncodeProposal

func DecodeProposal(bz []byte) (Proposal, int, error) {
var err error
var length int
var v Proposal
var n int
var total int
v.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.POLRound = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.BlockID.PartsHeader
// end of v.BlockID
v.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeProposal

func RandProposal(r RandSrc) Proposal {
var length int
var v Proposal
v.Type = SignedMsgType(r.GetUint8())
v.Height = r.GetInt64()
v.Round = r.GetInt()
v.POLRound = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.Hash = r.GetBytes(length)
v.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.BlockID.PartsHeader
// end of v.BlockID
v.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Signature = r.GetBytes(length)
return v
} //End of RandProposal

func DeepCopyProposal(in Proposal) (out Proposal) {
var length int
out.Type = in.Type
out.Height = in.Height
out.Round = in.Round
out.POLRound = in.POLRound
length = len(in.BlockID.Hash)
out.BlockID.Hash = make([]uint8, length)
copy(out.BlockID.Hash[:], in.BlockID.Hash[:])
out.BlockID.PartsHeader.Total = in.BlockID.PartsHeader.Total
length = len(in.BlockID.PartsHeader.Hash)
out.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.BlockID.PartsHeader.Hash[:], in.BlockID.PartsHeader.Hash[:])
// end of .BlockID.PartsHeader
// end of .BlockID
out.Timestamp = DeepCopyTime(in.Timestamp)
length = len(in.Signature)
out.Signature = make([]uint8, length)
copy(out.Signature[:], in.Signature[:])
return
} //End of DeepCopyProposal

// Non-Interface
func EncodePart(w *[]byte, v Part) {
codonEncodeVarint(w, int64(v.Index))
codonEncodeByteSlice(w, v.Bytes[:])
codonEncodeVarint(w, int64(v.Proof.Total))
codonEncodeVarint(w, int64(v.Proof.Index))
codonEncodeByteSlice(w, v.Proof.LeafHash[:])
codonEncodeVarint(w, int64(len(v.Proof.Aunts)))
for _0:=0; _0<len(v.Proof.Aunts); _0++ {
codonEncodeByteSlice(w, v.Proof.Aunts[_0][:])
}
// end of v.Proof
} //End of EncodePart

func DecodePart(bz []byte) (Part, int, error) {
var err error
var length int
var v Part
var n int
var total int
v.Index = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Bytes, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proof.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proof.Index = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proof.LeafHash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proof.Aunts = make([][]byte, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of slice
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Proof.Aunts[_0], n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.Proof
return v, total, nil
} //End of DecodePart

func RandPart(r RandSrc) Part {
var length int
var v Part
v.Index = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Bytes = r.GetBytes(length)
v.Proof.Total = r.GetInt()
v.Proof.Index = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proof.LeafHash = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proof.Aunts = make([][]byte, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of slice
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Proof.Aunts[_0] = r.GetBytes(length)
}
// end of v.Proof
return v
} //End of RandPart

func DeepCopyPart(in Part) (out Part) {
var length int
out.Index = in.Index
length = len(in.Bytes)
out.Bytes = make([]uint8, length)
copy(out.Bytes[:], in.Bytes[:])
out.Proof.Total = in.Proof.Total
out.Proof.Index = in.Proof.Index
length = len(in.Proof.LeafHash)
out.Proof.LeafHash = make([]uint8, length)
copy(out.Proof.LeafHash[:], in.Proof.LeafHash[:])
length = len(in.Proof.Aunts)
out.Proof.Aunts = make([][]byte, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of slice
length = len(in.Proof.Aunts[_0])
out.Proof.Aunts[_0] = make([]uint8, length)
copy(out.Proof.Aunts[_0][:], in.Proof.Aunts[_0][:])
}
// end of .Proof
return
} //End of DeepCopyPart

// Non-Interface
func EncodeProtocol(w *[]byte, v Protocol) {
codonEncodeUvarint(w, uint64(v))
} //End of EncodeProtocol

func DecodeProtocol(bz []byte) (Protocol, int, error) {
var err error
var v Protocol
var n int
var total int
v = Protocol(codonDecodeUint64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeProtocol

func RandProtocol(r RandSrc) Protocol {
var v Protocol
v = Protocol(r.GetUint64())
return v
} //End of RandProtocol

func DeepCopyProtocol(in Protocol) (out Protocol) {
out = in
return
} //End of DeepCopyProtocol

// Non-Interface
func EncodeTx(w *[]byte, v Tx) {
codonEncodeByteSlice(w, v[:])
} //End of EncodeTx

func DecodeTx(bz []byte) (Tx, int, error) {
var err error
var length int
var v Tx
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeTx

func RandTx(r RandSrc) Tx {
var length int
var v Tx
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v = r.GetBytes(length)
return v
} //End of RandTx

func DeepCopyTx(in Tx) (out Tx) {
var length int
length = len(in)
out = make([]uint8, length)
copy(out[:], in[:])
return
} //End of DeepCopyTx

// Non-Interface
func EncodeCommitSig(w *[]byte, v CommitSig) {
codonEncodeUint8(w, uint8(v.Type))
codonEncodeVarint(w, int64(v.Height))
codonEncodeVarint(w, int64(v.Round))
codonEncodeByteSlice(w, v.BlockID.Hash[:])
codonEncodeVarint(w, int64(v.BlockID.PartsHeader.Total))
codonEncodeByteSlice(w, v.BlockID.PartsHeader.Hash[:])
// end of v.BlockID.PartsHeader
// end of v.BlockID
EncodeTime(w, v.Timestamp)
codonEncodeByteSlice(w, v.ValidatorAddress[:])
codonEncodeVarint(w, int64(v.ValidatorIndex))
codonEncodeByteSlice(w, v.Signature[:])
} //End of EncodeCommitSig

func DecodeCommitSig(bz []byte) (CommitSig, int, error) {
var err error
var length int
var v CommitSig
var n int
var total int
v.Type = SignedMsgType(codonDecodeUint8(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Height = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Round = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Total = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.BlockID.PartsHeader.Hash, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.BlockID.PartsHeader
// end of v.BlockID
v.Timestamp, n, err = DecodeTime(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ValidatorAddress, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ValidatorIndex = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Signature, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeCommitSig

func RandCommitSig(r RandSrc) CommitSig {
var length int
var v CommitSig
v.Type = SignedMsgType(r.GetUint8())
v.Height = r.GetInt64()
v.Round = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.Hash = r.GetBytes(length)
v.BlockID.PartsHeader.Total = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.BlockID.PartsHeader.Hash = r.GetBytes(length)
// end of v.BlockID.PartsHeader
// end of v.BlockID
v.Timestamp = RandTime(r)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.ValidatorAddress = r.GetBytes(length)
v.ValidatorIndex = r.GetInt()
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Signature = r.GetBytes(length)
return v
} //End of RandCommitSig

func DeepCopyCommitSig(in CommitSig) (out CommitSig) {
var length int
out.Type = in.Type
out.Height = in.Height
out.Round = in.Round
length = len(in.BlockID.Hash)
out.BlockID.Hash = make([]uint8, length)
copy(out.BlockID.Hash[:], in.BlockID.Hash[:])
out.BlockID.PartsHeader.Total = in.BlockID.PartsHeader.Total
length = len(in.BlockID.PartsHeader.Hash)
out.BlockID.PartsHeader.Hash = make([]uint8, length)
copy(out.BlockID.PartsHeader.Hash[:], in.BlockID.PartsHeader.Hash[:])
// end of .BlockID.PartsHeader
// end of .BlockID
out.Timestamp = DeepCopyTime(in.Timestamp)
length = len(in.ValidatorAddress)
out.ValidatorAddress = make([]uint8, length)
copy(out.ValidatorAddress[:], in.ValidatorAddress[:])
out.ValidatorIndex = in.ValidatorIndex
length = len(in.Signature)
out.Signature = make([]uint8, length)
copy(out.Signature[:], in.Signature[:])
return
} //End of DeepCopyCommitSig

// Non-Interface
func EncodeValidator(w *[]byte, v Validator) {
codonEncodeByteSlice(w, v.Address[:])
EncodePubKey(w, v.PubKey) // interface_encode
codonEncodeVarint(w, int64(v.VotingPower))
codonEncodeVarint(w, int64(v.ProposerPriority))
} //End of EncodeValidator

func DecodeValidator(bz []byte) (Validator, int, error) {
var err error
var length int
var v Validator
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Address, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.PubKey, n, err = DecodePubKey(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n // interface_decode
v.VotingPower = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.ProposerPriority = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeValidator

func RandValidator(r RandSrc) Validator {
var length int
var v Validator
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Address = r.GetBytes(length)
v.PubKey = RandPubKey(r) // interface_decode
v.VotingPower = r.GetInt64()
v.ProposerPriority = r.GetInt64()
return v
} //End of RandValidator

func DeepCopyValidator(in Validator) (out Validator) {
var length int
length = len(in.Address)
out.Address = make([]uint8, length)
copy(out.Address[:], in.Address[:])
out.PubKey = DeepCopyPubKey(in.PubKey)
out.VotingPower = in.VotingPower
out.ProposerPriority = in.ProposerPriority
return
} //End of DeepCopyValidator

// Non-Interface
func EncodeEvent(w *[]byte, v Event) {
codonEncodeString(w, v.Type)
codonEncodeVarint(w, int64(len(v.Attributes)))
for _0:=0; _0<len(v.Attributes); _0++ {
codonEncodeByteSlice(w, v.Attributes[_0].Key[:])
codonEncodeByteSlice(w, v.Attributes[_0].Value[:])
// end of v.Attributes[_0].XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.Attributes[_0].XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.Attributes[_0].XXX_sizecache))
// end of v.Attributes[_0]
}
// end of v.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.XXX_sizecache))
} //End of EncodeEvent

func DecodeEvent(bz []byte) (Event, int, error) {
var err error
var length int
var v Event
var n int
var total int
v.Type = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Attributes = make([]KVPair, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.Attributes[_0], n, err = DecodeKVPair(bz)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeEvent

func RandEvent(r RandSrc) Event {
var length int
var v Event
v.Type = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Attributes = make([]KVPair, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
v.Attributes[_0] = RandKVPair(r)
}
// end of v.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.XXX_unrecognized = r.GetBytes(length)
v.XXX_sizecache = r.GetInt32()
return v
} //End of RandEvent

func DeepCopyEvent(in Event) (out Event) {
var length int
out.Type = in.Type
length = len(in.Attributes)
out.Attributes = make([]KVPair, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of struct
out.Attributes[_0] = DeepCopyKVPair(in.Attributes[_0])
}
// end of .XXX_NoUnkeyedLiteral
length = len(in.XXX_unrecognized)
out.XXX_unrecognized = make([]uint8, length)
copy(out.XXX_unrecognized[:], in.XXX_unrecognized[:])
out.XXX_sizecache = in.XXX_sizecache
return
} //End of DeepCopyEvent

// Non-Interface
func EncodeValidatorUpdate(w *[]byte, v ValidatorUpdate) {
codonEncodeString(w, v.PubKey.Type)
codonEncodeByteSlice(w, v.PubKey.Data[:])
// end of v.PubKey.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.PubKey.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.PubKey.XXX_sizecache))
// end of v.PubKey
codonEncodeVarint(w, int64(v.Power))
// end of v.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.XXX_sizecache))
} //End of EncodeValidatorUpdate

func DecodeValidatorUpdate(bz []byte) (ValidatorUpdate, int, error) {
var err error
var length int
var v ValidatorUpdate
var n int
var total int
v.PubKey.Type = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.PubKey.Data, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.PubKey.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.PubKey.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.PubKey.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.PubKey
v.Power = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeValidatorUpdate

func RandValidatorUpdate(r RandSrc) ValidatorUpdate {
var length int
var v ValidatorUpdate
v.PubKey.Type = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.PubKey.Data = r.GetBytes(length)
// end of v.PubKey.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.PubKey.XXX_unrecognized = r.GetBytes(length)
v.PubKey.XXX_sizecache = r.GetInt32()
// end of v.PubKey
v.Power = r.GetInt64()
// end of v.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.XXX_unrecognized = r.GetBytes(length)
v.XXX_sizecache = r.GetInt32()
return v
} //End of RandValidatorUpdate

func DeepCopyValidatorUpdate(in ValidatorUpdate) (out ValidatorUpdate) {
var length int
out.PubKey.Type = in.PubKey.Type
length = len(in.PubKey.Data)
out.PubKey.Data = make([]uint8, length)
copy(out.PubKey.Data[:], in.PubKey.Data[:])
// end of .PubKey.XXX_NoUnkeyedLiteral
length = len(in.PubKey.XXX_unrecognized)
out.PubKey.XXX_unrecognized = make([]uint8, length)
copy(out.PubKey.XXX_unrecognized[:], in.PubKey.XXX_unrecognized[:])
out.PubKey.XXX_sizecache = in.PubKey.XXX_sizecache
// end of .PubKey
out.Power = in.Power
// end of .XXX_NoUnkeyedLiteral
length = len(in.XXX_unrecognized)
out.XXX_unrecognized = make([]uint8, length)
copy(out.XXX_unrecognized[:], in.XXX_unrecognized[:])
out.XXX_sizecache = in.XXX_sizecache
return
} //End of DeepCopyValidatorUpdate

// Non-Interface
func EncodeConsensusParams(w *[]byte, v ConsensusParams) {
codonEncodeVarint(w, int64(v.Block.MaxBytes))
codonEncodeVarint(w, int64(v.Block.MaxGas))
// end of v.Block.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.Block.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.Block.XXX_sizecache))
// end of v.Block
codonEncodeVarint(w, int64(v.Evidence.MaxAge))
// end of v.Evidence.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.Evidence.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.Evidence.XXX_sizecache))
// end of v.Evidence
codonEncodeVarint(w, int64(len(v.Validator.PubKeyTypes)))
for _0:=0; _0<len(v.Validator.PubKeyTypes); _0++ {
codonEncodeString(w, v.Validator.PubKeyTypes[_0])
}
// end of v.Validator.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.Validator.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.Validator.XXX_sizecache))
// end of v.Validator
// end of v.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.XXX_sizecache))
} //End of EncodeConsensusParams

func DecodeConsensusParams(bz []byte) (ConsensusParams, int, error) {
var err error
var length int
var v ConsensusParams
var n int
var total int
v.Block = &BlockParams{}
v.Block.MaxBytes = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Block.MaxGas = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Block.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Block.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Block.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Block
v.Evidence = &EvidenceParams{}
v.Evidence.MaxAge = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Evidence.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Evidence.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Evidence.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Evidence
v.Validator = &ValidatorParams{}
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Validator.PubKeyTypes = make([]string, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of string
v.Validator.PubKeyTypes[_0] = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.Validator.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Validator.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Validator.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.Validator
// end of v.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeConsensusParams

func RandConsensusParams(r RandSrc) ConsensusParams {
var length int
var v ConsensusParams
v.Block = &BlockParams{}
v.Block.MaxBytes = r.GetInt64()
v.Block.MaxGas = r.GetInt64()
// end of v.Block.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Block.XXX_unrecognized = r.GetBytes(length)
v.Block.XXX_sizecache = r.GetInt32()
// end of v.Block
v.Evidence = &EvidenceParams{}
v.Evidence.MaxAge = r.GetInt64()
// end of v.Evidence.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Evidence.XXX_unrecognized = r.GetBytes(length)
v.Evidence.XXX_sizecache = r.GetInt32()
// end of v.Evidence
v.Validator = &ValidatorParams{}
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Validator.PubKeyTypes = make([]string, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of string
v.Validator.PubKeyTypes[_0] = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
}
// end of v.Validator.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Validator.XXX_unrecognized = r.GetBytes(length)
v.Validator.XXX_sizecache = r.GetInt32()
// end of v.Validator
// end of v.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.XXX_unrecognized = r.GetBytes(length)
v.XXX_sizecache = r.GetInt32()
return v
} //End of RandConsensusParams

func DeepCopyConsensusParams(in ConsensusParams) (out ConsensusParams) {
var length int
out.Block = &BlockParams{}
out.Block.MaxBytes = in.Block.MaxBytes
out.Block.MaxGas = in.Block.MaxGas
// end of .Block.XXX_NoUnkeyedLiteral
length = len(in.Block.XXX_unrecognized)
out.Block.XXX_unrecognized = make([]uint8, length)
copy(out.Block.XXX_unrecognized[:], in.Block.XXX_unrecognized[:])
out.Block.XXX_sizecache = in.Block.XXX_sizecache
// end of .Block
out.Evidence = &EvidenceParams{}
out.Evidence.MaxAge = in.Evidence.MaxAge
// end of .Evidence.XXX_NoUnkeyedLiteral
length = len(in.Evidence.XXX_unrecognized)
out.Evidence.XXX_unrecognized = make([]uint8, length)
copy(out.Evidence.XXX_unrecognized[:], in.Evidence.XXX_unrecognized[:])
out.Evidence.XXX_sizecache = in.Evidence.XXX_sizecache
// end of .Evidence
out.Validator = &ValidatorParams{}
length = len(in.Validator.PubKeyTypes)
out.Validator.PubKeyTypes = make([]string, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of string
out.Validator.PubKeyTypes[_0] = in.Validator.PubKeyTypes[_0]
}
// end of .Validator.XXX_NoUnkeyedLiteral
length = len(in.Validator.XXX_unrecognized)
out.Validator.XXX_unrecognized = make([]uint8, length)
copy(out.Validator.XXX_unrecognized[:], in.Validator.XXX_unrecognized[:])
out.Validator.XXX_sizecache = in.Validator.XXX_sizecache
// end of .Validator
// end of .XXX_NoUnkeyedLiteral
length = len(in.XXX_unrecognized)
out.XXX_unrecognized = make([]uint8, length)
copy(out.XXX_unrecognized[:], in.XXX_unrecognized[:])
out.XXX_sizecache = in.XXX_sizecache
return
} //End of DeepCopyConsensusParams

// Non-Interface
func EncodeBlockParams(w *[]byte, v BlockParams) {
codonEncodeVarint(w, int64(v.MaxBytes))
codonEncodeVarint(w, int64(v.MaxGas))
// end of v.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.XXX_sizecache))
} //End of EncodeBlockParams

func DecodeBlockParams(bz []byte) (BlockParams, int, error) {
var err error
var length int
var v BlockParams
var n int
var total int
v.MaxBytes = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.MaxGas = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeBlockParams

func RandBlockParams(r RandSrc) BlockParams {
var length int
var v BlockParams
v.MaxBytes = r.GetInt64()
v.MaxGas = r.GetInt64()
// end of v.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.XXX_unrecognized = r.GetBytes(length)
v.XXX_sizecache = r.GetInt32()
return v
} //End of RandBlockParams

func DeepCopyBlockParams(in BlockParams) (out BlockParams) {
var length int
out.MaxBytes = in.MaxBytes
out.MaxGas = in.MaxGas
// end of .XXX_NoUnkeyedLiteral
length = len(in.XXX_unrecognized)
out.XXX_unrecognized = make([]uint8, length)
copy(out.XXX_unrecognized[:], in.XXX_unrecognized[:])
out.XXX_sizecache = in.XXX_sizecache
return
} //End of DeepCopyBlockParams

// Non-Interface
func EncodeEvidenceParams(w *[]byte, v EvidenceParams) {
codonEncodeVarint(w, int64(v.MaxAge))
// end of v.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.XXX_sizecache))
} //End of EncodeEvidenceParams

func DecodeEvidenceParams(bz []byte) (EvidenceParams, int, error) {
var err error
var length int
var v EvidenceParams
var n int
var total int
v.MaxAge = int64(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeEvidenceParams

func RandEvidenceParams(r RandSrc) EvidenceParams {
var length int
var v EvidenceParams
v.MaxAge = r.GetInt64()
// end of v.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.XXX_unrecognized = r.GetBytes(length)
v.XXX_sizecache = r.GetInt32()
return v
} //End of RandEvidenceParams

func DeepCopyEvidenceParams(in EvidenceParams) (out EvidenceParams) {
var length int
out.MaxAge = in.MaxAge
// end of .XXX_NoUnkeyedLiteral
length = len(in.XXX_unrecognized)
out.XXX_unrecognized = make([]uint8, length)
copy(out.XXX_unrecognized[:], in.XXX_unrecognized[:])
out.XXX_sizecache = in.XXX_sizecache
return
} //End of DeepCopyEvidenceParams

// Non-Interface
func EncodeValidatorParams(w *[]byte, v ValidatorParams) {
codonEncodeVarint(w, int64(len(v.PubKeyTypes)))
for _0:=0; _0<len(v.PubKeyTypes); _0++ {
codonEncodeString(w, v.PubKeyTypes[_0])
}
// end of v.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.XXX_sizecache))
} //End of EncodeValidatorParams

func DecodeValidatorParams(bz []byte) (ValidatorParams, int, error) {
var err error
var length int
var v ValidatorParams
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.PubKeyTypes = make([]string, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of string
v.PubKeyTypes[_0] = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
}
// end of v.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeValidatorParams

func RandValidatorParams(r RandSrc) ValidatorParams {
var length int
var v ValidatorParams
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.PubKeyTypes = make([]string, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of string
v.PubKeyTypes[_0] = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
}
// end of v.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.XXX_unrecognized = r.GetBytes(length)
v.XXX_sizecache = r.GetInt32()
return v
} //End of RandValidatorParams

func DeepCopyValidatorParams(in ValidatorParams) (out ValidatorParams) {
var length int
length = len(in.PubKeyTypes)
out.PubKeyTypes = make([]string, length)
for _0, length_0 := 0, length; _0<length_0; _0++ { //slice of string
out.PubKeyTypes[_0] = in.PubKeyTypes[_0]
}
// end of .XXX_NoUnkeyedLiteral
length = len(in.XXX_unrecognized)
out.XXX_unrecognized = make([]uint8, length)
copy(out.XXX_unrecognized[:], in.XXX_unrecognized[:])
out.XXX_sizecache = in.XXX_sizecache
return
} //End of DeepCopyValidatorParams

// Non-Interface
func EncodeRemoteSignerError(w *[]byte, v RemoteSignerError) {
codonEncodeVarint(w, int64(v.Code))
codonEncodeString(w, v.Description)
} //End of EncodeRemoteSignerError

func DecodeRemoteSignerError(bz []byte) (RemoteSignerError, int, error) {
var err error
var v RemoteSignerError
var n int
var total int
v.Code = int(codonDecodeInt(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Description = string(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeRemoteSignerError

func RandRemoteSignerError(r RandSrc) RemoteSignerError {
var v RemoteSignerError
v.Code = r.GetInt()
v.Description = r.GetString(1+int(r.GetUint()%(MaxStringLength-1)))
return v
} //End of RandRemoteSignerError

func DeepCopyRemoteSignerError(in RemoteSignerError) (out RemoteSignerError) {
out.Code = in.Code
out.Description = in.Description
return
} //End of DeepCopyRemoteSignerError

// Non-Interface
func EncodeP2PID(w *[]byte, v P2PID) {
codonEncodeString(w, string(v))
} //End of EncodeP2PID

func DecodeP2PID(bz []byte) (P2PID, int, error) {
var err error
var v P2PID
var n int
var total int
v = P2PID(codonDecodeString(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeP2PID

func RandP2PID(r RandSrc) P2PID {
var v P2PID
v = P2PID(r.GetString(1+int(r.GetUint()%(MaxStringLength-1))))
return v
} //End of RandP2PID

func DeepCopyP2PID(in P2PID) (out P2PID) {
out = in
return
} //End of DeepCopyP2PID

// Non-Interface
func EncodeDuration(w *[]byte, v Duration) {
codonEncodeVarint(w, int64(v))
} //End of EncodeDuration

func DecodeDuration(bz []byte) (Duration, int, error) {
var err error
var v Duration
var n int
var total int
v = Duration(codonDecodeInt64(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeDuration

func RandDuration(r RandSrc) Duration {
var v Duration
v = Duration(r.GetInt64())
return v
} //End of RandDuration

func DeepCopyDuration(in Duration) (out Duration) {
out = in
return
} //End of DeepCopyDuration

// Non-Interface
func EncodeKVPair(w *[]byte, v KVPair) {
codonEncodeByteSlice(w, v.Key[:])
codonEncodeByteSlice(w, v.Value[:])
// end of v.XXX_NoUnkeyedLiteral
codonEncodeByteSlice(w, v.XXX_unrecognized[:])
codonEncodeVarint(w, int64(v.XXX_sizecache))
} //End of EncodeKVPair

func DecodeKVPair(bz []byte) (KVPair, int, error) {
var err error
var length int
var v KVPair
var n int
var total int
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Key, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.Value, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
// end of v.XXX_NoUnkeyedLiteral
length = codonDecodeInt(bz, &n, &err)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_unrecognized, n, err = codonGetByteSlice(bz, length)
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
v.XXX_sizecache = int32(codonDecodeInt32(bz, &n, &err))
if err != nil {return v, total, err}
bz = bz[n:]
total+=n
return v, total, nil
} //End of DecodeKVPair

func RandKVPair(r RandSrc) KVPair {
var length int
var v KVPair
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Key = r.GetBytes(length)
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.Value = r.GetBytes(length)
// end of v.XXX_NoUnkeyedLiteral
length = 1+int(r.GetUint()%(MaxSliceLength-1))
v.XXX_unrecognized = r.GetBytes(length)
v.XXX_sizecache = r.GetInt32()
return v
} //End of RandKVPair

func DeepCopyKVPair(in KVPair) (out KVPair) {
var length int
length = len(in.Key)
out.Key = make([]uint8, length)
copy(out.Key[:], in.Key[:])
length = len(in.Value)
out.Value = make([]uint8, length)
copy(out.Value[:], in.Value[:])
// end of .XXX_NoUnkeyedLiteral
length = len(in.XXX_unrecognized)
out.XXX_unrecognized = make([]uint8, length)
copy(out.XXX_unrecognized[:], in.XXX_unrecognized[:])
out.XXX_sizecache = in.XXX_sizecache
return
} //End of DeepCopyKVPair

// Interface
func DecodePubKey(bz []byte) (PubKey, int, error) {
var v PubKey
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{114,76,37,23}:
v, n, err := DecodePubKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{14,33,23,141}:
v, n, err := DecodePubKeyMultisigThreshold(bz[4:])
return v, n+4, err
case [4]byte{51,161,20,197}:
v, n, err := DecodePubKeySecp256k1(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodePubKey
func EncodePubKey(w *[]byte, x interface{}) {
switch v := x.(type) {
case PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, v)
case *PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, *v)
case PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, v)
case *PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, *v)
case PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, v)
case *PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandPubKey(r RandSrc) PubKey {
switch r.GetUint() % 2 {
case 0:
return RandPubKeyEd25519(r)
case 1:
return RandPubKeySecp256k1(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyPubKey(x PubKey) PubKey {
switch v := x.(type) {
case *PubKeyEd25519:
res := DeepCopyPubKeyEd25519(*v)
return &res
case PubKeyEd25519:
res := DeepCopyPubKeyEd25519(v)
return res
case *PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(*v)
return &res
case PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodePrivKey(bz []byte) (PrivKey, int, error) {
var v PrivKey
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{158,94,112,161}:
v, n, err := DecodePrivKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{83,16,177,42}:
v, n, err := DecodePrivKeySecp256k1(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodePrivKey
func EncodePrivKey(w *[]byte, x interface{}) {
switch v := x.(type) {
case PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, v)
case *PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, *v)
case PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, v)
case *PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandPrivKey(r RandSrc) PrivKey {
switch r.GetUint() % 2 {
case 0:
return RandPrivKeyEd25519(r)
case 1:
return RandPrivKeySecp256k1(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyPrivKey(x PrivKey) PrivKey {
switch v := x.(type) {
case *PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(*v)
return &res
case PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(v)
return res
case *PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(*v)
return &res
case PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodeConsensusMessage(bz []byte) (ConsensusMessage, int, error) {
var v ConsensusMessage
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeConsensusMessage
func EncodeConsensusMessage(w *[]byte, x interface{}) {
switch v := x.(type) {
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandConsensusMessage(r RandSrc) ConsensusMessage {
switch r.GetUint() % 2 {
case 0:
return RandMockBadEvidence(r)
case 1:
return RandMockGoodEvidence(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyConsensusMessage(x ConsensusMessage) ConsensusMessage {
switch v := x.(type) {
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodeWALMessage(bz []byte) (WALMessage, int, error) {
var v WALMessage
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{252,242,102,165}:
v, n, err := DecodeBc0BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{226,55,89,212}:
v, n, err := DecodeBc0NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{89,2,129,24}:
v, n, err := DecodeBc0StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{181,21,96,69}:
v, n, err := DecodeBc0StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{235,69,170,199}:
v, n, err := DecodeBc1BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{134,105,217,194}:
v, n, err := DecodeBc1NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{45,250,153,33}:
v, n, err := DecodeBc1StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{81,214,62,85}:
v, n, err := DecodeBc1StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{206,243,248,173}:
v, n, err := DecodeBitArray(bz[4:])
return v, n+4, err
case [4]byte{119,191,108,56}:
v, n, err := DecodeBlockParams(bz[4:])
return v, n+4, err
case [4]byte{79,243,233,203}:
v, n, err := DecodeBlockPartMessage(bz[4:])
return v, n+4, err
case [4]byte{163,200,208,51}:
v, n, err := DecodeCommitSig(bz[4:])
return v, n+4, err
case [4]byte{204,104,130,6}:
v, n, err := DecodeConsensusParams(bz[4:])
return v, n+4, err
case [4]byte{89,252,98,178}:
v, n, err := DecodeDuplicateVoteEvidence(bz[4:])
return v, n+4, err
case [4]byte{79,197,42,60}:
v, n, err := DecodeDuration(bz[4:])
return v, n+4, err
case [4]byte{143,131,137,178}:
v, n, err := DecodeEndHeightMessage(bz[4:])
return v, n+4, err
case [4]byte{78,31,73,169}:
v, n, err := DecodeEvent(bz[4:])
return v, n+4, err
case [4]byte{63,236,183,244}:
v, n, err := DecodeEventDataCompleteProposal(bz[4:])
return v, n+4, err
case [4]byte{10,8,210,215}:
v, n, err := DecodeEventDataNewBlockHeader(bz[4:])
return v, n+4, err
case [4]byte{220,234,247,95}:
v, n, err := DecodeEventDataNewRound(bz[4:])
return v, n+4, err
case [4]byte{224,234,115,110}:
v, n, err := DecodeEventDataRoundState(bz[4:])
return v, n+4, err
case [4]byte{197,0,12,147}:
v, n, err := DecodeEventDataString(bz[4:])
return v, n+4, err
case [4]byte{64,227,166,189}:
v, n, err := DecodeEventDataTx(bz[4:])
return v, n+4, err
case [4]byte{79,214,242,86}:
v, n, err := DecodeEventDataValidatorSetUpdates(bz[4:])
return v, n+4, err
case [4]byte{82,175,21,121}:
v, n, err := DecodeEventDataVote(bz[4:])
return v, n+4, err
case [4]byte{51,243,128,228}:
v, n, err := DecodeEvidenceListMessage(bz[4:])
return v, n+4, err
case [4]byte{169,170,199,49}:
v, n, err := DecodeEvidenceParams(bz[4:])
return v, n+4, err
case [4]byte{30,19,14,216}:
v, n, err := DecodeHasVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{63,28,146,21}:
v, n, err := DecodeKVPair(bz[4:])
return v, n+4, err
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
case [4]byte{72,21,63,101}:
v, n, err := DecodeMsgInfo(bz[4:])
return v, n+4, err
case [4]byte{193,38,23,92}:
v, n, err := DecodeNewRoundStepMessage(bz[4:])
return v, n+4, err
case [4]byte{132,190,161,107}:
v, n, err := DecodeNewValidBlockMessage(bz[4:])
return v, n+4, err
case [4]byte{162,154,17,21}:
v, n, err := DecodeP2PID(bz[4:])
return v, n+4, err
case [4]byte{28,180,161,19}:
v, n, err := DecodePacketMsg(bz[4:])
return v, n+4, err
case [4]byte{133,112,162,225}:
v, n, err := DecodePart(bz[4:])
return v, n+4, err
case [4]byte{158,94,112,161}:
v, n, err := DecodePrivKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{83,16,177,42}:
v, n, err := DecodePrivKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{93,66,118,108}:
v, n, err := DecodeProposal(bz[4:])
return v, n+4, err
case [4]byte{162,17,172,204}:
v, n, err := DecodeProposalMessage(bz[4:])
return v, n+4, err
case [4]byte{187,39,126,253}:
v, n, err := DecodeProposalPOLMessage(bz[4:])
return v, n+4, err
case [4]byte{207,8,131,52}:
v, n, err := DecodeProtocol(bz[4:])
return v, n+4, err
case [4]byte{114,76,37,23}:
v, n, err := DecodePubKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{14,33,23,141}:
v, n, err := DecodePubKeyMultisigThreshold(bz[4:])
return v, n+4, err
case [4]byte{1,41,197,74}:
v, n, err := DecodePubKeyResponse(bz[4:])
return v, n+4, err
case [4]byte{51,161,20,197}:
v, n, err := DecodePubKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{198,215,201,27}:
v, n, err := DecodeRemoteSignerError(bz[4:])
return v, n+4, err
case [4]byte{235,161,234,245}:
v, n, err := DecodeRoundStepType(bz[4:])
return v, n+4, err
case [4]byte{9,94,219,98}:
v, n, err := DecodeSignProposalRequest(bz[4:])
return v, n+4, err
case [4]byte{17,244,69,120}:
v, n, err := DecodeSignVoteRequest(bz[4:])
return v, n+4, err
case [4]byte{67,52,162,78}:
v, n, err := DecodeSignedMsgType(bz[4:])
return v, n+4, err
case [4]byte{173,178,196,98}:
v, n, err := DecodeSignedProposalResponse(bz[4:])
return v, n+4, err
case [4]byte{144,17,7,46}:
v, n, err := DecodeSignedVoteResponse(bz[4:])
return v, n+4, err
case [4]byte{47,153,154,201}:
v, n, err := DecodeTimeoutInfo(bz[4:])
return v, n+4, err
case [4]byte{91,158,141,124}:
v, n, err := DecodeTx(bz[4:])
return v, n+4, err
case [4]byte{138,192,166,6}:
v, n, err := DecodeTxMessage(bz[4:])
return v, n+4, err
case [4]byte{11,10,9,103}:
v, n, err := DecodeValidator(bz[4:])
return v, n+4, err
case [4]byte{228,6,69,233}:
v, n, err := DecodeValidatorParams(bz[4:])
return v, n+4, err
case [4]byte{60,45,116,56}:
v, n, err := DecodeValidatorUpdate(bz[4:])
return v, n+4, err
case [4]byte{205,85,136,219}:
v, n, err := DecodeVote(bz[4:])
return v, n+4, err
case [4]byte{239,126,120,100}:
v, n, err := DecodeVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{135,195,3,234}:
v, n, err := DecodeVoteSetBitsMessage(bz[4:])
return v, n+4, err
case [4]byte{232,73,18,176}:
v, n, err := DecodeVoteSetMaj23Message(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeWALMessage
func EncodeWALMessage(w *[]byte, x interface{}) {
switch v := x.(type) {
case Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, v)
case *Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, *v)
case Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, v)
case *Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, *v)
case Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, v)
case *Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, *v)
case Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, v)
case *Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, *v)
case Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, v)
case *Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, *v)
case Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, v)
case *Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, *v)
case Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, v)
case *Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, *v)
case Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, v)
case *Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, *v)
case BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, v)
case *BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, *v)
case BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, v)
case *BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, *v)
case BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, v)
case *BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, *v)
case CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, v)
case *CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, *v)
case ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, v)
case *ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, *v)
case DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, v)
case *DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, *v)
case Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, v)
case *Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, *v)
case EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, v)
case *EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, *v)
case Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, v)
case *Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, *v)
case EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, v)
case *EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, *v)
case EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, v)
case *EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, *v)
case EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, v)
case *EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, *v)
case EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, v)
case *EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, *v)
case EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, v)
case *EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, *v)
case EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, v)
case *EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, *v)
case EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, v)
case *EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, *v)
case EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, v)
case *EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, *v)
case EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, v)
case *EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, *v)
case EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, v)
case *EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, *v)
case HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, v)
case *HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, *v)
case KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, v)
case *KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, *v)
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
case MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, v)
case *MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, *v)
case NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, v)
case *NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, *v)
case NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, v)
case *NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, *v)
case P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, v)
case *P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, *v)
case PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, v)
case *PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, *v)
case Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, v)
case *Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, *v)
case PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, v)
case *PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, *v)
case PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, v)
case *PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, *v)
case Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, v)
case *Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, *v)
case ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, v)
case *ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, *v)
case ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, v)
case *ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, *v)
case Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, v)
case *Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, *v)
case PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, v)
case *PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, *v)
case PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, v)
case *PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, *v)
case PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, v)
case *PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, *v)
case PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, v)
case *PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, *v)
case RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, v)
case *RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, *v)
case RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, v)
case *RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, *v)
case SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, v)
case *SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, *v)
case SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, v)
case *SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, *v)
case SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, v)
case *SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, *v)
case SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, v)
case *SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, *v)
case SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, v)
case *SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, *v)
case TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, v)
case *TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, *v)
case Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, v)
case *Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, *v)
case TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, v)
case *TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, *v)
case Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, v)
case *Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, *v)
case ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, v)
case *ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, *v)
case ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, v)
case *ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, *v)
case Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, v)
case *Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, *v)
case VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, v)
case *VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, *v)
case VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, v)
case *VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, *v)
case VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, v)
case *VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandWALMessage(r RandSrc) WALMessage {
switch r.GetUint() % 64 {
case 0:
return RandBc0BlockRequestMessage(r)
case 1:
return RandBc0NoBlockResponseMessage(r)
case 2:
return RandBc0StatusRequestMessage(r)
case 3:
return RandBc0StatusResponseMessage(r)
case 4:
return RandBc1BlockRequestMessage(r)
case 5:
return RandBc1NoBlockResponseMessage(r)
case 6:
return RandBc1StatusRequestMessage(r)
case 7:
return RandBc1StatusResponseMessage(r)
case 8:
return RandBitArray(r)
case 9:
return RandBlockParams(r)
case 10:
return RandBlockPartMessage(r)
case 11:
return RandCommitSig(r)
case 12:
return RandConsensusParams(r)
case 13:
return RandDuplicateVoteEvidence(r)
case 14:
return RandDuration(r)
case 15:
return RandEndHeightMessage(r)
case 16:
return RandEvent(r)
case 17:
return RandEventDataCompleteProposal(r)
case 18:
return RandEventDataNewBlockHeader(r)
case 19:
return RandEventDataNewRound(r)
case 20:
return RandEventDataRoundState(r)
case 21:
return RandEventDataString(r)
case 22:
return RandEventDataTx(r)
case 23:
return RandEventDataValidatorSetUpdates(r)
case 24:
return RandEventDataVote(r)
case 25:
return RandEvidenceListMessage(r)
case 26:
return RandEvidenceParams(r)
case 27:
return RandHasVoteMessage(r)
case 28:
return RandKVPair(r)
case 29:
return RandMockBadEvidence(r)
case 30:
return RandMockGoodEvidence(r)
case 31:
return RandMsgInfo(r)
case 32:
return RandNewRoundStepMessage(r)
case 33:
return RandNewValidBlockMessage(r)
case 34:
return RandP2PID(r)
case 35:
return RandPacketMsg(r)
case 36:
return RandPart(r)
case 37:
return RandPrivKeyEd25519(r)
case 38:
return RandPrivKeySecp256k1(r)
case 39:
return RandProposal(r)
case 40:
return RandProposalMessage(r)
case 41:
return RandProposalPOLMessage(r)
case 42:
return RandProtocol(r)
case 43:
return RandPubKeyEd25519(r)
case 44:
return RandPubKeyMultisigThreshold(r)
case 45:
return RandPubKeyResponse(r)
case 46:
return RandPubKeySecp256k1(r)
case 47:
return RandRemoteSignerError(r)
case 48:
return RandRoundStepType(r)
case 49:
return RandSignProposalRequest(r)
case 50:
return RandSignVoteRequest(r)
case 51:
return RandSignedMsgType(r)
case 52:
return RandSignedProposalResponse(r)
case 53:
return RandSignedVoteResponse(r)
case 54:
return RandTimeoutInfo(r)
case 55:
return RandTx(r)
case 56:
return RandTxMessage(r)
case 57:
return RandValidator(r)
case 58:
return RandValidatorParams(r)
case 59:
return RandValidatorUpdate(r)
case 60:
return RandVote(r)
case 61:
return RandVoteMessage(r)
case 62:
return RandVoteSetBitsMessage(r)
case 63:
return RandVoteSetMaj23Message(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyWALMessage(x WALMessage) WALMessage {
switch v := x.(type) {
case *Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(*v)
return &res
case Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(v)
return res
case *Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(*v)
return &res
case Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(v)
return res
case *Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(*v)
return &res
case Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(v)
return res
case *Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(*v)
return &res
case Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(v)
return res
case *Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(*v)
return &res
case Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(v)
return res
case *Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(*v)
return &res
case Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(v)
return res
case *Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(*v)
return &res
case Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(v)
return res
case *Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(*v)
return &res
case Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(v)
return res
case *BitArray:
res := DeepCopyBitArray(*v)
return &res
case BitArray:
res := DeepCopyBitArray(v)
return res
case *BlockParams:
res := DeepCopyBlockParams(*v)
return &res
case BlockParams:
res := DeepCopyBlockParams(v)
return res
case *BlockPartMessage:
res := DeepCopyBlockPartMessage(*v)
return &res
case BlockPartMessage:
res := DeepCopyBlockPartMessage(v)
return res
case *CommitSig:
res := DeepCopyCommitSig(*v)
return &res
case CommitSig:
res := DeepCopyCommitSig(v)
return res
case *ConsensusParams:
res := DeepCopyConsensusParams(*v)
return &res
case ConsensusParams:
res := DeepCopyConsensusParams(v)
return res
case *DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(*v)
return &res
case DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(v)
return res
case *Duration:
res := DeepCopyDuration(*v)
return &res
case Duration:
res := DeepCopyDuration(v)
return res
case *EndHeightMessage:
res := DeepCopyEndHeightMessage(*v)
return &res
case EndHeightMessage:
res := DeepCopyEndHeightMessage(v)
return res
case *Event:
res := DeepCopyEvent(*v)
return &res
case Event:
res := DeepCopyEvent(v)
return res
case *EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(*v)
return &res
case EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(v)
return res
case *EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(*v)
return &res
case EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(v)
return res
case *EventDataNewRound:
res := DeepCopyEventDataNewRound(*v)
return &res
case EventDataNewRound:
res := DeepCopyEventDataNewRound(v)
return res
case *EventDataRoundState:
res := DeepCopyEventDataRoundState(*v)
return &res
case EventDataRoundState:
res := DeepCopyEventDataRoundState(v)
return res
case *EventDataString:
res := DeepCopyEventDataString(*v)
return &res
case EventDataString:
res := DeepCopyEventDataString(v)
return res
case *EventDataTx:
res := DeepCopyEventDataTx(*v)
return &res
case EventDataTx:
res := DeepCopyEventDataTx(v)
return res
case *EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(*v)
return &res
case EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(v)
return res
case *EventDataVote:
res := DeepCopyEventDataVote(*v)
return &res
case EventDataVote:
res := DeepCopyEventDataVote(v)
return res
case *EvidenceListMessage:
res := DeepCopyEvidenceListMessage(*v)
return &res
case EvidenceListMessage:
res := DeepCopyEvidenceListMessage(v)
return res
case *EvidenceParams:
res := DeepCopyEvidenceParams(*v)
return &res
case EvidenceParams:
res := DeepCopyEvidenceParams(v)
return res
case *HasVoteMessage:
res := DeepCopyHasVoteMessage(*v)
return &res
case HasVoteMessage:
res := DeepCopyHasVoteMessage(v)
return res
case *KVPair:
res := DeepCopyKVPair(*v)
return &res
case KVPair:
res := DeepCopyKVPair(v)
return res
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
case *MsgInfo:
res := DeepCopyMsgInfo(*v)
return &res
case MsgInfo:
res := DeepCopyMsgInfo(v)
return res
case *NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(*v)
return &res
case NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(v)
return res
case *NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(*v)
return &res
case NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(v)
return res
case *P2PID:
res := DeepCopyP2PID(*v)
return &res
case P2PID:
res := DeepCopyP2PID(v)
return res
case *PacketMsg:
res := DeepCopyPacketMsg(*v)
return &res
case PacketMsg:
res := DeepCopyPacketMsg(v)
return res
case *Part:
res := DeepCopyPart(*v)
return &res
case Part:
res := DeepCopyPart(v)
return res
case *PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(*v)
return &res
case PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(v)
return res
case *PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(*v)
return &res
case PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(v)
return res
case *Proposal:
res := DeepCopyProposal(*v)
return &res
case Proposal:
res := DeepCopyProposal(v)
return res
case *ProposalMessage:
res := DeepCopyProposalMessage(*v)
return &res
case ProposalMessage:
res := DeepCopyProposalMessage(v)
return res
case *ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(*v)
return &res
case ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(v)
return res
case *Protocol:
res := DeepCopyProtocol(*v)
return &res
case Protocol:
res := DeepCopyProtocol(v)
return res
case *PubKeyEd25519:
res := DeepCopyPubKeyEd25519(*v)
return &res
case PubKeyEd25519:
res := DeepCopyPubKeyEd25519(v)
return res
case *PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(*v)
return &res
case PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(v)
return res
case *PubKeyResponse:
res := DeepCopyPubKeyResponse(*v)
return &res
case PubKeyResponse:
res := DeepCopyPubKeyResponse(v)
return res
case *PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(*v)
return &res
case PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(v)
return res
case *RemoteSignerError:
res := DeepCopyRemoteSignerError(*v)
return &res
case RemoteSignerError:
res := DeepCopyRemoteSignerError(v)
return res
case *RoundStepType:
res := DeepCopyRoundStepType(*v)
return &res
case RoundStepType:
res := DeepCopyRoundStepType(v)
return res
case *SignProposalRequest:
res := DeepCopySignProposalRequest(*v)
return &res
case SignProposalRequest:
res := DeepCopySignProposalRequest(v)
return res
case *SignVoteRequest:
res := DeepCopySignVoteRequest(*v)
return &res
case SignVoteRequest:
res := DeepCopySignVoteRequest(v)
return res
case *SignedMsgType:
res := DeepCopySignedMsgType(*v)
return &res
case SignedMsgType:
res := DeepCopySignedMsgType(v)
return res
case *SignedProposalResponse:
res := DeepCopySignedProposalResponse(*v)
return &res
case SignedProposalResponse:
res := DeepCopySignedProposalResponse(v)
return res
case *SignedVoteResponse:
res := DeepCopySignedVoteResponse(*v)
return &res
case SignedVoteResponse:
res := DeepCopySignedVoteResponse(v)
return res
case *TimeoutInfo:
res := DeepCopyTimeoutInfo(*v)
return &res
case TimeoutInfo:
res := DeepCopyTimeoutInfo(v)
return res
case *Tx:
res := DeepCopyTx(*v)
return &res
case Tx:
res := DeepCopyTx(v)
return res
case *TxMessage:
res := DeepCopyTxMessage(*v)
return &res
case TxMessage:
res := DeepCopyTxMessage(v)
return res
case *Validator:
res := DeepCopyValidator(*v)
return &res
case Validator:
res := DeepCopyValidator(v)
return res
case *ValidatorParams:
res := DeepCopyValidatorParams(*v)
return &res
case ValidatorParams:
res := DeepCopyValidatorParams(v)
return res
case *ValidatorUpdate:
res := DeepCopyValidatorUpdate(*v)
return &res
case ValidatorUpdate:
res := DeepCopyValidatorUpdate(v)
return res
case *Vote:
res := DeepCopyVote(*v)
return &res
case Vote:
res := DeepCopyVote(v)
return res
case *VoteMessage:
res := DeepCopyVoteMessage(*v)
return &res
case VoteMessage:
res := DeepCopyVoteMessage(v)
return res
case *VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(*v)
return &res
case VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(v)
return res
case *VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(*v)
return &res
case VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodeTMEventData(bz []byte) (TMEventData, int, error) {
var v TMEventData
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{252,242,102,165}:
v, n, err := DecodeBc0BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{226,55,89,212}:
v, n, err := DecodeBc0NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{89,2,129,24}:
v, n, err := DecodeBc0StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{181,21,96,69}:
v, n, err := DecodeBc0StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{235,69,170,199}:
v, n, err := DecodeBc1BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{134,105,217,194}:
v, n, err := DecodeBc1NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{45,250,153,33}:
v, n, err := DecodeBc1StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{81,214,62,85}:
v, n, err := DecodeBc1StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{206,243,248,173}:
v, n, err := DecodeBitArray(bz[4:])
return v, n+4, err
case [4]byte{119,191,108,56}:
v, n, err := DecodeBlockParams(bz[4:])
return v, n+4, err
case [4]byte{79,243,233,203}:
v, n, err := DecodeBlockPartMessage(bz[4:])
return v, n+4, err
case [4]byte{163,200,208,51}:
v, n, err := DecodeCommitSig(bz[4:])
return v, n+4, err
case [4]byte{204,104,130,6}:
v, n, err := DecodeConsensusParams(bz[4:])
return v, n+4, err
case [4]byte{89,252,98,178}:
v, n, err := DecodeDuplicateVoteEvidence(bz[4:])
return v, n+4, err
case [4]byte{79,197,42,60}:
v, n, err := DecodeDuration(bz[4:])
return v, n+4, err
case [4]byte{143,131,137,178}:
v, n, err := DecodeEndHeightMessage(bz[4:])
return v, n+4, err
case [4]byte{78,31,73,169}:
v, n, err := DecodeEvent(bz[4:])
return v, n+4, err
case [4]byte{63,236,183,244}:
v, n, err := DecodeEventDataCompleteProposal(bz[4:])
return v, n+4, err
case [4]byte{10,8,210,215}:
v, n, err := DecodeEventDataNewBlockHeader(bz[4:])
return v, n+4, err
case [4]byte{220,234,247,95}:
v, n, err := DecodeEventDataNewRound(bz[4:])
return v, n+4, err
case [4]byte{224,234,115,110}:
v, n, err := DecodeEventDataRoundState(bz[4:])
return v, n+4, err
case [4]byte{197,0,12,147}:
v, n, err := DecodeEventDataString(bz[4:])
return v, n+4, err
case [4]byte{64,227,166,189}:
v, n, err := DecodeEventDataTx(bz[4:])
return v, n+4, err
case [4]byte{79,214,242,86}:
v, n, err := DecodeEventDataValidatorSetUpdates(bz[4:])
return v, n+4, err
case [4]byte{82,175,21,121}:
v, n, err := DecodeEventDataVote(bz[4:])
return v, n+4, err
case [4]byte{51,243,128,228}:
v, n, err := DecodeEvidenceListMessage(bz[4:])
return v, n+4, err
case [4]byte{169,170,199,49}:
v, n, err := DecodeEvidenceParams(bz[4:])
return v, n+4, err
case [4]byte{30,19,14,216}:
v, n, err := DecodeHasVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{63,28,146,21}:
v, n, err := DecodeKVPair(bz[4:])
return v, n+4, err
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
case [4]byte{72,21,63,101}:
v, n, err := DecodeMsgInfo(bz[4:])
return v, n+4, err
case [4]byte{193,38,23,92}:
v, n, err := DecodeNewRoundStepMessage(bz[4:])
return v, n+4, err
case [4]byte{132,190,161,107}:
v, n, err := DecodeNewValidBlockMessage(bz[4:])
return v, n+4, err
case [4]byte{162,154,17,21}:
v, n, err := DecodeP2PID(bz[4:])
return v, n+4, err
case [4]byte{28,180,161,19}:
v, n, err := DecodePacketMsg(bz[4:])
return v, n+4, err
case [4]byte{133,112,162,225}:
v, n, err := DecodePart(bz[4:])
return v, n+4, err
case [4]byte{158,94,112,161}:
v, n, err := DecodePrivKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{83,16,177,42}:
v, n, err := DecodePrivKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{93,66,118,108}:
v, n, err := DecodeProposal(bz[4:])
return v, n+4, err
case [4]byte{162,17,172,204}:
v, n, err := DecodeProposalMessage(bz[4:])
return v, n+4, err
case [4]byte{187,39,126,253}:
v, n, err := DecodeProposalPOLMessage(bz[4:])
return v, n+4, err
case [4]byte{207,8,131,52}:
v, n, err := DecodeProtocol(bz[4:])
return v, n+4, err
case [4]byte{114,76,37,23}:
v, n, err := DecodePubKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{14,33,23,141}:
v, n, err := DecodePubKeyMultisigThreshold(bz[4:])
return v, n+4, err
case [4]byte{1,41,197,74}:
v, n, err := DecodePubKeyResponse(bz[4:])
return v, n+4, err
case [4]byte{51,161,20,197}:
v, n, err := DecodePubKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{198,215,201,27}:
v, n, err := DecodeRemoteSignerError(bz[4:])
return v, n+4, err
case [4]byte{235,161,234,245}:
v, n, err := DecodeRoundStepType(bz[4:])
return v, n+4, err
case [4]byte{9,94,219,98}:
v, n, err := DecodeSignProposalRequest(bz[4:])
return v, n+4, err
case [4]byte{17,244,69,120}:
v, n, err := DecodeSignVoteRequest(bz[4:])
return v, n+4, err
case [4]byte{67,52,162,78}:
v, n, err := DecodeSignedMsgType(bz[4:])
return v, n+4, err
case [4]byte{173,178,196,98}:
v, n, err := DecodeSignedProposalResponse(bz[4:])
return v, n+4, err
case [4]byte{144,17,7,46}:
v, n, err := DecodeSignedVoteResponse(bz[4:])
return v, n+4, err
case [4]byte{47,153,154,201}:
v, n, err := DecodeTimeoutInfo(bz[4:])
return v, n+4, err
case [4]byte{91,158,141,124}:
v, n, err := DecodeTx(bz[4:])
return v, n+4, err
case [4]byte{138,192,166,6}:
v, n, err := DecodeTxMessage(bz[4:])
return v, n+4, err
case [4]byte{11,10,9,103}:
v, n, err := DecodeValidator(bz[4:])
return v, n+4, err
case [4]byte{228,6,69,233}:
v, n, err := DecodeValidatorParams(bz[4:])
return v, n+4, err
case [4]byte{60,45,116,56}:
v, n, err := DecodeValidatorUpdate(bz[4:])
return v, n+4, err
case [4]byte{205,85,136,219}:
v, n, err := DecodeVote(bz[4:])
return v, n+4, err
case [4]byte{239,126,120,100}:
v, n, err := DecodeVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{135,195,3,234}:
v, n, err := DecodeVoteSetBitsMessage(bz[4:])
return v, n+4, err
case [4]byte{232,73,18,176}:
v, n, err := DecodeVoteSetMaj23Message(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeTMEventData
func EncodeTMEventData(w *[]byte, x interface{}) {
switch v := x.(type) {
case Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, v)
case *Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, *v)
case Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, v)
case *Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, *v)
case Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, v)
case *Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, *v)
case Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, v)
case *Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, *v)
case Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, v)
case *Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, *v)
case Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, v)
case *Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, *v)
case Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, v)
case *Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, *v)
case Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, v)
case *Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, *v)
case BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, v)
case *BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, *v)
case BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, v)
case *BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, *v)
case BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, v)
case *BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, *v)
case CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, v)
case *CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, *v)
case ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, v)
case *ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, *v)
case DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, v)
case *DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, *v)
case Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, v)
case *Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, *v)
case EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, v)
case *EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, *v)
case Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, v)
case *Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, *v)
case EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, v)
case *EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, *v)
case EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, v)
case *EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, *v)
case EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, v)
case *EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, *v)
case EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, v)
case *EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, *v)
case EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, v)
case *EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, *v)
case EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, v)
case *EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, *v)
case EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, v)
case *EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, *v)
case EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, v)
case *EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, *v)
case EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, v)
case *EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, *v)
case EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, v)
case *EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, *v)
case HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, v)
case *HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, *v)
case KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, v)
case *KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, *v)
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
case MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, v)
case *MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, *v)
case NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, v)
case *NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, *v)
case NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, v)
case *NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, *v)
case P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, v)
case *P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, *v)
case PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, v)
case *PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, *v)
case Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, v)
case *Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, *v)
case PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, v)
case *PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, *v)
case PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, v)
case *PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, *v)
case Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, v)
case *Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, *v)
case ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, v)
case *ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, *v)
case ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, v)
case *ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, *v)
case Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, v)
case *Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, *v)
case PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, v)
case *PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, *v)
case PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, v)
case *PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, *v)
case PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, v)
case *PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, *v)
case PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, v)
case *PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, *v)
case RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, v)
case *RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, *v)
case RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, v)
case *RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, *v)
case SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, v)
case *SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, *v)
case SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, v)
case *SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, *v)
case SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, v)
case *SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, *v)
case SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, v)
case *SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, *v)
case SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, v)
case *SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, *v)
case TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, v)
case *TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, *v)
case Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, v)
case *Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, *v)
case TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, v)
case *TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, *v)
case Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, v)
case *Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, *v)
case ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, v)
case *ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, *v)
case ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, v)
case *ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, *v)
case Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, v)
case *Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, *v)
case VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, v)
case *VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, *v)
case VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, v)
case *VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, *v)
case VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, v)
case *VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandTMEventData(r RandSrc) TMEventData {
switch r.GetUint() % 64 {
case 0:
return RandBc0BlockRequestMessage(r)
case 1:
return RandBc0NoBlockResponseMessage(r)
case 2:
return RandBc0StatusRequestMessage(r)
case 3:
return RandBc0StatusResponseMessage(r)
case 4:
return RandBc1BlockRequestMessage(r)
case 5:
return RandBc1NoBlockResponseMessage(r)
case 6:
return RandBc1StatusRequestMessage(r)
case 7:
return RandBc1StatusResponseMessage(r)
case 8:
return RandBitArray(r)
case 9:
return RandBlockParams(r)
case 10:
return RandBlockPartMessage(r)
case 11:
return RandCommitSig(r)
case 12:
return RandConsensusParams(r)
case 13:
return RandDuplicateVoteEvidence(r)
case 14:
return RandDuration(r)
case 15:
return RandEndHeightMessage(r)
case 16:
return RandEvent(r)
case 17:
return RandEventDataCompleteProposal(r)
case 18:
return RandEventDataNewBlockHeader(r)
case 19:
return RandEventDataNewRound(r)
case 20:
return RandEventDataRoundState(r)
case 21:
return RandEventDataString(r)
case 22:
return RandEventDataTx(r)
case 23:
return RandEventDataValidatorSetUpdates(r)
case 24:
return RandEventDataVote(r)
case 25:
return RandEvidenceListMessage(r)
case 26:
return RandEvidenceParams(r)
case 27:
return RandHasVoteMessage(r)
case 28:
return RandKVPair(r)
case 29:
return RandMockBadEvidence(r)
case 30:
return RandMockGoodEvidence(r)
case 31:
return RandMsgInfo(r)
case 32:
return RandNewRoundStepMessage(r)
case 33:
return RandNewValidBlockMessage(r)
case 34:
return RandP2PID(r)
case 35:
return RandPacketMsg(r)
case 36:
return RandPart(r)
case 37:
return RandPrivKeyEd25519(r)
case 38:
return RandPrivKeySecp256k1(r)
case 39:
return RandProposal(r)
case 40:
return RandProposalMessage(r)
case 41:
return RandProposalPOLMessage(r)
case 42:
return RandProtocol(r)
case 43:
return RandPubKeyEd25519(r)
case 44:
return RandPubKeyMultisigThreshold(r)
case 45:
return RandPubKeyResponse(r)
case 46:
return RandPubKeySecp256k1(r)
case 47:
return RandRemoteSignerError(r)
case 48:
return RandRoundStepType(r)
case 49:
return RandSignProposalRequest(r)
case 50:
return RandSignVoteRequest(r)
case 51:
return RandSignedMsgType(r)
case 52:
return RandSignedProposalResponse(r)
case 53:
return RandSignedVoteResponse(r)
case 54:
return RandTimeoutInfo(r)
case 55:
return RandTx(r)
case 56:
return RandTxMessage(r)
case 57:
return RandValidator(r)
case 58:
return RandValidatorParams(r)
case 59:
return RandValidatorUpdate(r)
case 60:
return RandVote(r)
case 61:
return RandVoteMessage(r)
case 62:
return RandVoteSetBitsMessage(r)
case 63:
return RandVoteSetMaj23Message(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyTMEventData(x TMEventData) TMEventData {
switch v := x.(type) {
case *Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(*v)
return &res
case Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(v)
return res
case *Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(*v)
return &res
case Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(v)
return res
case *Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(*v)
return &res
case Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(v)
return res
case *Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(*v)
return &res
case Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(v)
return res
case *Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(*v)
return &res
case Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(v)
return res
case *Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(*v)
return &res
case Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(v)
return res
case *Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(*v)
return &res
case Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(v)
return res
case *Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(*v)
return &res
case Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(v)
return res
case *BitArray:
res := DeepCopyBitArray(*v)
return &res
case BitArray:
res := DeepCopyBitArray(v)
return res
case *BlockParams:
res := DeepCopyBlockParams(*v)
return &res
case BlockParams:
res := DeepCopyBlockParams(v)
return res
case *BlockPartMessage:
res := DeepCopyBlockPartMessage(*v)
return &res
case BlockPartMessage:
res := DeepCopyBlockPartMessage(v)
return res
case *CommitSig:
res := DeepCopyCommitSig(*v)
return &res
case CommitSig:
res := DeepCopyCommitSig(v)
return res
case *ConsensusParams:
res := DeepCopyConsensusParams(*v)
return &res
case ConsensusParams:
res := DeepCopyConsensusParams(v)
return res
case *DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(*v)
return &res
case DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(v)
return res
case *Duration:
res := DeepCopyDuration(*v)
return &res
case Duration:
res := DeepCopyDuration(v)
return res
case *EndHeightMessage:
res := DeepCopyEndHeightMessage(*v)
return &res
case EndHeightMessage:
res := DeepCopyEndHeightMessage(v)
return res
case *Event:
res := DeepCopyEvent(*v)
return &res
case Event:
res := DeepCopyEvent(v)
return res
case *EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(*v)
return &res
case EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(v)
return res
case *EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(*v)
return &res
case EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(v)
return res
case *EventDataNewRound:
res := DeepCopyEventDataNewRound(*v)
return &res
case EventDataNewRound:
res := DeepCopyEventDataNewRound(v)
return res
case *EventDataRoundState:
res := DeepCopyEventDataRoundState(*v)
return &res
case EventDataRoundState:
res := DeepCopyEventDataRoundState(v)
return res
case *EventDataString:
res := DeepCopyEventDataString(*v)
return &res
case EventDataString:
res := DeepCopyEventDataString(v)
return res
case *EventDataTx:
res := DeepCopyEventDataTx(*v)
return &res
case EventDataTx:
res := DeepCopyEventDataTx(v)
return res
case *EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(*v)
return &res
case EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(v)
return res
case *EventDataVote:
res := DeepCopyEventDataVote(*v)
return &res
case EventDataVote:
res := DeepCopyEventDataVote(v)
return res
case *EvidenceListMessage:
res := DeepCopyEvidenceListMessage(*v)
return &res
case EvidenceListMessage:
res := DeepCopyEvidenceListMessage(v)
return res
case *EvidenceParams:
res := DeepCopyEvidenceParams(*v)
return &res
case EvidenceParams:
res := DeepCopyEvidenceParams(v)
return res
case *HasVoteMessage:
res := DeepCopyHasVoteMessage(*v)
return &res
case HasVoteMessage:
res := DeepCopyHasVoteMessage(v)
return res
case *KVPair:
res := DeepCopyKVPair(*v)
return &res
case KVPair:
res := DeepCopyKVPair(v)
return res
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
case *MsgInfo:
res := DeepCopyMsgInfo(*v)
return &res
case MsgInfo:
res := DeepCopyMsgInfo(v)
return res
case *NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(*v)
return &res
case NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(v)
return res
case *NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(*v)
return &res
case NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(v)
return res
case *P2PID:
res := DeepCopyP2PID(*v)
return &res
case P2PID:
res := DeepCopyP2PID(v)
return res
case *PacketMsg:
res := DeepCopyPacketMsg(*v)
return &res
case PacketMsg:
res := DeepCopyPacketMsg(v)
return res
case *Part:
res := DeepCopyPart(*v)
return &res
case Part:
res := DeepCopyPart(v)
return res
case *PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(*v)
return &res
case PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(v)
return res
case *PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(*v)
return &res
case PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(v)
return res
case *Proposal:
res := DeepCopyProposal(*v)
return &res
case Proposal:
res := DeepCopyProposal(v)
return res
case *ProposalMessage:
res := DeepCopyProposalMessage(*v)
return &res
case ProposalMessage:
res := DeepCopyProposalMessage(v)
return res
case *ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(*v)
return &res
case ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(v)
return res
case *Protocol:
res := DeepCopyProtocol(*v)
return &res
case Protocol:
res := DeepCopyProtocol(v)
return res
case *PubKeyEd25519:
res := DeepCopyPubKeyEd25519(*v)
return &res
case PubKeyEd25519:
res := DeepCopyPubKeyEd25519(v)
return res
case *PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(*v)
return &res
case PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(v)
return res
case *PubKeyResponse:
res := DeepCopyPubKeyResponse(*v)
return &res
case PubKeyResponse:
res := DeepCopyPubKeyResponse(v)
return res
case *PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(*v)
return &res
case PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(v)
return res
case *RemoteSignerError:
res := DeepCopyRemoteSignerError(*v)
return &res
case RemoteSignerError:
res := DeepCopyRemoteSignerError(v)
return res
case *RoundStepType:
res := DeepCopyRoundStepType(*v)
return &res
case RoundStepType:
res := DeepCopyRoundStepType(v)
return res
case *SignProposalRequest:
res := DeepCopySignProposalRequest(*v)
return &res
case SignProposalRequest:
res := DeepCopySignProposalRequest(v)
return res
case *SignVoteRequest:
res := DeepCopySignVoteRequest(*v)
return &res
case SignVoteRequest:
res := DeepCopySignVoteRequest(v)
return res
case *SignedMsgType:
res := DeepCopySignedMsgType(*v)
return &res
case SignedMsgType:
res := DeepCopySignedMsgType(v)
return res
case *SignedProposalResponse:
res := DeepCopySignedProposalResponse(*v)
return &res
case SignedProposalResponse:
res := DeepCopySignedProposalResponse(v)
return res
case *SignedVoteResponse:
res := DeepCopySignedVoteResponse(*v)
return &res
case SignedVoteResponse:
res := DeepCopySignedVoteResponse(v)
return res
case *TimeoutInfo:
res := DeepCopyTimeoutInfo(*v)
return &res
case TimeoutInfo:
res := DeepCopyTimeoutInfo(v)
return res
case *Tx:
res := DeepCopyTx(*v)
return &res
case Tx:
res := DeepCopyTx(v)
return res
case *TxMessage:
res := DeepCopyTxMessage(*v)
return &res
case TxMessage:
res := DeepCopyTxMessage(v)
return res
case *Validator:
res := DeepCopyValidator(*v)
return &res
case Validator:
res := DeepCopyValidator(v)
return res
case *ValidatorParams:
res := DeepCopyValidatorParams(*v)
return &res
case ValidatorParams:
res := DeepCopyValidatorParams(v)
return res
case *ValidatorUpdate:
res := DeepCopyValidatorUpdate(*v)
return &res
case ValidatorUpdate:
res := DeepCopyValidatorUpdate(v)
return res
case *Vote:
res := DeepCopyVote(*v)
return &res
case Vote:
res := DeepCopyVote(v)
return res
case *VoteMessage:
res := DeepCopyVoteMessage(*v)
return &res
case VoteMessage:
res := DeepCopyVoteMessage(v)
return res
case *VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(*v)
return &res
case VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(v)
return res
case *VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(*v)
return &res
case VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodeEvidence(bz []byte) (Evidence, int, error) {
var v Evidence
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeEvidence
func EncodeEvidence(w *[]byte, x interface{}) {
switch v := x.(type) {
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandEvidence(r RandSrc) Evidence {
switch r.GetUint() % 2 {
case 0:
return RandMockBadEvidence(r)
case 1:
return RandMockGoodEvidence(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyEvidence(x Evidence) Evidence {
switch v := x.(type) {
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodeEvidenceMessage(bz []byte) (EvidenceMessage, int, error) {
var v EvidenceMessage
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeEvidenceMessage
func EncodeEvidenceMessage(w *[]byte, x interface{}) {
switch v := x.(type) {
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandEvidenceMessage(r RandSrc) EvidenceMessage {
switch r.GetUint() % 2 {
case 0:
return RandMockBadEvidence(r)
case 1:
return RandMockGoodEvidence(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyEvidenceMessage(x EvidenceMessage) EvidenceMessage {
switch v := x.(type) {
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodeSignerMessage(bz []byte) (SignerMessage, int, error) {
var v SignerMessage
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{252,242,102,165}:
v, n, err := DecodeBc0BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{226,55,89,212}:
v, n, err := DecodeBc0NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{89,2,129,24}:
v, n, err := DecodeBc0StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{181,21,96,69}:
v, n, err := DecodeBc0StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{235,69,170,199}:
v, n, err := DecodeBc1BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{134,105,217,194}:
v, n, err := DecodeBc1NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{45,250,153,33}:
v, n, err := DecodeBc1StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{81,214,62,85}:
v, n, err := DecodeBc1StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{206,243,248,173}:
v, n, err := DecodeBitArray(bz[4:])
return v, n+4, err
case [4]byte{119,191,108,56}:
v, n, err := DecodeBlockParams(bz[4:])
return v, n+4, err
case [4]byte{79,243,233,203}:
v, n, err := DecodeBlockPartMessage(bz[4:])
return v, n+4, err
case [4]byte{163,200,208,51}:
v, n, err := DecodeCommitSig(bz[4:])
return v, n+4, err
case [4]byte{204,104,130,6}:
v, n, err := DecodeConsensusParams(bz[4:])
return v, n+4, err
case [4]byte{89,252,98,178}:
v, n, err := DecodeDuplicateVoteEvidence(bz[4:])
return v, n+4, err
case [4]byte{79,197,42,60}:
v, n, err := DecodeDuration(bz[4:])
return v, n+4, err
case [4]byte{143,131,137,178}:
v, n, err := DecodeEndHeightMessage(bz[4:])
return v, n+4, err
case [4]byte{78,31,73,169}:
v, n, err := DecodeEvent(bz[4:])
return v, n+4, err
case [4]byte{63,236,183,244}:
v, n, err := DecodeEventDataCompleteProposal(bz[4:])
return v, n+4, err
case [4]byte{10,8,210,215}:
v, n, err := DecodeEventDataNewBlockHeader(bz[4:])
return v, n+4, err
case [4]byte{220,234,247,95}:
v, n, err := DecodeEventDataNewRound(bz[4:])
return v, n+4, err
case [4]byte{224,234,115,110}:
v, n, err := DecodeEventDataRoundState(bz[4:])
return v, n+4, err
case [4]byte{197,0,12,147}:
v, n, err := DecodeEventDataString(bz[4:])
return v, n+4, err
case [4]byte{64,227,166,189}:
v, n, err := DecodeEventDataTx(bz[4:])
return v, n+4, err
case [4]byte{79,214,242,86}:
v, n, err := DecodeEventDataValidatorSetUpdates(bz[4:])
return v, n+4, err
case [4]byte{82,175,21,121}:
v, n, err := DecodeEventDataVote(bz[4:])
return v, n+4, err
case [4]byte{51,243,128,228}:
v, n, err := DecodeEvidenceListMessage(bz[4:])
return v, n+4, err
case [4]byte{169,170,199,49}:
v, n, err := DecodeEvidenceParams(bz[4:])
return v, n+4, err
case [4]byte{30,19,14,216}:
v, n, err := DecodeHasVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{63,28,146,21}:
v, n, err := DecodeKVPair(bz[4:])
return v, n+4, err
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
case [4]byte{72,21,63,101}:
v, n, err := DecodeMsgInfo(bz[4:])
return v, n+4, err
case [4]byte{193,38,23,92}:
v, n, err := DecodeNewRoundStepMessage(bz[4:])
return v, n+4, err
case [4]byte{132,190,161,107}:
v, n, err := DecodeNewValidBlockMessage(bz[4:])
return v, n+4, err
case [4]byte{162,154,17,21}:
v, n, err := DecodeP2PID(bz[4:])
return v, n+4, err
case [4]byte{28,180,161,19}:
v, n, err := DecodePacketMsg(bz[4:])
return v, n+4, err
case [4]byte{133,112,162,225}:
v, n, err := DecodePart(bz[4:])
return v, n+4, err
case [4]byte{158,94,112,161}:
v, n, err := DecodePrivKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{83,16,177,42}:
v, n, err := DecodePrivKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{93,66,118,108}:
v, n, err := DecodeProposal(bz[4:])
return v, n+4, err
case [4]byte{162,17,172,204}:
v, n, err := DecodeProposalMessage(bz[4:])
return v, n+4, err
case [4]byte{187,39,126,253}:
v, n, err := DecodeProposalPOLMessage(bz[4:])
return v, n+4, err
case [4]byte{207,8,131,52}:
v, n, err := DecodeProtocol(bz[4:])
return v, n+4, err
case [4]byte{114,76,37,23}:
v, n, err := DecodePubKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{14,33,23,141}:
v, n, err := DecodePubKeyMultisigThreshold(bz[4:])
return v, n+4, err
case [4]byte{1,41,197,74}:
v, n, err := DecodePubKeyResponse(bz[4:])
return v, n+4, err
case [4]byte{51,161,20,197}:
v, n, err := DecodePubKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{198,215,201,27}:
v, n, err := DecodeRemoteSignerError(bz[4:])
return v, n+4, err
case [4]byte{235,161,234,245}:
v, n, err := DecodeRoundStepType(bz[4:])
return v, n+4, err
case [4]byte{9,94,219,98}:
v, n, err := DecodeSignProposalRequest(bz[4:])
return v, n+4, err
case [4]byte{17,244,69,120}:
v, n, err := DecodeSignVoteRequest(bz[4:])
return v, n+4, err
case [4]byte{67,52,162,78}:
v, n, err := DecodeSignedMsgType(bz[4:])
return v, n+4, err
case [4]byte{173,178,196,98}:
v, n, err := DecodeSignedProposalResponse(bz[4:])
return v, n+4, err
case [4]byte{144,17,7,46}:
v, n, err := DecodeSignedVoteResponse(bz[4:])
return v, n+4, err
case [4]byte{47,153,154,201}:
v, n, err := DecodeTimeoutInfo(bz[4:])
return v, n+4, err
case [4]byte{91,158,141,124}:
v, n, err := DecodeTx(bz[4:])
return v, n+4, err
case [4]byte{138,192,166,6}:
v, n, err := DecodeTxMessage(bz[4:])
return v, n+4, err
case [4]byte{11,10,9,103}:
v, n, err := DecodeValidator(bz[4:])
return v, n+4, err
case [4]byte{228,6,69,233}:
v, n, err := DecodeValidatorParams(bz[4:])
return v, n+4, err
case [4]byte{60,45,116,56}:
v, n, err := DecodeValidatorUpdate(bz[4:])
return v, n+4, err
case [4]byte{205,85,136,219}:
v, n, err := DecodeVote(bz[4:])
return v, n+4, err
case [4]byte{239,126,120,100}:
v, n, err := DecodeVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{135,195,3,234}:
v, n, err := DecodeVoteSetBitsMessage(bz[4:])
return v, n+4, err
case [4]byte{232,73,18,176}:
v, n, err := DecodeVoteSetMaj23Message(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeSignerMessage
func EncodeSignerMessage(w *[]byte, x interface{}) {
switch v := x.(type) {
case Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, v)
case *Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, *v)
case Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, v)
case *Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, *v)
case Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, v)
case *Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, *v)
case Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, v)
case *Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, *v)
case Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, v)
case *Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, *v)
case Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, v)
case *Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, *v)
case Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, v)
case *Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, *v)
case Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, v)
case *Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, *v)
case BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, v)
case *BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, *v)
case BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, v)
case *BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, *v)
case BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, v)
case *BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, *v)
case CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, v)
case *CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, *v)
case ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, v)
case *ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, *v)
case DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, v)
case *DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, *v)
case Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, v)
case *Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, *v)
case EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, v)
case *EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, *v)
case Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, v)
case *Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, *v)
case EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, v)
case *EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, *v)
case EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, v)
case *EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, *v)
case EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, v)
case *EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, *v)
case EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, v)
case *EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, *v)
case EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, v)
case *EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, *v)
case EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, v)
case *EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, *v)
case EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, v)
case *EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, *v)
case EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, v)
case *EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, *v)
case EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, v)
case *EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, *v)
case EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, v)
case *EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, *v)
case HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, v)
case *HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, *v)
case KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, v)
case *KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, *v)
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
case MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, v)
case *MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, *v)
case NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, v)
case *NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, *v)
case NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, v)
case *NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, *v)
case P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, v)
case *P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, *v)
case PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, v)
case *PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, *v)
case Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, v)
case *Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, *v)
case PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, v)
case *PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, *v)
case PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, v)
case *PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, *v)
case Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, v)
case *Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, *v)
case ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, v)
case *ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, *v)
case ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, v)
case *ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, *v)
case Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, v)
case *Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, *v)
case PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, v)
case *PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, *v)
case PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, v)
case *PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, *v)
case PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, v)
case *PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, *v)
case PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, v)
case *PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, *v)
case RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, v)
case *RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, *v)
case RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, v)
case *RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, *v)
case SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, v)
case *SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, *v)
case SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, v)
case *SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, *v)
case SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, v)
case *SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, *v)
case SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, v)
case *SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, *v)
case SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, v)
case *SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, *v)
case TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, v)
case *TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, *v)
case Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, v)
case *Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, *v)
case TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, v)
case *TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, *v)
case Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, v)
case *Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, *v)
case ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, v)
case *ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, *v)
case ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, v)
case *ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, *v)
case Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, v)
case *Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, *v)
case VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, v)
case *VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, *v)
case VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, v)
case *VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, *v)
case VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, v)
case *VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandSignerMessage(r RandSrc) SignerMessage {
switch r.GetUint() % 64 {
case 0:
return RandBc0BlockRequestMessage(r)
case 1:
return RandBc0NoBlockResponseMessage(r)
case 2:
return RandBc0StatusRequestMessage(r)
case 3:
return RandBc0StatusResponseMessage(r)
case 4:
return RandBc1BlockRequestMessage(r)
case 5:
return RandBc1NoBlockResponseMessage(r)
case 6:
return RandBc1StatusRequestMessage(r)
case 7:
return RandBc1StatusResponseMessage(r)
case 8:
return RandBitArray(r)
case 9:
return RandBlockParams(r)
case 10:
return RandBlockPartMessage(r)
case 11:
return RandCommitSig(r)
case 12:
return RandConsensusParams(r)
case 13:
return RandDuplicateVoteEvidence(r)
case 14:
return RandDuration(r)
case 15:
return RandEndHeightMessage(r)
case 16:
return RandEvent(r)
case 17:
return RandEventDataCompleteProposal(r)
case 18:
return RandEventDataNewBlockHeader(r)
case 19:
return RandEventDataNewRound(r)
case 20:
return RandEventDataRoundState(r)
case 21:
return RandEventDataString(r)
case 22:
return RandEventDataTx(r)
case 23:
return RandEventDataValidatorSetUpdates(r)
case 24:
return RandEventDataVote(r)
case 25:
return RandEvidenceListMessage(r)
case 26:
return RandEvidenceParams(r)
case 27:
return RandHasVoteMessage(r)
case 28:
return RandKVPair(r)
case 29:
return RandMockBadEvidence(r)
case 30:
return RandMockGoodEvidence(r)
case 31:
return RandMsgInfo(r)
case 32:
return RandNewRoundStepMessage(r)
case 33:
return RandNewValidBlockMessage(r)
case 34:
return RandP2PID(r)
case 35:
return RandPacketMsg(r)
case 36:
return RandPart(r)
case 37:
return RandPrivKeyEd25519(r)
case 38:
return RandPrivKeySecp256k1(r)
case 39:
return RandProposal(r)
case 40:
return RandProposalMessage(r)
case 41:
return RandProposalPOLMessage(r)
case 42:
return RandProtocol(r)
case 43:
return RandPubKeyEd25519(r)
case 44:
return RandPubKeyMultisigThreshold(r)
case 45:
return RandPubKeyResponse(r)
case 46:
return RandPubKeySecp256k1(r)
case 47:
return RandRemoteSignerError(r)
case 48:
return RandRoundStepType(r)
case 49:
return RandSignProposalRequest(r)
case 50:
return RandSignVoteRequest(r)
case 51:
return RandSignedMsgType(r)
case 52:
return RandSignedProposalResponse(r)
case 53:
return RandSignedVoteResponse(r)
case 54:
return RandTimeoutInfo(r)
case 55:
return RandTx(r)
case 56:
return RandTxMessage(r)
case 57:
return RandValidator(r)
case 58:
return RandValidatorParams(r)
case 59:
return RandValidatorUpdate(r)
case 60:
return RandVote(r)
case 61:
return RandVoteMessage(r)
case 62:
return RandVoteSetBitsMessage(r)
case 63:
return RandVoteSetMaj23Message(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopySignerMessage(x SignerMessage) SignerMessage {
switch v := x.(type) {
case *Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(*v)
return &res
case Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(v)
return res
case *Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(*v)
return &res
case Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(v)
return res
case *Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(*v)
return &res
case Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(v)
return res
case *Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(*v)
return &res
case Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(v)
return res
case *Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(*v)
return &res
case Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(v)
return res
case *Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(*v)
return &res
case Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(v)
return res
case *Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(*v)
return &res
case Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(v)
return res
case *Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(*v)
return &res
case Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(v)
return res
case *BitArray:
res := DeepCopyBitArray(*v)
return &res
case BitArray:
res := DeepCopyBitArray(v)
return res
case *BlockParams:
res := DeepCopyBlockParams(*v)
return &res
case BlockParams:
res := DeepCopyBlockParams(v)
return res
case *BlockPartMessage:
res := DeepCopyBlockPartMessage(*v)
return &res
case BlockPartMessage:
res := DeepCopyBlockPartMessage(v)
return res
case *CommitSig:
res := DeepCopyCommitSig(*v)
return &res
case CommitSig:
res := DeepCopyCommitSig(v)
return res
case *ConsensusParams:
res := DeepCopyConsensusParams(*v)
return &res
case ConsensusParams:
res := DeepCopyConsensusParams(v)
return res
case *DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(*v)
return &res
case DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(v)
return res
case *Duration:
res := DeepCopyDuration(*v)
return &res
case Duration:
res := DeepCopyDuration(v)
return res
case *EndHeightMessage:
res := DeepCopyEndHeightMessage(*v)
return &res
case EndHeightMessage:
res := DeepCopyEndHeightMessage(v)
return res
case *Event:
res := DeepCopyEvent(*v)
return &res
case Event:
res := DeepCopyEvent(v)
return res
case *EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(*v)
return &res
case EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(v)
return res
case *EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(*v)
return &res
case EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(v)
return res
case *EventDataNewRound:
res := DeepCopyEventDataNewRound(*v)
return &res
case EventDataNewRound:
res := DeepCopyEventDataNewRound(v)
return res
case *EventDataRoundState:
res := DeepCopyEventDataRoundState(*v)
return &res
case EventDataRoundState:
res := DeepCopyEventDataRoundState(v)
return res
case *EventDataString:
res := DeepCopyEventDataString(*v)
return &res
case EventDataString:
res := DeepCopyEventDataString(v)
return res
case *EventDataTx:
res := DeepCopyEventDataTx(*v)
return &res
case EventDataTx:
res := DeepCopyEventDataTx(v)
return res
case *EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(*v)
return &res
case EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(v)
return res
case *EventDataVote:
res := DeepCopyEventDataVote(*v)
return &res
case EventDataVote:
res := DeepCopyEventDataVote(v)
return res
case *EvidenceListMessage:
res := DeepCopyEvidenceListMessage(*v)
return &res
case EvidenceListMessage:
res := DeepCopyEvidenceListMessage(v)
return res
case *EvidenceParams:
res := DeepCopyEvidenceParams(*v)
return &res
case EvidenceParams:
res := DeepCopyEvidenceParams(v)
return res
case *HasVoteMessage:
res := DeepCopyHasVoteMessage(*v)
return &res
case HasVoteMessage:
res := DeepCopyHasVoteMessage(v)
return res
case *KVPair:
res := DeepCopyKVPair(*v)
return &res
case KVPair:
res := DeepCopyKVPair(v)
return res
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
case *MsgInfo:
res := DeepCopyMsgInfo(*v)
return &res
case MsgInfo:
res := DeepCopyMsgInfo(v)
return res
case *NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(*v)
return &res
case NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(v)
return res
case *NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(*v)
return &res
case NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(v)
return res
case *P2PID:
res := DeepCopyP2PID(*v)
return &res
case P2PID:
res := DeepCopyP2PID(v)
return res
case *PacketMsg:
res := DeepCopyPacketMsg(*v)
return &res
case PacketMsg:
res := DeepCopyPacketMsg(v)
return res
case *Part:
res := DeepCopyPart(*v)
return &res
case Part:
res := DeepCopyPart(v)
return res
case *PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(*v)
return &res
case PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(v)
return res
case *PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(*v)
return &res
case PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(v)
return res
case *Proposal:
res := DeepCopyProposal(*v)
return &res
case Proposal:
res := DeepCopyProposal(v)
return res
case *ProposalMessage:
res := DeepCopyProposalMessage(*v)
return &res
case ProposalMessage:
res := DeepCopyProposalMessage(v)
return res
case *ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(*v)
return &res
case ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(v)
return res
case *Protocol:
res := DeepCopyProtocol(*v)
return &res
case Protocol:
res := DeepCopyProtocol(v)
return res
case *PubKeyEd25519:
res := DeepCopyPubKeyEd25519(*v)
return &res
case PubKeyEd25519:
res := DeepCopyPubKeyEd25519(v)
return res
case *PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(*v)
return &res
case PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(v)
return res
case *PubKeyResponse:
res := DeepCopyPubKeyResponse(*v)
return &res
case PubKeyResponse:
res := DeepCopyPubKeyResponse(v)
return res
case *PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(*v)
return &res
case PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(v)
return res
case *RemoteSignerError:
res := DeepCopyRemoteSignerError(*v)
return &res
case RemoteSignerError:
res := DeepCopyRemoteSignerError(v)
return res
case *RoundStepType:
res := DeepCopyRoundStepType(*v)
return &res
case RoundStepType:
res := DeepCopyRoundStepType(v)
return res
case *SignProposalRequest:
res := DeepCopySignProposalRequest(*v)
return &res
case SignProposalRequest:
res := DeepCopySignProposalRequest(v)
return res
case *SignVoteRequest:
res := DeepCopySignVoteRequest(*v)
return &res
case SignVoteRequest:
res := DeepCopySignVoteRequest(v)
return res
case *SignedMsgType:
res := DeepCopySignedMsgType(*v)
return &res
case SignedMsgType:
res := DeepCopySignedMsgType(v)
return res
case *SignedProposalResponse:
res := DeepCopySignedProposalResponse(*v)
return &res
case SignedProposalResponse:
res := DeepCopySignedProposalResponse(v)
return res
case *SignedVoteResponse:
res := DeepCopySignedVoteResponse(*v)
return &res
case SignedVoteResponse:
res := DeepCopySignedVoteResponse(v)
return res
case *TimeoutInfo:
res := DeepCopyTimeoutInfo(*v)
return &res
case TimeoutInfo:
res := DeepCopyTimeoutInfo(v)
return res
case *Tx:
res := DeepCopyTx(*v)
return &res
case Tx:
res := DeepCopyTx(v)
return res
case *TxMessage:
res := DeepCopyTxMessage(*v)
return &res
case TxMessage:
res := DeepCopyTxMessage(v)
return res
case *Validator:
res := DeepCopyValidator(*v)
return &res
case Validator:
res := DeepCopyValidator(v)
return res
case *ValidatorParams:
res := DeepCopyValidatorParams(*v)
return &res
case ValidatorParams:
res := DeepCopyValidatorParams(v)
return res
case *ValidatorUpdate:
res := DeepCopyValidatorUpdate(*v)
return &res
case ValidatorUpdate:
res := DeepCopyValidatorUpdate(v)
return res
case *Vote:
res := DeepCopyVote(*v)
return &res
case Vote:
res := DeepCopyVote(v)
return res
case *VoteMessage:
res := DeepCopyVoteMessage(*v)
return &res
case VoteMessage:
res := DeepCopyVoteMessage(v)
return res
case *VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(*v)
return &res
case VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(v)
return res
case *VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(*v)
return &res
case VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodeBlockchainMessageV1(bz []byte) (BlockchainMessageV1, int, error) {
var v BlockchainMessageV1
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeBlockchainMessageV1
func EncodeBlockchainMessageV1(w *[]byte, x interface{}) {
switch v := x.(type) {
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandBlockchainMessageV1(r RandSrc) BlockchainMessageV1 {
switch r.GetUint() % 2 {
case 0:
return RandMockBadEvidence(r)
case 1:
return RandMockGoodEvidence(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyBlockchainMessageV1(x BlockchainMessageV1) BlockchainMessageV1 {
switch v := x.(type) {
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodeBlockchainMessageV0(bz []byte) (BlockchainMessageV0, int, error) {
var v BlockchainMessageV0
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeBlockchainMessageV0
func EncodeBlockchainMessageV0(w *[]byte, x interface{}) {
switch v := x.(type) {
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandBlockchainMessageV0(r RandSrc) BlockchainMessageV0 {
switch r.GetUint() % 2 {
case 0:
return RandMockBadEvidence(r)
case 1:
return RandMockGoodEvidence(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyBlockchainMessageV0(x BlockchainMessageV0) BlockchainMessageV0 {
switch v := x.(type) {
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodePexMessage(bz []byte) (PexMessage, int, error) {
var v PexMessage
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{252,242,102,165}:
v, n, err := DecodeBc0BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{226,55,89,212}:
v, n, err := DecodeBc0NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{89,2,129,24}:
v, n, err := DecodeBc0StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{181,21,96,69}:
v, n, err := DecodeBc0StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{235,69,170,199}:
v, n, err := DecodeBc1BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{134,105,217,194}:
v, n, err := DecodeBc1NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{45,250,153,33}:
v, n, err := DecodeBc1StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{81,214,62,85}:
v, n, err := DecodeBc1StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{206,243,248,173}:
v, n, err := DecodeBitArray(bz[4:])
return v, n+4, err
case [4]byte{119,191,108,56}:
v, n, err := DecodeBlockParams(bz[4:])
return v, n+4, err
case [4]byte{79,243,233,203}:
v, n, err := DecodeBlockPartMessage(bz[4:])
return v, n+4, err
case [4]byte{163,200,208,51}:
v, n, err := DecodeCommitSig(bz[4:])
return v, n+4, err
case [4]byte{204,104,130,6}:
v, n, err := DecodeConsensusParams(bz[4:])
return v, n+4, err
case [4]byte{89,252,98,178}:
v, n, err := DecodeDuplicateVoteEvidence(bz[4:])
return v, n+4, err
case [4]byte{79,197,42,60}:
v, n, err := DecodeDuration(bz[4:])
return v, n+4, err
case [4]byte{143,131,137,178}:
v, n, err := DecodeEndHeightMessage(bz[4:])
return v, n+4, err
case [4]byte{78,31,73,169}:
v, n, err := DecodeEvent(bz[4:])
return v, n+4, err
case [4]byte{63,236,183,244}:
v, n, err := DecodeEventDataCompleteProposal(bz[4:])
return v, n+4, err
case [4]byte{10,8,210,215}:
v, n, err := DecodeEventDataNewBlockHeader(bz[4:])
return v, n+4, err
case [4]byte{220,234,247,95}:
v, n, err := DecodeEventDataNewRound(bz[4:])
return v, n+4, err
case [4]byte{224,234,115,110}:
v, n, err := DecodeEventDataRoundState(bz[4:])
return v, n+4, err
case [4]byte{197,0,12,147}:
v, n, err := DecodeEventDataString(bz[4:])
return v, n+4, err
case [4]byte{64,227,166,189}:
v, n, err := DecodeEventDataTx(bz[4:])
return v, n+4, err
case [4]byte{79,214,242,86}:
v, n, err := DecodeEventDataValidatorSetUpdates(bz[4:])
return v, n+4, err
case [4]byte{82,175,21,121}:
v, n, err := DecodeEventDataVote(bz[4:])
return v, n+4, err
case [4]byte{51,243,128,228}:
v, n, err := DecodeEvidenceListMessage(bz[4:])
return v, n+4, err
case [4]byte{169,170,199,49}:
v, n, err := DecodeEvidenceParams(bz[4:])
return v, n+4, err
case [4]byte{30,19,14,216}:
v, n, err := DecodeHasVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{63,28,146,21}:
v, n, err := DecodeKVPair(bz[4:])
return v, n+4, err
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
case [4]byte{72,21,63,101}:
v, n, err := DecodeMsgInfo(bz[4:])
return v, n+4, err
case [4]byte{193,38,23,92}:
v, n, err := DecodeNewRoundStepMessage(bz[4:])
return v, n+4, err
case [4]byte{132,190,161,107}:
v, n, err := DecodeNewValidBlockMessage(bz[4:])
return v, n+4, err
case [4]byte{162,154,17,21}:
v, n, err := DecodeP2PID(bz[4:])
return v, n+4, err
case [4]byte{28,180,161,19}:
v, n, err := DecodePacketMsg(bz[4:])
return v, n+4, err
case [4]byte{133,112,162,225}:
v, n, err := DecodePart(bz[4:])
return v, n+4, err
case [4]byte{158,94,112,161}:
v, n, err := DecodePrivKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{83,16,177,42}:
v, n, err := DecodePrivKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{93,66,118,108}:
v, n, err := DecodeProposal(bz[4:])
return v, n+4, err
case [4]byte{162,17,172,204}:
v, n, err := DecodeProposalMessage(bz[4:])
return v, n+4, err
case [4]byte{187,39,126,253}:
v, n, err := DecodeProposalPOLMessage(bz[4:])
return v, n+4, err
case [4]byte{207,8,131,52}:
v, n, err := DecodeProtocol(bz[4:])
return v, n+4, err
case [4]byte{114,76,37,23}:
v, n, err := DecodePubKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{14,33,23,141}:
v, n, err := DecodePubKeyMultisigThreshold(bz[4:])
return v, n+4, err
case [4]byte{1,41,197,74}:
v, n, err := DecodePubKeyResponse(bz[4:])
return v, n+4, err
case [4]byte{51,161,20,197}:
v, n, err := DecodePubKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{198,215,201,27}:
v, n, err := DecodeRemoteSignerError(bz[4:])
return v, n+4, err
case [4]byte{235,161,234,245}:
v, n, err := DecodeRoundStepType(bz[4:])
return v, n+4, err
case [4]byte{9,94,219,98}:
v, n, err := DecodeSignProposalRequest(bz[4:])
return v, n+4, err
case [4]byte{17,244,69,120}:
v, n, err := DecodeSignVoteRequest(bz[4:])
return v, n+4, err
case [4]byte{67,52,162,78}:
v, n, err := DecodeSignedMsgType(bz[4:])
return v, n+4, err
case [4]byte{173,178,196,98}:
v, n, err := DecodeSignedProposalResponse(bz[4:])
return v, n+4, err
case [4]byte{144,17,7,46}:
v, n, err := DecodeSignedVoteResponse(bz[4:])
return v, n+4, err
case [4]byte{47,153,154,201}:
v, n, err := DecodeTimeoutInfo(bz[4:])
return v, n+4, err
case [4]byte{91,158,141,124}:
v, n, err := DecodeTx(bz[4:])
return v, n+4, err
case [4]byte{138,192,166,6}:
v, n, err := DecodeTxMessage(bz[4:])
return v, n+4, err
case [4]byte{11,10,9,103}:
v, n, err := DecodeValidator(bz[4:])
return v, n+4, err
case [4]byte{228,6,69,233}:
v, n, err := DecodeValidatorParams(bz[4:])
return v, n+4, err
case [4]byte{60,45,116,56}:
v, n, err := DecodeValidatorUpdate(bz[4:])
return v, n+4, err
case [4]byte{205,85,136,219}:
v, n, err := DecodeVote(bz[4:])
return v, n+4, err
case [4]byte{239,126,120,100}:
v, n, err := DecodeVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{135,195,3,234}:
v, n, err := DecodeVoteSetBitsMessage(bz[4:])
return v, n+4, err
case [4]byte{232,73,18,176}:
v, n, err := DecodeVoteSetMaj23Message(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodePexMessage
func EncodePexMessage(w *[]byte, x interface{}) {
switch v := x.(type) {
case Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, v)
case *Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, *v)
case Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, v)
case *Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, *v)
case Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, v)
case *Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, *v)
case Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, v)
case *Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, *v)
case Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, v)
case *Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, *v)
case Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, v)
case *Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, *v)
case Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, v)
case *Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, *v)
case Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, v)
case *Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, *v)
case BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, v)
case *BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, *v)
case BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, v)
case *BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, *v)
case BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, v)
case *BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, *v)
case CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, v)
case *CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, *v)
case ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, v)
case *ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, *v)
case DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, v)
case *DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, *v)
case Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, v)
case *Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, *v)
case EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, v)
case *EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, *v)
case Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, v)
case *Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, *v)
case EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, v)
case *EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, *v)
case EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, v)
case *EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, *v)
case EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, v)
case *EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, *v)
case EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, v)
case *EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, *v)
case EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, v)
case *EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, *v)
case EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, v)
case *EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, *v)
case EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, v)
case *EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, *v)
case EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, v)
case *EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, *v)
case EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, v)
case *EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, *v)
case EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, v)
case *EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, *v)
case HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, v)
case *HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, *v)
case KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, v)
case *KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, *v)
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
case MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, v)
case *MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, *v)
case NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, v)
case *NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, *v)
case NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, v)
case *NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, *v)
case P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, v)
case *P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, *v)
case PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, v)
case *PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, *v)
case Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, v)
case *Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, *v)
case PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, v)
case *PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, *v)
case PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, v)
case *PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, *v)
case Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, v)
case *Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, *v)
case ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, v)
case *ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, *v)
case ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, v)
case *ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, *v)
case Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, v)
case *Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, *v)
case PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, v)
case *PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, *v)
case PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, v)
case *PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, *v)
case PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, v)
case *PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, *v)
case PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, v)
case *PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, *v)
case RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, v)
case *RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, *v)
case RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, v)
case *RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, *v)
case SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, v)
case *SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, *v)
case SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, v)
case *SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, *v)
case SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, v)
case *SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, *v)
case SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, v)
case *SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, *v)
case SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, v)
case *SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, *v)
case TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, v)
case *TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, *v)
case Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, v)
case *Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, *v)
case TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, v)
case *TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, *v)
case Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, v)
case *Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, *v)
case ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, v)
case *ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, *v)
case ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, v)
case *ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, *v)
case Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, v)
case *Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, *v)
case VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, v)
case *VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, *v)
case VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, v)
case *VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, *v)
case VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, v)
case *VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandPexMessage(r RandSrc) PexMessage {
switch r.GetUint() % 64 {
case 0:
return RandBc0BlockRequestMessage(r)
case 1:
return RandBc0NoBlockResponseMessage(r)
case 2:
return RandBc0StatusRequestMessage(r)
case 3:
return RandBc0StatusResponseMessage(r)
case 4:
return RandBc1BlockRequestMessage(r)
case 5:
return RandBc1NoBlockResponseMessage(r)
case 6:
return RandBc1StatusRequestMessage(r)
case 7:
return RandBc1StatusResponseMessage(r)
case 8:
return RandBitArray(r)
case 9:
return RandBlockParams(r)
case 10:
return RandBlockPartMessage(r)
case 11:
return RandCommitSig(r)
case 12:
return RandConsensusParams(r)
case 13:
return RandDuplicateVoteEvidence(r)
case 14:
return RandDuration(r)
case 15:
return RandEndHeightMessage(r)
case 16:
return RandEvent(r)
case 17:
return RandEventDataCompleteProposal(r)
case 18:
return RandEventDataNewBlockHeader(r)
case 19:
return RandEventDataNewRound(r)
case 20:
return RandEventDataRoundState(r)
case 21:
return RandEventDataString(r)
case 22:
return RandEventDataTx(r)
case 23:
return RandEventDataValidatorSetUpdates(r)
case 24:
return RandEventDataVote(r)
case 25:
return RandEvidenceListMessage(r)
case 26:
return RandEvidenceParams(r)
case 27:
return RandHasVoteMessage(r)
case 28:
return RandKVPair(r)
case 29:
return RandMockBadEvidence(r)
case 30:
return RandMockGoodEvidence(r)
case 31:
return RandMsgInfo(r)
case 32:
return RandNewRoundStepMessage(r)
case 33:
return RandNewValidBlockMessage(r)
case 34:
return RandP2PID(r)
case 35:
return RandPacketMsg(r)
case 36:
return RandPart(r)
case 37:
return RandPrivKeyEd25519(r)
case 38:
return RandPrivKeySecp256k1(r)
case 39:
return RandProposal(r)
case 40:
return RandProposalMessage(r)
case 41:
return RandProposalPOLMessage(r)
case 42:
return RandProtocol(r)
case 43:
return RandPubKeyEd25519(r)
case 44:
return RandPubKeyMultisigThreshold(r)
case 45:
return RandPubKeyResponse(r)
case 46:
return RandPubKeySecp256k1(r)
case 47:
return RandRemoteSignerError(r)
case 48:
return RandRoundStepType(r)
case 49:
return RandSignProposalRequest(r)
case 50:
return RandSignVoteRequest(r)
case 51:
return RandSignedMsgType(r)
case 52:
return RandSignedProposalResponse(r)
case 53:
return RandSignedVoteResponse(r)
case 54:
return RandTimeoutInfo(r)
case 55:
return RandTx(r)
case 56:
return RandTxMessage(r)
case 57:
return RandValidator(r)
case 58:
return RandValidatorParams(r)
case 59:
return RandValidatorUpdate(r)
case 60:
return RandVote(r)
case 61:
return RandVoteMessage(r)
case 62:
return RandVoteSetBitsMessage(r)
case 63:
return RandVoteSetMaj23Message(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyPexMessage(x PexMessage) PexMessage {
switch v := x.(type) {
case *Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(*v)
return &res
case Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(v)
return res
case *Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(*v)
return &res
case Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(v)
return res
case *Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(*v)
return &res
case Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(v)
return res
case *Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(*v)
return &res
case Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(v)
return res
case *Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(*v)
return &res
case Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(v)
return res
case *Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(*v)
return &res
case Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(v)
return res
case *Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(*v)
return &res
case Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(v)
return res
case *Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(*v)
return &res
case Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(v)
return res
case *BitArray:
res := DeepCopyBitArray(*v)
return &res
case BitArray:
res := DeepCopyBitArray(v)
return res
case *BlockParams:
res := DeepCopyBlockParams(*v)
return &res
case BlockParams:
res := DeepCopyBlockParams(v)
return res
case *BlockPartMessage:
res := DeepCopyBlockPartMessage(*v)
return &res
case BlockPartMessage:
res := DeepCopyBlockPartMessage(v)
return res
case *CommitSig:
res := DeepCopyCommitSig(*v)
return &res
case CommitSig:
res := DeepCopyCommitSig(v)
return res
case *ConsensusParams:
res := DeepCopyConsensusParams(*v)
return &res
case ConsensusParams:
res := DeepCopyConsensusParams(v)
return res
case *DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(*v)
return &res
case DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(v)
return res
case *Duration:
res := DeepCopyDuration(*v)
return &res
case Duration:
res := DeepCopyDuration(v)
return res
case *EndHeightMessage:
res := DeepCopyEndHeightMessage(*v)
return &res
case EndHeightMessage:
res := DeepCopyEndHeightMessage(v)
return res
case *Event:
res := DeepCopyEvent(*v)
return &res
case Event:
res := DeepCopyEvent(v)
return res
case *EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(*v)
return &res
case EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(v)
return res
case *EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(*v)
return &res
case EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(v)
return res
case *EventDataNewRound:
res := DeepCopyEventDataNewRound(*v)
return &res
case EventDataNewRound:
res := DeepCopyEventDataNewRound(v)
return res
case *EventDataRoundState:
res := DeepCopyEventDataRoundState(*v)
return &res
case EventDataRoundState:
res := DeepCopyEventDataRoundState(v)
return res
case *EventDataString:
res := DeepCopyEventDataString(*v)
return &res
case EventDataString:
res := DeepCopyEventDataString(v)
return res
case *EventDataTx:
res := DeepCopyEventDataTx(*v)
return &res
case EventDataTx:
res := DeepCopyEventDataTx(v)
return res
case *EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(*v)
return &res
case EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(v)
return res
case *EventDataVote:
res := DeepCopyEventDataVote(*v)
return &res
case EventDataVote:
res := DeepCopyEventDataVote(v)
return res
case *EvidenceListMessage:
res := DeepCopyEvidenceListMessage(*v)
return &res
case EvidenceListMessage:
res := DeepCopyEvidenceListMessage(v)
return res
case *EvidenceParams:
res := DeepCopyEvidenceParams(*v)
return &res
case EvidenceParams:
res := DeepCopyEvidenceParams(v)
return res
case *HasVoteMessage:
res := DeepCopyHasVoteMessage(*v)
return &res
case HasVoteMessage:
res := DeepCopyHasVoteMessage(v)
return res
case *KVPair:
res := DeepCopyKVPair(*v)
return &res
case KVPair:
res := DeepCopyKVPair(v)
return res
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
case *MsgInfo:
res := DeepCopyMsgInfo(*v)
return &res
case MsgInfo:
res := DeepCopyMsgInfo(v)
return res
case *NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(*v)
return &res
case NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(v)
return res
case *NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(*v)
return &res
case NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(v)
return res
case *P2PID:
res := DeepCopyP2PID(*v)
return &res
case P2PID:
res := DeepCopyP2PID(v)
return res
case *PacketMsg:
res := DeepCopyPacketMsg(*v)
return &res
case PacketMsg:
res := DeepCopyPacketMsg(v)
return res
case *Part:
res := DeepCopyPart(*v)
return &res
case Part:
res := DeepCopyPart(v)
return res
case *PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(*v)
return &res
case PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(v)
return res
case *PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(*v)
return &res
case PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(v)
return res
case *Proposal:
res := DeepCopyProposal(*v)
return &res
case Proposal:
res := DeepCopyProposal(v)
return res
case *ProposalMessage:
res := DeepCopyProposalMessage(*v)
return &res
case ProposalMessage:
res := DeepCopyProposalMessage(v)
return res
case *ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(*v)
return &res
case ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(v)
return res
case *Protocol:
res := DeepCopyProtocol(*v)
return &res
case Protocol:
res := DeepCopyProtocol(v)
return res
case *PubKeyEd25519:
res := DeepCopyPubKeyEd25519(*v)
return &res
case PubKeyEd25519:
res := DeepCopyPubKeyEd25519(v)
return res
case *PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(*v)
return &res
case PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(v)
return res
case *PubKeyResponse:
res := DeepCopyPubKeyResponse(*v)
return &res
case PubKeyResponse:
res := DeepCopyPubKeyResponse(v)
return res
case *PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(*v)
return &res
case PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(v)
return res
case *RemoteSignerError:
res := DeepCopyRemoteSignerError(*v)
return &res
case RemoteSignerError:
res := DeepCopyRemoteSignerError(v)
return res
case *RoundStepType:
res := DeepCopyRoundStepType(*v)
return &res
case RoundStepType:
res := DeepCopyRoundStepType(v)
return res
case *SignProposalRequest:
res := DeepCopySignProposalRequest(*v)
return &res
case SignProposalRequest:
res := DeepCopySignProposalRequest(v)
return res
case *SignVoteRequest:
res := DeepCopySignVoteRequest(*v)
return &res
case SignVoteRequest:
res := DeepCopySignVoteRequest(v)
return res
case *SignedMsgType:
res := DeepCopySignedMsgType(*v)
return &res
case SignedMsgType:
res := DeepCopySignedMsgType(v)
return res
case *SignedProposalResponse:
res := DeepCopySignedProposalResponse(*v)
return &res
case SignedProposalResponse:
res := DeepCopySignedProposalResponse(v)
return res
case *SignedVoteResponse:
res := DeepCopySignedVoteResponse(*v)
return &res
case SignedVoteResponse:
res := DeepCopySignedVoteResponse(v)
return res
case *TimeoutInfo:
res := DeepCopyTimeoutInfo(*v)
return &res
case TimeoutInfo:
res := DeepCopyTimeoutInfo(v)
return res
case *Tx:
res := DeepCopyTx(*v)
return &res
case Tx:
res := DeepCopyTx(v)
return res
case *TxMessage:
res := DeepCopyTxMessage(*v)
return &res
case TxMessage:
res := DeepCopyTxMessage(v)
return res
case *Validator:
res := DeepCopyValidator(*v)
return &res
case Validator:
res := DeepCopyValidator(v)
return res
case *ValidatorParams:
res := DeepCopyValidatorParams(*v)
return &res
case ValidatorParams:
res := DeepCopyValidatorParams(v)
return res
case *ValidatorUpdate:
res := DeepCopyValidatorUpdate(*v)
return &res
case ValidatorUpdate:
res := DeepCopyValidatorUpdate(v)
return res
case *Vote:
res := DeepCopyVote(*v)
return &res
case Vote:
res := DeepCopyVote(v)
return res
case *VoteMessage:
res := DeepCopyVoteMessage(*v)
return &res
case VoteMessage:
res := DeepCopyVoteMessage(v)
return res
case *VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(*v)
return &res
case VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(v)
return res
case *VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(*v)
return &res
case VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodePacket(bz []byte) (Packet, int, error) {
var v Packet
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{28,180,161,19}:
v, n, err := DecodePacketMsg(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodePacket
func EncodePacket(w *[]byte, x interface{}) {
switch v := x.(type) {
case PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, v)
case *PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandPacket(r RandSrc) Packet {
switch r.GetUint() % 1 {
case 0:
return RandPacketMsg(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyPacket(x Packet) Packet {
switch v := x.(type) {
case *PacketMsg:
res := DeepCopyPacketMsg(*v)
return &res
case PacketMsg:
res := DeepCopyPacketMsg(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
// Interface
func DecodeMempoolMessage(bz []byte) (MempoolMessage, int, error) {
var v MempoolMessage
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{252,242,102,165}:
v, n, err := DecodeBc0BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{226,55,89,212}:
v, n, err := DecodeBc0NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{89,2,129,24}:
v, n, err := DecodeBc0StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{181,21,96,69}:
v, n, err := DecodeBc0StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{235,69,170,199}:
v, n, err := DecodeBc1BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{134,105,217,194}:
v, n, err := DecodeBc1NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{45,250,153,33}:
v, n, err := DecodeBc1StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{81,214,62,85}:
v, n, err := DecodeBc1StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{206,243,248,173}:
v, n, err := DecodeBitArray(bz[4:])
return v, n+4, err
case [4]byte{119,191,108,56}:
v, n, err := DecodeBlockParams(bz[4:])
return v, n+4, err
case [4]byte{79,243,233,203}:
v, n, err := DecodeBlockPartMessage(bz[4:])
return v, n+4, err
case [4]byte{163,200,208,51}:
v, n, err := DecodeCommitSig(bz[4:])
return v, n+4, err
case [4]byte{204,104,130,6}:
v, n, err := DecodeConsensusParams(bz[4:])
return v, n+4, err
case [4]byte{89,252,98,178}:
v, n, err := DecodeDuplicateVoteEvidence(bz[4:])
return v, n+4, err
case [4]byte{79,197,42,60}:
v, n, err := DecodeDuration(bz[4:])
return v, n+4, err
case [4]byte{143,131,137,178}:
v, n, err := DecodeEndHeightMessage(bz[4:])
return v, n+4, err
case [4]byte{78,31,73,169}:
v, n, err := DecodeEvent(bz[4:])
return v, n+4, err
case [4]byte{63,236,183,244}:
v, n, err := DecodeEventDataCompleteProposal(bz[4:])
return v, n+4, err
case [4]byte{10,8,210,215}:
v, n, err := DecodeEventDataNewBlockHeader(bz[4:])
return v, n+4, err
case [4]byte{220,234,247,95}:
v, n, err := DecodeEventDataNewRound(bz[4:])
return v, n+4, err
case [4]byte{224,234,115,110}:
v, n, err := DecodeEventDataRoundState(bz[4:])
return v, n+4, err
case [4]byte{197,0,12,147}:
v, n, err := DecodeEventDataString(bz[4:])
return v, n+4, err
case [4]byte{64,227,166,189}:
v, n, err := DecodeEventDataTx(bz[4:])
return v, n+4, err
case [4]byte{79,214,242,86}:
v, n, err := DecodeEventDataValidatorSetUpdates(bz[4:])
return v, n+4, err
case [4]byte{82,175,21,121}:
v, n, err := DecodeEventDataVote(bz[4:])
return v, n+4, err
case [4]byte{51,243,128,228}:
v, n, err := DecodeEvidenceListMessage(bz[4:])
return v, n+4, err
case [4]byte{169,170,199,49}:
v, n, err := DecodeEvidenceParams(bz[4:])
return v, n+4, err
case [4]byte{30,19,14,216}:
v, n, err := DecodeHasVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{63,28,146,21}:
v, n, err := DecodeKVPair(bz[4:])
return v, n+4, err
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
case [4]byte{72,21,63,101}:
v, n, err := DecodeMsgInfo(bz[4:])
return v, n+4, err
case [4]byte{193,38,23,92}:
v, n, err := DecodeNewRoundStepMessage(bz[4:])
return v, n+4, err
case [4]byte{132,190,161,107}:
v, n, err := DecodeNewValidBlockMessage(bz[4:])
return v, n+4, err
case [4]byte{162,154,17,21}:
v, n, err := DecodeP2PID(bz[4:])
return v, n+4, err
case [4]byte{28,180,161,19}:
v, n, err := DecodePacketMsg(bz[4:])
return v, n+4, err
case [4]byte{133,112,162,225}:
v, n, err := DecodePart(bz[4:])
return v, n+4, err
case [4]byte{158,94,112,161}:
v, n, err := DecodePrivKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{83,16,177,42}:
v, n, err := DecodePrivKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{93,66,118,108}:
v, n, err := DecodeProposal(bz[4:])
return v, n+4, err
case [4]byte{162,17,172,204}:
v, n, err := DecodeProposalMessage(bz[4:])
return v, n+4, err
case [4]byte{187,39,126,253}:
v, n, err := DecodeProposalPOLMessage(bz[4:])
return v, n+4, err
case [4]byte{207,8,131,52}:
v, n, err := DecodeProtocol(bz[4:])
return v, n+4, err
case [4]byte{114,76,37,23}:
v, n, err := DecodePubKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{14,33,23,141}:
v, n, err := DecodePubKeyMultisigThreshold(bz[4:])
return v, n+4, err
case [4]byte{1,41,197,74}:
v, n, err := DecodePubKeyResponse(bz[4:])
return v, n+4, err
case [4]byte{51,161,20,197}:
v, n, err := DecodePubKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{198,215,201,27}:
v, n, err := DecodeRemoteSignerError(bz[4:])
return v, n+4, err
case [4]byte{235,161,234,245}:
v, n, err := DecodeRoundStepType(bz[4:])
return v, n+4, err
case [4]byte{9,94,219,98}:
v, n, err := DecodeSignProposalRequest(bz[4:])
return v, n+4, err
case [4]byte{17,244,69,120}:
v, n, err := DecodeSignVoteRequest(bz[4:])
return v, n+4, err
case [4]byte{67,52,162,78}:
v, n, err := DecodeSignedMsgType(bz[4:])
return v, n+4, err
case [4]byte{173,178,196,98}:
v, n, err := DecodeSignedProposalResponse(bz[4:])
return v, n+4, err
case [4]byte{144,17,7,46}:
v, n, err := DecodeSignedVoteResponse(bz[4:])
return v, n+4, err
case [4]byte{47,153,154,201}:
v, n, err := DecodeTimeoutInfo(bz[4:])
return v, n+4, err
case [4]byte{91,158,141,124}:
v, n, err := DecodeTx(bz[4:])
return v, n+4, err
case [4]byte{138,192,166,6}:
v, n, err := DecodeTxMessage(bz[4:])
return v, n+4, err
case [4]byte{11,10,9,103}:
v, n, err := DecodeValidator(bz[4:])
return v, n+4, err
case [4]byte{228,6,69,233}:
v, n, err := DecodeValidatorParams(bz[4:])
return v, n+4, err
case [4]byte{60,45,116,56}:
v, n, err := DecodeValidatorUpdate(bz[4:])
return v, n+4, err
case [4]byte{205,85,136,219}:
v, n, err := DecodeVote(bz[4:])
return v, n+4, err
case [4]byte{239,126,120,100}:
v, n, err := DecodeVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{135,195,3,234}:
v, n, err := DecodeVoteSetBitsMessage(bz[4:])
return v, n+4, err
case [4]byte{232,73,18,176}:
v, n, err := DecodeVoteSetMaj23Message(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeMempoolMessage
func EncodeMempoolMessage(w *[]byte, x interface{}) {
switch v := x.(type) {
case Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, v)
case *Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, *v)
case Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, v)
case *Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, *v)
case Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, v)
case *Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, *v)
case Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, v)
case *Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, *v)
case Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, v)
case *Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, *v)
case Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, v)
case *Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, *v)
case Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, v)
case *Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, *v)
case Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, v)
case *Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, *v)
case BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, v)
case *BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, *v)
case BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, v)
case *BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, *v)
case BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, v)
case *BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, *v)
case CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, v)
case *CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, *v)
case ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, v)
case *ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, *v)
case DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, v)
case *DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, *v)
case Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, v)
case *Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, *v)
case EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, v)
case *EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, *v)
case Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, v)
case *Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, *v)
case EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, v)
case *EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, *v)
case EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, v)
case *EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, *v)
case EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, v)
case *EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, *v)
case EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, v)
case *EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, *v)
case EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, v)
case *EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, *v)
case EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, v)
case *EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, *v)
case EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, v)
case *EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, *v)
case EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, v)
case *EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, *v)
case EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, v)
case *EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, *v)
case EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, v)
case *EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, *v)
case HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, v)
case *HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, *v)
case KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, v)
case *KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, *v)
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
case MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, v)
case *MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, *v)
case NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, v)
case *NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, *v)
case NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, v)
case *NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, *v)
case P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, v)
case *P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, *v)
case PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, v)
case *PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, *v)
case Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, v)
case *Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, *v)
case PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, v)
case *PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, *v)
case PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, v)
case *PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, *v)
case Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, v)
case *Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, *v)
case ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, v)
case *ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, *v)
case ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, v)
case *ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, *v)
case Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, v)
case *Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, *v)
case PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, v)
case *PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, *v)
case PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, v)
case *PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, *v)
case PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, v)
case *PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, *v)
case PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, v)
case *PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, *v)
case RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, v)
case *RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, *v)
case RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, v)
case *RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, *v)
case SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, v)
case *SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, *v)
case SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, v)
case *SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, *v)
case SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, v)
case *SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, *v)
case SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, v)
case *SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, *v)
case SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, v)
case *SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, *v)
case TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, v)
case *TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, *v)
case Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, v)
case *Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, *v)
case TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, v)
case *TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, *v)
case Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, v)
case *Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, *v)
case ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, v)
case *ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, *v)
case ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, v)
case *ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, *v)
case Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, v)
case *Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, *v)
case VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, v)
case *VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, *v)
case VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, v)
case *VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, *v)
case VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, v)
case *VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func RandMempoolMessage(r RandSrc) MempoolMessage {
switch r.GetUint() % 64 {
case 0:
return RandBc0BlockRequestMessage(r)
case 1:
return RandBc0NoBlockResponseMessage(r)
case 2:
return RandBc0StatusRequestMessage(r)
case 3:
return RandBc0StatusResponseMessage(r)
case 4:
return RandBc1BlockRequestMessage(r)
case 5:
return RandBc1NoBlockResponseMessage(r)
case 6:
return RandBc1StatusRequestMessage(r)
case 7:
return RandBc1StatusResponseMessage(r)
case 8:
return RandBitArray(r)
case 9:
return RandBlockParams(r)
case 10:
return RandBlockPartMessage(r)
case 11:
return RandCommitSig(r)
case 12:
return RandConsensusParams(r)
case 13:
return RandDuplicateVoteEvidence(r)
case 14:
return RandDuration(r)
case 15:
return RandEndHeightMessage(r)
case 16:
return RandEvent(r)
case 17:
return RandEventDataCompleteProposal(r)
case 18:
return RandEventDataNewBlockHeader(r)
case 19:
return RandEventDataNewRound(r)
case 20:
return RandEventDataRoundState(r)
case 21:
return RandEventDataString(r)
case 22:
return RandEventDataTx(r)
case 23:
return RandEventDataValidatorSetUpdates(r)
case 24:
return RandEventDataVote(r)
case 25:
return RandEvidenceListMessage(r)
case 26:
return RandEvidenceParams(r)
case 27:
return RandHasVoteMessage(r)
case 28:
return RandKVPair(r)
case 29:
return RandMockBadEvidence(r)
case 30:
return RandMockGoodEvidence(r)
case 31:
return RandMsgInfo(r)
case 32:
return RandNewRoundStepMessage(r)
case 33:
return RandNewValidBlockMessage(r)
case 34:
return RandP2PID(r)
case 35:
return RandPacketMsg(r)
case 36:
return RandPart(r)
case 37:
return RandPrivKeyEd25519(r)
case 38:
return RandPrivKeySecp256k1(r)
case 39:
return RandProposal(r)
case 40:
return RandProposalMessage(r)
case 41:
return RandProposalPOLMessage(r)
case 42:
return RandProtocol(r)
case 43:
return RandPubKeyEd25519(r)
case 44:
return RandPubKeyMultisigThreshold(r)
case 45:
return RandPubKeyResponse(r)
case 46:
return RandPubKeySecp256k1(r)
case 47:
return RandRemoteSignerError(r)
case 48:
return RandRoundStepType(r)
case 49:
return RandSignProposalRequest(r)
case 50:
return RandSignVoteRequest(r)
case 51:
return RandSignedMsgType(r)
case 52:
return RandSignedProposalResponse(r)
case 53:
return RandSignedVoteResponse(r)
case 54:
return RandTimeoutInfo(r)
case 55:
return RandTx(r)
case 56:
return RandTxMessage(r)
case 57:
return RandValidator(r)
case 58:
return RandValidatorParams(r)
case 59:
return RandValidatorUpdate(r)
case 60:
return RandVote(r)
case 61:
return RandVoteMessage(r)
case 62:
return RandVoteSetBitsMessage(r)
case 63:
return RandVoteSetMaj23Message(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyMempoolMessage(x MempoolMessage) MempoolMessage {
switch v := x.(type) {
case *Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(*v)
return &res
case Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(v)
return res
case *Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(*v)
return &res
case Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(v)
return res
case *Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(*v)
return &res
case Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(v)
return res
case *Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(*v)
return &res
case Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(v)
return res
case *Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(*v)
return &res
case Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(v)
return res
case *Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(*v)
return &res
case Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(v)
return res
case *Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(*v)
return &res
case Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(v)
return res
case *Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(*v)
return &res
case Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(v)
return res
case *BitArray:
res := DeepCopyBitArray(*v)
return &res
case BitArray:
res := DeepCopyBitArray(v)
return res
case *BlockParams:
res := DeepCopyBlockParams(*v)
return &res
case BlockParams:
res := DeepCopyBlockParams(v)
return res
case *BlockPartMessage:
res := DeepCopyBlockPartMessage(*v)
return &res
case BlockPartMessage:
res := DeepCopyBlockPartMessage(v)
return res
case *CommitSig:
res := DeepCopyCommitSig(*v)
return &res
case CommitSig:
res := DeepCopyCommitSig(v)
return res
case *ConsensusParams:
res := DeepCopyConsensusParams(*v)
return &res
case ConsensusParams:
res := DeepCopyConsensusParams(v)
return res
case *DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(*v)
return &res
case DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(v)
return res
case *Duration:
res := DeepCopyDuration(*v)
return &res
case Duration:
res := DeepCopyDuration(v)
return res
case *EndHeightMessage:
res := DeepCopyEndHeightMessage(*v)
return &res
case EndHeightMessage:
res := DeepCopyEndHeightMessage(v)
return res
case *Event:
res := DeepCopyEvent(*v)
return &res
case Event:
res := DeepCopyEvent(v)
return res
case *EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(*v)
return &res
case EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(v)
return res
case *EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(*v)
return &res
case EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(v)
return res
case *EventDataNewRound:
res := DeepCopyEventDataNewRound(*v)
return &res
case EventDataNewRound:
res := DeepCopyEventDataNewRound(v)
return res
case *EventDataRoundState:
res := DeepCopyEventDataRoundState(*v)
return &res
case EventDataRoundState:
res := DeepCopyEventDataRoundState(v)
return res
case *EventDataString:
res := DeepCopyEventDataString(*v)
return &res
case EventDataString:
res := DeepCopyEventDataString(v)
return res
case *EventDataTx:
res := DeepCopyEventDataTx(*v)
return &res
case EventDataTx:
res := DeepCopyEventDataTx(v)
return res
case *EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(*v)
return &res
case EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(v)
return res
case *EventDataVote:
res := DeepCopyEventDataVote(*v)
return &res
case EventDataVote:
res := DeepCopyEventDataVote(v)
return res
case *EvidenceListMessage:
res := DeepCopyEvidenceListMessage(*v)
return &res
case EvidenceListMessage:
res := DeepCopyEvidenceListMessage(v)
return res
case *EvidenceParams:
res := DeepCopyEvidenceParams(*v)
return &res
case EvidenceParams:
res := DeepCopyEvidenceParams(v)
return res
case *HasVoteMessage:
res := DeepCopyHasVoteMessage(*v)
return &res
case HasVoteMessage:
res := DeepCopyHasVoteMessage(v)
return res
case *KVPair:
res := DeepCopyKVPair(*v)
return &res
case KVPair:
res := DeepCopyKVPair(v)
return res
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
case *MsgInfo:
res := DeepCopyMsgInfo(*v)
return &res
case MsgInfo:
res := DeepCopyMsgInfo(v)
return res
case *NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(*v)
return &res
case NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(v)
return res
case *NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(*v)
return &res
case NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(v)
return res
case *P2PID:
res := DeepCopyP2PID(*v)
return &res
case P2PID:
res := DeepCopyP2PID(v)
return res
case *PacketMsg:
res := DeepCopyPacketMsg(*v)
return &res
case PacketMsg:
res := DeepCopyPacketMsg(v)
return res
case *Part:
res := DeepCopyPart(*v)
return &res
case Part:
res := DeepCopyPart(v)
return res
case *PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(*v)
return &res
case PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(v)
return res
case *PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(*v)
return &res
case PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(v)
return res
case *Proposal:
res := DeepCopyProposal(*v)
return &res
case Proposal:
res := DeepCopyProposal(v)
return res
case *ProposalMessage:
res := DeepCopyProposalMessage(*v)
return &res
case ProposalMessage:
res := DeepCopyProposalMessage(v)
return res
case *ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(*v)
return &res
case ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(v)
return res
case *Protocol:
res := DeepCopyProtocol(*v)
return &res
case Protocol:
res := DeepCopyProtocol(v)
return res
case *PubKeyEd25519:
res := DeepCopyPubKeyEd25519(*v)
return &res
case PubKeyEd25519:
res := DeepCopyPubKeyEd25519(v)
return res
case *PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(*v)
return &res
case PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(v)
return res
case *PubKeyResponse:
res := DeepCopyPubKeyResponse(*v)
return &res
case PubKeyResponse:
res := DeepCopyPubKeyResponse(v)
return res
case *PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(*v)
return &res
case PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(v)
return res
case *RemoteSignerError:
res := DeepCopyRemoteSignerError(*v)
return &res
case RemoteSignerError:
res := DeepCopyRemoteSignerError(v)
return res
case *RoundStepType:
res := DeepCopyRoundStepType(*v)
return &res
case RoundStepType:
res := DeepCopyRoundStepType(v)
return res
case *SignProposalRequest:
res := DeepCopySignProposalRequest(*v)
return &res
case SignProposalRequest:
res := DeepCopySignProposalRequest(v)
return res
case *SignVoteRequest:
res := DeepCopySignVoteRequest(*v)
return &res
case SignVoteRequest:
res := DeepCopySignVoteRequest(v)
return res
case *SignedMsgType:
res := DeepCopySignedMsgType(*v)
return &res
case SignedMsgType:
res := DeepCopySignedMsgType(v)
return res
case *SignedProposalResponse:
res := DeepCopySignedProposalResponse(*v)
return &res
case SignedProposalResponse:
res := DeepCopySignedProposalResponse(v)
return res
case *SignedVoteResponse:
res := DeepCopySignedVoteResponse(*v)
return &res
case SignedVoteResponse:
res := DeepCopySignedVoteResponse(v)
return res
case *TimeoutInfo:
res := DeepCopyTimeoutInfo(*v)
return &res
case TimeoutInfo:
res := DeepCopyTimeoutInfo(v)
return res
case *Tx:
res := DeepCopyTx(*v)
return &res
case Tx:
res := DeepCopyTx(v)
return res
case *TxMessage:
res := DeepCopyTxMessage(*v)
return &res
case TxMessage:
res := DeepCopyTxMessage(v)
return res
case *Validator:
res := DeepCopyValidator(*v)
return &res
case Validator:
res := DeepCopyValidator(v)
return res
case *ValidatorParams:
res := DeepCopyValidatorParams(*v)
return &res
case ValidatorParams:
res := DeepCopyValidatorParams(v)
return res
case *ValidatorUpdate:
res := DeepCopyValidatorUpdate(*v)
return &res
case ValidatorUpdate:
res := DeepCopyValidatorUpdate(v)
return res
case *Vote:
res := DeepCopyVote(*v)
return &res
case Vote:
res := DeepCopyVote(v)
return res
case *VoteMessage:
res := DeepCopyVoteMessage(*v)
return &res
case VoteMessage:
res := DeepCopyVoteMessage(v)
return res
case *VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(*v)
return &res
case VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(v)
return res
case *VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(*v)
return &res
case VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func getMagicBytes(name string) []byte {
switch name {
case "Bc0BlockRequestMessage":
return []byte{252,242,102,165}
case "Bc0NoBlockResponseMessage":
return []byte{226,55,89,212}
case "Bc0StatusRequestMessage":
return []byte{89,2,129,24}
case "Bc0StatusResponseMessage":
return []byte{181,21,96,69}
case "Bc1BlockRequestMessage":
return []byte{235,69,170,199}
case "Bc1NoBlockResponseMessage":
return []byte{134,105,217,194}
case "Bc1StatusRequestMessage":
return []byte{45,250,153,33}
case "Bc1StatusResponseMessage":
return []byte{81,214,62,85}
case "BitArray":
return []byte{206,243,248,173}
case "BlockParams":
return []byte{119,191,108,56}
case "BlockPartMessage":
return []byte{79,243,233,203}
case "CommitSig":
return []byte{163,200,208,51}
case "ConsensusParams":
return []byte{204,104,130,6}
case "DuplicateVoteEvidence":
return []byte{89,252,98,178}
case "Duration":
return []byte{79,197,42,60}
case "EndHeightMessage":
return []byte{143,131,137,178}
case "Event":
return []byte{78,31,73,169}
case "EventDataCompleteProposal":
return []byte{63,236,183,244}
case "EventDataNewBlockHeader":
return []byte{10,8,210,215}
case "EventDataNewRound":
return []byte{220,234,247,95}
case "EventDataRoundState":
return []byte{224,234,115,110}
case "EventDataString":
return []byte{197,0,12,147}
case "EventDataTx":
return []byte{64,227,166,189}
case "EventDataValidatorSetUpdates":
return []byte{79,214,242,86}
case "EventDataVote":
return []byte{82,175,21,121}
case "EvidenceListMessage":
return []byte{51,243,128,228}
case "EvidenceParams":
return []byte{169,170,199,49}
case "HasVoteMessage":
return []byte{30,19,14,216}
case "KVPair":
return []byte{63,28,146,21}
case "MockBadEvidence":
return []byte{63,69,132,79}
case "MockGoodEvidence":
return []byte{248,174,56,150}
case "MsgInfo":
return []byte{72,21,63,101}
case "NewRoundStepMessage":
return []byte{193,38,23,92}
case "NewValidBlockMessage":
return []byte{132,190,161,107}
case "P2PID":
return []byte{162,154,17,21}
case "PacketMsg":
return []byte{28,180,161,19}
case "Part":
return []byte{133,112,162,225}
case "PrivKeyEd25519":
return []byte{158,94,112,161}
case "PrivKeySecp256k1":
return []byte{83,16,177,42}
case "Proposal":
return []byte{93,66,118,108}
case "ProposalMessage":
return []byte{162,17,172,204}
case "ProposalPOLMessage":
return []byte{187,39,126,253}
case "Protocol":
return []byte{207,8,131,52}
case "PubKeyEd25519":
return []byte{114,76,37,23}
case "PubKeyMultisigThreshold":
return []byte{14,33,23,141}
case "PubKeyResponse":
return []byte{1,41,197,74}
case "PubKeySecp256k1":
return []byte{51,161,20,197}
case "RemoteSignerError":
return []byte{198,215,201,27}
case "RoundStepType":
return []byte{235,161,234,245}
case "SignProposalRequest":
return []byte{9,94,219,98}
case "SignVoteRequest":
return []byte{17,244,69,120}
case "SignedMsgType":
return []byte{67,52,162,78}
case "SignedProposalResponse":
return []byte{173,178,196,98}
case "SignedVoteResponse":
return []byte{144,17,7,46}
case "TimeoutInfo":
return []byte{47,153,154,201}
case "Tx":
return []byte{91,158,141,124}
case "TxMessage":
return []byte{138,192,166,6}
case "Validator":
return []byte{11,10,9,103}
case "ValidatorParams":
return []byte{228,6,69,233}
case "ValidatorUpdate":
return []byte{60,45,116,56}
case "Vote":
return []byte{205,85,136,219}
case "VoteMessage":
return []byte{239,126,120,100}
case "VoteSetBitsMessage":
return []byte{135,195,3,234}
case "VoteSetMaj23Message":
return []byte{232,73,18,176}
} // end of switch
panic("Should not reach here")
return []byte{}
} // end of getMagicBytes
func getMagicBytesOfVar(x interface{}) ([4]byte, error) {
switch x.(type) {
case *Bc0BlockRequestMessage, Bc0BlockRequestMessage:
return [4]byte{252,242,102,165}, nil
case *Bc0NoBlockResponseMessage, Bc0NoBlockResponseMessage:
return [4]byte{226,55,89,212}, nil
case *Bc0StatusRequestMessage, Bc0StatusRequestMessage:
return [4]byte{89,2,129,24}, nil
case *Bc0StatusResponseMessage, Bc0StatusResponseMessage:
return [4]byte{181,21,96,69}, nil
case *Bc1BlockRequestMessage, Bc1BlockRequestMessage:
return [4]byte{235,69,170,199}, nil
case *Bc1NoBlockResponseMessage, Bc1NoBlockResponseMessage:
return [4]byte{134,105,217,194}, nil
case *Bc1StatusRequestMessage, Bc1StatusRequestMessage:
return [4]byte{45,250,153,33}, nil
case *Bc1StatusResponseMessage, Bc1StatusResponseMessage:
return [4]byte{81,214,62,85}, nil
case *BitArray, BitArray:
return [4]byte{206,243,248,173}, nil
case *BlockParams, BlockParams:
return [4]byte{119,191,108,56}, nil
case *BlockPartMessage, BlockPartMessage:
return [4]byte{79,243,233,203}, nil
case *CommitSig, CommitSig:
return [4]byte{163,200,208,51}, nil
case *ConsensusParams, ConsensusParams:
return [4]byte{204,104,130,6}, nil
case *DuplicateVoteEvidence, DuplicateVoteEvidence:
return [4]byte{89,252,98,178}, nil
case *Duration, Duration:
return [4]byte{79,197,42,60}, nil
case *EndHeightMessage, EndHeightMessage:
return [4]byte{143,131,137,178}, nil
case *Event, Event:
return [4]byte{78,31,73,169}, nil
case *EventDataCompleteProposal, EventDataCompleteProposal:
return [4]byte{63,236,183,244}, nil
case *EventDataNewBlockHeader, EventDataNewBlockHeader:
return [4]byte{10,8,210,215}, nil
case *EventDataNewRound, EventDataNewRound:
return [4]byte{220,234,247,95}, nil
case *EventDataRoundState, EventDataRoundState:
return [4]byte{224,234,115,110}, nil
case *EventDataString, EventDataString:
return [4]byte{197,0,12,147}, nil
case *EventDataTx, EventDataTx:
return [4]byte{64,227,166,189}, nil
case *EventDataValidatorSetUpdates, EventDataValidatorSetUpdates:
return [4]byte{79,214,242,86}, nil
case *EventDataVote, EventDataVote:
return [4]byte{82,175,21,121}, nil
case *EvidenceListMessage, EvidenceListMessage:
return [4]byte{51,243,128,228}, nil
case *EvidenceParams, EvidenceParams:
return [4]byte{169,170,199,49}, nil
case *HasVoteMessage, HasVoteMessage:
return [4]byte{30,19,14,216}, nil
case *KVPair, KVPair:
return [4]byte{63,28,146,21}, nil
case *MockBadEvidence, MockBadEvidence:
return [4]byte{63,69,132,79}, nil
case *MockGoodEvidence, MockGoodEvidence:
return [4]byte{248,174,56,150}, nil
case *MsgInfo, MsgInfo:
return [4]byte{72,21,63,101}, nil
case *NewRoundStepMessage, NewRoundStepMessage:
return [4]byte{193,38,23,92}, nil
case *NewValidBlockMessage, NewValidBlockMessage:
return [4]byte{132,190,161,107}, nil
case *P2PID, P2PID:
return [4]byte{162,154,17,21}, nil
case *PacketMsg, PacketMsg:
return [4]byte{28,180,161,19}, nil
case *Part, Part:
return [4]byte{133,112,162,225}, nil
case *PrivKeyEd25519, PrivKeyEd25519:
return [4]byte{158,94,112,161}, nil
case *PrivKeySecp256k1, PrivKeySecp256k1:
return [4]byte{83,16,177,42}, nil
case *Proposal, Proposal:
return [4]byte{93,66,118,108}, nil
case *ProposalMessage, ProposalMessage:
return [4]byte{162,17,172,204}, nil
case *ProposalPOLMessage, ProposalPOLMessage:
return [4]byte{187,39,126,253}, nil
case *Protocol, Protocol:
return [4]byte{207,8,131,52}, nil
case *PubKeyEd25519, PubKeyEd25519:
return [4]byte{114,76,37,23}, nil
case *PubKeyMultisigThreshold, PubKeyMultisigThreshold:
return [4]byte{14,33,23,141}, nil
case *PubKeyResponse, PubKeyResponse:
return [4]byte{1,41,197,74}, nil
case *PubKeySecp256k1, PubKeySecp256k1:
return [4]byte{51,161,20,197}, nil
case *RemoteSignerError, RemoteSignerError:
return [4]byte{198,215,201,27}, nil
case *RoundStepType, RoundStepType:
return [4]byte{235,161,234,245}, nil
case *SignProposalRequest, SignProposalRequest:
return [4]byte{9,94,219,98}, nil
case *SignVoteRequest, SignVoteRequest:
return [4]byte{17,244,69,120}, nil
case *SignedMsgType, SignedMsgType:
return [4]byte{67,52,162,78}, nil
case *SignedProposalResponse, SignedProposalResponse:
return [4]byte{173,178,196,98}, nil
case *SignedVoteResponse, SignedVoteResponse:
return [4]byte{144,17,7,46}, nil
case *TimeoutInfo, TimeoutInfo:
return [4]byte{47,153,154,201}, nil
case *Tx, Tx:
return [4]byte{91,158,141,124}, nil
case *TxMessage, TxMessage:
return [4]byte{138,192,166,6}, nil
case *Validator, Validator:
return [4]byte{11,10,9,103}, nil
case *ValidatorParams, ValidatorParams:
return [4]byte{228,6,69,233}, nil
case *ValidatorUpdate, ValidatorUpdate:
return [4]byte{60,45,116,56}, nil
case *Vote, Vote:
return [4]byte{205,85,136,219}, nil
case *VoteMessage, VoteMessage:
return [4]byte{239,126,120,100}, nil
case *VoteSetBitsMessage, VoteSetBitsMessage:
return [4]byte{135,195,3,234}, nil
case *VoteSetMaj23Message, VoteSetMaj23Message:
return [4]byte{232,73,18,176}, nil
default:
return [4]byte{0,0,0,0}, errors.New("Unknown Type")
} // end of switch
} // end of func
func EncodeAny(w *[]byte, x interface{}) {
switch v := x.(type) {
case Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, v)
case *Bc0BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc0BlockRequestMessage")...)
EncodeBc0BlockRequestMessage(w, *v)
case Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, v)
case *Bc0NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc0NoBlockResponseMessage")...)
EncodeBc0NoBlockResponseMessage(w, *v)
case Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, v)
case *Bc0StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc0StatusRequestMessage")...)
EncodeBc0StatusRequestMessage(w, *v)
case Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, v)
case *Bc0StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc0StatusResponseMessage")...)
EncodeBc0StatusResponseMessage(w, *v)
case Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, v)
case *Bc1BlockRequestMessage:
*w = append(*w, getMagicBytes("Bc1BlockRequestMessage")...)
EncodeBc1BlockRequestMessage(w, *v)
case Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, v)
case *Bc1NoBlockResponseMessage:
*w = append(*w, getMagicBytes("Bc1NoBlockResponseMessage")...)
EncodeBc1NoBlockResponseMessage(w, *v)
case Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, v)
case *Bc1StatusRequestMessage:
*w = append(*w, getMagicBytes("Bc1StatusRequestMessage")...)
EncodeBc1StatusRequestMessage(w, *v)
case Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, v)
case *Bc1StatusResponseMessage:
*w = append(*w, getMagicBytes("Bc1StatusResponseMessage")...)
EncodeBc1StatusResponseMessage(w, *v)
case BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, v)
case *BitArray:
*w = append(*w, getMagicBytes("BitArray")...)
EncodeBitArray(w, *v)
case BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, v)
case *BlockParams:
*w = append(*w, getMagicBytes("BlockParams")...)
EncodeBlockParams(w, *v)
case BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, v)
case *BlockPartMessage:
*w = append(*w, getMagicBytes("BlockPartMessage")...)
EncodeBlockPartMessage(w, *v)
case CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, v)
case *CommitSig:
*w = append(*w, getMagicBytes("CommitSig")...)
EncodeCommitSig(w, *v)
case ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, v)
case *ConsensusParams:
*w = append(*w, getMagicBytes("ConsensusParams")...)
EncodeConsensusParams(w, *v)
case DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, v)
case *DuplicateVoteEvidence:
*w = append(*w, getMagicBytes("DuplicateVoteEvidence")...)
EncodeDuplicateVoteEvidence(w, *v)
case Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, v)
case *Duration:
*w = append(*w, getMagicBytes("Duration")...)
EncodeDuration(w, *v)
case EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, v)
case *EndHeightMessage:
*w = append(*w, getMagicBytes("EndHeightMessage")...)
EncodeEndHeightMessage(w, *v)
case Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, v)
case *Event:
*w = append(*w, getMagicBytes("Event")...)
EncodeEvent(w, *v)
case EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, v)
case *EventDataCompleteProposal:
*w = append(*w, getMagicBytes("EventDataCompleteProposal")...)
EncodeEventDataCompleteProposal(w, *v)
case EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, v)
case *EventDataNewBlockHeader:
*w = append(*w, getMagicBytes("EventDataNewBlockHeader")...)
EncodeEventDataNewBlockHeader(w, *v)
case EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, v)
case *EventDataNewRound:
*w = append(*w, getMagicBytes("EventDataNewRound")...)
EncodeEventDataNewRound(w, *v)
case EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, v)
case *EventDataRoundState:
*w = append(*w, getMagicBytes("EventDataRoundState")...)
EncodeEventDataRoundState(w, *v)
case EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, v)
case *EventDataString:
*w = append(*w, getMagicBytes("EventDataString")...)
EncodeEventDataString(w, *v)
case EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, v)
case *EventDataTx:
*w = append(*w, getMagicBytes("EventDataTx")...)
EncodeEventDataTx(w, *v)
case EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, v)
case *EventDataValidatorSetUpdates:
*w = append(*w, getMagicBytes("EventDataValidatorSetUpdates")...)
EncodeEventDataValidatorSetUpdates(w, *v)
case EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, v)
case *EventDataVote:
*w = append(*w, getMagicBytes("EventDataVote")...)
EncodeEventDataVote(w, *v)
case EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, v)
case *EvidenceListMessage:
*w = append(*w, getMagicBytes("EvidenceListMessage")...)
EncodeEvidenceListMessage(w, *v)
case EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, v)
case *EvidenceParams:
*w = append(*w, getMagicBytes("EvidenceParams")...)
EncodeEvidenceParams(w, *v)
case HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, v)
case *HasVoteMessage:
*w = append(*w, getMagicBytes("HasVoteMessage")...)
EncodeHasVoteMessage(w, *v)
case KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, v)
case *KVPair:
*w = append(*w, getMagicBytes("KVPair")...)
EncodeKVPair(w, *v)
case MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, v)
case *MockBadEvidence:
*w = append(*w, getMagicBytes("MockBadEvidence")...)
EncodeMockBadEvidence(w, *v)
case MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, v)
case *MockGoodEvidence:
*w = append(*w, getMagicBytes("MockGoodEvidence")...)
EncodeMockGoodEvidence(w, *v)
case MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, v)
case *MsgInfo:
*w = append(*w, getMagicBytes("MsgInfo")...)
EncodeMsgInfo(w, *v)
case NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, v)
case *NewRoundStepMessage:
*w = append(*w, getMagicBytes("NewRoundStepMessage")...)
EncodeNewRoundStepMessage(w, *v)
case NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, v)
case *NewValidBlockMessage:
*w = append(*w, getMagicBytes("NewValidBlockMessage")...)
EncodeNewValidBlockMessage(w, *v)
case P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, v)
case *P2PID:
*w = append(*w, getMagicBytes("P2PID")...)
EncodeP2PID(w, *v)
case PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, v)
case *PacketMsg:
*w = append(*w, getMagicBytes("PacketMsg")...)
EncodePacketMsg(w, *v)
case Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, v)
case *Part:
*w = append(*w, getMagicBytes("Part")...)
EncodePart(w, *v)
case PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, v)
case *PrivKeyEd25519:
*w = append(*w, getMagicBytes("PrivKeyEd25519")...)
EncodePrivKeyEd25519(w, *v)
case PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, v)
case *PrivKeySecp256k1:
*w = append(*w, getMagicBytes("PrivKeySecp256k1")...)
EncodePrivKeySecp256k1(w, *v)
case Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, v)
case *Proposal:
*w = append(*w, getMagicBytes("Proposal")...)
EncodeProposal(w, *v)
case ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, v)
case *ProposalMessage:
*w = append(*w, getMagicBytes("ProposalMessage")...)
EncodeProposalMessage(w, *v)
case ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, v)
case *ProposalPOLMessage:
*w = append(*w, getMagicBytes("ProposalPOLMessage")...)
EncodeProposalPOLMessage(w, *v)
case Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, v)
case *Protocol:
*w = append(*w, getMagicBytes("Protocol")...)
EncodeProtocol(w, *v)
case PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, v)
case *PubKeyEd25519:
*w = append(*w, getMagicBytes("PubKeyEd25519")...)
EncodePubKeyEd25519(w, *v)
case PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, v)
case *PubKeyMultisigThreshold:
*w = append(*w, getMagicBytes("PubKeyMultisigThreshold")...)
EncodePubKeyMultisigThreshold(w, *v)
case PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, v)
case *PubKeyResponse:
*w = append(*w, getMagicBytes("PubKeyResponse")...)
EncodePubKeyResponse(w, *v)
case PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, v)
case *PubKeySecp256k1:
*w = append(*w, getMagicBytes("PubKeySecp256k1")...)
EncodePubKeySecp256k1(w, *v)
case RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, v)
case *RemoteSignerError:
*w = append(*w, getMagicBytes("RemoteSignerError")...)
EncodeRemoteSignerError(w, *v)
case RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, v)
case *RoundStepType:
*w = append(*w, getMagicBytes("RoundStepType")...)
EncodeRoundStepType(w, *v)
case SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, v)
case *SignProposalRequest:
*w = append(*w, getMagicBytes("SignProposalRequest")...)
EncodeSignProposalRequest(w, *v)
case SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, v)
case *SignVoteRequest:
*w = append(*w, getMagicBytes("SignVoteRequest")...)
EncodeSignVoteRequest(w, *v)
case SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, v)
case *SignedMsgType:
*w = append(*w, getMagicBytes("SignedMsgType")...)
EncodeSignedMsgType(w, *v)
case SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, v)
case *SignedProposalResponse:
*w = append(*w, getMagicBytes("SignedProposalResponse")...)
EncodeSignedProposalResponse(w, *v)
case SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, v)
case *SignedVoteResponse:
*w = append(*w, getMagicBytes("SignedVoteResponse")...)
EncodeSignedVoteResponse(w, *v)
case TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, v)
case *TimeoutInfo:
*w = append(*w, getMagicBytes("TimeoutInfo")...)
EncodeTimeoutInfo(w, *v)
case Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, v)
case *Tx:
*w = append(*w, getMagicBytes("Tx")...)
EncodeTx(w, *v)
case TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, v)
case *TxMessage:
*w = append(*w, getMagicBytes("TxMessage")...)
EncodeTxMessage(w, *v)
case Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, v)
case *Validator:
*w = append(*w, getMagicBytes("Validator")...)
EncodeValidator(w, *v)
case ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, v)
case *ValidatorParams:
*w = append(*w, getMagicBytes("ValidatorParams")...)
EncodeValidatorParams(w, *v)
case ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, v)
case *ValidatorUpdate:
*w = append(*w, getMagicBytes("ValidatorUpdate")...)
EncodeValidatorUpdate(w, *v)
case Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, v)
case *Vote:
*w = append(*w, getMagicBytes("Vote")...)
EncodeVote(w, *v)
case VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, v)
case *VoteMessage:
*w = append(*w, getMagicBytes("VoteMessage")...)
EncodeVoteMessage(w, *v)
case VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, v)
case *VoteSetBitsMessage:
*w = append(*w, getMagicBytes("VoteSetBitsMessage")...)
EncodeVoteSetBitsMessage(w, *v)
case VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, v)
case *VoteSetMaj23Message:
*w = append(*w, getMagicBytes("VoteSetMaj23Message")...)
EncodeVoteSetMaj23Message(w, *v)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DecodeAny(bz []byte) (interface{}, int, error) {
var v interface{}
var magicBytes [4]byte
var n int
for i:=0; i<4; i++ {magicBytes[i] = bz[i]}
switch magicBytes {
case [4]byte{252,242,102,165}:
v, n, err := DecodeBc0BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{226,55,89,212}:
v, n, err := DecodeBc0NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{89,2,129,24}:
v, n, err := DecodeBc0StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{181,21,96,69}:
v, n, err := DecodeBc0StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{235,69,170,199}:
v, n, err := DecodeBc1BlockRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{134,105,217,194}:
v, n, err := DecodeBc1NoBlockResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{45,250,153,33}:
v, n, err := DecodeBc1StatusRequestMessage(bz[4:])
return v, n+4, err
case [4]byte{81,214,62,85}:
v, n, err := DecodeBc1StatusResponseMessage(bz[4:])
return v, n+4, err
case [4]byte{206,243,248,173}:
v, n, err := DecodeBitArray(bz[4:])
return v, n+4, err
case [4]byte{119,191,108,56}:
v, n, err := DecodeBlockParams(bz[4:])
return v, n+4, err
case [4]byte{79,243,233,203}:
v, n, err := DecodeBlockPartMessage(bz[4:])
return v, n+4, err
case [4]byte{163,200,208,51}:
v, n, err := DecodeCommitSig(bz[4:])
return v, n+4, err
case [4]byte{204,104,130,6}:
v, n, err := DecodeConsensusParams(bz[4:])
return v, n+4, err
case [4]byte{89,252,98,178}:
v, n, err := DecodeDuplicateVoteEvidence(bz[4:])
return v, n+4, err
case [4]byte{79,197,42,60}:
v, n, err := DecodeDuration(bz[4:])
return v, n+4, err
case [4]byte{143,131,137,178}:
v, n, err := DecodeEndHeightMessage(bz[4:])
return v, n+4, err
case [4]byte{78,31,73,169}:
v, n, err := DecodeEvent(bz[4:])
return v, n+4, err
case [4]byte{63,236,183,244}:
v, n, err := DecodeEventDataCompleteProposal(bz[4:])
return v, n+4, err
case [4]byte{10,8,210,215}:
v, n, err := DecodeEventDataNewBlockHeader(bz[4:])
return v, n+4, err
case [4]byte{220,234,247,95}:
v, n, err := DecodeEventDataNewRound(bz[4:])
return v, n+4, err
case [4]byte{224,234,115,110}:
v, n, err := DecodeEventDataRoundState(bz[4:])
return v, n+4, err
case [4]byte{197,0,12,147}:
v, n, err := DecodeEventDataString(bz[4:])
return v, n+4, err
case [4]byte{64,227,166,189}:
v, n, err := DecodeEventDataTx(bz[4:])
return v, n+4, err
case [4]byte{79,214,242,86}:
v, n, err := DecodeEventDataValidatorSetUpdates(bz[4:])
return v, n+4, err
case [4]byte{82,175,21,121}:
v, n, err := DecodeEventDataVote(bz[4:])
return v, n+4, err
case [4]byte{51,243,128,228}:
v, n, err := DecodeEvidenceListMessage(bz[4:])
return v, n+4, err
case [4]byte{169,170,199,49}:
v, n, err := DecodeEvidenceParams(bz[4:])
return v, n+4, err
case [4]byte{30,19,14,216}:
v, n, err := DecodeHasVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{63,28,146,21}:
v, n, err := DecodeKVPair(bz[4:])
return v, n+4, err
case [4]byte{63,69,132,79}:
v, n, err := DecodeMockBadEvidence(bz[4:])
return v, n+4, err
case [4]byte{248,174,56,150}:
v, n, err := DecodeMockGoodEvidence(bz[4:])
return v, n+4, err
case [4]byte{72,21,63,101}:
v, n, err := DecodeMsgInfo(bz[4:])
return v, n+4, err
case [4]byte{193,38,23,92}:
v, n, err := DecodeNewRoundStepMessage(bz[4:])
return v, n+4, err
case [4]byte{132,190,161,107}:
v, n, err := DecodeNewValidBlockMessage(bz[4:])
return v, n+4, err
case [4]byte{162,154,17,21}:
v, n, err := DecodeP2PID(bz[4:])
return v, n+4, err
case [4]byte{28,180,161,19}:
v, n, err := DecodePacketMsg(bz[4:])
return v, n+4, err
case [4]byte{133,112,162,225}:
v, n, err := DecodePart(bz[4:])
return v, n+4, err
case [4]byte{158,94,112,161}:
v, n, err := DecodePrivKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{83,16,177,42}:
v, n, err := DecodePrivKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{93,66,118,108}:
v, n, err := DecodeProposal(bz[4:])
return v, n+4, err
case [4]byte{162,17,172,204}:
v, n, err := DecodeProposalMessage(bz[4:])
return v, n+4, err
case [4]byte{187,39,126,253}:
v, n, err := DecodeProposalPOLMessage(bz[4:])
return v, n+4, err
case [4]byte{207,8,131,52}:
v, n, err := DecodeProtocol(bz[4:])
return v, n+4, err
case [4]byte{114,76,37,23}:
v, n, err := DecodePubKeyEd25519(bz[4:])
return v, n+4, err
case [4]byte{14,33,23,141}:
v, n, err := DecodePubKeyMultisigThreshold(bz[4:])
return v, n+4, err
case [4]byte{1,41,197,74}:
v, n, err := DecodePubKeyResponse(bz[4:])
return v, n+4, err
case [4]byte{51,161,20,197}:
v, n, err := DecodePubKeySecp256k1(bz[4:])
return v, n+4, err
case [4]byte{198,215,201,27}:
v, n, err := DecodeRemoteSignerError(bz[4:])
return v, n+4, err
case [4]byte{235,161,234,245}:
v, n, err := DecodeRoundStepType(bz[4:])
return v, n+4, err
case [4]byte{9,94,219,98}:
v, n, err := DecodeSignProposalRequest(bz[4:])
return v, n+4, err
case [4]byte{17,244,69,120}:
v, n, err := DecodeSignVoteRequest(bz[4:])
return v, n+4, err
case [4]byte{67,52,162,78}:
v, n, err := DecodeSignedMsgType(bz[4:])
return v, n+4, err
case [4]byte{173,178,196,98}:
v, n, err := DecodeSignedProposalResponse(bz[4:])
return v, n+4, err
case [4]byte{144,17,7,46}:
v, n, err := DecodeSignedVoteResponse(bz[4:])
return v, n+4, err
case [4]byte{47,153,154,201}:
v, n, err := DecodeTimeoutInfo(bz[4:])
return v, n+4, err
case [4]byte{91,158,141,124}:
v, n, err := DecodeTx(bz[4:])
return v, n+4, err
case [4]byte{138,192,166,6}:
v, n, err := DecodeTxMessage(bz[4:])
return v, n+4, err
case [4]byte{11,10,9,103}:
v, n, err := DecodeValidator(bz[4:])
return v, n+4, err
case [4]byte{228,6,69,233}:
v, n, err := DecodeValidatorParams(bz[4:])
return v, n+4, err
case [4]byte{60,45,116,56}:
v, n, err := DecodeValidatorUpdate(bz[4:])
return v, n+4, err
case [4]byte{205,85,136,219}:
v, n, err := DecodeVote(bz[4:])
return v, n+4, err
case [4]byte{239,126,120,100}:
v, n, err := DecodeVoteMessage(bz[4:])
return v, n+4, err
case [4]byte{135,195,3,234}:
v, n, err := DecodeVoteSetBitsMessage(bz[4:])
return v, n+4, err
case [4]byte{232,73,18,176}:
v, n, err := DecodeVoteSetMaj23Message(bz[4:])
return v, n+4, err
default:
panic("Unknown type")
} // end of switch
return v, n, nil
} // end of DecodeAny
func RandAny(r RandSrc) interface{} {
switch r.GetUint() % 64 {
case 0:
return RandBc0BlockRequestMessage(r)
case 1:
return RandBc0NoBlockResponseMessage(r)
case 2:
return RandBc0StatusRequestMessage(r)
case 3:
return RandBc0StatusResponseMessage(r)
case 4:
return RandBc1BlockRequestMessage(r)
case 5:
return RandBc1NoBlockResponseMessage(r)
case 6:
return RandBc1StatusRequestMessage(r)
case 7:
return RandBc1StatusResponseMessage(r)
case 8:
return RandBitArray(r)
case 9:
return RandBlockParams(r)
case 10:
return RandBlockPartMessage(r)
case 11:
return RandCommitSig(r)
case 12:
return RandConsensusParams(r)
case 13:
return RandDuplicateVoteEvidence(r)
case 14:
return RandDuration(r)
case 15:
return RandEndHeightMessage(r)
case 16:
return RandEvent(r)
case 17:
return RandEventDataCompleteProposal(r)
case 18:
return RandEventDataNewBlockHeader(r)
case 19:
return RandEventDataNewRound(r)
case 20:
return RandEventDataRoundState(r)
case 21:
return RandEventDataString(r)
case 22:
return RandEventDataTx(r)
case 23:
return RandEventDataValidatorSetUpdates(r)
case 24:
return RandEventDataVote(r)
case 25:
return RandEvidenceListMessage(r)
case 26:
return RandEvidenceParams(r)
case 27:
return RandHasVoteMessage(r)
case 28:
return RandKVPair(r)
case 29:
return RandMockBadEvidence(r)
case 30:
return RandMockGoodEvidence(r)
case 31:
return RandMsgInfo(r)
case 32:
return RandNewRoundStepMessage(r)
case 33:
return RandNewValidBlockMessage(r)
case 34:
return RandP2PID(r)
case 35:
return RandPacketMsg(r)
case 36:
return RandPart(r)
case 37:
return RandPrivKeyEd25519(r)
case 38:
return RandPrivKeySecp256k1(r)
case 39:
return RandProposal(r)
case 40:
return RandProposalMessage(r)
case 41:
return RandProposalPOLMessage(r)
case 42:
return RandProtocol(r)
case 43:
return RandPubKeyEd25519(r)
case 44:
return RandPubKeyMultisigThreshold(r)
case 45:
return RandPubKeyResponse(r)
case 46:
return RandPubKeySecp256k1(r)
case 47:
return RandRemoteSignerError(r)
case 48:
return RandRoundStepType(r)
case 49:
return RandSignProposalRequest(r)
case 50:
return RandSignVoteRequest(r)
case 51:
return RandSignedMsgType(r)
case 52:
return RandSignedProposalResponse(r)
case 53:
return RandSignedVoteResponse(r)
case 54:
return RandTimeoutInfo(r)
case 55:
return RandTx(r)
case 56:
return RandTxMessage(r)
case 57:
return RandValidator(r)
case 58:
return RandValidatorParams(r)
case 59:
return RandValidatorUpdate(r)
case 60:
return RandVote(r)
case 61:
return RandVoteMessage(r)
case 62:
return RandVoteSetBitsMessage(r)
case 63:
return RandVoteSetMaj23Message(r)
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func DeepCopyAny(x interface{}) interface{} {
switch v := x.(type) {
case *Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(*v)
return &res
case Bc0BlockRequestMessage:
res := DeepCopyBc0BlockRequestMessage(v)
return res
case *Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(*v)
return &res
case Bc0NoBlockResponseMessage:
res := DeepCopyBc0NoBlockResponseMessage(v)
return res
case *Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(*v)
return &res
case Bc0StatusRequestMessage:
res := DeepCopyBc0StatusRequestMessage(v)
return res
case *Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(*v)
return &res
case Bc0StatusResponseMessage:
res := DeepCopyBc0StatusResponseMessage(v)
return res
case *Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(*v)
return &res
case Bc1BlockRequestMessage:
res := DeepCopyBc1BlockRequestMessage(v)
return res
case *Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(*v)
return &res
case Bc1NoBlockResponseMessage:
res := DeepCopyBc1NoBlockResponseMessage(v)
return res
case *Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(*v)
return &res
case Bc1StatusRequestMessage:
res := DeepCopyBc1StatusRequestMessage(v)
return res
case *Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(*v)
return &res
case Bc1StatusResponseMessage:
res := DeepCopyBc1StatusResponseMessage(v)
return res
case *BitArray:
res := DeepCopyBitArray(*v)
return &res
case BitArray:
res := DeepCopyBitArray(v)
return res
case *BlockParams:
res := DeepCopyBlockParams(*v)
return &res
case BlockParams:
res := DeepCopyBlockParams(v)
return res
case *BlockPartMessage:
res := DeepCopyBlockPartMessage(*v)
return &res
case BlockPartMessage:
res := DeepCopyBlockPartMessage(v)
return res
case *CommitSig:
res := DeepCopyCommitSig(*v)
return &res
case CommitSig:
res := DeepCopyCommitSig(v)
return res
case *ConsensusParams:
res := DeepCopyConsensusParams(*v)
return &res
case ConsensusParams:
res := DeepCopyConsensusParams(v)
return res
case *DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(*v)
return &res
case DuplicateVoteEvidence:
res := DeepCopyDuplicateVoteEvidence(v)
return res
case *Duration:
res := DeepCopyDuration(*v)
return &res
case Duration:
res := DeepCopyDuration(v)
return res
case *EndHeightMessage:
res := DeepCopyEndHeightMessage(*v)
return &res
case EndHeightMessage:
res := DeepCopyEndHeightMessage(v)
return res
case *Event:
res := DeepCopyEvent(*v)
return &res
case Event:
res := DeepCopyEvent(v)
return res
case *EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(*v)
return &res
case EventDataCompleteProposal:
res := DeepCopyEventDataCompleteProposal(v)
return res
case *EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(*v)
return &res
case EventDataNewBlockHeader:
res := DeepCopyEventDataNewBlockHeader(v)
return res
case *EventDataNewRound:
res := DeepCopyEventDataNewRound(*v)
return &res
case EventDataNewRound:
res := DeepCopyEventDataNewRound(v)
return res
case *EventDataRoundState:
res := DeepCopyEventDataRoundState(*v)
return &res
case EventDataRoundState:
res := DeepCopyEventDataRoundState(v)
return res
case *EventDataString:
res := DeepCopyEventDataString(*v)
return &res
case EventDataString:
res := DeepCopyEventDataString(v)
return res
case *EventDataTx:
res := DeepCopyEventDataTx(*v)
return &res
case EventDataTx:
res := DeepCopyEventDataTx(v)
return res
case *EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(*v)
return &res
case EventDataValidatorSetUpdates:
res := DeepCopyEventDataValidatorSetUpdates(v)
return res
case *EventDataVote:
res := DeepCopyEventDataVote(*v)
return &res
case EventDataVote:
res := DeepCopyEventDataVote(v)
return res
case *EvidenceListMessage:
res := DeepCopyEvidenceListMessage(*v)
return &res
case EvidenceListMessage:
res := DeepCopyEvidenceListMessage(v)
return res
case *EvidenceParams:
res := DeepCopyEvidenceParams(*v)
return &res
case EvidenceParams:
res := DeepCopyEvidenceParams(v)
return res
case *HasVoteMessage:
res := DeepCopyHasVoteMessage(*v)
return &res
case HasVoteMessage:
res := DeepCopyHasVoteMessage(v)
return res
case *KVPair:
res := DeepCopyKVPair(*v)
return &res
case KVPair:
res := DeepCopyKVPair(v)
return res
case *MockBadEvidence:
res := DeepCopyMockBadEvidence(*v)
return &res
case MockBadEvidence:
res := DeepCopyMockBadEvidence(v)
return res
case *MockGoodEvidence:
res := DeepCopyMockGoodEvidence(*v)
return &res
case MockGoodEvidence:
res := DeepCopyMockGoodEvidence(v)
return res
case *MsgInfo:
res := DeepCopyMsgInfo(*v)
return &res
case MsgInfo:
res := DeepCopyMsgInfo(v)
return res
case *NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(*v)
return &res
case NewRoundStepMessage:
res := DeepCopyNewRoundStepMessage(v)
return res
case *NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(*v)
return &res
case NewValidBlockMessage:
res := DeepCopyNewValidBlockMessage(v)
return res
case *P2PID:
res := DeepCopyP2PID(*v)
return &res
case P2PID:
res := DeepCopyP2PID(v)
return res
case *PacketMsg:
res := DeepCopyPacketMsg(*v)
return &res
case PacketMsg:
res := DeepCopyPacketMsg(v)
return res
case *Part:
res := DeepCopyPart(*v)
return &res
case Part:
res := DeepCopyPart(v)
return res
case *PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(*v)
return &res
case PrivKeyEd25519:
res := DeepCopyPrivKeyEd25519(v)
return res
case *PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(*v)
return &res
case PrivKeySecp256k1:
res := DeepCopyPrivKeySecp256k1(v)
return res
case *Proposal:
res := DeepCopyProposal(*v)
return &res
case Proposal:
res := DeepCopyProposal(v)
return res
case *ProposalMessage:
res := DeepCopyProposalMessage(*v)
return &res
case ProposalMessage:
res := DeepCopyProposalMessage(v)
return res
case *ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(*v)
return &res
case ProposalPOLMessage:
res := DeepCopyProposalPOLMessage(v)
return res
case *Protocol:
res := DeepCopyProtocol(*v)
return &res
case Protocol:
res := DeepCopyProtocol(v)
return res
case *PubKeyEd25519:
res := DeepCopyPubKeyEd25519(*v)
return &res
case PubKeyEd25519:
res := DeepCopyPubKeyEd25519(v)
return res
case *PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(*v)
return &res
case PubKeyMultisigThreshold:
res := DeepCopyPubKeyMultisigThreshold(v)
return res
case *PubKeyResponse:
res := DeepCopyPubKeyResponse(*v)
return &res
case PubKeyResponse:
res := DeepCopyPubKeyResponse(v)
return res
case *PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(*v)
return &res
case PubKeySecp256k1:
res := DeepCopyPubKeySecp256k1(v)
return res
case *RemoteSignerError:
res := DeepCopyRemoteSignerError(*v)
return &res
case RemoteSignerError:
res := DeepCopyRemoteSignerError(v)
return res
case *RoundStepType:
res := DeepCopyRoundStepType(*v)
return &res
case RoundStepType:
res := DeepCopyRoundStepType(v)
return res
case *SignProposalRequest:
res := DeepCopySignProposalRequest(*v)
return &res
case SignProposalRequest:
res := DeepCopySignProposalRequest(v)
return res
case *SignVoteRequest:
res := DeepCopySignVoteRequest(*v)
return &res
case SignVoteRequest:
res := DeepCopySignVoteRequest(v)
return res
case *SignedMsgType:
res := DeepCopySignedMsgType(*v)
return &res
case SignedMsgType:
res := DeepCopySignedMsgType(v)
return res
case *SignedProposalResponse:
res := DeepCopySignedProposalResponse(*v)
return &res
case SignedProposalResponse:
res := DeepCopySignedProposalResponse(v)
return res
case *SignedVoteResponse:
res := DeepCopySignedVoteResponse(*v)
return &res
case SignedVoteResponse:
res := DeepCopySignedVoteResponse(v)
return res
case *TimeoutInfo:
res := DeepCopyTimeoutInfo(*v)
return &res
case TimeoutInfo:
res := DeepCopyTimeoutInfo(v)
return res
case *Tx:
res := DeepCopyTx(*v)
return &res
case Tx:
res := DeepCopyTx(v)
return res
case *TxMessage:
res := DeepCopyTxMessage(*v)
return &res
case TxMessage:
res := DeepCopyTxMessage(v)
return res
case *Validator:
res := DeepCopyValidator(*v)
return &res
case Validator:
res := DeepCopyValidator(v)
return res
case *ValidatorParams:
res := DeepCopyValidatorParams(*v)
return &res
case ValidatorParams:
res := DeepCopyValidatorParams(v)
return res
case *ValidatorUpdate:
res := DeepCopyValidatorUpdate(*v)
return &res
case ValidatorUpdate:
res := DeepCopyValidatorUpdate(v)
return res
case *Vote:
res := DeepCopyVote(*v)
return &res
case Vote:
res := DeepCopyVote(v)
return res
case *VoteMessage:
res := DeepCopyVoteMessage(*v)
return &res
case VoteMessage:
res := DeepCopyVoteMessage(v)
return res
case *VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(*v)
return &res
case VoteSetBitsMessage:
res := DeepCopyVoteSetBitsMessage(v)
return res
case *VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(*v)
return &res
case VoteSetMaj23Message:
res := DeepCopyVoteSetMaj23Message(v)
return res
default:
panic("Unknown Type.")
} // end of switch
} // end of func
func GetSupportList() []string {
return []string {
"github.com/tendermint/tendermint/abci/types.BlockParams",
"github.com/tendermint/tendermint/abci/types.ConsensusParams",
"github.com/tendermint/tendermint/abci/types.Event",
"github.com/tendermint/tendermint/abci/types.EvidenceParams",
"github.com/tendermint/tendermint/abci/types.ValidatorParams",
"github.com/tendermint/tendermint/abci/types.ValidatorUpdate",
"github.com/tendermint/tendermint/blockchain/v0.BlockchainMessage",
"github.com/tendermint/tendermint/blockchain/v0.bcBlockRequestMessage",
"github.com/tendermint/tendermint/blockchain/v0.bcNoBlockResponseMessage",
"github.com/tendermint/tendermint/blockchain/v0.bcStatusRequestMessage",
"github.com/tendermint/tendermint/blockchain/v0.bcStatusResponseMessage",
"github.com/tendermint/tendermint/blockchain/v1.BlockchainMessage",
"github.com/tendermint/tendermint/blockchain/v1.bcBlockRequestMessage",
"github.com/tendermint/tendermint/blockchain/v1.bcNoBlockResponseMessage",
"github.com/tendermint/tendermint/blockchain/v1.bcStatusRequestMessage",
"github.com/tendermint/tendermint/blockchain/v1.bcStatusResponseMessage",
"github.com/tendermint/tendermint/consensus.BlockPartMessage",
"github.com/tendermint/tendermint/consensus.ConsensusMessage",
"github.com/tendermint/tendermint/consensus.EndHeightMessage",
"github.com/tendermint/tendermint/consensus.HasVoteMessage",
"github.com/tendermint/tendermint/consensus.NewRoundStepMessage",
"github.com/tendermint/tendermint/consensus.NewValidBlockMessage",
"github.com/tendermint/tendermint/consensus.ProposalMessage",
"github.com/tendermint/tendermint/consensus.ProposalPOLMessage",
"github.com/tendermint/tendermint/consensus.VoteMessage",
"github.com/tendermint/tendermint/consensus.VoteSetBitsMessage",
"github.com/tendermint/tendermint/consensus.VoteSetMaj23Message",
"github.com/tendermint/tendermint/consensus.WALMessage",
"github.com/tendermint/tendermint/consensus.msgInfo",
"github.com/tendermint/tendermint/consensus.timeoutInfo",
"github.com/tendermint/tendermint/consensus/types.RoundStepType",
"github.com/tendermint/tendermint/crypto.PrivKey",
"github.com/tendermint/tendermint/crypto.PubKey",
"github.com/tendermint/tendermint/crypto/ed25519.PrivKeyEd25519",
"github.com/tendermint/tendermint/crypto/ed25519.PubKeyEd25519",
"github.com/tendermint/tendermint/crypto/multisig.PubKeyMultisigThreshold",
"github.com/tendermint/tendermint/crypto/secp256k1.PrivKeySecp256k1",
"github.com/tendermint/tendermint/crypto/secp256k1.PubKeySecp256k1",
"github.com/tendermint/tendermint/evidence.EvidenceListMessage",
"github.com/tendermint/tendermint/evidence.EvidenceMessage",
"github.com/tendermint/tendermint/libs/common.BitArray",
"github.com/tendermint/tendermint/libs/common.KVPair",
"github.com/tendermint/tendermint/mempool.MempoolMessage",
"github.com/tendermint/tendermint/mempool.TxMessage",
"github.com/tendermint/tendermint/p2p.ID",
"github.com/tendermint/tendermint/p2p/conn.Packet",
"github.com/tendermint/tendermint/p2p/conn.PacketMsg",
"github.com/tendermint/tendermint/p2p/pex.PexMessage",
"github.com/tendermint/tendermint/privval.PubKeyResponse",
"github.com/tendermint/tendermint/privval.RemoteSignerError",
"github.com/tendermint/tendermint/privval.SignProposalRequest",
"github.com/tendermint/tendermint/privval.SignVoteRequest",
"github.com/tendermint/tendermint/privval.SignedProposalResponse",
"github.com/tendermint/tendermint/privval.SignedVoteResponse",
"github.com/tendermint/tendermint/privval.SignerMessage",
"github.com/tendermint/tendermint/types.CommitSig",
"github.com/tendermint/tendermint/types.DuplicateVoteEvidence",
"github.com/tendermint/tendermint/types.EventDataCompleteProposal",
"github.com/tendermint/tendermint/types.EventDataNewBlockHeader",
"github.com/tendermint/tendermint/types.EventDataNewRound",
"github.com/tendermint/tendermint/types.EventDataRoundState",
"github.com/tendermint/tendermint/types.EventDataString",
"github.com/tendermint/tendermint/types.EventDataTx",
"github.com/tendermint/tendermint/types.EventDataValidatorSetUpdates",
"github.com/tendermint/tendermint/types.EventDataVote",
"github.com/tendermint/tendermint/types.Evidence",
"github.com/tendermint/tendermint/types.MockBadEvidence",
"github.com/tendermint/tendermint/types.MockGoodEvidence",
"github.com/tendermint/tendermint/types.Part",
"github.com/tendermint/tendermint/types.Proposal",
"github.com/tendermint/tendermint/types.SignedMsgType",
"github.com/tendermint/tendermint/types.TMEventData",
"github.com/tendermint/tendermint/types.Tx",
"github.com/tendermint/tendermint/types.Validator",
"github.com/tendermint/tendermint/types.Vote",
"github.com/tendermint/tendermint/version.Protocol",
"time.Duration",
}
} // end of GetSupportList
