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

// Composes a 'data item'
type Composer struct {
	header     byte
	w          io.Writer
	indefinite bool
}

// Create a new Composer with the given
// io.Writer and returns back it's address
func NewComposer(w io.Writer) *Composer {
	return &Composer{w: w}
}

func (c *Composer) composeInformation(major Major, info byte) error {
	c.header = (byte(major) << 5) | info
	if _, err := c.w.Write([]byte{c.header}); err != nil {
		return fmt.Errorf("while composing inforamtion byte: %s", err)
	}
	return nil
}

// Write bytes into the io.Writer, returns the
// number of bytes written and an error in case of any
func (c *Composer) write(buf []byte) (n int, err error) {
	if len(buf) == 0 || buf == nil {
		return 0, nil
	}

	n, err = c.w.Write(buf)
	if err != nil {
		return n, err
	}
	if n != len(buf) {
		err = fmt.Errorf("buf was %d bytes length but only %d bytes were written", len(buf), n)
	}
	return n, err
}

// Writes a single byte into the io.Writer
func (c *Composer) write1(b byte) error {
	if _, err := c.write([]byte{b}); err != nil {
		return err
	}
	return nil
}

// Handle unsigned integers writing
func (c *Composer) composeUint(i uint64, infoType ...Major) (n int, err error) {
	var t Major = cborUnsignedInt
	if len(infoType) > 0 {
		t = infoType[0]
	}
	if i < 24 {
		if err := c.composeInformation(t, byte(i)); err != nil {
			return 0, err
		}
		return 1, nil
	}
	if i <= math.MaxUint8 {
		if err := c.composeInformation(t, cborUint8); err != nil {
			return 0, err
		}
		return c.composeUint8(byte(i))
	}
	if i <= math.MaxUint16 {
		if err := c.composeInformation(t, cborUint16); err != nil {
			return 0, err
		}
		return c.composeUint16(uint16(i))
	}
	if i <= math.MaxUint32 {
		if err := c.composeInformation(t, cborUint32); err != nil {
			return 0, err
		}
		return c.composeUint32(uint32(i))
	}
	if i <= math.MaxUint64 {
		if err := c.composeInformation(t, cborUint64); err != nil {
			return 0, err
		}
		return c.composeUin64(i)
	}
	return 0, fmt.Errorf("universe paradox detected, we are gonna die!")
}

// Handle signed integers writing
func (c *Composer) composeInt(i int64) (n int, err error) {
	if i < 0 {
		return c.composeUint(uint64(^i), cborUnsignedInt)
	}
	return c.composeUint(uint64(i))
}

// Write a single byte into the io.Writer
// as an encoded CBOR unsigned int of 8 bits
func (c *Composer) composeUint8(b uint8) (int, error) {
	if b < 24 {
		return 0, fmt.Errorf("%d must be send in a single byte 0x%x\n", b, b)
	}
	if err := c.write1(b); err != nil {
		return 0, err
	}
	return 1, nil
}

// Write two bytes into the io.Writer
// as an encoded CBOR unsigned int of 16 bits
func (c *Composer) composeUint16(i uint16) (int, error) {
	if err := binary.Write(c.w, binary.BigEndian, i); err != nil {
		return 0, err
	}
	return 2, nil
}

// Write two bytes into the io.Writer
// as an encoded CBOR unsigned int of 32 bits
func (c *Composer) composeUint32(i uint32) (int, error) {
	if err := binary.Write(c.w, binary.BigEndian, i); err != nil {
		return 0, err
	}
	return 4, nil
}

// Write two bytes into the io.Writer
// as an encoded CBOR unsigned int of 64 bits
func (c *Composer) composeUin64(i uint64) (int, error) {
	if err := binary.Write(c.w, binary.BigEndian, i); err != nil {
		return 0, err
	}
	return 8, nil
}
