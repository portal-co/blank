package blank

import (
	"bytes"
	"io"
)

type BlockReader interface {
	ReadBlock() ([]byte, error)
}
type BlockRWC interface {
	BlockReader
	io.WriteCloser
}
type RWWC struct {
	io.Reader
	io.WriteCloser
}
type BWWC struct {
	BlockReader
	io.WriteCloser
}
type Blockinator struct {
	io.Reader
	N int
}

func (b Blockinator) ReadBlock() ([]byte, error) {
	x := make([]byte, b.N)
	n, err := b.Read(x)
	if err != nil {
		return nil, err
	}
	return x[:n], nil
}

type Unblocker struct {
	BlockReader
	Buf bytes.Buffer
}

func (b *Unblocker) Read(x []byte) (int, error) {
	for len(b.Buf.Bytes()) < len(x) {
		r, err := b.ReadBlock()
		if err != nil {
			return 0, err
		}
		b.Buf.Write(r)
	}
	return b.Buf.Read(x)
}

func MakeBlockRWC(x io.ReadWriteCloser, n int) BlockRWC {
	return BWWC{Blockinator{x, n}, x}
}
func MakeUnblockRWC(x BlockRWC) io.ReadWriteCloser {
	return RWWC{&Unblocker{x, bytes.Buffer{}}, x}
}
