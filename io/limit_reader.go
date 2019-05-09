package ioutil

import (
	"errors"
	"io"
)

var (
	// LimitedReaderEOF is an error returned by LimitedReader to give feedback to the fact that we did not hit an
	// EOF of the Reader but hit the limit imposed by the LimitedReader.
	LimitedReaderEOF = errors.New("LimitedReader EOF: limit reached")
)

// LimitReader returns a Reader that reads from r
// but stops with LimitedReaderEOF after n bytes.
// The underlying implementation is a *LimitedReader.
func LimitReader(r io.Reader, n int64) io.Reader {
	return &LimitedReader{R: r, N: n}
}

// A LimitedReader reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns LimitedReaderEOF when N <= 0 or when the underlying R returns EOF.
// Unlike the std io.LimitedReader this provides feedback
// that the limit was reached through the returned error.
type LimitedReader struct {
	R io.Reader
	N int64 // bytes alloted
}

func (l *LimitedReader) Read(p []byte) (n int, err error) {
	if int64(len(p)) > l.N {
		p = p[0 : l.N+1]
	}
	n, err = l.R.Read(p)
	l.N -= int64(n)
	if err != nil {
		return
	}
	if l.N < 0 {
		return n, LimitedReaderEOF
	}
	return
}
