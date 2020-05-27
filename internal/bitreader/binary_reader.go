package bitreader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	sled     = 8
	sledMask = sled - 1
	sledBits = sled << 3
)

type stack []int

func (s stack) push(v int) stack {
	return append(s, v)
}

func (s stack) pop() (stack, int) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (s stack) top() int {
	return s[len(s)-1]
}

type BinaryReader struct {
	underlying   io.Reader
	buffer       []byte
	offset       int
	bitsInBuffer int
	lazyPosition int
	chunkTargets stack
	endReached   bool
}

func (r *BinaryReader) LazyPosition() int {
	return r.lazyPosition
}

func (r *BinaryReader) ActualPosition() int {
	return r.lazyPosition + r.offset
}

func (r *BinaryReader) Open(underlying io.Reader, bufferSize int) {
	r.OpenWithBuffer(underlying, make([]byte, bufferSize))
}

func (r *BinaryReader) OpenWithBuffer(underlying io.Reader, buffer []byte) {
	if len(buffer)&sledMask != 0 {
		panic(fmt.Sprintf("Buffer must be a multiple of %d", sled))
	}
	if len(buffer) <= sled<<1 {
		panic(fmt.Sprintf("Buffer must be larger than %d bytes", sled<<1))
	}

	r.endReached = false
	r.underlying = underlying
	r.buffer = buffer

	bytes, err := r.underlying.Read(r.buffer)
	if err != nil {
		panic(err)
	}

	r.bitsInBuffer = (bytes << 3) - sledBits
	if bytes < len(r.buffer)-sled {
		r.bitsInBuffer += sledBits
	}
}

func (r *BinaryReader) Close() {
	r.underlying = nil
	r.buffer = nil
	r.offset = 0
	r.bitsInBuffer = 0
	r.chunkTargets = stack{}
	r.lazyPosition = 0
}

func (r *BinaryReader) ReadBit() bool {
	res := (r.buffer[r.offset>>3] & (1 << uint(r.offset&7))) != 0
	r.advance(1)
	return res
}

func (r *BinaryReader) ReadBits(n int) []byte {
	b := make([]byte, (n+7)>>3)
	bitLevel := r.offset&7 != 0
	for i := 0; i < n>>3; i++ {
		b[i] = r.readByteInternal(bitLevel)
	}
	if n&7 != 0 {
		b[n>>3] = r.ReadBitsToByte(n & 7)
	}
	return b
}

func (r *BinaryReader) ReadSingleByte() byte {
	return r.readByteInternal(r.offset&7 != 0)
}

func (r *BinaryReader) readByteInternal(bitLevel bool) byte {
	if !bitLevel {
		res := r.buffer[r.offset>>3]
		r.advance(8)
		return res
	}
	return r.ReadBitsToByte(8)
}

func (r *BinaryReader) ReadBitsToByte(n int) byte {
	return byte(r.ReadInt(n))
}

func (r *BinaryReader) ReadInt(n int) uint {
	val := binary.LittleEndian.Uint64(r.buffer[r.offset>>3&^3:])
	res := uint(val << uint(64-(r.offset&31)-n) >> (64 - uint(n)))
	r.advance(n)
	return res
}

func (r *BinaryReader) ReadBytes(n int) []byte {
	res := make([]byte, 0, n)
	r.ReadBytesInto(&res, n)
	return res
}

func (r *BinaryReader) ReadBytesInto(out *[]byte, n int) {
	bitLevel := r.offset&7 != 0
	if !bitLevel && r.offset+(n<<3) <= r.bitsInBuffer {
		*out = append(*out, r.buffer[r.offset>>3:(r.offset>>3)+n]...)
		r.advance(n << 3)
	} else {
		for i := 0; i < n; i++ {
			*out = append(*out, r.readByteInternal(bitLevel))
		}
	}
}

func (r *BinaryReader) ReadCString(n int) string {
	b := r.ReadBytes(n)
	end := bytes.IndexByte(b, 0)
	if end < 0 {
		end = n
	}
	return string(b[:end])
}

func (r *BinaryReader) ReadSignedInt(n int) int {
	val := binary.LittleEndian.Uint64(r.buffer[r.offset>>3&^3:])
	res := int(int64(val<<uint(64-(r.offset&31)-n)) >> (64 - uint(n)))
	r.advance(n)
	return res
}

func (r *BinaryReader) BeginChunk(n int) {
	r.chunkTargets = r.chunkTargets.push(r.ActualPosition() + n)
}

func (r *BinaryReader) EndChunk() {
	var target int
	r.chunkTargets, target = r.chunkTargets.pop()
	delta := target - r.ActualPosition()
	if delta < 0 {
		panic("Someone read beyond a chunk boundary, what a dick")
	} else if delta > 0 {
		r.Skip(delta)
	}
	if target != r.ActualPosition() {
		panic(fmt.Sprintf("Skipping data failed, expected position %d got %d", target, r.ActualPosition()))
	}
}

func (r *BinaryReader) ChunkFinished() bool {
	return r.chunkTargets.top() <= r.ActualPosition()
}

func (r *BinaryReader) Skip(n int) {
	bufferBits := r.bitsInBuffer - r.offset
	seeker, ok := r.underlying.(io.Seeker)
	if n > bufferBits+sledBits && ok {
		unbufferedSkipBits := n - bufferBits
		globalOffset, err := seeker.Seek(int64((unbufferedSkipBits>>3)-sled), io.SeekCurrent)
		if err != nil {
			panic(err)
		}
		r.lazyPosition = int(globalOffset) << 3

		newBytes, err := r.underlying.Read(r.buffer)
		if err != nil {
			panic(err)
		}
		r.offset = unbufferedSkipBits & sledMask

		r.bitsInBuffer = (newBytes << 3) - sledBits
		if newBytes <= sled {
			r.bitsInBuffer += sledBits
		}
	} else {
		r.advance(n)
	}
}

func (r *BinaryReader) advance(bits int) {
	r.offset += bits
	for r.offset > r.bitsInBuffer {
		r.refillBuffer()
	}
}

func (r *BinaryReader) refillBuffer() {
	copy(r.buffer[0:sled], r.buffer[r.bitsInBuffer>>3:(r.bitsInBuffer>>3)+sled])

	r.offset -= r.bitsInBuffer
	r.lazyPosition += r.bitsInBuffer

	newBytes, err := r.underlying.Read(r.buffer[sled:])
	r.bitsInBuffer = newBytes << 3

	if err != nil {
		if err == io.EOF {
			if r.endReached {
				panic(io.ErrUnexpectedEOF)
			}

			r.bitsInBuffer += sledBits
			r.endReached = true
		} else {
			panic(err)
		}
	}
}
