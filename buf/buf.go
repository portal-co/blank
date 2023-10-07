package buf

import (
	"bytes"
	"io"
)

func Buffer(x io.Writer) (io.Writer, func() error) {
	b := &bytes.Buffer{}

	return b, func() error {
		_, err := x.Write(b.Bytes())
		return err
	}
}
