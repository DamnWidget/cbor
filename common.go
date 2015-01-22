//go:generate stringer -type=Major

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
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

var someInt int

// int size in the running platform
const intSize = unsafe.Sizeof(someInt)

type Major byte

// RFC7049 defines eight "Major Types" that are contained in the
// higher-order 3 bits in the initial byte of data of a 'data item'
const (
	cborUnsignedInt Major = iota // Unsigned integers
	cborNegativeInt              // Negative integers
	cborByteString               // String of bytes
	cborTextString               // String of text UTF-8 encoded
	cborDataArray                // Array of arbitrary data
	cborDataMap                  // Map of arbitrary data
	cborTag                      // Semantic tag
	cborNC                       // Other types that needs no content like "break"
)

// Additional information contained in the 5 low-order bits of
// the header byte have an specific meaning in general and a
// special meaning in case of the Major 7
const (
	cborSmallInt   byte = 0x17
	cborUint8           = 0x18
	cborUint16          = 0x19
	cborUint32          = 0x1a
	cborUint64          = 0x1b
	cborIndefinite      = 0x1f
)

// Additional information values for Major 7
const (
	cborFalse byte = 0x14 + iota
	cborTrue
	cborNil
	cborUndef
	cborSimple
	cborFloat16
	cborFloat32
	cborFloat64
)

// Additional tags on RFC7049
const (
	cborTextDateTime  byte = 0x00
	cborUnixTimestamp      = 0x01
	cborBigNum             = 0x02
	cborBigNegNum          = 0x03
	cborFraction           = 0x04
	cborBigFloat           = 0x05
	cborBase64Url          = 0x15
	cborBase64             = 0x16
	cborBase16             = 0x17
	cborEnc                = 0x18
	cborURI                = 0x20
	cborTextBase64Url      = 0x21
	cborTextBase64         = 0x22
	cborRegexp             = 0x23
	cborMime               = 0x24
	cborSelfDescribe       = 0xd9d9f7
)

// this is being used to break indefinite streams
const cborBreak byte = 0xff

// convenience contants to help to the blind decoder
const (
	absoluteFalse byte = 0xf4 + iota
	absoluteTrue
	absoluteNil
	absoluteUndef
	absoluteSimple
	absoluteFloat16
	absoluteFloat32
	absoluteFloat64
)

const (
	absoluteIndefiniteBytes  byte = 0x5f
	absoluteIndefiniteString      = 0x7f
	absoluteIndefiniteArray       = 0x9f
	absoluteIndefiniteMap         = 0xbf
	absoluteUint                  = 0x00
	absoluteInt                   = 0x20
	absoluteBytes                 = 0x40
	absoluteString                = 0x60
	absoluteArray                 = 0x80
	absoluteMap                   = 0xa0
	absoluteTag                   = 0xc0
	absoluteStringDateTime        = absoluteTag
	absoluteEpochDateTime         = 0xc1
	absolutePositiveBigNum        = 0xc2
	absoluteNegativeBigNum        = 0xc3
	absoluteDecimalFraction       = 0xc4
	absoluteBigFloat              = 0xc5
	absoluteBase64Url             = 0xd5
	absoluteNoContent             = 0xe0
)

// type constants used on blind decode
const (
	byteString reflect.Kind = reflect.UnsafePointer + 1 + iota
	stringDateTime
	epochDateTime
	bigNum
	decimalFraction
	bigFloat
	base64Url
)

type float16 float32

// taken from OGRE 3D rendering engine
func float16toUint32(yy uint16) (d uint32) {
	y := uint32(yy)
	s := (y >> 15) & 0x00000001
	e := (y >> 10) & 0x0000001f
	m := y & 0x000003ff

	if e == 0 {
		if m == 0 { // Plus or minus zero
			return s << 31
		} else { // Denormalized number -- renormalize it
			for (m & 0x00000400) == 0 {
				m <<= 1
				e -= 1
			}
			e += 1
			m &= ^uint32(0x00000400)
		}
	} else if e == 31 {
		if m == 0 { // Inf
			return (s << 31) | 0x7f800000
		} else { // NaN
			return (s << 31) | 0x7f800000 | (m << 13)
		}
	}
	e = e + (127 - 15)
	m = m << 13
	return (s << 31) | (e << 23) | m
}

// convert a mantissa and an exponent into a float32
func decimalFractionToFloat(m, e int64) float32 {
	be := math.Pow10(int(e))
	return float32(float64(m) * be)
}

// convert a float32 to an exponent and a mantissa
func floatToDecimalFraction(f float32) (int64, int64) {
	fs := strconv.FormatFloat(float64(f), 'f', -1, 32)
	l := len(fs) - (strings.Index(fs, ".") + 1)
	return int64(l), int64(f * float32(math.Pow10(l)))
}

// convert a mantissa and an exponent into a *big.Rat
func bigFloatToRatFromInt64(m, e int64) *big.Rat {
	be := math.Pow(2, float64(e))
	f := float32(float64(m) * be)
	r := &big.Rat{}
	r.SetFloat64(float64(f))
	return r
}

// convert a mantissa and an exponent into a *big.Tar from a *big.Int
func bigFloatToRatFromBigInt(m *big.Int, e int64) *big.Rat {
	multiplier := big.NewInt(2 * int64(math.Abs(float64(e))))
	r := &big.Rat{}
	return r.SetFrac(m, multiplier)
}

// convert a *big.Rat to an exponent and a mantissa
func ratToBigFloat(r *big.Rat) (int64, *big.Int) {
	f, _ := r.Float64()
	fs := strconv.FormatFloat(f, 'f', -1, 32)
	l := len(fs) - (strings.Index(fs, ".") + 1)
	return int64(l), big.NewInt(int64(f * float64(math.Pow(2, float64(l)))))
}
