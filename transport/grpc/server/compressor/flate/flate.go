package flate

import (
	"compress/flate"
	"google.golang.org/grpc/encoding"
	"io"
	"sync"
)

const Name = "deflate"

type compressor struct {
	poolCompressor   sync.Pool
	poolDecompressor sync.Pool
}

type writer struct {
	*flate.Writer
	pool *sync.Pool
}

type reader struct {
	reader io.Reader
	pool   *sync.Pool
}

func New() encoding.Compressor {
	c := &compressor{}
	c.poolCompressor.New = func() interface{} {
		w, err := flate.NewWriter(io.Discard, flate.DefaultCompression)
		if err != nil {
			panic(err)
		}
		return &writer{Writer: w, pool: &c.poolCompressor}
	}
	return c
}

func (c *compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	z := c.poolCompressor.Get().(*writer)
	z.Writer.Reset(w)
	return z, nil
}

func (c *compressor) Decompress(r io.Reader) (io.Reader, error) {
	z, inPool := c.poolDecompressor.Get().(*reader)
	if !inPool {
		newR := flate.NewReader(r)
		return &reader{reader: newR, pool: &c.poolCompressor}, nil
	}
	if err := z.reader.(flate.Resetter).Reset(r, nil); err != nil {
		c.poolDecompressor.Put(z)
		return nil, err
	}
	return z, nil
}

func (c *compressor) Name() string {
	return Name
}

func (z *writer) Close() error {
	err := z.Writer.Close()
	z.pool.Put(z)
	return err
}

func (z *reader) Read(p []byte) (n int, err error) {
	n, err = z.reader.Read(p)
	if err == io.EOF {
		z.pool.Put(z)
	}
	return n, err
}
