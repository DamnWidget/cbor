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
	"math/big"
	"time"
	"unsafe"
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

// Write a single byte into the io.Writer
// as an encoded CBOR Null value
func (c *Composer) composeNil() error {
	if err := c.write1(absoluteNil); err != nil {
		return fmt.Errorf("while writting nil value: %s", err.Error())
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
	return 0, fmt.Errorf("totally unexpected error, Uint size is unknown %v!", i)
}

// Handle signed integers writing
func (c *Composer) composeInt(i int64) (n int, err error) {
	if i < 0 {
		return c.composeUint(uint64(^i), cborNegativeInt)
	}
	return c.composeUint(uint64(i))
}

// Write a single byte into the io.Writer
// as an encoded CBOR unsigned int of 8 bits
func (c *Composer) composeUint8(i uint8) (int, error) {
	if i < 24 {
		return 0, NewCanonicalModeError(fmt.Sprintf("%d must be send in a single byte 0x%x\n", i, i))
	}
	if err := binary.Write(c.w, binary.BigEndian, i); err != nil {
		return 0, err
	}
	return 1, nil
}

// Write two bytes into the io.Writer
// as an encoded CBOR unsigned int of 16 bits
func (c *Composer) composeUint16(i uint16) (int, error) {
	buf := []byte{byte(i >> 8), byte(i)}
	if _, err := c.write(buf); err != nil {
		return 0, err
	}
	return 2, nil
}

// Write two bytes into the io.Writer
// as an encoded CBOR unsigned int of 32 bits
func (c *Composer) composeUint32(i uint32) (int, error) {
	buf := []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
	if _, err := c.write(buf); err != nil {
		return 0, err
	}
	return 4, nil
}

// Write two bytes into the io.Writer
// as an encoded CBOR unsigned int of 64 bits
func (c *Composer) composeUin64(i uint64) (int, error) {
	buf := []byte{
		byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32),
		byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i),
	}
	if _, err := c.write(buf); err != nil {
		return 0, err
	}
	return 8, nil
}

// Write one byte into the io.Writer
// as an encoded CBOR boolean value
func (c *Composer) composeBoolean(v bool) error {
	b := absoluteFalse
	if v {
		b = absoluteTrue
	}
	if err := c.write1(b); err != nil {
		return fmt.Errorf("while writting boolean %v value: %s", v, err.Error())
	}
	return nil
}

// Write two bytes into the io.Writer
// as an encoded CBOR float16
func (c *Composer) composeFloat16(f float16) error {
	if err := c.write1(absoluteFloat16); err != nil {
		return err
	}
	f16 := uint32toFloat16(*(*uint32)(unsafe.Pointer(&f)))
	buf := []byte{byte(f16 >> 8), byte(f16)}
	if _, err := c.write(buf); err != nil {
		return err
	}
	return nil
}

// Write four bytes into the io.Writer
// as an encoded CBOR float32
func (c *Composer) composeFloat32(f float32) error {
	if err := c.write1(absoluteFloat32); err != nil {
		return err
	}
	i := math.Float32bits(f)
	buf := []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
	if _, err := c.write(buf); err != nil {
		return err
	}
	return nil
}

// Write eight bytes into the io.Writer
// as an encoded CBOR float64
func (c *Composer) composeFloat64(f float64) error {
	if err := c.write1(absoluteFloat64); err != nil {
		return err
	}
	i := math.Float64bits(f)
	buf := []byte{
		byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32),
		byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i),
	}
	if _, err := c.write(buf); err != nil {
		return err
	}
	return nil
}

// Write len(b) + 1 bytes into the
// io.Writer as a sequence of bytes
func (c *Composer) composeBytes(b []byte, major ...Major) (err error) {
	m := cborByteString
	if len(major) != 0 {
		m = major[0]
	}
	l := uint(len(b))
	if l <= 24 {
		err = c.composeInformation(m, byte(l))
	} else {
		info, err := infoHelper(l)
		if err != nil {
			return err
		}
		err = c.composeInformation(m, info)
	}
	if err != nil {
		return err
	}
	if _, err := c.write(b); err != nil {
		return err
	}
	return nil
}

// Write N bytes into the io.Writer
// as an encoded CBOR positive big.Int
func (c *Composer) composeBigUint(n big.Int) error {
	if err := c.write1(absolutePositiveBigNum); err != nil {
		return err
	}
	return c.composeBytes(n.Bytes())
}

// Write N bytes into the io.Writer
// as an encoded CBOR negative big.Int
func (c *Composer) composeBigInt(n big.Int) error {
	if err := c.write1(absoluteNegativeBigNum); err != nil {
		return err
	}
	buf := n.Bytes()
	buf[len(buf)-1]--
	return c.composeBytes(buf)
}

// Write N bytes into the io.Writer
// as an encoded CBOR epoch-based datetime
func (c *Composer) composeEpochDateTime(t time.Time) error {
	if err := c.write1(absoluteEpochDateTime); err != nil {
		return err
	}
	_, err := c.composeInt(t.Unix())
	return err
}

// Write N bytes into the io.Writer
// as an encoded CBOR Big Float
func (c *Composer) composeBigFloat(r big.Rat) error {
	if _, err := c.write([]byte{absoluteBigFloat, byte(0x82)}); err != nil {
		return err
	}
	f, _ := r.Float64()
	m, e := math.Frexp(f)
	if _, err := c.composeInt(int64(e)); err != nil {
		return err
	}
	if err := c.composeFloat64(m); err != nil {
		return err
	}
	return nil
}

// Write len(s) + 1 bytes into the
// io.Writer as an UTF-8 string
func (c *Composer) composeString(s string) error {
	return c.composeBytes([]byte(s), cborTextString)
}

// Write 5 bytes into the
// io.Writer as a CBOR NaN value
func (c *Composer) composeNaN() error {
	if _, err := c.write([]byte{0xfa, 0x7f, 0xc0, 0x00, 0x00}); err != nil {
		return err
	}
	return nil
}

// Write 5 bytes into the
// io.Writer as a CBOR Infinity value
func (c *Composer) composeInfinity(neg ...bool) error {
	data := []byte{0xfa, 0x7f, 0x80, 0x00, 0x00}
	if len(neg) > 0 && neg[0] {
		data = []byte{0xfa, 0xff, 0x80, 0x00, 0x00}
	}
	if _, err := c.write(data); err != nil {
		return err
	}
	return nil
}

// Write 9 bytes into the io.Writer as a
// CBOR NaN value for double precission
func (c *Composer) composeDoublePrecissionNaN() error {
	if _, err := c.write([]byte{0xfb, 0x7f, 0xf8, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}); err != nil {
		return err
	}
	return nil
}

// Write 9 bytes into the io.Writer as a
// CBOR Infinity value for double precission
func (c *Composer) composeDoublePrecissionInfinity(neg ...bool) error {
	data := []byte{0xfb, 0x7f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	if len(neg) > 0 && neg[0] {
		data = []byte{0xfb, 0xff, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	}
	if _, err := c.write(data); err != nil {
		return err
	}
	return nil
}

// Write 3 bytes into the io.Writer
// as a CBOR NaN canonicalized float16 vlaue
func (c *Composer) composeCanonicalNaN() error {
	if _, err := c.write([]byte{0xf9, 0x7e, 0x00}); err != nil {
		return err
	}
	return nil
}

// Write 3 bytes into the io.Writer
// as a CBOR Infinity canonicalized float16 value
func (c *Composer) composeCanonicalInfinity(neg ...bool) error {
	data := []byte{0xf9, 0x7c, 0x00}
	if len(neg) > 0 && neg[0] {
		data = []byte{0xf9, 0xfc, 0x00}
	}
	if _, err := c.write(data); err != nil {
		return err
	}
	return nil
}

// get the info code depending of the size of l
func infoHelper(l uint) (byte, error) {
	var info byte
	if l <= math.MaxUint8 {
		info = cborUint8
	} else if l <= math.MaxUint16 {
		info = cborUint16
	} else if l <= math.MaxUint32 {
		info = cborUint32
	} else if l <= math.MaxUint64 {
		info = cborUint64
	} else {
		return 0, fmt.Errorf("totally unexpected error, []byte buf length size is unkwnown %v!", l)
	}
	return info, nil
}
