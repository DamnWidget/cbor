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
	"io/ioutil"
	"log"

	"math/big"

	"os"
	"reflect"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard)
	os.Exit(m.Run())
}

func TestDecodeUint8(t *testing.T) {
	buf := []byte{0x18, 0x6f}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a uint8
	check(d.Decode(&a))
	expect(uint8(111), a, t)

	buf = []byte{0x38, 0x6f}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)

	buf = []byte{0x19, 0x6f, 0x00}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeInt8(t *testing.T) {
	buf := []byte{0x38, 0x6f}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a int8
	check(d.Decode(&a))
	expect(int8(-112), a, t)

	buf = []byte{0x18, 0x6f}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)

	buf = []byte{0x39, 0x6f, 0x00}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeUint16(t *testing.T) {
	buf := []byte{0x19, 0x45, 0xab}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a uint16
	check(d.Decode(&a))
	expect(uint16(17835), a, t)

	buf = []byte{0x39, 0x45, 0xab}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)

	buf = []byte{0x1a, 0x45, 0xab, 0x00, 0x00}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeInt16(t *testing.T) {
	buf := []byte{0x39, 0x45, 0xab}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a int16
	check(d.Decode(&a))
	expect(int16(-17836), a, t)

	buf = []byte{0x19, 0x45, 0xab}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)

	buf = []byte{0x38, 0x45, 0xab}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeUint32(t *testing.T) {
	buf := []byte{0x1a, 0x45, 0xab, 0x23, 0x00}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a uint32
	check(d.Decode(&a))
	expect(uint32(1168843520), a, t)

	buf = []byte{0x3a, 0x45, 0xab, 0x23, 0x00}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)

	buf = []byte{0x19, 0x45, 0xab, 0x23, 0x00}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeInt32(t *testing.T) {
	buf := []byte{0x3a, 0x45, 0xab, 0x23, 0x00}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a int32
	check(d.Decode(&a))
	expect(int32(-1168843521), a, t)

	buf = []byte{0x1a, 0x45, 0xab, 0x23, 0x00}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)

	buf = []byte{0x39, 0x45, 0xab, 0x23, 0x00}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeUint64(t *testing.T) {
	buf := []byte{0x1b, 0x45, 0xab, 0x23, 0x00, 0x10, 0x11, 0x12, 0x13}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a uint64
	check(d.Decode(&a))
	expect(uint64(5020144692811076115), a, t)

	buf = []byte{0x3b, 0x45, 0xab, 0x23, 0x00, 0x10, 0x11, 0x12, 0x13}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)

	buf = []byte{0x19, 0x45, 0xab, 0x23, 0x00, 0x10, 0x11, 0x12, 0x13}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeInt64(t *testing.T) {
	buf := []byte{0x3b, 0x45, 0xab, 0x23, 0x00, 0x10, 0x11, 0x12, 0x13}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a int64
	check(d.Decode(&a))
	expect(int64(-5020144692811076116), a, t)

	buf = []byte{0x1b, 0x45, 0xab, 0x23, 0x00, 0x10, 0x11, 0x12, 0x13}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)

	buf = []byte{0x39, 0x45, 0xab, 0x23, 0x00, 0x10, 0x11, 0x12, 0x13}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeFloat16(t *testing.T) {
	buf := []byte{0xf9, 0x3f, 0xe0}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a float16
	check(d.Decode(&a))
	expect(float16(1.96875), a, t)
}

func TestDecodeFloat32(t *testing.T) {
	buf := []byte{0xfa, 0x3f, 0x66, 0x66, 0x66}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a float32
	check(d.Decode(&a))
	expect(float32(0.9), a, t)

	buf = []byte{0xfb, 0x40, 0x63, 0x8e, 0xa6, 0xb7, 0x23, 0xee, 0x1c}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeFloat64(t *testing.T) {
	buf := []byte{0xfb, 0x40, 0x63, 0x8e, 0xa6, 0xb7, 0x23, 0xee, 0x1c}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a float64
	check(d.Decode(&a))
	expect(float64(156.457851), a, t)

	buf = []byte{0xfa, 0x40, 0x63, 0x8e, 0xa6, 0xb7, 0x23, 0xee, 0x1c}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(&a) != nil, true, t)
}

func TestDecodeBytes(t *testing.T) {
	buf := []byte{0x4c, 0x62, 0x79, 0x74, 0x65, 0x73, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a []byte
	check(d.Decode(&a))
	expect("bytes string", string(a), t)
}

func TestDecodeString(t *testing.T) {
	buf := []byte{0x67, 0x65, 0x73, 0x70, 0x61, 0xc3, 0xb1, 0x61}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a string
	check(d.Decode(&a))
	expect("españa", a, t)
}

func TestDecodeBool(t *testing.T) {
	buf := []byte{0xf4}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a bool
	check(d.Decode(&a))
	expect(false, a, t)

	buf = []byte{0xf5}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	check(d.Decode(&a))
	expect(true, a, t)
}

func TestDecodeIndefiniteBytes(t *testing.T) {
	buf := []byte{0x5f, 0x4c, 0x62, 0x79, 0x74, 0x65, 0x73, 0x20, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x43, 0x20, 0x78, 0x44, 0xff}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a []byte
	check(d.Decode(&a))
	expect("bytes string xD", string(a), t)
}

func TestDecodeIndefiniteString(t *testing.T) {
	buf := []byte{0x7f, 0x63, 0xe4, 0xb8, 0x96, 0x63, 0xe7, 0x95, 0x8c, 0xff}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a string
	check(d.Decode(&a))
	expect("世界", a, t)
}

func TestDecodeKInt(t *testing.T) {
	buf := []byte{0x3a, 0x45, 0xab, 0x23, 0x00}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a int32
	check(d.Decode(reflect.ValueOf(&a)))
	expect(int32(-1168843521), a, t)

	buf = []byte{0x39, 0x45, 0xab, 0x23, 0x00}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(reflect.ValueOf(&a)) != nil, true, t)
}

func TestDecodeKUint(t *testing.T) {
	buf := []byte{0x1a, 0x45, 0xab, 0x23, 0x00}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a uint32
	check(d.Decode(reflect.ValueOf(&a)))
	expect(uint32(1168843520), a, t)

	buf = []byte{0x19, 0x45, 0xab, 0x23, 0x00}
	r = bytes.NewReader(buf)
	d = NewDecoder(r)
	expect(d.Decode(reflect.ValueOf(&a)) != nil, true, t)
}

func TestDecodeUnsignedIntsArray(t *testing.T) {
	buf := []byte{0x84, 0x04, 0x09, 0x19, 0x04, 0x00, 0x10}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a []uint
	check(d.Decode(&a))
	expected := []uint{4, 9, 1024, 16}
	for i, e := range a {
		expect(expected[i], e, t)
	}
}

func TestDecodeUnsignedIntsIndefiniteArray(t *testing.T) {
	buf := []byte{0x9f, 0x04, 0x09, 0x19, 0x04, 0x00, 0x10, 0xff}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a []uint
	check(d.Decode(&a))
	expected := []uint{4, 9, 1024, 16}
	for i, e := range a {
		expect(expected[i], e, t)
	}
}

func TestDecodeArrayOfUin32(t *testing.T) {
	buf := []byte{0x81, 0x1a, 0x45, 0xab, 0x23, 0x00}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a [1]uint32
	check(d.Decode(&a))
	expected := [1]uint32{1168843520}
	expect(expected[0], a[0], t)
	expect(expected, a, t)
}

func TestDecodeInterface(t *testing.T) {
	buf := []byte{0x85, 0x04, 0x09, 0x19, 0x04, 0x00, 0x10, 0x83, 0x01, 0x02, 0x67, 0x65, 0x73, 0x70, 0x61, 0xc3, 0xb1, 0x61}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	av := *a.(*[]interface{})
	expected := []interface{}{uint8(4), uint8(9), uint16(1024), uint8(16)}
	for i := 0; i < 4; i++ {
		expect(expected[i], av[i], t)
	}
	aiv := *av[4].(*[]interface{})
	expect(aiv[0], uint8(1), t)
	expect(aiv[1], uint8(2), t)
	expect(aiv[2], "españa", t)
}

func TestDecodeMap(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a map[string]interface{}
	check(d.Decode(&a))
	v1, ok := a["Fun"]
	expect(ok, true, t)
	expect(v1, true, t)
	v2, ok := a["Amt"]
	expect(ok, true, t)
	expect(v2, int8(-2), t)
}

func TestDecodeStrictMap(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x46, 0x75, 0x6e, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r, func(dec *Decoder) { dec.strict = true })
	var a map[string]interface{}
	err := d.Decode(&a)
	expect(err != nil, true, t)

	buf = []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x46, 0x75, 0x6e, 0x21}
	r = bytes.NewReader(buf)
	d = NewDecoder(r, func(dec *Decoder) { dec.strict = false })
	var a2 map[string]interface{}
	err = d.Decode(&a2)
	expect(len(a2), 1, t)
	_, ok := a["Fun"]
	expect(ok, true, t)
}

func TestDecodeIndefiniteMap(t *testing.T) {
	buf := []byte{0xbf, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21, 0xff}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a map[string]interface{}
	check(d.Decode(&a))
	v1, ok := a["Fun"]
	expect(ok, true, t)
	expect(v1, true, t)
	v2, ok := a["Amt"]
	expect(ok, true, t)
	expect(v2, int8(-2), t)
}

func TestDecodeInterfaceKeyInterfaceValueMap(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a map[interface{}]interface{}
	check(d.Decode(&a))
	v1, ok := a["Fun"]
	expect(ok, true, t)
	expect(v1, true, t)
	v2, ok := a["Amt"]
	expect(ok, true, t)
	expect(v2, int8(-2), t)
}

func TestDecodeMapIntoInterface(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	av := *a.(*map[interface{}]interface{})
	v1, ok := av["Fun"]
	expect(ok, true, t)
	expect(v1, true, t)
	v2, ok := av["Amt"]
	expect(ok, true, t)
	expect(v2, int8(-2), t)
}

func TestDecodeMapIntoStruct(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Amt, int8(-2), t)
}

func TestDecodeIndefiniteMapIntoStruct(t *testing.T) {
	buf := []byte{0xbf, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21, 0xff}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Amt, int8(-2), t)
}

func TestDecodeDuplicateKeysMapIntoStruct(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x46, 0x75, 0x6e, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Amt, int8(0), t) // is not set
}

func TestDecodeDuplicateKeysMapIntoStructStrictMode(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x46, 0x75, 0x6e, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r, func(dec *Decoder) { dec.strict = true })
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	err := d.Decode(&a)
	expect(err != nil, true, t)
	expect(fmt.Sprint(err), "strict-mode: duplicated key Fun in map", t)
}

func TestDecodeMapIntoStructNonStringKeys(t *testing.T) {
	buf := []byte{0xa2, 0x10, 0xf5, 0x11, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	err := d.Decode(&a)
	expect(err != nil, true, t)
	expect(fmt.Sprint(err), "map keys must be string, cborUnsignedInt received", t)
}

func TestDecodeMapNonFieldIntoStruct(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
}

func TestDecodeMapNonFieldIntoStructWithValidTag(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun   bool
		Other int8 `cbor:"Amt"`
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Other, int8(-2), t)
}

func TestDecodeIndefiniteMapNonFieldIntoStruct(t *testing.T) {
	buf := []byte{0xbf, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21, 0xff}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
}

func TestDecodeMapNonFieldIntoStructStrictMode(t *testing.T) {
	buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r, func(dec *Decoder) { dec.strict = true })
	type MyType struct {
		Fun  bool
		None int8
	}
	var a MyType
	err := d.Decode(&a)
	expect(err != nil, true, t)
	expect(fmt.Sprint(err), "strict-mode: key Amt doesn't match with any field", t)
}

func TestDecodeMapOutboundsIntoStruct(t *testing.T) {
	buf := []byte{0xa3, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21, 0x64, 0x46, 0x61, 0x69, 0x6c, 0x04}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Amt, int8(-2), t)
	expect(r.Len(), 0, t)
}

func TestDecodeMapOutboundsIntoStructStrictMode(t *testing.T) {
	buf := []byte{0xa3, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21, 0x64, 0x46, 0x61, 0x69, 0x6c, 0x04}
	r := bytes.NewReader(buf)
	d := NewDecoder(r, func(dec *Decoder) { dec.strict = true })
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	err := d.Decode(&a)
	expect(err != nil, true, t)
	expect(fmt.Sprint(err), "strict-mode: destination struct fields num 2 doesn't match map length 3", t)
}

func TestDecodeArrayIntoStruct(t *testing.T) {
	buf := []byte{0x84, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Amt, int8(-2), t)
	expect(r.Len(), 0, t)
}

func TestDecodeIndefiniteArrayIntoStruct(t *testing.T) {
	buf := []byte{0x8f, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21, 0xff}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Amt, int8(-2), t)
}

func TestDecodeDuplicateKeysArrayIntoStruct(t *testing.T) {
	buf := []byte{0x84, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x46, 0x75, 0x6e, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Amt, int8(0), t) // is not set
}

func TestDecodeDuplicateKeysArrayIntoStructStrictMode(t *testing.T) {
	buf := []byte{0x84, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x46, 0x75, 0x6e, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r, func(dec *Decoder) { dec.strict = true })
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	err := d.Decode(&a)
	expect(err != nil, true, t)
	expect(fmt.Sprint(err), "strict-mode: duplicated key Fun in map", t)
}

func TestDecodeArrayIntoStructNonStringKeys(t *testing.T) {
	buf := []byte{0x84, 0x10, 0xf5, 0x11, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	err := d.Decode(&a)
	expect(err != nil, true, t)
	expect(fmt.Sprint(err), "array keys must be string, cborUnsignedInt received", t)
}

func TestDecodeArrayNonFieldIntoStruct(t *testing.T) {
	buf := []byte{0x84, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
}

func TestDecodeArrayNonFieldIntoStructWithValidTag(t *testing.T) {
	buf := []byte{0x84, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun   bool
		Other int8 `cbor:"Amt"`
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Other, int8(-2), t)
}

func TestDecodeIndefiniteArrayNonFieldIntoStruct(t *testing.T) {
	buf := []byte{0x8f, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21, 0xff}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
}

func TestDecodeArrayNonFieldIntoStructStrictMode(t *testing.T) {
	buf := []byte{0x84, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
	r := bytes.NewReader(buf)
	d := NewDecoder(r, func(dec *Decoder) { dec.strict = true })
	type MyType struct {
		Fun  bool
		None int8
	}
	var a MyType
	err := d.Decode(&a)
	expect(err != nil, true, t)
	expect(fmt.Sprint(err), "strict-mode: key Amt doesn't match with any field", t)
}

func TestDecodeArrayOutboundsIntoStruct(t *testing.T) {
	buf := []byte{0x86, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21, 0x64, 0x46, 0x61, 0x69, 0x6c, 0x04}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Amt, int8(-2), t)
	expect(r.Len(), 0, t)
}

func TestDecodeArrayOutboundsIntoStructStrictMode(t *testing.T) {
	buf := []byte{0x86, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21, 0x64, 0x46, 0x61, 0x69, 0x6c, 0x04}
	r := bytes.NewReader(buf)
	d := NewDecoder(r, func(dec *Decoder) { dec.strict = true })
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	err := d.Decode(&a)
	expect(err != nil, true, t)
	expect(fmt.Sprint(err), "strict-mode: destination struct fields num 2 doesn't match map length 3", t)
}

func TestDecodeArrayIntoStructWithNilValue(t *testing.T) {
	buf := []byte{0x84, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0xf6}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	type MyType struct {
		Fun bool
		Amt int8
	}
	var a MyType
	check(d.Decode(&a))
	expect(a.Fun, true, t)
	expect(a.Amt, int8(0), t)
	expect(r.Len(), 0, t)
}

func TestDecodePositiveBigNum(t *testing.T) {
	buf := []byte{0xc2, 0x49, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	a := new(big.Int)
	check(d.Decode(a))
	expect(fmt.Sprint(a), "18446744073709551616", t)
}

func TestDecodePositiveBigNumFromInterface(t *testing.T) {
	buf := []byte{0xc2, 0x49, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	expect(fmt.Sprint(a), "18446744073709551616", t)
}

func TestDecodeNegativeBigNum(t *testing.T) {
	buf := []byte{0xc3, 0x49, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	a := big.NewInt(-1)
	check(d.Decode(a))
	expect(fmt.Sprint(a), "-18446744073709551616", t)
}

func TestDecodeNegativeBigNumFromInterface(t *testing.T) {
	buf := []byte{0xc3, 0x49, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	expect(fmt.Sprint(a), "-18446744073709551616", t)
}

func TestDecodeUtf8DAteTimeFromInterface(t *testing.T) {
	buf := []byte{0xc0, 0x74, 0x32, 0x30, 0x30, 0x33, 0x2d, 0x31, 0x32, 0x2d, 0x31, 0x33, 0x54, 0x31, 0x38, 0x3a, 0x33, 0x30, 0x3a, 0x30, 0x32, 0x5a}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	expect(a.(time.Time).Year(), 2003, t)
	expect(a.(time.Time).Month(), time.December, t)
	expect(a.(time.Time).Day(), 13, t)
	expect(a.(time.Time).Hour(), 18, t)
	expect(a.(time.Time).Minute(), 30, t)
	expect(a.(time.Time).Nanosecond(), 0, t)
	expect(a.(time.Time).Location(), time.UTC, t)
}

func TestDecodeEpochDateTimeFromInterface(t *testing.T) {
	buf := []byte{0xc1, 0x1a, 0x3f, 0xdb, 0x5a, 0xaa}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	expect(a.(time.Time).Year(), 2003, t)
	expect(a.(time.Time).Month(), time.December, t)
	expect(a.(time.Time).Day(), 13, t)
	expect(a.(time.Time).Hour(), 18, t)
	expect(a.(time.Time).Minute(), 30, t)
	expect(a.(time.Time).Nanosecond(), 0, t)
	expect(a.(time.Time).Location(), time.Local, t)
}

func TestDecodeNegativeEpochDateTimeFromInterface(t *testing.T) {
	buf := []byte{0xc1, 0x3a, 0x01, 0x93, 0xa9, 0x4b}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	expect(a.(time.Time).Year(), 1969, t)
	expect(a.(time.Time).Month(), time.February, t)
	expect(a.(time.Time).Day(), 28, t)
	expect(a.(time.Time).Hour(), 20, t)
	expect(a.(time.Time).Minute(), 34, t)
	expect(a.(time.Time).Second(), 12, t)
	expect(a.(time.Time).Nanosecond(), 0, t)
	expect(a.(time.Time).Location(), time.Local, t)
}

func TestDecodeDecimalFraction(t *testing.T) {
	buf := []byte{0xc4, 0x82, 0x21, 0x19, 0x6a, 0xb3}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	expect(a, float32(273.15), t)
}

func TestDecodeBigFloat(t *testing.T) {
	buf := []byte{0xc5, 0x82, 0x20, 0x03}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	expect(a.(*big.Rat).String(), big.NewRat(3, 2).String(), t)
}

func TestDecodeBigFloatFromBigInt(t *testing.T) {
	buf := []byte{0xc5, 0x82, 0x20, 0xc2, 0x41, 0x03}
	r := bytes.NewReader(buf)
	d := NewDecoder(r)
	var a interface{}
	check(d.Decode(&a))
	expect(a.(*big.Rat).String(), big.NewRat(3, 2).String(), t)
}

// Some benchmarks
func BenchmarkDecodeUint8(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0x18, 0x6f}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a uint8
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeFLoat16(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0xf9, 0x3f, 0xe0}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a float16
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeFLoat32(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0xfa, 0x3f, 0x66, 0x66, 0x66}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a float32
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeUnsignedIntsIndefiniteArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0x9f, 0x04, 0x09, 0x19, 0x04, 0x00, 0x10, 0xff}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a []uint
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeUnsignedIntsArray(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0x84, 0x04, 0x09, 0x19, 0x04, 0x00, 0x10}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a []uint
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0x85, 0x04, 0x09, 0x19, 0x04, 0x00, 0x10, 0x83, 0x01, 0x02, 0x67, 0x65, 0x73, 0x70, 0x61, 0xc3, 0xb1, 0x61}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a interface{}
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeMapInterfaceKeyInterfaceValues(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a map[interface{}]interface{}
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeMapInterfaceValues(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a map[string]interface{}
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeMapInt8Values(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0x23, 0x63, 0x41, 0x6d, 0x74, 0x21}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a map[string]int8
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeMapIntoInterface(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		var a interface{}
		check(d.Decode(&a))
	}
}

func BenchmarkDecodeMapIntoStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := []byte{0xa2, 0x63, 0x46, 0x75, 0x6e, 0xf5, 0x63, 0x41, 0x6d, 0x74, 0x21}
		r := bytes.NewReader(buf)
		d := NewDecoder(r)
		type MyType struct {
			Fun bool
			Amt int8
		}
		var a MyType
		check(d.Decode(&a))
	}
}
