package buffer

import (
	"github.com/yates-z/easel/core/pool"
	"strconv"
)

const _size = 1024

type Buffer []byte

var bufPool = pool.New(func() *Buffer {
	b := make(Buffer, 0, _size)
	return &b
})

func New() *Buffer {
	return bufPool.Get()
}

func (b *Buffer) Free() {
	// To reduce peak allocation, return only smaller buffers to the pool.
	const maxBufferSize = 16 << 10
	if cap(*b) <= maxBufferSize {
		*b = (*b)[:0]
		bufPool.Put(b)
	}
}

func (b *Buffer) Reset() {
	b.SetLen(0)
}

func (b *Buffer) Write(p []byte) (int, error) {
	*b = append(*b, p...)
	return len(p), nil
}

func (b *Buffer) WriteString(s string) (int, error) {
	*b = append(*b, s...)
	return len(s), nil
}

func (b *Buffer) WriteByte(c byte) error {
	*b = append(*b, c)
	return nil
}

// WriteInt appends an integer to the underlying buffer (assuming base 10).
func (b *Buffer) WriteInt(i int64) {
	*b = strconv.AppendInt(*b, i, 10)
}

func (b *Buffer) String() string {
	return string(*b)
}

func (b *Buffer) Len() int {
	return len(*b)
}

func (b *Buffer) SetLen(n int) {
	*b = (*b)[:n]
}

func (b *Buffer) Replace(n int, c byte) {
	*b = append((*b)[:n], c)
	*b = append(*b, (*b)[n+1:]...)
}
