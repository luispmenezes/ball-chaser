package bitreader

import (
	"bytes"
	"io"
	"log"
	"math"
)

const bufferSize = 2048

type Reader struct {
	BinaryReader
	buffer *[]byte
}

func (r *Reader) ReadString(length uint) string {
	if length == 0 {
		length = r.ReadInt(32)
	}
	return string(bytes.Trim(r.ReadBytes(int(length)), "\x00"))
}

func (r *Reader) ReadFloat() float32 {
	return math.Float32frombits(uint32(r.ReadInt(32)))
}

func (r *Reader) ReadInt32() uint {
	return r.ReadInt(32)
}

func (r *Reader) ReadUInt32Max(maxValue uint) uint {
	maxBits := math.Floor(math.Log10(float64(maxValue))/math.Log10(2)) + 1

	value := r.ReadInt(int(maxBits))

	if value > maxValue {
		log.Fatal("failed reading UInt32Max")
	}

	return value
}

func NewLargeBitReader(underlying io.Reader) *Reader {
	b := make([]byte, bufferSize)
	br := new(Reader)
	br.buffer = &b
	br.OpenWithBuffer(underlying, *br.buffer)
	return br
}
