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
	"math/big"
	"reflect"
	"time"
	"unsafe"
)

// Type of function that handler encoding of extensions
type handleEncFn handleDecFn

// An Encoder writes and encode CBOR objects to an output stream
type Encoder struct {
	composer *Composer
	strict   bool
}

// NewEncoder returns a new encoder that write to w
func NewEncoder(w io.Writer, options ...func(*Encoder)) *Encoder {
	e := &Encoder{composer: &Composer{w: w}, strict: false}
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
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	// fast path encoding for simple values
	switch t := v.(type) {
	case nil:
		enc.encodeNil()
	case bool:
		enc.encodeBool(t)
	case uint8:
		enc.encodeUint(uint64(t))
	case int8:
		enc.encodeInt(int64(t))
	case uint16:
		enc.encodeUint(uint64(t))
	case int16:
		enc.encodeInt(int64(t))
	case uint32:
		enc.encodeUint(uint64(t))
	case int32:
		enc.encodeInt(int64(t))
	case uint64:
		enc.encodeUint(t)
	case int64:
		enc.encodeInt(t)
	case uint:
		enc.encodeUint(uint64(t))
	case int:
		enc.encodeInt(int64(t))
	case float16:
		enc.encodeFloat16(t)
	case float32:
		enc.encodeFloat32(t)
	case float64:
		enc.encodeFloat64(t)
	case big.Int:
		if t.Sign() < 0 {
			enc.encodeBigInt(t)
		} else {
			enc.encodeBigUint(t)
		}
	case time.Time:
		enc.encodeEpochDateTime(t)
	case big.Rat:
		enc.encodeBigFloat(t)
	case []uint8:
		enc.encodeByteString(t)
	case string:
		enc.encodeTextString(t)
	case *bool:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeBool(*t)
		}
	case *uint8:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeUint(uint64(*t))
		}
	case *int8:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeInt(int64(*t))
		}
	case *uint16:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeUint(uint64(*t))
		}
	case *int16:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeInt(int64(*t))
		}
	case *uint32:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeUint(uint64(*t))
		}
	case *int32:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeInt(int64(*t))
		}
	case *uint64:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeUint(*t)
		}
	case *int64:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeInt(*t)
		}
	case *uint:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeUint(uint64(*t))
		}
	case *int:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeInt(int64(*t))
		}
	case *float16:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeFloat16(*t)
		}
	case *float32:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeFloat32(*t)
		}
	case *float64:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeFloat64(*t)
		}
	case *big.Int:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			if t.Sign() < 0 {
				enc.encodeBigInt(*t)
			} else {
				enc.encodeBigUint(*t)
			}
		}
	case *time.Time:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeEpochDateTime(*t)
		}
	case *big.Rat:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeBigFloat(*t)
		}
	case *[]uint8:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeByteString(*t)
		}
	case *string:
		if enc.isValidPointer(unsafe.Pointer(t)) {
			enc.encodeTextString(*t)
		}
	case reflect.Value:
		enc.encode(t, v)
	default:
		enc.encode(reflect.ValueOf(v), v)
	}

	return nil
}

// encode is being used when the type of the supplier of the encode
// operation is a slice, a map an interface or any other custom type
func (enc *Encoder) encode(rv reflect.Value, vs ...interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New(fmt.Sprint(r))
		}
	}()

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
	var v interface{} = rv.Interface()
	if len(vs) > 0 {
		v = vs[0]
	}

	switch rv.Type().Kind() {
	case reflect.Bool:
		err = enc.composer.composeBoolean(v.(bool))
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		_, err = enc.composer.composeUint(v.(uint64))
	case reflect.Int32:
		_, err = enc.composer.composeInt(int64(v.(int32)))
	case reflect.Int:
		_, err = enc.composer.composeInt(int64(v.(int)))
	case reflect.Int8, reflect.Int16, reflect.Int64:
		_, err = enc.composer.composeInt(v.(int64))
	case reflect.Float32:
		err = enc.composer.composeFloat32(v.(float32))
	case reflect.Float64:
		err = enc.composer.composeFloat64(v.(float64))
	case reflect.String:
		enc.encodeTextString(v.(string))
	case reflect.Invalid:
		err = enc.composer.composeNil()
	case reflect.Slice, reflect.Array:
		enc.encodeSlice(rv)
	case reflect.Map:
		enc.encodeMap(rv)
		// case reflect.Struct:
		// 	err = enc.encodeStruct()
		// case reflect.Interface:
		// 	err = enc.encodeInterface()
		// default:
		// 	err = enc.lookupExtension(rv)
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

// Encode a float16
func (enc *Encoder) encodeFloat16(v float16) {
	if err := enc.composer.composeFloat16(v); err != nil {
		panic(err)
	}
}

// Encode a float32
func (enc *Encoder) encodeFloat32(v float32) {
	if err := enc.composer.composeFloat32(v); err != nil {
		panic(err)
	}
}

// Encode a float64
func (enc *Encoder) encodeFloat64(v float64) {
	if err := enc.composer.composeFloat64(v); err != nil {
		panic(err)
	}
}

// Encode a bytes string
func (enc *Encoder) encodeByteString(v []byte) {
	if err := enc.composer.composeBytes(v); err != nil {
		panic(err)
	}
}

// Encode a positive big.Int
func (enc *Encoder) encodeBigUint(v big.Int) {
	if err := enc.composer.composeBigUint(v); err != nil {
		panic(err)
	}
}

// Encode a negative big.Int
func (enc *Encoder) encodeBigInt(v big.Int) {
	if err := enc.composer.composeBigInt(v); err != nil {
		panic(err)
	}
}

// Encode a datetime as epoch
func (enc *Encoder) encodeEpochDateTime(v time.Time) {
	if err := enc.composer.composeEpochDateTime(v); err != nil {
		panic(err)
	}
}

// Encode a big float
func (enc *Encoder) encodeBigFloat(v big.Rat) {
	if err := enc.composer.composeBigFloat(v); err != nil {
		panic(err)
	}
}

// Encode a Text String (UTF-8)
func (enc *Encoder) encodeTextString(v string) {
	if err := enc.composer.composeString(v); err != nil {
		panic(err)
	}
}

// Encode an Slice
func (enc *Encoder) encodeSlice(rv reflect.Value) {
	etp := rv.Type().Elem()
	if etp.Kind() == reflect.Uint8 {
		// Bytes String
		enc.encodeByteString(rv.Bytes())
		return
	}
	l := rv.Len()
	info, err := calculateInfoFromIntLength(l)
	if err != nil {
		panic(err)
	}
	if err := enc.composer.composeInformation(cborDataArray, info); err != nil {
		panic(err)
	}
	if info > cborSmallInt {
		enc.encodeUint(uint64(l))
	}
	for i := 0; i < l; i++ {
		if err := enc.encode(rv.Index(i)); err != nil {
			panic(err)
		}
	}
}

// Encode a Map
func (enc *Encoder) encodeMap(rv reflect.Value) {
	l := rv.Len()
	info, err := calculateInfoFromIntLength(l)
	if err != nil {
		panic(err)
	}
	if err := enc.composer.composeInformation(cborDataMap, info); err != nil {
		panic(err)
	}
	if info > cborSmallInt {
		enc.encodeUint(uint64(l))
	}
	for _, key := range rv.MapKeys() {
		if err := enc.encode(key); err != nil {
			panic(err)
		}
		if err := enc.encode(rv.MapIndex(key)); err != nil {
			panic(err)
		}
	}

}

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
