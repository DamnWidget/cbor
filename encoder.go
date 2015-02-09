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
	"io"
	"math/big"
	"time"
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
	// case string:
	// 	enc.encodeStringTest(v)
	case *bool:
		enc.encodeBool(*t)
	case *uint8:
		enc.encodeUint(uint64(*t))
	case *int8:
		enc.encodeInt(int64(*t))
	case *uint16:
		enc.encodeUint(uint64(*t))
	case *int16:
		enc.encodeInt(int64(*t))
	case *uint32:
		enc.encodeUint(uint64(*t))
	case *int32:
		enc.encodeInt(int64(*t))
	case *uint64:
		enc.encodeUint(*t)
	case *int64:
		enc.encodeInt(*t)
	case *uint:
		enc.encodeUint(uint64(*t))
	case *int:
		enc.encodeInt(int64(*t))
	case *float16:
		enc.encodeFloat16(*t)
	case *float32:
		enc.encodeFloat32(*t)
	case *float64:
		enc.encodeFloat64(*t)
	case *big.Int:
		if t.Sign() < 0 {
			enc.encodeBigInt(*t)
		} else {
			enc.encodeBigUint(*t)
		}
	case *time.Time:
		enc.encodeEpochDateTime(*t)
	case *big.Rat:
		enc.encodeBigFloat(*t)
	case *[]uint8:
		enc.encodeByteString(*t)
		// case *string:
		// 	enc.encodeStringTest(*v)
		// case reflect.Value:
		// 	enc.encode(v)
		// default:
		// 	rv := reflect.ValueOf(*v)
		// 	enc.encode(rv)
	}

	return nil
}

// encode is being used when the type of the supplier of the encode
// operation is a slice, a map an interface or any other custom type
// func (enc *Encoder) encode(rv reflect.Value) (err error) {
// 	// If rv is a pointer, get the value it's references
// 	for rv.Kind() == reflect.Ptr {
// 		// Lets encode nil values if present
// 		if rv.IsNil() {
// 			enc.encodeNil()
// 			return
// 		}
// 		rv = rv.Elem()
// 	}
// 	if !rv.IsValid() {
// 		// Encode as nil if the value is not valid
// 		enc.encodeNil()
// 		return
// 	}
// 	var handler handleEncFn
// 	handler, err = enc.lookupFn(rv)
// 	if err != nil {
// 		return err
// 	}
// 	defer func() {
// 		if r := recover(); r != nil {
// 			err = errors.New(fmt.Sprint(r))
// 		}
// 	}()
// 	return handler(enc, rv)
// }

// lookup for an encode function based on the value type
// func (enc *Encoder) lookupFn(rv reflect.Value) (handler handleEncFn, e error) {
// 	switch rv.Type().Kind() {
// 	case reflect.Bool:
// 		handler = (*Encoder).encodekBool
// 	case reflect.Uint8:
// 		handler = (*Encoder).encodekUint8
// 	case reflect.Int8:
// 		handler = (*Encoder).encodekInt8
// 	case reflect.Uint16:
// 		handler = (*Encoder).encodekUint16
// 	case reflect.Int16:
// 		handler = (*Encoder).encodekInt16
// 	case reflect.Uint32:
// 		handler = (*Encoder).encodekUin32
// 	case reflect.Int32:
// 		handler = (*Encoder).encodekInt32
// 	case reflect.Uint64:
// 		handler = (*Encoder).encodekUin64
// 	case reflect.Int64:
// 		handler = (*Encoder).encodekInt64
// 	case reflect.Uint:
// 		handler = (*Encoder).encodekUint
// 	case reflect.Int:
// 		handler = (*Encoder).encodekInt
// 	case reflect.Float32:
// 		handler = (*Encoder).encodekFloat32
// 	case reflect.Float64:
// 		handler = (*Encoder).encodekFloat64
// 	case reflect.Invalid:
// 		handler = (*Encoder).encodekInvalid
// 	case reflect.Slice, reflect.Array:
// 		handler = (*Encoder).encodekSlice
// 	case reflect.Map:
// 		handler = (*Encoder).encodekMap
// 	case reflect.Struct:
// 		handler = (*Encoder).encodekStruct
// 	case reflect.Interface:
// 		handler = (*Encoder).encodekInterface
// 	default:
// 		handler, e = LookupExtensionFn(rv.Type())
// 	}
// 	return handler, e
// }

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
