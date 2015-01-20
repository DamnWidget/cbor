// A Golang RFC7049 implementation
// Copyright (C) 2015 Oscar Campos

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cbor

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

// errors returned by the parser
type ParserErr struct {
	Msg string
}

// implements the Error interface
func (pe ParserErr) Error() string {
	return pe.Msg
}

// creates a new ParseErr component and return it back
func NewParseErr(msg string) ParserErr {
	return ParserErr{msg}
}

// Parses a 'data item' and checks it's well-formedness
//
// It defines an internal buffer that is used to check
// the well-formedness of the 'data item' and to store
// data to be processed later
type Parser struct {
	header     byte
	r          io.Reader
	indefinite bool
	buf        []byte
	off        int // the offset inside the buf
}

// Create a new Parser with the given
// io.Reader and resturns back it's address
func NewParser(r io.Reader) *Parser {
	return &Parser{r: r}
}

// Returns true if the header is the
// break opcode, returns false otherwise
func (p *Parser) isBreak() bool {
	return p.header == cborBreak
}

// Returns true if the header is the
// nil opcode, returns false otherwise
func (p *Parser) isNil() bool {
	return p.header == absoluteNil
}

// Returns true if the header is the
// undef opcode, returns false otherwise
func (p *Parser) isUndef() bool {
	return p.header == absoluteUndef
}

// Parses the information part of a 'data item' for any Major type
//
// It also populates the internal buffer if major is not Tag (6) and the
// additional information is not an undefinite (streamed data) type (31)
func (p *Parser) parseInformation() (major Major, info byte, err error) {
	p.header, err = p.scan1()
	if err != nil {
		return 0, 0, err
	}
	major, infotype := p.parseHeader()
	if infotype <= cborSmallInt {
		p.buf = []byte{infotype}
		return major, infotype, nil
	}
	if infotype == cborIndefinite {
		if major < cborByteString || major > cborDataMap && major != cborNC {
			return major, info, NewParseErr(fmt.Sprintf(
				"received additional info 31 (indefinite) for wrong major %d\n", major))
		}
		p.indefinite = true
		return major, infotype, nil
	}
	if major == cborTag {
		return p.parseTagInformation(infotype)
	}
	if (infotype >= 28 && infotype <= 30) || infotype > 31 {
		return major, info, NewParseErr(
			fmt.Sprintf("invalid additional info %d", infotype))
	}
	bytes := 1 << uint(3-(0x1b-uint(infotype)))
	_, p.buf, err = p.scan(bytes)
	return major, infotype, err
}

// TODO: Parses tag information
func (p *Parser) parseTagInformation(infotype byte) (major Major, info byte, err error) {
	return major, info, nil
}

// Parses the header returning back major and additional information
func (p *Parser) parseHeader() (Major, byte) {
	return Major(p.header >> 5), p.header & 0x1f
}

// returns back the lenght of the buffer
func (p *Parser) buflen() uint64 {
	var v uint64
	info := p.header & 0x1f
	if info <= cborSmallInt {
		v = uint64(info)
	} else {
		switch len(p.buf) {
		case 1:
			v = uint64(p.parseUint8())
		case 2:
			v = uint64(p.parseUint16())
		case 4:
			v = uint64(p.parseUint32())
		case 8:
			v = uint64(p.parseUint64())
		}
	}
	return v
}

// Read N bytes from the internal buffer
// If the buffer doesn't contains that many
// bytes, the function just panic (as it had
// to be checked already during the scanning)
func (p *Parser) read(n int) []byte {
	a := (len(p.buf) - p.off)
	if n > a {
		panic(fmt.Sprintf(
			"can't read %d bytes from buffer as only %d are available\n", n, a))
	}
	oldOff := p.off
	p.off += n
	return p.buf[oldOff:p.off]
}

// Reads N bytes from the parser io.Reader
//
// Returns the number of bytes readed or zero when errors and a bytes slice
// containing the data that has been readed from the io.Reader
func (p *Parser) scan(n int) (numbytes int, data []byte, err error) {
	if n <= 0 {
		return
	}
	data = make([]byte, n)
	if numbytes, err = p.r.Read(data); err != nil {
		return 0, nil, err
	}
	if numbytes < n {
		return 0, nil, NewParseErr(fmt.Sprintf(
			"can't scan %d bytes from buffer as only %d are available\n", n, numbytes))
	}
	p.off = 0
	return numbytes, data, nil
}

// Reads a single byte from the parser io.Reader
func (p *Parser) scan1() (byte, error) {
	_, tmpdata, err := p.scan(1)
	if err != nil {
		return 0, err
	}
	return tmpdata[0], nil
}

// Read a single byte from the internal
// buffer and returns it back as an uint8
func (p *Parser) parseUint8() uint8 {
	return uint8(p.read(1)[0])
}

// Read two bytes from the internal
// buffer and returns it back as uint16
func (p *Parser) parseUint16() uint16 {
	return binary.BigEndian.Uint16(p.read(2))
}

// Read four bytes from the internal
// buffer and returns it back as uint32
func (p *Parser) parseUint32() uint32 {
	return binary.BigEndian.Uint32(p.read(4))
}

// Read eight bytes from the internal
// buffer and returns it back as uint64
func (p *Parser) parseUint64() uint64 {
	return binary.BigEndian.Uint64(p.read(8))
}

// Read two bytes from the internal
// buffer and returns it back as float16
func (p *Parser) parseFloat16() float16 {
	return float16(
		math.Float32frombits(float16toUint32(binary.BigEndian.Uint16(p.read(2)))))
}

// Read four bytes from the internal
// buffer and returns it back as float32
func (p *Parser) parseFloat32() float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(p.read(4)))
}

// Read eight bytes from the internal
// buffer and returns it back as float64
func (p *Parser) parseFloat64() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(p.read(8)))
}

// Read a boolean value from the internal buffer
func (p *Parser) parseBool() bool {
	v := true
	if uint8(p.buflen()) == cborFalse {
		v = false
	}
	return v
}
