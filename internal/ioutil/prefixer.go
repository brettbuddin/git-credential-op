package ioutil

import (
	"bytes"
	"fmt"
	"io"
)

// PrefixWriter returns a [io.Reader] prepends a string to all lines written to the w.
func PrefixWriter(w io.Writer, prefix string) io.Writer {
	return &prefixer{
		w:      w,
		buf:    bytes.NewBuffer(nil),
		prefix: prefix,
	}
}

type prefixer struct {
	w      io.Writer
	buf    *bytes.Buffer
	prefix string
}

func (p *prefixer) Write(b []byte) (int, error) {
	const newline = '\n'

	var written int
	for _, b := range b {
		p.buf.WriteByte(b)

		if b == newline {
			// Write prefix
			n, err := fmt.Fprint(p.w, p.prefix)
			if err != nil {
				return 0, err
			}
			written += n

			// Flush buffer to output
			m, err := p.buf.WriteTo(p.w)
			if err != nil {
				return n, err
			}
			written += int(m)
			p.buf.Truncate(0)
		}
	}

	return written, nil
}

func (p *prefixer) Close() error {
	_, err := p.buf.WriteTo(p.w)
	return err
}
