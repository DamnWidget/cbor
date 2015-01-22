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
	"bytes"
	"io"
	"testing"
)

func TestScan(t *testing.T) {
	buf := []byte{'t', 'e', 's', 't'}
	r := bytes.NewBuffer(buf)

	p := NewParser(r)
	n, value, err := p.scan(1)
	check(err)
	expect(1, n, t)
	expect(byte('t'), value[0], t, "TestScan")

	r = bytes.NewBuffer(buf)
	p = NewParser(r)
	n, value, err = p.scan(3)
	check(err)
	expect(byte('t'), value[0], t, "TestScan")
	expect(byte('e'), value[1], t, "TestScan")
	expect(byte('s'), value[2], t, "TestScan")

	r = bytes.NewBuffer(buf)
	p = NewParser(r)
	n, value, err = p.scan(4)
	check(err)
	expect(string(buf), string(value), t, "TestScan")

	r = bytes.NewBuffer(buf)
	p = NewParser(r)
	n, _, err = p.scan(5)
	expect(err, NewParseErr("can't scan 5 bytes from buffer as only 4 are available\n"), t, "TestScan")
}

func TestScan1(t *testing.T) {
	buf := []byte{'t', 'e', 's', 't'}
	r := bytes.NewBuffer(buf)
	p := NewParser(r)

	value, err := p.scan1()
	check(err)
	expect(byte('t'), value, t, "TestScan1")
	value, err = p.scan1()
	check(err)
	expect(byte('e'), value, t, "TestScan1")
	value, err = p.scan1()
	check(err)
	expect(byte('s'), value, t, "TestScan1")
	value, err = p.scan1()
	check(err)
	expect(byte('t'), value, t, "TestScan1")

	// read beyond limits returns io.EOF
	_, err = p.scan1()
	expect(err, io.EOF, t, "TestScan1")
}

func TestParseUint8(t *testing.T) {
	buf := []byte{0x6f}
	p := new(Parser)
	p.buf = buf
	expect(uint8(111), p.parseUint8(), t, "TestParseUint8")
}

func TestParseUint16(t *testing.T) {
	buf := []byte{0x45, 0xab}
	p := new(Parser)
	p.buf = buf
	expect(uint16(17835), p.parseUint16(), t, "TestParseUint16")
}

func TestParseUint32(t *testing.T) {
	buf := []byte{0x8c, 0x7e, 0xe1, 0x38}
	p := new(Parser)
	p.buf = buf
	expect(uint32(2357125432), p.parseUint32(), t, "TestParseUint32")
}

func TestParseUint64(t *testing.T) {
	buf := []byte{0xb6, 0x70, 0x0f, 0xa8, 0xcd, 0x99, 0x87, 0x8d}
	p := new(Parser)
	p.buf = buf
	expect(uint64(13146024529972791181), p.parseUint64(), t, "TestParseUint64")
}

func TestParseBool(t *testing.T) {
	p := new(Parser)
	p.header = byte(0xf4)
	expect(false, p.parseBool(), t, "TestParseBool")
}

func TestParseInformation(t *testing.T) {
	buf := []byte{0x19, 0x10, 0x23}
	r := bytes.NewBuffer(buf)
	p := NewParser(r)
	_, info, err := p.parseInformation()
	check(err)
	expect(byte(cborUint16), info, t, "TestParseInformation")
	d := p.read(2)
	for i, b := range buf[1:] {
		expect(b, d[i], t, "TestParseInformation")
	}

	buf = []byte{0x1b, 0xb6, 0x70, 0x0f, 0xa8, 0xcd, 0x99, 0x87, 0x8d}
	r = bytes.NewBuffer(buf)
	p = NewParser(r)
	_, info, err = p.parseInformation()
	check(err)
	expect(byte(cborUint64), info, t, "TestParseInformation")
	d = p.read(8)
	for i, b := range buf[1:] {
		expect(b, d[i], t, "TestParseInformation")
	}

	buf = []byte{0x5f}
	r = bytes.NewBuffer(buf)
	p = NewParser(r)
	_, info, err = p.parseInformation()
	check(err)
	expect(byte(cborIndefinite), info, t, "TestParseInformation")

	buf = []byte{0x3f}
	r = bytes.NewBuffer(buf)
	p = NewParser(r)
	_, info, err = p.parseInformation()
	expect(err, NewParseErr("received additional info 31 (indefinite) for wrong major 1\n"), t, "TestParseInformation")
}
