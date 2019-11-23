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
	"io"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"time"
	"unicode"
	"unsafe"
)

// Marshaler interface
type Marshaler interface {
	MarshalCBOR() ([]byte, error)
}

// An Encoder writes and encode CBOR objects to an output stream
type Encoder struct {
	composer  *composer
	canonical bool
	strict    bool
}

// NewEncoder returns a new encoder that write to w
func NewEncoder(w io.Writer, options ...func(*Encoder)) *Encoder {
	e := &Encoder{composer: &composer{w: w}, strict: false}
	if len(options) > 0 {
		for _, option := range options {
			option(e)
		}
	}
	return e
}

// Check if the pointer passed to Encode
// is nil and then call enc.encodeNil()
func (enc *Encoder) isValidPointer(t unsafe.Pointer) bool {
	if t == nil {
		enc.encodeNil()
		return false
	}
	return true
}

// Encoder takes any object passed as parameter and
// writes it into a io.Writer using the C.B.O.R encoding format.
func (enc *Encoder) Encode(v interface{}) (err error) {

	if ok, err := enc.fastPath(v); !ok {
		return err
	}

	return err
}

func (enc *Encoder) fastPath(v interface{}) (ok bool, err error) {
	// fast path encoding for builting and simple values
	switch t := v.(type) {
	case nil:
		err = enc.composer.composeNil()
	case bool:
		err = enc.composer.composeBoolean(t)
	case int, int8, int16, int32, int64:
		_, err = enc.composer.composeInt(reflect.ValueOf(t).Int())
	case uint, uint8, uint16, uint32, uint64, uintptr:
		_, err = enc.composer.composeUint(reflect.ValueOf(t).Uint())
	case float16:
		err = encodeFloat16(enc.composer, reflect.ValueOf(t))
	case float32:
		err = encodeFloat32(enc.composer, reflect.ValueOf(t))
	case float64:
		err = encodeFloat64(enc.composer, reflect.ValueOf(t))
	case big.Int:
		if t.Sign() < 0 {
			err = enc.composer.composeBigInt(&t)
		} else {
			err = enc.composer.composeBigUint(&t)
		}
	case time.Time:
		err = enc.composer.composeEpochDateTime(&t)
	case big.Rat:
		err = enc.composer.composeBigFloat(&t)
	case CBORMIME:
		err = enc.composer.composeCBORMIME(&t)
	case []uint8:
		err = enc.composer.composeBytes(t)
	case string:
		err = enc.composer.composeString(t)
	case *bool:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			err = enc.composer.composeBoolean(*t)
		}
	case *int, *int8, *int16, *int32, *int64:
		if enc.isValidPointer(unsafe.Pointer(reflect.ValueOf(v).Pointer())) {
			_, err = enc.composer.composeInt(reflect.ValueOf(v).Elem().Int())
		}
	case *uint, *uint8, *uint16, *uint32, *uint64, *uintptr:
		if enc.isValidPointer(unsafe.Pointer(reflect.ValueOf(v).Pointer())) {
			_, err = enc.composer.composeUint(reflect.ValueOf(v).Elem().Uint())
		}
	case *float16:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			err = encodeFloat16(enc.composer, reflect.ValueOf(t).Elem())
		}
	case *float32:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			err = encodeFloat32(enc.composer, reflect.ValueOf(t).Elem())
		}
	case *float64:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			err = encodeFloat64(enc.composer, reflect.ValueOf(t).Elem())
		}
	case *big.Int:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			if t.Sign() < 0 {
				err = enc.composer.composeBigInt(t)
			} else {
				err = enc.composer.composeBigUint(t)
			}
		}
	case *time.Time:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			err = enc.composer.composeEpochDateTime(t)
		}
	case *big.Rat:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			err = enc.composer.composeBigFloat(t)
		}
	case *CBORMIME:
		err = enc.composer.composeCBORMIME(t)
	case *url.URL:

	case *[]uint8:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			err = enc.composer.composeBytes(*t)
		}
	case *string:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			err = enc.composer.composeString(*t)
		}
	case reflect.Value:
		err = enc.encode(t)
	default:
		err = enc.encode(reflect.ValueOf(v))
	}

	if err == nil {
		return true, nil
	}
	return false, err
}

// encode is being used when the type of the supplier of the encode
// operation is a slice, a map an interface or any other custom type
func (enc *Encoder) encode(rv reflect.Value) (err error) {

	// If rv is a pointer, get the value it's references
	for rv.Kind() == reflect.Ptr {
		// Lets encode nil values if present
		if rv.IsNil() {
			enc.encodeNil()
			return
		}
		rv = rv.Elem()
	}
	if !rv.IsValid() {
		// Encode as nil if the value is not valid
		enc.encodeNil()
		return
	}

	rt := rv.Type()
	switch rt {
	case bigNumType:
		t := rv.Interface().(big.Int)
		if t.Sign() < 0 {
			err = enc.composer.composeBigInt(&t)
		} else {
			err = enc.composer.composeBigUint(&t)
		}
	case bigFloatType:
		r := rv.Interface().(big.Rat)
		err = enc.composer.composeBigFloat(&r)
	case epochTimeType:
		t := rv.Interface().(time.Time)
		err = enc.composer.composeEpochDateTime(&t)
	case cborMimeType:
		t := rv.Interface().(CBORMIME)
		err = enc.composer.composeCBORMIME(&t)
	case float16Type:
		err = enc.composer.composeFloat16(rv.Interface().(float16))
	default:
		switch rt.Kind() {
		case reflect.Bool:
			err = enc.composer.composeBoolean(rv.Bool())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			_, err = enc.composer.composeUint(rv.Uint())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			_, err = enc.composer.composeInt(rv.Int())
		case reflect.Float32:
			err = enc.composer.composeFloat32(rv.Interface().(float32))
		case reflect.Float64:
			err = enc.composer.composeFloat64(rv.Float())
		case reflect.String:
			err = enc.composer.composeString(rv.String())
		case reflect.Invalid:
			err = enc.composer.composeNil()
		case reflect.Slice, reflect.Array:
			err = enc.encodeSlice(rv)
		case reflect.Map:
			err = enc.encodeMap(rv)
		case reflect.Struct:
			err = enc.encodeStruct(rv)
		case reflect.Interface:
			err = enc.encodeInterface(rv)
		case reflect.Ptr:
			err = enc.encode(rv.Elem())
			// default:
			// 	err = enc.lookupExtension(rv)
		}
	}

	return err
}

// Encode a Nil value
func (enc *Encoder) encodeNil() {
	if err := enc.composer.composeNil(); err != nil {
		panic(err)
	}
}

// Encode a boolean value
func (enc *Encoder) encodeBool(v bool) {
	if err := enc.composer.composeBoolean(v); err != nil {
		panic(err)
	}
}

// Encode a signed in of any size
func (enc *Encoder) encodeInt(v int64) {
	if _, err := enc.composer.composeInt(v); err != nil {
		panic(err)
	}
}

// Encode an unsigned int of any size
func (enc *Encoder) encodeUint(v uint64) {
	if _, err := enc.composer.composeUint(v); err != nil {
		panic(err)
	}
}

type floatEncoder int // number of bits

func (bits floatEncoder) encode(c *composer, v reflect.Value) (err error) {
	f := v.Float()
	if math.IsInf(f, 0) {
		err = bits.encodeNewInfinity(c, v)
	} else if math.IsNaN(f) {
		err = bits.encodeNewNaN(c, v)
	} else {
		b := int(bits)
		if b == 16 || v.Type() == float16Type {
			err = c.composeFloat16(v.Interface().(float16))
		} else if b == 32 {
			err = c.composeFloat32(v.Interface().(float32))
		} else {
			err = c.composeFloat64(f)
		}
	}
	return err
}

func (bits floatEncoder) encodeNewInfinity(c *composer, v reflect.Value) (err error) {
	if c.isCanonical() {
		err = c.composeCanonicalInfinity()
	} else {
		switch int(bits) {
		case 16:
			err = c.composeCanonicalInfinity()
		case 32:
			err = c.composeInfinity()
		case 64:
			err = c.composeDoublePrecissionInfinity()
		default:
			err = &UnsupportedValueError{v, fmt.Sprintf("%#v", v)}
		}
	}
	return err
}

func (bits floatEncoder) encodeNewNaN(c *composer, v reflect.Value) (err error) {
	if c.isCanonical() {
		err = c.composeCanonicalNaN()
	} else {
		switch int(bits) {
		case 16:
			err = c.composeCanonicalNaN()
		case 32:
			err = c.composeNaN()
		case 64:
			err = c.composeDoublePrecissionNaN()
		default:
			err = &UnsupportedValueError{v, v.Type().String()}
		}
	}
	return err
}

var (
	encodeFloat16 = (floatEncoder(16)).encode
	encodeFloat32 = (floatEncoder(32)).encode
	encodeFloat64 = (floatEncoder(64)).encode
)

// Encode an Slice
func (enc *Encoder) encodeSlice(rv reflect.Value) error {
	etp := rv.Type().Elem()
	if etp.Kind() == reflect.Uint8 {
		// Bytes String
		enc.composer.composeBytes(rv.Bytes())
		return nil
	}
	l := rv.Len()
	if _, err := enc.composer.composeUint(uint64(l), cborDataArray); err != nil {
		return fmt.Errorf("while enoding slice %v: %s", rv.Type(), err.Error())
	}

	for i := 0; i < l; i++ {
		if err := enc.encode(rv.Index(i)); err != nil {
			return fmt.Errorf("while enoding slice %v: %s", rv.Type(), err.Error())
		}
	}
	return nil
}

// Encode a Map
func (enc *Encoder) encodeMap(rv reflect.Value) error {
	l := rv.Len()
	if _, err := enc.composer.composeUint(uint64(l), cborDataMap); err != nil {
		return fmt.Errorf("while enoding map %v: %s", rv.Type(), err.Error())
	}

	for _, key := range rv.MapKeys() {
		if err := enc.encode(key); err != nil {
			return fmt.Errorf("while enoding map %v: %s", rv.Type(), err.Error())
		}
		if err := enc.encode(rv.MapIndex(key)); err != nil {
			return fmt.Errorf("while enoding map %v: %s", rv.Type(), err.Error())
		}
	}
	return nil
}

// Encode a Struct
func (enc *Encoder) encodeStruct(rv reflect.Value, array ...bool) error {
	// buffer the fields encoding
	buf := bytes.NewBuffer(nil)
	w := enc.composer.w
	enc.composer.w = buf

	exportedFields := 0
	numfields := rv.NumField()
	for i := 0; i < numfields; i++ {
		field := rv.Type().Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		key := field.Name
		if unicode.IsUpper(rune(key[0])) {
			tag := field.Tag.Get("cbor")
			if tag != "" {
				if tag == "-" {
					continue
				}
				key = tag
			}
			exportedFields++
			if err := enc.composer.composeString(key); err != nil {
				return fmt.Errorf("while enoding struct %v: %s", rv.Type(), err.Error())
			}
			if err := enc.encode(rv.Field(i)); err != nil {
				return fmt.Errorf("while enoding struct %v: %s", rv.Type(), err.Error())
			}
		}
	}

	enc.composer.w = w
	var info byte
	if len(array) > 0 && array[0] {
		info, _ = calculateInfoFromIntLength(exportedFields * 2)
	} else {
		info, _ = calculateInfoFromIntLength(exportedFields)
	}
	if err := enc.composer.composeInformation(cborDataMap, info); err != nil {
		return fmt.Errorf("while enoding struct %v: %s", rv.Type(), err.Error())
	}
	if _, err := enc.composer.write(buf.Bytes()); err != nil {
		return fmt.Errorf("while enoding struct %v: %s", rv.Type(), err.Error())
	}
	return nil
}

func (enc *Encoder) encodeInterface(rv reflect.Value) error { return enc.encode(rv.Elem()) }

// helper function that calculates the size
// of the info byte depending on the given length
func calculateInfoFromIntLength(l int) (info byte, err error) {
	if l < int(cborSmallInt) {
		info = byte(l)
	} else {
		if info, err = infoHelper(uint(l)); err != nil {
			return 0, err
		}
	}
	return info, nil
}
