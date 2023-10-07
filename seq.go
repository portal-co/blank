package blank

import (
	"bytes"
	"encoding/gob"
	"io"

	"github.com/portal-co/blank/buf"
)

type SeqFrame struct {
	N byte
	X []byte
}
type UnSeq struct {
	BlockReader
	F map[byte][]byte
	C byte
}

func (u *UnSeq) ReadBlock() ([]byte, error) {
	c := u.C
	u.C += 1
	for _, ok := u.F[c]; !ok; {
		b, err := u.BlockReader.ReadBlock()
		if err != nil {
			return nil, err
		}
		var f SeqFrame
		err = gob.NewDecoder(bytes.NewBuffer(b)).Decode(&f)
		if err != nil {
			return nil, err
		}
		u.F[f.N] = f.X
	}
	f := u.F[c]
	delete(u.F, c)
	return f, nil
}

type Seq struct {
	C byte
	io.WriteCloser
}

func (s *Seq) Write(x []byte) (int, error) {
	b, d := buf.Buffer(s.WriteCloser)
	defer d()
	err := gob.NewEncoder(b).Encode(SeqFrame{s.C, x})
	if err != nil {
		return 0, err
	}
	s.C += 1
	return len(x), nil
}
func MakeSeq(x BlockRWC) BlockRWC {
	return BWWC{&UnSeq{x, map[byte][]byte{}, 0}, &Seq{0, x}}
}
