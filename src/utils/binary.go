package utils

import (
	"math"
	"bytes"
	"io/ioutil"
	"encoding/binary"

	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/japanese"
)

type BinaryReader struct {
	Position		int
	Bytes			[]byte
}

func NewBinaryReader(data []byte) *BinaryReader {
	return &BinaryReader{0, data}
}

func (r *BinaryReader) Seek(n int) {
	r.Position = n
}

func (r *BinaryReader) Skip(n int) {
	r.Position = r.Position + n
}

func (r *BinaryReader) ReadByte() byte {
	b := r.Bytes[r.Position]
	r.Position = r.Position + 1
	return b
}

func (r *BinaryReader) ReadBytes(sz int) []byte {
	bs := r.Bytes[r.Position:r.Position + sz]
	r.Position = r.Position + sz
	return bs
}

func (r *BinaryReader) ReadString(sz int) string {
	bs := r.ReadBytes(sz)
	return toShiftJIS(bs)
}

func (r *BinaryReader) ReadStringUntilNull() string {
	i := int(0)
	for i = r.Position; r.Bytes[i] != 0; i++ {}
	bs := r.ReadBytes(i - r.Position)
	r.Position = r.Position + 1

	return toShiftJIS(bs)
}

func (r *BinaryReader) ReadInt16() int16 {
	bs := r.ReadBytes(2)
	return int16(binary.LittleEndian.Uint16(bs))
}

func (r *BinaryReader) ReadInt32() int {
	bs := r.ReadBytes(4)
	return int(binary.LittleEndian.Uint32(bs))
}

func (r *BinaryReader) ReadInt64() int64 {
	bs := r.ReadBytes(8)
	return int64(binary.LittleEndian.Uint64(bs))
}

func (r *BinaryReader) ReadFloat() float32 {
	bs := r.ReadBytes(4)
	tmp := binary.LittleEndian.Uint32(bs)
	return float32(math.Float32frombits(tmp))
}

func (r *BinaryReader) ReadDouble() float64 {
	bs := r.ReadBytes(8)
	tmp := binary.LittleEndian.Uint64(bs)
	return float64(math.Float64frombits(tmp))
}

func toShiftJIS(data []byte) string {
	reader := bytes.NewReader(data)
	encoder := japanese.ShiftJIS.NewEncoder()
	ret, _ := ioutil.ReadAll(transform.NewReader(reader, encoder))

	return string(ret)
}