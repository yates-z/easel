package zlib

import (
	"compress/zlib"
	"google.golang.org/grpc/encoding"
	"io"
	"sync"
)

const Name = "zlib"

func New() encoding.Compressor {
	return &compressor{}
}

type writer struct {
	*zlib.Writer
	pool *sync.Pool
}

func (c *compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	z := c.poolCompressor.Get().(*writer)
	z.Writer.Reset(w)
	return z, nil
}

func (z *writer) Close() error {
	defer z.pool.Put(z)
	return z.Writer.Close()
}

type reader struct {
	r    io.Reader
	pool *sync.Pool
}

func (c *compressor) Decompress(r io.Reader) (io.Reader, error) {
	z, inPool := c.poolDecompressor.Get().(*reader)
	if !inPool {
		newZ, err := zlib.NewReader(r)
		if err != nil {
			return nil, err
		}
		return &reader{r: newZ, pool: &c.poolDecompressor}, nil
	}
	if err := z.r.(zlib.Resetter).Reset(r, nil); err != nil {
		c.poolDecompressor.Put(z)
		return nil, err
	}
	return z, nil
}

func (z *reader) Read(p []byte) (n int, err error) {
	n, err = z.r.Read(p)
	if err == io.EOF {
		z.pool.Put(z)
	}
	return n, err
}

func (c *compressor) Name() string {
	return Name
}

type compressor struct {
	poolCompressor   sync.Pool
	poolDecompressor sync.Pool
}
