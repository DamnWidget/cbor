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
	"fmt"
	"math/big"
	"testing"
	"time"
)

func TestEncodeNil(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(nil))
	expect(buf.Bytes()[0], absoluteNil, t, "TestEncodeNil")
}

func TestEncodeBool(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(false))
	expect(buf.Bytes()[0], absoluteFalse, t, "TestEncodeBool")
	check(e.Encode(true))
	expect(buf.Bytes()[1], absoluteTrue, t, "TestEncodeBool")
}

func TestEncodePointerToBool(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v bool = false
	check(e.Encode(&v))
	expect(buf.Bytes()[0], absoluteFalse, t, "TestEncodePointerToBool")
	v = true
	check(e.Encode(&v))
	expect(buf.Bytes()[1], absoluteTrue, t, "TestEncodePointerToBool")
}

func TestEncodeUint8(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(uint8(200)))
	expect(buf.Bytes()[0], byte(cborUint8), t, "TestEncodeUint8")
	expect(buf.Bytes()[1], uint8(200), t, "TestEncodeUint8")
	check(e.Encode(uint8(10)))
	expect(buf.Bytes()[2], byte(cborUnsignedInt+Major(0x0a)), t, "TestEncodeUint8")
}

func TestEncodePointerToUint8(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v uint8 = 200
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(cborUint8), t, "TestEncodePointerToUint8")
	expect(buf.Bytes()[1], uint8(200), t, "TestEncodePointerToUint8")
}

func TestEncodeUint16(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(uint16(65000)))
	expect(buf.Bytes()[0], byte(cborUint16), t, "TestEncodeUint16")
	expect(buf.Bytes()[1], byte(0xfd), t, "TestEncodeUint16")
	expect(buf.Bytes()[2], byte(0xe8), t, "TestEncodeUint16")
}

func TestEncodePointerToUint16(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v uint16 = 65000
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(cborUint16), t, "TestEncodePointerToUint16")
	expect(buf.Bytes()[1], byte(0xfd), t, "TestEncodePointerToUint16")
	expect(buf.Bytes()[2], byte(0xe8), t, "TestEncodePointerToUint16")
}

func TestEncodeUint32(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(uint32(650000)))
	expect(buf.Bytes()[0], byte(cborUint32), t, "TestEncodeUint32")
	expect(buf.Bytes()[1], byte(0x00), t, "TestEncodeUint32")
	expect(buf.Bytes()[2], byte(0x09), t, "TestEncodeUint32")
	expect(buf.Bytes()[3], byte(0xeb), t, "TestEncodeUint32")
	expect(buf.Bytes()[4], byte(0x10), t, "TestEncodeUint32")
}

func TestEncodePointerToUint32(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v uint32 = 650000
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(cborUint32), t, "TestEncodePointerToUint32")
	expect(buf.Bytes()[1], byte(0x00), t, "TestEncodePointerToUint32")
	expect(buf.Bytes()[2], byte(0x09), t, "TestEncodePointerToUint32")
	expect(buf.Bytes()[3], byte(0xeb), t, "TestEncodePointerToUint32")
	expect(buf.Bytes()[4], byte(0x10), t, "TestEncodePointerToUint32")
}

func TestEncodeUint64(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(uint64(6500000000)))
	expect(buf.Bytes()[0], byte(cborUint64), t, "TestEncodeUint64")
	expect(buf.Bytes()[1], byte(0x00), t, "TestEncodeUint64")
	expect(buf.Bytes()[2], byte(0x00), t, "TestEncodeUint64")
	expect(buf.Bytes()[3], byte(0x00), t, "TestEncodeUint64")
	expect(buf.Bytes()[4], byte(0x01), t, "TestEncodeUint64")
	expect(buf.Bytes()[5], byte(0x83), t, "TestEncodeUint64")
	expect(buf.Bytes()[6], byte(0x6e), t, "TestEncodeUint64")
	expect(buf.Bytes()[7], byte(0x21), t, "TestEncodeUint64")
	expect(buf.Bytes()[8], byte(0x00), t, "TestEncodeUint64")
}

func TestEncodePointerToUint64(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v uint64 = 6500000000
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(cborUint64), t, "TestEncodePointerToUint64")
	expect(buf.Bytes()[1], byte(0x00), t, "TestEncodePointerToUint64")
	expect(buf.Bytes()[2], byte(0x00), t, "TestEncodePointerToUint64")
	expect(buf.Bytes()[3], byte(0x00), t, "TestEncodePointerToUint64")
	expect(buf.Bytes()[4], byte(0x01), t, "TestEncodePointerToUint64")
	expect(buf.Bytes()[5], byte(0x83), t, "TestEncodePointerToUint64")
	expect(buf.Bytes()[6], byte(0x6e), t, "TestEncodePointerToUint64")
	expect(buf.Bytes()[7], byte(0x21), t, "TestEncodePointerToUint64")
	expect(buf.Bytes()[8], byte(0x00), t, "TestEncodePointerToUint64")
}

func TestEncodeInt8(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(int8(-16)))
	expect(Major(buf.Bytes()[0]>>5), cborNegativeInt, t, "TestEncodeInt")
	expect(buf.Bytes()[0]&0x1f, byte(0x0f), t, "TestEncodeInt8")
	check(e.Encode(int8(-100)))
	expect(Major(buf.Bytes()[1]>>5), cborNegativeInt, t, "TestEncodeInt8")
	expect(buf.Bytes()[1]&0x1f, uint8(cborUint8), t, "TestEncodeInt8")
	expect(buf.Bytes()[2], byte(0x63), t, "TestEncodeInt8")
	check(e.Encode(int8(16)))
	expect(buf.Bytes()[3], byte(cborUnsignedInt+Major(0x10)), t, "TestEncodeInt8")
	check(e.Encode(int8(100)))
	expect(buf.Bytes()[4], byte(cborUint8), t, "TestEncodeInt8")
	expect(buf.Bytes()[5], byte(0x64), t, "TestEncodeInt8")
}

func TestEncodePointerToInt8(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v int8 = -16
	check(e.Encode(&v))
	expect(Major(buf.Bytes()[0]>>5), cborNegativeInt, t, "TestEncodeInt")
	expect(buf.Bytes()[0]&0x1f, byte(0x0f), t, "TestEncodePointerToInt8")
	v = -100
	check(e.Encode(&v))
	expect(Major(buf.Bytes()[1]>>5), cborNegativeInt, t, "TestEncodePointerToInt8")
	expect(buf.Bytes()[1]&0x1f, uint8(cborUint8), t, "TestEncodePointerToInt8")
	expect(buf.Bytes()[2], byte(0x63), t, "TestEncodePointerToInt8")
	v = 16
	check(e.Encode(&v))
	expect(buf.Bytes()[3], byte(cborUnsignedInt+Major(0x10)), t, "TestEncodePointerToInt8")
	v = 100
	check(e.Encode(&v))
	expect(buf.Bytes()[4], byte(cborUint8), t, "TestEncodePointerToInt8")
	expect(buf.Bytes()[5], byte(0x64), t, "TestEncodePointerToInt8")
}

func TestEncodeInt16(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(int16(-32000)))
	expect(Major(buf.Bytes()[0]>>5), cborNegativeInt, t, "TestEncodeInt16")
	expect(buf.Bytes()[0]&0x1f, uint8(cborUint16), t, "TestEncodeInt16")
	expect(buf.Bytes()[1], byte(0x7c), t, "TestEncodeInt16")
	expect(buf.Bytes()[2], byte(0xff), t, "TestEncodeInt16")
	check(e.Encode(int16(32000)))
	expect(buf.Bytes()[3], byte(cborUint16), t, "TestEncodeInt16")
	expect(buf.Bytes()[4], byte(0x7d), t, "TestEncodeInt16")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodeInt16")
}

func TestEncodePointerToInt16(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v int16 = -32000
	check(e.Encode(&v))
	expect(Major(buf.Bytes()[0]>>5), cborNegativeInt, t, "TestEncodePointerToInt16")
	expect(buf.Bytes()[0]&0x1f, uint8(cborUint16), t, "TestEncodePointerToInt16")
	expect(buf.Bytes()[1], byte(0x7c), t, "TestEncodePointerToInt16")
	expect(buf.Bytes()[2], byte(0xff), t, "TestEncodePointerToInt16")
	v = 32000
	check(e.Encode(&v))
	expect(buf.Bytes()[3], byte(cborUint16), t, "TestEncodePointerToInt16")
	expect(buf.Bytes()[4], byte(0x7d), t, "TestEncodePointerToInt16")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodePointerToInt16")
}

func TestEncodeInt32(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(int32(-2147483647)))
	expect(Major(buf.Bytes()[0]>>5), cborNegativeInt, t, "TestEncodeInt32")
	expect(buf.Bytes()[0]&0x1f, uint8(cborUint32), t, "TestEncodeInt32")
	expect(buf.Bytes()[1], byte(0x7f), t, "TestEncodeInt32")
	expect(buf.Bytes()[2], byte(0xff), t, "TestEncodeInt32")
	expect(buf.Bytes()[3], byte(0xff), t, "TestEncodeInt32")
	expect(buf.Bytes()[4], byte(0xfe), t, "TestEncodeInt32")
	check(e.Encode(int32(2147483647)))
	expect(buf.Bytes()[5], byte(cborUint32), t, "TestEncodeInt32")
	expect(buf.Bytes()[6], byte(0x7f), t, "TestEncodeInt32")
	expect(buf.Bytes()[7], byte(0xff), t, "TestEncodeInt32")
	expect(buf.Bytes()[8], byte(0xff), t, "TestEncodeInt32")
	expect(buf.Bytes()[9], byte(0xff), t, "TestEncodeInt32")
}

func TestEncodePointerToInt32(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v int32 = -2147483647
	check(e.Encode(&v))
	expect(Major(buf.Bytes()[0]>>5), cborNegativeInt, t, "TestEncodePointerToInt32")
	expect(buf.Bytes()[0]&0x1f, uint8(cborUint32), t, "TestEncodePointerToInt32")
	expect(buf.Bytes()[1], byte(0x7f), t, "TestEncodePointerToInt32")
	expect(buf.Bytes()[2], byte(0xff), t, "TestEncodePointerToInt32")
	expect(buf.Bytes()[3], byte(0xff), t, "TestEncodePointerToInt32")
	expect(buf.Bytes()[4], byte(0xfe), t, "TestEncodePointerToInt32")
	v = 2147483647
	check(e.Encode(&v))
	expect(buf.Bytes()[5], byte(cborUint32), t, "TestEncodeInt32")
	expect(buf.Bytes()[6], byte(0x7f), t, "TestEncodeInt32")
	expect(buf.Bytes()[7], byte(0xff), t, "TestEncodeInt32")
	expect(buf.Bytes()[8], byte(0xff), t, "TestEncodeInt32")
	expect(buf.Bytes()[9], byte(0xff), t, "TestEncodeInt32")
}

func TestEncodeInt64(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(int64(-184467440737095516)))
	expect(Major(buf.Bytes()[0]>>5), cborNegativeInt, t, "TestEncodeInt64")
	expect(buf.Bytes()[0]&0x1f, uint8(cborUint64), t, "TestEncodeInt64")
	expect(buf.Bytes()[1], byte(0x02), t, "TestEncodeInt64")
	expect(buf.Bytes()[2], byte(0x8f), t, "TestEncodeInt64")
	expect(buf.Bytes()[3], byte(0x5c), t, "TestEncodeInt64")
	expect(buf.Bytes()[4], byte(0x28), t, "TestEncodeInt64")
	expect(buf.Bytes()[5], byte(0xf5), t, "TestEncodeInt64")
	expect(buf.Bytes()[6], byte(0xc2), t, "TestEncodeInt64")
	expect(buf.Bytes()[7], byte(0x8f), t, "TestEncodeInt64")
	expect(buf.Bytes()[8], byte(0x5b), t, "TestEncodeInt64")
	check(e.Encode(int64(184467440737095516)))
	expect(buf.Bytes()[9], byte(cborUint64), t, "TestEncodeInt64")
	expect(buf.Bytes()[10], byte(0x02), t, "TestEncodeInt64")
	expect(buf.Bytes()[11], byte(0x8f), t, "TestEncodeInt64")
	expect(buf.Bytes()[12], byte(0x5c), t, "TestEncodeInt64")
	expect(buf.Bytes()[13], byte(0x28), t, "TestEncodeInt64")
	expect(buf.Bytes()[14], byte(0xf5), t, "TestEncodeInt64")
	expect(buf.Bytes()[15], byte(0xc2), t, "TestEncodeInt64")
	expect(buf.Bytes()[16], byte(0x8f), t, "TestEncodeInt64")
	expect(buf.Bytes()[17], byte(0x5c), t, "TestEncodeInt64")
}

func TestEncodePointerToInt64(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v int64 = -184467440737095516
	check(e.Encode(&v))
	expect(Major(buf.Bytes()[0]>>5), cborNegativeInt, t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[0]&0x1f, uint8(cborUint64), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[1], byte(0x02), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[2], byte(0x8f), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[3], byte(0x5c), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[4], byte(0x28), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[5], byte(0xf5), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[6], byte(0xc2), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[7], byte(0x8f), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[8], byte(0x5b), t, "TestEncodePointerToInt64")
	v = 184467440737095516
	check(e.Encode(&v))
	expect(buf.Bytes()[9], byte(cborUint64), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[10], byte(0x02), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[11], byte(0x8f), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[12], byte(0x5c), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[13], byte(0x28), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[14], byte(0xf5), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[15], byte(0xc2), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[16], byte(0x8f), t, "TestEncodePointerToInt64")
	expect(buf.Bytes()[17], byte(0x5c), t, "TestEncodePointerToInt64")
}

func TestEncodeUint(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(uint(23)))
	expect(buf.Bytes()[0], byte(0x17), t, "TestEncodeUint")
	check(e.Encode(uint(1000000)))
	expect(buf.Bytes()[1], byte(0x1a), t, "TestEncodeUint")
	expect(buf.Bytes()[2], byte(0x00), t, "TestEncodeUint")
	expect(buf.Bytes()[3], byte(0x0f), t, "TestEncodeUint")
	expect(buf.Bytes()[4], byte(0x42), t, "TestEncodeUint")
	expect(buf.Bytes()[5], byte(0x40), t, "TestEncodeUint")
}

func TestEncodePointerToUint(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v uint = 23
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(0x17), t, "TestEncodePointerToUint")
	v = 1000000
	check(e.Encode(&v))
	expect(buf.Bytes()[1], byte(0x1a), t, "TestEncodePointerToUint")
	expect(buf.Bytes()[2], byte(0x00), t, "TestEncodePointerToUint")
	expect(buf.Bytes()[3], byte(0x0f), t, "TestEncodePointerToUint")
	expect(buf.Bytes()[4], byte(0x42), t, "TestEncodePointerToUint")
	expect(buf.Bytes()[5], byte(0x40), t, "TestEncodePointerToUint")
}

func TestEncodeInt(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(int(-10)))
	expect(buf.Bytes()[0], byte(0x29), t, "TestEncodeInt")
	check(e.Encode(int(-1000)))
	expect(buf.Bytes()[1], byte(0x39), t, "TestEncodeInt")
	expect(buf.Bytes()[2], byte(0x03), t, "TestEncodeInt")
	expect(buf.Bytes()[3], byte(0xe7), t, "TestEncodeInt")
	check(e.Encode(int(23)))
	expect(buf.Bytes()[4], byte(0x17), t, "TestEncodeInt")
	check(e.Encode(int(24)))
	expect(buf.Bytes()[5], byte(0x18), t, "TestEncodeInt")
	expect(buf.Bytes()[6], byte(0x18), t, "TestEncodeInt")
}

func TestEncodePointerToInt(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v int = -10
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(0x29), t, "TestEncodePointerToInt")
	v = -1000
	check(e.Encode(&v))
	expect(buf.Bytes()[1], byte(0x39), t, "TestEncodePointerToInt")
	expect(buf.Bytes()[2], byte(0x03), t, "TestEncodePointerToInt")
	expect(buf.Bytes()[3], byte(0xe7), t, "TestEncodePointerToInt")
	v = 23
	check(e.Encode(&v))
	expect(buf.Bytes()[4], byte(0x17), t, "TestEncodeInt")
	v = 24
	check(e.Encode(&v))
	expect(buf.Bytes()[5], byte(0x18), t, "TestEncodeInt")
	expect(buf.Bytes()[6], byte(0x18), t, "TestEncodeInt")
}

func TestEncodeFloat16(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(float16(1.5)))
	expect(buf.Bytes()[0], byte(0xf9), t, "TestEncodeFloat16")
	expect(buf.Bytes()[1], byte(0x3e), t, "TestEncodeFloat16")
	expect(buf.Bytes()[2], byte(0x00), t, "TestEncodeFloat16")
	check(e.Encode(float16(1.0)))
	expect(buf.Bytes()[3], byte(0xf9), t, "TestEncodeFloat16")
	expect(buf.Bytes()[4], byte(0x3c), t, "TestEncodeFloat16")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodeFloat16")
	check(e.Encode(float16(65504.0)))
	expect(buf.Bytes()[6], byte(0xf9), t, "TestEncodeFloat16")
	expect(buf.Bytes()[7], byte(0x7b), t, "TestEncodeFloat16")
	expect(buf.Bytes()[8], byte(0xff), t, "TestEncodeFloat16")
	check(e.Encode(float16(0.00006103515625)))
	expect(buf.Bytes()[9], byte(0xf9), t, "TestEncodeFloat16")
	expect(buf.Bytes()[10], byte(0x04), t, "TestEncodeFloat16")
	expect(buf.Bytes()[11], byte(0x00), t, "TestEncodeFloat16")
}

func TestEncodePointerToFloat16(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v float16 = 1.5
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(0xf9), t, "TestEncodePointerToFloat16")
	expect(buf.Bytes()[1], byte(0x3e), t, "TestEncodePointerToFloat16")
	expect(buf.Bytes()[2], byte(0x00), t, "TestEncodePointerToFloat16")
	v = 1.0
	check(e.Encode(&v))
	expect(buf.Bytes()[3], byte(0xf9), t, "TestEncodePointerToFloat16")
	expect(buf.Bytes()[4], byte(0x3c), t, "TestEncodePointerToFloat16")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodePointerToFloat16")
	v = 65504.0
	check(e.Encode(&v))
	expect(buf.Bytes()[6], byte(0xf9), t, "TestEncodePointerToFloat16")
	expect(buf.Bytes()[7], byte(0x7b), t, "TestEncodePointerToFloat16")
	expect(buf.Bytes()[8], byte(0xff), t, "TestEncodePointerToFloat16")
	v = 0.00006103515625
	check(e.Encode(&v))
	expect(buf.Bytes()[9], byte(0xf9), t, "TestEncodePointerToFloat16")
	expect(buf.Bytes()[10], byte(0x04), t, "TestEncodePointerToFloat16")
	expect(buf.Bytes()[11], byte(0x00), t, "TestEncodePointerToFloat16")
}

func TestEncodeFloat32(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(float32(100000.0)))
	expect(buf.Bytes()[0], byte(0xfa), t, "TestEncodeFloat32")
	expect(buf.Bytes()[1], byte(0x47), t, "TestEncodeFloat32")
	expect(buf.Bytes()[2], byte(0xc3), t, "TestEncodeFloat32")
	expect(buf.Bytes()[3], byte(0x50), t, "TestEncodeFloat32")
	expect(buf.Bytes()[4], byte(0x00), t, "TestEncodeFloat32")
	check(e.Encode(float32(3.4028234663852886e+38)))
	expect(buf.Bytes()[5], byte(0xfa), t, "TestEncodeFloat32")
	expect(buf.Bytes()[6], byte(0x7f), t, "TestEncodeFloat32")
	expect(buf.Bytes()[7], byte(0x7f), t, "TestEncodeFloat32")
	expect(buf.Bytes()[8], byte(0xff), t, "TestEncodeFloat32")
	expect(buf.Bytes()[9], byte(0xff), t, "TestEncodeFloat32")
}

func TestEncodePointerToFloat32(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v float32 = 100000.0
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(0xfa), t, "TestEncodePointerToFloat32")
	expect(buf.Bytes()[1], byte(0x47), t, "TestEncodePointerToFloat32")
	expect(buf.Bytes()[2], byte(0xc3), t, "TestEncodePointerToFloat32")
	expect(buf.Bytes()[3], byte(0x50), t, "TestEncodePointerToFloat32")
	expect(buf.Bytes()[4], byte(0x00), t, "TestEncodePointerToFloat32")
	v = 3.4028234663852886e+38
	check(e.Encode(&v))
	expect(buf.Bytes()[5], byte(0xfa), t, "TestEncodePointerToFloat32")
	expect(buf.Bytes()[6], byte(0x7f), t, "TestEncodePointerToFloat32")
	expect(buf.Bytes()[7], byte(0x7f), t, "TestEncodePointerToFloat32")
	expect(buf.Bytes()[8], byte(0xff), t, "TestEncodePointerToFloat32")
	expect(buf.Bytes()[9], byte(0xff), t, "TestEncodePointerToFloat32")
}

func TestEncodeFloat64(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(float64(1.1)))
	expect(buf.Bytes()[0], byte(0xfb), t, "TestEncodeFloat64")
	expect(buf.Bytes()[1], byte(0x3f), t, "TestEncodeFloat64")
	expect(buf.Bytes()[2], byte(0xf1), t, "TestEncodeFloat64")
	expect(buf.Bytes()[3], byte(0x99), t, "TestEncodeFloat64")
	expect(buf.Bytes()[4], byte(0x99), t, "TestEncodeFloat64")
	expect(buf.Bytes()[5], byte(0x99), t, "TestEncodeFloat64")
	expect(buf.Bytes()[6], byte(0x99), t, "TestEncodeFloat64")
	expect(buf.Bytes()[7], byte(0x99), t, "TestEncodeFloat64")
	expect(buf.Bytes()[8], byte(0x9a), t, "TestEncodeFloat64")
	check(e.Encode(float64(-4.1)))
	expect(buf.Bytes()[9], byte(0xfb), t, "TestEncodeFloat64")
	expect(buf.Bytes()[10], byte(0xc0), t, "TestEncodeFloat64")
	expect(buf.Bytes()[11], byte(0x10), t, "TestEncodeFloat64")
	expect(buf.Bytes()[12], byte(0x66), t, "TestEncodeFloat64")
	expect(buf.Bytes()[13], byte(0x66), t, "TestEncodeFloat64")
	expect(buf.Bytes()[14], byte(0x66), t, "TestEncodeFloat64")
	expect(buf.Bytes()[15], byte(0x66), t, "TestEncodeFloat64")
	expect(buf.Bytes()[16], byte(0x66), t, "TestEncodeFloat64")
	expect(buf.Bytes()[17], byte(0x66), t, "TestEncodeFloat64")
	check(e.Encode(float64(1.0e+300)))
	expect(buf.Bytes()[18], byte(0xfb), t, "TestEncodeFloat64")
	expect(buf.Bytes()[19], byte(0x7e), t, "TestEncodeFloat64")
	expect(buf.Bytes()[20], byte(0x37), t, "TestEncodeFloat64")
	expect(buf.Bytes()[21], byte(0xe4), t, "TestEncodeFloat64")
	expect(buf.Bytes()[22], byte(0x3c), t, "TestEncodeFloat64")
	expect(buf.Bytes()[23], byte(0x88), t, "TestEncodeFloat64")
	expect(buf.Bytes()[24], byte(0x00), t, "TestEncodeFloat64")
	expect(buf.Bytes()[25], byte(0x75), t, "TestEncodeFloat64")
	expect(buf.Bytes()[26], byte(0x9c), t, "TestEncodeFloat64")
}

func TestEncodePointerToFloat64(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v float64 = 1.1
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(0xfb), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[1], byte(0x3f), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[2], byte(0xf1), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[3], byte(0x99), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[4], byte(0x99), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[5], byte(0x99), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[6], byte(0x99), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[7], byte(0x99), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[8], byte(0x9a), t, "TestEncodePointerToFloat64")
	v = -4.1
	check(e.Encode(&v))
	expect(buf.Bytes()[9], byte(0xfb), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[10], byte(0xc0), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[11], byte(0x10), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[12], byte(0x66), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[13], byte(0x66), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[14], byte(0x66), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[15], byte(0x66), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[16], byte(0x66), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[17], byte(0x66), t, "TestEncodePointerToFloat64")
	v = 1.0e+300
	check(e.Encode(&v))
	expect(buf.Bytes()[18], byte(0xfb), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[19], byte(0x7e), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[20], byte(0x37), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[21], byte(0xe4), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[22], byte(0x3c), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[23], byte(0x88), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[24], byte(0x00), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[25], byte(0x75), t, "TestEncodePointerToFloat64")
	expect(buf.Bytes()[26], byte(0x9c), t, "TestEncodePointerToFloat64")
}

func TestEncodeByteString(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode([]byte("byte string")))
	b := []byte{0x4b, 0x62, 0x79, 0x74, 0x65, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}
	for i, c := range b {
		expect(buf.Bytes()[i], c, t, "TestEncodeByteString")
	}
}

func TestEncodePointerToByteString(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v []byte = []byte("byte string")
	check(e.Encode(&v))
	b := []byte{0x4b, 0x62, 0x79, 0x74, 0x65, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}
	for i, c := range b {
		expect(buf.Bytes()[i], c, t, "TestEncodePointerToByteString")
	}
}

func TestEncodePositiveBigNum(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	bn := new(big.Int)
	bn.SetString("18446744073709551616", 10)
	check(e.Encode(*bn))
	expect(buf.Bytes()[0], byte(0xc2), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[1], byte(0x49), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[2], byte(0x01), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[3], byte(0x00), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[4], byte(0x00), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[6], byte(0x00), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[7], byte(0x00), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[8], byte(0x00), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[9], byte(0x00), t, "TestEncodePositiveBigNum")
	expect(buf.Bytes()[10], byte(0x00), t, "TestEncodePositiveBigNum")
}

func TestEncodePointerToPositiveBigNum(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	bn := new(big.Int)
	bn.SetString("18446744073709551616", 10)
	check(e.Encode(bn))
	expect(buf.Bytes()[0], byte(0xc2), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[1], byte(0x49), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[2], byte(0x01), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[3], byte(0x00), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[4], byte(0x00), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[6], byte(0x00), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[7], byte(0x00), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[8], byte(0x00), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[9], byte(0x00), t, "TestEncodePointerToPositiveBigNum")
	expect(buf.Bytes()[10], byte(0x00), t, "TestEncodePointerToPositiveBigNum")
}

func TestEncodeNegativeBigNum(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	bn := new(big.Int)
	bn.SetString("-18446744073709551617", 10)
	check(e.Encode(*bn))
	expect(buf.Bytes()[0], byte(0xc3), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[1], byte(0x49), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[2], byte(0x01), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[3], byte(0x00), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[4], byte(0x00), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[6], byte(0x00), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[7], byte(0x00), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[8], byte(0x00), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[9], byte(0x00), t, "TestEncodeNegativeBigNum")
	expect(buf.Bytes()[10], byte(0x00), t, "TestEncodeNegativeBigNum")
}

func TestEncodePoiinterToNegativeBigNum(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	bn := new(big.Int)
	bn.SetString("-18446744073709551617", 10)
	check(e.Encode(bn))
	expect(buf.Bytes()[0], byte(0xc3), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[1], byte(0x49), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[2], byte(0x01), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[3], byte(0x00), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[4], byte(0x00), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[6], byte(0x00), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[7], byte(0x00), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[8], byte(0x00), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[9], byte(0x00), t, "TestEncodePoiinterToNegativeBigNum")
	expect(buf.Bytes()[10], byte(0x00), t, "TestEncodePoiinterToNegativeBigNum")
}

func TestEncodeEpochDateTime(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode(time.Unix(1363896240, int64(0))))
	expect(buf.Bytes()[0], byte(0xc1), t, "TestEncodeEpochDateTime")
	expect(buf.Bytes()[1], byte(0x1a), t, "TestEncodeEpochDateTime")
	expect(buf.Bytes()[2], byte(0x51), t, "TestEncodeEpochDateTime")
	expect(buf.Bytes()[3], byte(0x4b), t, "TestEncodeEpochDateTime")
	expect(buf.Bytes()[4], byte(0x67), t, "TestEncodeEpochDateTime")
	expect(buf.Bytes()[5], byte(0xb0), t, "TestEncodeEpochDateTime")
}

func TestEncodePointerToEpochDateTime(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v time.Time = time.Unix(1363896240, int64(0))
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(0xc1), t, "TestEncodePointerToEpochDateTime")
	expect(buf.Bytes()[1], byte(0x1a), t, "TestEncodePointerToEpochDateTime")
	expect(buf.Bytes()[2], byte(0x51), t, "TestEncodePointerToEpochDateTime")
	expect(buf.Bytes()[3], byte(0x4b), t, "TestEncodePointerToEpochDateTime")
	expect(buf.Bytes()[4], byte(0x67), t, "TestEncodePointerToEpochDateTime")
	expect(buf.Bytes()[5], byte(0xb0), t, "TestEncodePointerToEpochDateTime")
}

func TestEncodeBigFloat(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v *big.Rat = big.NewRat(3, 2)
	check(e.Encode(*v))
	expect(buf.Bytes()[0], byte(0xc5), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[1], byte(0x82), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[2], byte(0x01), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[3], byte(0xfb), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[4], byte(0x3f), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[5], byte(0xe8), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[6], byte(0x00), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[7], byte(0x00), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[8], byte(0x00), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[9], byte(0x00), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[10], byte(0x00), t, "TestEncodeBigFloat")
	expect(buf.Bytes()[11], byte(0x00), t, "TestEncodeBigFloat")
}

func TestEncodePointerToBigFloat(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v *big.Rat = big.NewRat(3, 2)
	check(e.Encode(v))
	expect(buf.Bytes()[0], byte(0xc5), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[1], byte(0x82), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[2], byte(0x01), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[3], byte(0xfb), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[4], byte(0x3f), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[5], byte(0xe8), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[6], byte(0x00), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[7], byte(0x00), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[8], byte(0x00), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[9], byte(0x00), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[10], byte(0x00), t, "TestEncodePointerToBigFloat")
	expect(buf.Bytes()[11], byte(0x00), t, "TestEncodePointerToBigFloat")
}

func TestEncodeString(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode("水"))
	expect(buf.Bytes()[0], byte(0x63), t, "TestEncodeString")
	expect(buf.Bytes()[1], byte(0xe6), t, "TestEncodeString")
	expect(buf.Bytes()[2], byte(0xb0), t, "TestEncodeString")
	expect(buf.Bytes()[3], byte(0xb4), t, "TestEncodeString")
	check(e.Encode("𐅑"))
	expect(buf.Bytes()[4], byte(0x64), t, "TestEncodeString")
	expect(buf.Bytes()[5], byte(0xf0), t, "TestEncodeString")
	expect(buf.Bytes()[6], byte(0x90), t, "TestEncodeString")
	expect(buf.Bytes()[7], byte(0x85), t, "TestEncodeString")
	expect(buf.Bytes()[8], byte(0x91), t, "TestEncodeString")
}

func TestEncodePointerToString(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v string = "水"
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(0x63), t, "TestEncodePointerToString")
	expect(buf.Bytes()[1], byte(0xe6), t, "TestEncodePointerToString")
	expect(buf.Bytes()[2], byte(0xb0), t, "TestEncodePointerToString")
	expect(buf.Bytes()[3], byte(0xb4), t, "TestEncodePointerToString")
	v = "𐅑"
	check(e.Encode(&v))
	expect(buf.Bytes()[4], byte(0x64), t, "TestEncodePointerToString")
	expect(buf.Bytes()[5], byte(0xf0), t, "TestEncodePointerToString")
	expect(buf.Bytes()[6], byte(0x90), t, "TestEncodePointerToString")
	expect(buf.Bytes()[7], byte(0x85), t, "TestEncodePointerToString")
	expect(buf.Bytes()[8], byte(0x91), t, "TestEncodePointerToString")
}

func TestEncodeNilPointer(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v *int32 = nil
	check(e.Encode(v))
	expect(buf.Bytes()[0], absoluteNil, t, "TestEncodeNil")
	type MyType struct{ a int32 }
	var m *MyType
	check(e.Encode(m))
	expect(buf.Bytes()[1], absoluteNil, t, "TestEncodeNil")
}

func TestEncodeBoolInterface(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v interface{} = false
	check(e.Encode(v))
	expect(buf.Bytes()[0], absoluteFalse, t, "TestEncodeBoolInterface")
	v = true
	check(e.Encode(v))
	expect(buf.Bytes()[1], absoluteTrue, t, "TestEncodeBoolInterface")
}

func TestEncodeUintInterface(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v interface{} = uint8(10)
	check(e.Encode(v))
	expect(buf.Bytes()[0], byte(0x0a), t, "TestEncodeUintInterface")
	v = uint16(1000)
	check(e.Encode(v))
	expect(buf.Bytes()[1], byte(0x19), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[2], byte(0x03), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[3], byte(0xe8), t, "TestEncodeUintInterface")
	v = uint32(1000000)
	check(e.Encode(v))
	expect(buf.Bytes()[4], byte(0x1a), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[6], byte(0x0f), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[7], byte(0x42), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[8], byte(0x40), t, "TestEncodeUintInterface")
	v = uint64(18446744073709551615)
	check(e.Encode(v))
	expect(buf.Bytes()[9], byte(0x1b), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[10], byte(0xff), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[11], byte(0xff), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[12], byte(0xff), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[13], byte(0xff), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[14], byte(0xff), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[15], byte(0xff), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[16], byte(0xff), t, "TestEncodeUintInterface")
	expect(buf.Bytes()[17], byte(0xff), t, "TestEncodeUintInterface")
}

func TestEncodeIntInterface(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v interface{} = int8(-10)
	check(e.Encode(v))
	expect(buf.Bytes()[0], byte(0x29), t, "TestEncodeIntInterface")
	v = int16(-1000)
	check(e.Encode(v))
	expect(buf.Bytes()[1], byte(0x39), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[2], byte(0x03), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[3], byte(0xe7), t, "TestEncodeIntInterface")
	v = int32(-1000000)
	check(e.Encode(v))
	expect(buf.Bytes()[4], byte(0x3a), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[5], byte(0x00), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[6], byte(0x0f), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[7], byte(0x42), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[8], byte(0x3f), t, "TestEncodeIntInterface")
	v = int64(-18446744073709551)
	check(e.Encode(v))
	expect(buf.Bytes()[9], byte(0x3b), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[10], byte(0x00), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[11], byte(0x41), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[12], byte(0x89), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[13], byte(0x37), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[14], byte(0x4b), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[15], byte(0xc6), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[16], byte(0xa7), t, "TestEncodeIntInterface")
	expect(buf.Bytes()[17], byte(0xee), t, "TestEncodeIntInterface")
}

func TestEncodeFloat32Interface(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v interface{} = float32(3.4028234663852886e+38)
	check(e.Encode(v))
	expect(buf.Bytes()[0], byte(0xfa), t, "TestEncodeFloat32Interface")
	expect(buf.Bytes()[1], byte(0x7f), t, "TestEncodeFloat32Interface")
	expect(buf.Bytes()[2], byte(0x7f), t, "TestEncodeFloat32Interface")
	expect(buf.Bytes()[3], byte(0xff), t, "TestEncodeFloat32Interface")
	expect(buf.Bytes()[4], byte(0xff), t, "TestEncodeFloat32Interface")
}

func TestEncodeFloat64Interface(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v interface{} = float64(1.0e+300)
	check(e.Encode(v))
	expect(buf.Bytes()[0], byte(0xfb), t, "TestEncodeFloat64Interface")
	expect(buf.Bytes()[1], byte(0x7e), t, "TestEncodeFloat64Interface")
	expect(buf.Bytes()[2], byte(0x37), t, "TestEncodeFloat64Interface")
	expect(buf.Bytes()[3], byte(0xe4), t, "TestEncodeFloat64Interface")
	expect(buf.Bytes()[4], byte(0x3c), t, "TestEncodeFloat64Interface")
	expect(buf.Bytes()[5], byte(0x88), t, "TestEncodeFloat64Interface")
	expect(buf.Bytes()[6], byte(0x00), t, "TestEncodeFloat64Interface")
	expect(buf.Bytes()[7], byte(0x75), t, "TestEncodeFloat64Interface")
	expect(buf.Bytes()[8], byte(0x9c), t, "TestEncodeFloat64Interface")
}

func TestEncodeSlice(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	check(e.Encode([]int32{-10, 1000, 10, -1000}))
	expect(buf.Bytes()[0], byte(0x84), t, "TestEncodeSlice")
	expect(buf.Bytes()[1], byte(0x29), t, "TestEncodeSlice")
	expect(buf.Bytes()[2], byte(0x19), t, "TestEncodeSlice")
	expect(buf.Bytes()[3], byte(0x03), t, "TestEncodeSlice")
	expect(buf.Bytes()[4], byte(0xe8), t, "TestEncodeSlice")
	expect(buf.Bytes()[5], byte(0x0a), t, "TestEncodeSlice")
	expect(buf.Bytes()[6], byte(0x39), t, "TestEncodeSlice")
	expect(buf.Bytes()[7], byte(0x03), t, "TestEncodeSlice")
	expect(buf.Bytes()[8], byte(0xe7), t, "TestEncodeSlice")
}

func TestEncodePointerToSlice(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	v := []int32{-10, 1000, 10, -1000}
	check(e.Encode(&v))
	expect(buf.Bytes()[0], byte(0x84), t, "TestEncodePointerToSlice")
	expect(buf.Bytes()[1], byte(0x29), t, "TestEncodePointerToSlice")
	expect(buf.Bytes()[2], byte(0x19), t, "TestEncodePointerToSlice")
	expect(buf.Bytes()[3], byte(0x03), t, "TestEncodePointerToSlice")
	expect(buf.Bytes()[4], byte(0xe8), t, "TestEncodePointerToSlice")
	expect(buf.Bytes()[5], byte(0x0a), t, "TestEncodePointerToSlice")
	expect(buf.Bytes()[6], byte(0x39), t, "TestEncodePointerToSlice")
	expect(buf.Bytes()[7], byte(0x03), t, "TestEncodePointerToSlice")
	expect(buf.Bytes()[8], byte(0xe7), t, "TestEncodePointerToSlice")
}

func TestEncodeSliceOfSlicesOfBools(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	v := [][]bool{
		[]bool{true, false, false, false},
		[]bool{},
		[]bool{false, true},
	}
	check(e.Encode(v))
	expect(buf.Bytes()[0], byte(0x83), t, "TestEncodeSliceOfSliceOfBools")
	expect(buf.Bytes()[1], byte(0x84), t, "TestEncodeSliceOfSliceOfBools")
	expect(buf.Bytes()[2], absoluteTrue, t, "TestEncodeSliceOfSliceOfBools")
	expect(buf.Bytes()[3], absoluteFalse, t, "TestEncodeSliceOfSliceOfBools")
	expect(buf.Bytes()[4], absoluteFalse, t, "TestEncodeSliceOfSliceOfBools")
	expect(buf.Bytes()[5], absoluteFalse, t, "TestEncodeSliceOfSliceOfBools")
	expect(buf.Bytes()[6], byte(0x80), t, "TestEncodeSliceOfSliceOfBools")
	expect(buf.Bytes()[7], byte(0x82), t, "TestEncodeSliceOfSliceOfBools")
	expect(buf.Bytes()[8], absoluteFalse, t, "TestEncodeSliceOfSliceOfBools")
	expect(buf.Bytes()[9], absoluteTrue, t, "TestEncodeSliceOfSliceOfBools")
}

func TestEncodeMapOfStringInt(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	v := map[string]int{"One": 1}
	check(e.Encode(v))
	expect(buf.Bytes()[0], byte(0xa1), t, "TestEncodeMapOfStringInt")
	expect(buf.Bytes()[1], byte(0x63), t, "TestEncodeMapOfStringInt")
	expect(buf.Bytes()[2], byte(0x4f), t, "TestEncodeMapOfStringInt")
	expect(buf.Bytes()[3], byte(0x6e), t, "TestEncodeMapOfStringInt")
	expect(buf.Bytes()[4], byte(0x65), t, "TestEncodeMapOfStringInt")
	expect(buf.Bytes()[5], byte(0x01), t, "TestEncodeMapOfStringInt")
}

func TestEncodeStruct(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	type MyType struct {
		polla    int
		Name     string
		Age      uint8
		Address1 []byte
		Address2 []byte
		Married  bool
		Height   float64
	}
	v := MyType{
		Name:     "Test Person",
		Age:      34,
		Address1: []byte("4 CBOR St"),
		Married:  false,
		Height:   1.77,
	}
	check(e.Encode(v))
	fmt.Printf("%#v\n", buf.Bytes())
	expect(buf.Bytes()[0], byte(0xa6), t, "TestEncodeStruct")
	expect(buf.Bytes()[1], byte(0x64), t, "TestEncodeStruct")
	name := []byte{0x4e, 0x61, 0x6d, 0x65}
	for i := 0; i < len(name); i++ {
		expect(buf.Bytes()[i+2], name[i], t, "TestEncodeStruct")
	}
	expect(buf.Bytes()[6], byte(0x6b), t, "TestEncodeStruct")
	test_person := []byte{0x54, 0x65, 0x73, 0x74, 0x20, 0x50, 0x65, 0x72, 0x73, 0x6f, 0x6e}
	for i := 0; i < len(test_person); i++ {
		expect(buf.Bytes()[i+7], test_person[i], t, "TestEncodeStruct")
	}
	// age := []byte{0x41, 0x67, 0x65}
}

// benchmarks
func BenchmarkEncodeBool(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	for i := 0; i < b.N; i++ {
		e.Encode(false)
	}
}

func BenchmarkEncodeUint(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	for i := 0; i < b.N; i++ {
		e.Encode(uint64(6500000000))
	}
}

func BenchmarkEncodeInt(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	v := int32(-650000)
	for i := 0; i < b.N; i++ {
		e.Encode(&v)
	}
}

func BenchmarkEncodeFloat16(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		e.Encode(float16(1.5))
	}
}

func BenchmarkEncodeFloat32(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		e.Encode(float32(100000.0))
	}
}

func BenchmarkEncodeFloat64(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		check(e.Encode(float64(1.0e+300)))
	}
}

func BenchmarkEncodeBytes(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		e.Encode([]byte("byte string"))
	}
}

func BenchmarkEncodePositiveBigNum(b *testing.B) {
	bn := new(big.Int)
	bn.SetString("18446744073709551616", 10)
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		e.Encode(*bn)
	}
}

func BenchmarkEncodeNegativeBigNum(b *testing.B) {
	bn := new(big.Int)
	bn.SetString("-18446744073709551617", 10)
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)

	for i := 0; i < b.N; i++ {
		e.Encode(*bn)
	}
}

func BenchmarkEncodeEpochDateTime(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v time.Time = time.Unix(1363896240, int64(0))

	for i := 0; i < b.N; i++ {
		e.Encode(v)
	}
}

func BenchmarkEncodeBigFloat(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v *big.Rat = big.NewRat(3, 2)

	for i := 0; i < b.N; i++ {
		e.Encode(v)
	}
}

func BenchmarkEncodeString(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	for i := 0; i < b.N; i++ {
		e.Encode("水")
	}
}

func BenchmarkEncodeBoolInterface(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	var v interface{} = false
	for i := 0; i < b.N; i++ {
		e.Encode(v)
	}
}

func becnhmarkInterfaceHelper(b *testing.B, v interface{}) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	for i := 0; i < b.N; i++ {
		e.Encode(v)
	}
}

func BenchmarkEncodeUintInterface(b *testing.B) {
	var v interface{} = uint8(10)
	becnhmarkInterfaceHelper(b, v)
}

func BenchmarkEncodeIntInterface(b *testing.B) {
	var v interface{} = int8(-10)
	becnhmarkInterfaceHelper(b, v)
}

func BenchmarkEncodeFloat32Interface(b *testing.B) {
	var v interface{} = float32(3.4028234663852886e+38)
	becnhmarkInterfaceHelper(b, v)
}

func BenchmarkEncodeFloat64Interface(b *testing.B) {
	var v interface{} = float64(1.0e+300)
	becnhmarkInterfaceHelper(b, v)
}

func BenchmarkEncodeSliceFourInts32(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	v := []int32{-10, 1000, 10, -1000}
	for i := 0; i < b.N; i++ {
		e.Encode(v)
	}
}

func BenchmarkEncodeSliceOfSlicesOfBools(b *testing.B) {
	buf := bytes.NewBuffer(nil)
	e := NewEncoder(buf)
	v := [][]bool{
		[]bool{true, false, false, false},
		[]bool{},
		[]bool{false, true},
	}
	for i := 0; i < b.N; i++ {
		e.Encode(v)
	}
}
