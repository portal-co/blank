package weave

import (
	"encoding/json"
	"io"
	"sync"

	"github.com/portal-co/blank"
	"github.com/portal-co/blank/buf"
)

type WvFrame struct {
	Name string
	Body []byte
}
type WvReader struct {
	blank.BlockReader
	Frames []WvFrame
	Mtx    sync.Mutex
}

func (w *WvReader) In(a string) blank.BlockReader {
	return wvBlocker{w, a}
}

func (w *WvReader) readBlock(x string) ([]byte, error) {
	w.Mtx.Lock()
	defer w.Mtx.Unlock()
	for {
		var f WvFrame
		if len(w.Frames) != 0 {
			f = w.Frames[0]
			w.Frames = w.Frames[1:]
		} else {
			b, err := w.BlockReader.ReadBlock()
			if err != nil {
				return nil, err
			}
			// var f SeqFrame
			err = json.Unmarshal(b, &f)
			if err != nil {
				return nil, err
			}
		}
		if f.Name == x {
			return f.Body, nil
		}
		defer func() {
			w.Frames = append(w.Frames, f)
		}()
	}
}

type wvBlocker struct {
	*WvReader
	Target string
}

func (w wvBlocker) ReadBlock() ([]byte, error) {
	return w.readBlock(w.Target)
}

type WvWriter struct {
	io.WriteCloser
	Target string
}

func (w WvWriter) Write(p []byte) (int, error) {
	b, d := buf.Buffer(w.WriteCloser)
	defer d()
	err := json.NewEncoder(b).Encode(WvFrame{w.Target, p})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

type WvReadWriter interface {
	In(a string) blank.BlockRWC
}

type wvReadWriter struct {
	*WvReader
	io.WriteCloser
}

func (w wvReadWriter) In(a string) blank.BlockRWC {
	return blank.BWWC{BlockReader: w.WvReader.In(a), WriteCloser: WvWriter{w.WriteCloser, a}}
}

func NewRW(x blank.BlockRWC) WvReadWriter {
	return wvReadWriter{&WvReader{BlockReader: x, Frames: []WvFrame{}}, x}
}
