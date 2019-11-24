package parser

import (
	"fmt"
	"io"
)

type reader struct {
	r      io.Reader // underlying reader
	buf    []byte    // data buffer
	offset int       // head offset
}

// take reads n bytes from buf and moves the offset accordingly
// it returns back the number of bytes read and the buffer, if
// the number of bytes to read exceeds the upper boundary limit
// of the buffer just the pending bytes are read, this is
// reflected in the number of bytes readed return parameter.
// If take is run in an already read header, 0 and nil is returned
func (r *reader) take(n int) (int, []byte) {

	// if n is zero just return empty
	if n == 0 {
		return 0, nil
	}

	// get number of bytes left on buffer
	bytesLeft := len(r.buf) - (r.offset)

	// if there are no more bytes left to read return empty
	if bytesLeft == 0 {
		return 0, nil
	}

	// normalize n
	if n > bytesLeft {
		n = bytesLeft
	}

	// read the buffer and advance the offset
	offset := r.offset
	r.offset += n
	return n, r.buf[offset:r.offset]
}

// NewReader creates a new reader value with the given reader and returns it back
func NewReader(r io.Reader) reader {

	// create new reader value with the given buffer
	newReader := reader{r: r}
	return newReader
}

// scan reads n bytes from the r io.Reader reader
// returns the number of bytes read into the giuven buf or zero on errors
func scan(r io.Reader, buf []byte) (int, error) {

	n := len(buf)
	if n == 0 {
		return 0, nil
	}

	numbytes, err := r.Read(buf)
	if err != nil {
		return 0, fmt.Errorf("could not read from buffer: %w", err)
	}

	return numbytes, nil
}
