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
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"reflect"
	"time"
)

// Type of function that handles decoding of extensions
type handleDecFn func(*Decoder, reflect.Value) error

// tag maps, used by user code to register custom extensions
// using the major type 6 Optional Semantic Tagging for more
// information refer to http://tools.ietf.org/html/rfc7049#section-2.4
type extensionTagDecMap map[uintptr]handleDecFn

// global extensions register
var tagExtensionsDec extensionTagDecMap

// Registers a new extension in the extensions tags register
func (e *extensionTagDecMap) register(t reflect.Type, fn handleDecFn) error {
	tid := reflect.ValueOf(t).Pointer()
	if _, ok := tagExtensionsDec[tid]; ok {
		return errors.New(fmt.Sprintf("%s type is already registered\n", t))
	}
	tagExtensionsDec[tid] = fn
	return nil
}

// Look for a function registered to handle a given type
func (e *extensionTagDecMap) lookup(t reflect.Type) (handleDecFn, error) {
	tid := reflect.ValueOf(t).Pointer()
	fn, ok := tagExtensionsDec[tid]
	if !ok {
		return nil, errors.New(fmt.Sprintf(
			"%s not matched as registered extension handler", t))
	}
	return fn, nil
}

// Registers a new function to hanle decode of extensions
func RegisterExtensionFn(t reflect.Type, fn handleDecFn) error {
	return tagExtensionsDec.register(t, fn)
}

// Lookup for a registered function that handles the given type decode
func LookupExtensionFn(t reflect.Type) (handleDecFn, error) {
	return tagExtensionsDec.lookup(t)
}

// A Decoder reads and decode CBOR objects from an input stream.
type Decoder struct {
	parser *Parser
	strict bool
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader, options ...func(*Decoder)) *Decoder {
	d := &Decoder{parser: &Parser{r: r}, strict: false}
	if len(options) > 0 {
		for _, option := range options {
			option(d)
		}
	}
	return d
}

// Decode reads the next CBOR-encoded value from its
// input and stores it in the value pointed to by v.
// It also checks for the well-formedness of the 'data item'
func (dec *Decoder) Decode(v interface{}) error {
	major, info, err := dec.parser.parseInformation()
	if err != nil {
		return err
	}
	if err := dec.checkTypes(reflect.TypeOf(v), major, info); err != nil {
		return err
	}
	switch t := v.(type) {
	case nil:
		return errors.New("can't decode a value into nil")
	case *uint8:
		*t = dec.decodeUint8()
	case *int8:
		*t = dec.decodeInt8()
	case *uint16:
		*t = dec.decodeUint16()
	case *int16:
		*t = dec.decodeInt16()
	case *uint32:
		*t = dec.decodeUint32()
	case *int32:
		*t = dec.decodeInt32()
	case *uint64:
		*t = dec.decodeUint64()
	case *int64:
		*t = dec.decodeInt64()
	case *float16:
		*t = dec.decodeFloat16()
	case *float32:
		*t = dec.decodeFloat32()
	case *float64:
		*t = dec.decodeFloat64()
	case *big.Int:
		n := dec.decodeBigNum()
		if v.(*big.Int).Sign() < 0 {
			*t = *n.Neg(n)
		} else {
			*t = *n
		}
	case *time.Time:
		if major == cborTextString {
			*t = dec.decodeStringDateTime()
		} else {
			*t = dec.decodeEpochDateTime()
		}
	case *[]byte:
		*t = dec.decodeBytes()
	case *string:
		*t = dec.decodeString()
	case *bool:
		*t = dec.decodeBool()
	case *interface{}:
		return dec.decode(reflect.ValueOf(v).Elem())
	case reflect.Value:
		return dec.decode(t.Elem())
	default:
		rv := reflect.ValueOf(v)
		if rv.Kind() == reflect.Ptr && !rv.IsNil() || !rv.IsValid() {
			return dec.decode(rv.Elem())
		}
		return &InvalidDecodeError{rv.Type()}
	}

	return nil
}

// decode is being used when the type of the receiver of the decode
// operation is a slice, a map an interface or any type of custom type
func (dec *Decoder) decode(rv reflect.Value) (err error) {
	// Decode nil and undef into zero values
	if dec.parser.isNil() || dec.parser.isUndef() {
		if rv.Kind() == reflect.Ptr {
			if !rv.IsNil() {
				rv.Set(reflect.Zero(rv.Type()))
			}
			return nil
		}
		if rv.IsValid() && rv.CanSet() {
			rv.Set(reflect.Zero(rv.Type()))
		}
		return nil
	}
	var handler handleDecFn
	handler, err = dec.lookupFn(rv)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()
	return handler(dec, rv)
}

// lookup for decode function based on type Kind
func (dec *Decoder) lookupFn(rv reflect.Value) (handler handleDecFn, e error) {
	rk := rv.Kind()
	switch rk {
	case reflect.Map:
		handler = (*Decoder).decodekMap
	case reflect.Struct:
		handler = (*Decoder).decodekStruct
	case reflect.Interface:
		handler = (*Decoder).decodekInterface
	case reflect.String:
		handler = (*Decoder).decodekString
	case reflect.Int:
		handler = (*Decoder).decodekInt
	case reflect.Int8:
		handler = (*Decoder).decodekInt8
	case reflect.Int16:
		handler = (*Decoder).decodekInt16
	case reflect.Int32:
		handler = (*Decoder).decodekInt32
	case reflect.Int64:
		handler = (*Decoder).decodekInt64
	case reflect.Uint:
		handler = (*Decoder).decodekUint
	case reflect.Uint8:
		handler = (*Decoder).decodekUint8
	case reflect.Uint16:
		handler = (*Decoder).decodekUint16
	case reflect.Uint32:
		handler = (*Decoder).decodekUint32
	case reflect.Uint64:
		handler = (*Decoder).decodekUint64
	case reflect.Bool:
		handler = (*Decoder).decodekBool
	case reflect.Float32:
		handler = (*Decoder).decodekFloat32
	case reflect.Float64:
		handler = (*Decoder).decodekFloat64
	case reflect.Slice:
		handler = (*Decoder).decodekSlice
	case reflect.Array:
		handler = (*Decoder).decodekArray
	default:
		handler, e = LookupExtensionFn(rv.Type())
	}
	return handler, e
}

// check if the major and info types are the expected for decode and return
// an error in case of the encoded data doesn't match or well-formedness errors
func (dec *Decoder) checkTypes(t reflect.Type, major Major, info byte) error {
	if major == cborTag || major == cborDataArray || major == cborDataMap || t == reflect.TypeOf(reflect.Value{}) {
		return nil
	}
	msg := "expected %s, got %s (major %d, info %d [%#v])\n"
	e, ok := expectedTypesMap[major][info]
	if !ok {
		switch major {
		case cborUnsignedInt:
			if info <= cborSmallInt {
				e = reflect.PtrTo(reflect.TypeOf(uint8(0)))
				break
			}
			return errors.New(fmt.Sprintf("Unknown info %d for major 1", info))
		case cborByteString:
			if info <= cborSmallInt || info == cborIndefinite {
				e = reflect.TypeOf([]byte{})
				break
			}
			return errors.New(fmt.Sprintf("Unknown info %d for major 2", info))
		case cborTextString:
			if info <= cborSmallInt || info == cborIndefinite {
				e = reflect.TypeOf("")
				break
			}
			return errors.New(fmt.Sprintf("Unknown info %d for major 3", info))
		case cborNC:
			if info == cborFalse || info == cborTrue {
				e = reflect.TypeOf(false)
				break
			}
			if info == cborNil || info == cborUndef {
				e = reflect.TypeOf(reflect.Interface)
				break
			}
			return errors.New(fmt.Sprintf("Unknown info %d for major 7", info))
		}
	}
	e = reflect.PtrTo(e)
	header := byte((major << 5)) | info
	if e != t {
		return errors.New(fmt.Sprintf(msg, t, e, major, info, header))
	}
	return nil
}

// Decode into an unsigned int
// of any size between 8 and 64 bits
func (dec *Decoder) decodeUint() uint64 {
	return dec.parser.buflen()
}

// Decode into an signed int
// of any size between 8 and 64 bits
func (dec *Decoder) decodeInt() int64 {
	return ^int64(dec.parser.buflen())
}

// Decodes into an unsigned integer of 8 bits
func (dec *Decoder) decodeUint8() uint8 {
	return dec.parser.parseUint8()
}

// Decodes into an unsigned integer of 16 bits
func (dec *Decoder) decodeUint16() uint16 {
	return dec.parser.parseUint16()
}

// Decodes into an unsigend integer of 32 bits
func (dec *Decoder) decodeUint32() uint32 {
	return dec.parser.parseUint32()
}

// Decodes into an unsigned integer of 64 bits
func (dec *Decoder) decodeUint64() uint64 {
	return dec.parser.parseUint64()
}

// Decodes into a signed integer of 8 bits
func (dec *Decoder) decodeInt8() int8 {
	return int8(^dec.decodeUint8())
}

// Decodes into a signed integer of 16 bits
func (dec *Decoder) decodeInt16() int16 {
	return int16(^dec.decodeUint16())
}

// Decodes into a signed integer of 32 bits
func (dec *Decoder) decodeInt32() int32 {
	return int32(^dec.decodeUint32())
}

// Decodes into a signed integer of 64 bits
func (dec *Decoder) decodeInt64() int64 {
	return int64(^dec.decodeUint64())
}

// Decode into a float16
func (dec *Decoder) decodeFloat16() float16 {
	return dec.parser.parseFloat16()
}

// Decode into a float32
func (dec *Decoder) decodeFloat32() float32 {
	return dec.parser.parseFloat32()
}

// Decode into a float64
func (dec *Decoder) decodeFloat64() float64 {
	return dec.parser.parseFloat64()
}

// Decode a string date representation
// that follows the standard format defined in
// RFC3339 with RFC4287 Section 3.3 additions
func (dec *Decoder) decodeStringDateTime() time.Time {
	major, info, err := dec.parser.parseInformation()
	checkErr(err)

	if major != cborTextString {
		log.Fatal(fmt.Sprintf("expected UTF-8 string, found %s", info))
	}
	t, err := time.Parse(time.RFC3339, dec.decodeString())
	checkErr(err)
	return t
}

// Decode a positive or negative
// integer or floating point with
// additional information a time.Time
func (dec *Decoder) decodeEpochDateTime() time.Time {
	major, _, err := dec.parser.parseInformation()
	checkErr(err)
	var n int64
	switch major {
	case cborUnsignedInt:
		n = int64(dec.decodeUint())
	case cborNegativeInt:
		n = dec.decodeInt()
	default:
		switch dec.parser.header {
		case absoluteFloat16:
			n = int64(int(dec.decodeFloat16()))
		case absoluteFloat32:
			n = int64(int(dec.decodeFloat32()))
		case absoluteFloat64:
			n = int64(int(dec.decodeFloat64()))
		default:
			log.Fatal("can't decode Epoch timestamp %#v", dec.parser.header)
		}
	}
	return time.Unix(n, int64(0))
}

// Decode a decimal fraction as defined in Section 2.4.3 of RFC7049
// http://tools.ietf.org/html/rfc7049#section-2.4.3
func (dec *Decoder) decodeDecimalFraction() float32 {
	major, _, err := dec.parser.parseInformation()
	checkErr(err)
	if major != cborDataArray {
		log.Fatal("Decimal Fraction must be represented as an array of two elements")
	}

	major, _, err = dec.parser.parseInformation()
	checkErr(err)
	if major > cborNegativeInt {
		log.Fatal(fmt.Sprintf("Can't decode %s as decimal fraction exponent", major))
	}
	e := dec.decodeInt()
	major, _, err = dec.parser.parseInformation()
	checkErr(err)
	if major > cborNegativeInt {
		log.Fatal(fmt.Sprintf("Can't decode %s as decimal fraction mantissa", major))
	}
	var m int64
	if major == cborUnsignedInt {
		m = int64(dec.decodeUint())
	} else {
		m = dec.decodeInt()
	}
	return decimalFractionToFloat(m, e)
}

// Decode a big float a defined in Section 2.3.4 of RFC7049
// http://tools.ietf.org/html/rfc7049#section-2.4.3
func (dec *Decoder) decodeBigFloat() *big.Rat {
	major, _, err := dec.parser.parseInformation()
	checkErr(err)
	if major != cborDataArray {
		log.Fatal("Decimal Fraction must be represented as an array of two elements")
	}

	major, _, err = dec.parser.parseInformation()
	checkErr(err)
	if major > cborNegativeInt {
		log.Fatal(fmt.Sprintf("Can't decode %s as decimal fraction exponent", major))
	}
	e := dec.decodeInt()
	major, info, err := dec.parser.parseInformation()
	checkErr(err)
	if major > cborNegativeInt && (major != cborTag && info != cborBigNum) {
		log.Fatal(fmt.Sprintf("Can't decode %s as decimal fraction mantissa", major))
	}
	switch major {
	case cborUnsignedInt:
		m := int64(dec.decodeUint())
		return bigFloatToRatFromInt64(m, e)
	case cborNegativeInt:
		m := int64(dec.decodeInt())
		return bigFloatToRatFromInt64(m, e)
	case cborTag:
		m := dec.decodeBigNum()
		return bigFloatToRatFromBigInt(m, e)
	}
	return big.NewRat(0, 0)
}

// Decode big num
func (dec *Decoder) decodeBigNum() *big.Int {
	major, info, err := dec.parser.parseInformation()
	checkErr(err)

	if major != cborByteString {
		log.Fatal(fmt.Sprintf("expected bytes found 0x%x", info))
	}
	i := new(big.Int)
	i.SetBytes(dec.decodeBytes())
	return i
}

// Decode into a byte string
func (dec *Decoder) decodeBytes() []byte {
	_, info := dec.parser.parseHeader()
	if dec.parser.isNil() || dec.parser.isUndef() {
		return nil
	}

	if info != cborIndefinite {
		_, d, err := dec.parser.scan(int(dec.parser.buflen()))
		checkErr(err)
		return d
	}

	return dec.decodeIndefiniteBytes(nil)
}

// Decode an UTF8 text string
func (dec *Decoder) decodeString() string {
	return string(dec.decodeBytes())
}

// decode an indefinite stream of bytes
// it doesn't really decode it, just read it and returns it back
func (dec *Decoder) decodeIndefiniteBytes(buf []byte) []byte {
	for {
		if dec.parser.isBreak() {
			break
		}
		buflen := int(dec.parser.buflen())
		n, d, err := dec.parser.scan(buflen)
		checkErr(err)
		if n < buflen {
			panic(fmt.Sprintf("expected %d bytes in buffer, got %d", buflen, n))
		}
		buf = append(buf, d...)
		if _, _, err := dec.parser.parseInformation(); err != nil {
			panic(err)
		}
	}
	return buf
}

// Decode into a boolean value
func (dec *Decoder) decodeBool() bool {
	return dec.parser.parseBool()
}

// helper function that logs and exists if err is not nil
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
