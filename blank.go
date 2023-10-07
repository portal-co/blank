package blank

import (
	"bytes"
	"encoding/gob"
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

type GobBlockinator struct {
	io.Reader
}

func (b GobBlockinator) ReadBlock() ([]byte, error) {
	var x []byte
	err := gob.NewDecoder(b.Reader).Decode(&x)
	return x, err
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

type GobUnblockinator struct {
	io.WriteCloser
}

func (g GobUnblockinator) Write(p []byte) (n int, err error) {
	err = gob.NewEncoder(g.WriteCloser).Encode(p)
	n = len(p)
	return
}

func MakeBlockRWC(x io.ReadWriteCloser, n int) BlockRWC {
	return BWWC{Blockinator{x, n}, x}
}
func MakeGobBlockRWC(x io.ReadWriteCloser) BlockRWC {
	return BWWC{GobBlockinator{x}, GobUnblockinator{x}}
}
func MakeUnblockRWC(x BlockRWC) io.ReadWriteCloser {
	return RWWC{&Unblocker{x, bytes.Buffer{}}, x}
}
