package fs

import (
	"errors"
	"fmt"
	"io"

	"github.com/hack-pad/hackpadfs"
	"github.com/portal-co/blank"
)

type FSReader struct {
	I   uint64
	Fmt string
	hackpadfs.FS
}

// Close implements io.WriteCloser.
func (f *FSReader) Close() error {
	return nil
}

// Write implements io.WriteCloser.
func (f *FSReader) Write(p []byte) (n int, err error) {
	err = hackpadfs.WriteFullFile(f.FS, fmt.Sprintf(f.Fmt, f.I), p, 0777)
	if err == nil {
		f.I++
	}
	n = len(p)
	return
}

// ReadBlock implements blank.BlockReader.
func (f *FSReader) ReadBlock() ([]byte, error) {
	for {
		x, err := hackpadfs.ReadFile(f.FS, fmt.Sprintf(f.Fmt, f.I))
		if err != nil {
			if errors.Is(err, hackpadfs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		return x, nil
	}
}

var _ blank.BlockReader = &FSReader{}
var _ io.WriteCloser = &FSReader{}
