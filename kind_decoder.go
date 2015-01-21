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
	"reflect"
	"strconv"
	"strings"
)

// magic error to force the decoder to continue in non strict mode
var forceContinueError = errors.New("")

const (
	d_NOP uint = iota
	d_BREAK
	d_CONTINUE
)

func (dec *Decoder) decodekInt(rv reflect.Value) error {
	rv.SetInt(^int64(dec.parser.buflen()))
	return nil
}

func (dec *Decoder) decodekUint(rv reflect.Value) error {
	rv.SetUint(dec.parser.buflen())
	return nil
}

func (dec *Decoder) decodekInt8(rv reflect.Value) error {
	rv.SetInt(int64(dec.decodeInt8()))
	return nil
}

func (dec *Decoder) decodekUint8(rv reflect.Value) error {
	rv.SetUint(uint64(dec.decodeUint8()))
	return nil
}

func (dec *Decoder) decodekInt16(rv reflect.Value) error {
	rv.SetInt(int64(dec.decodeInt16()))
	return nil
}

func (dec *Decoder) decodekUint16(rv reflect.Value) error {
	rv.SetUint(uint64(dec.decodeUint16()))
	return nil
}

func (dec *Decoder) decodekInt32(rv reflect.Value) error {
	rv.SetInt(int64(dec.decodeInt32()))
	return nil
}

func (dec *Decoder) decodekUint32(rv reflect.Value) error {
	rv.SetUint(uint64(dec.decodeUint32()))
	return nil
}

func (dec *Decoder) decodekInt64(rv reflect.Value) error {
	rv.SetInt(int64(dec.decodeInt64()))
	return nil
}

func (dec *Decoder) decodekUint64(rv reflect.Value) error {
	rv.SetUint(uint64(dec.decodeUint64()))
	return nil
}

func (dec *Decoder) decodekFloat32(rv reflect.Value) error {
	rv.SetFloat(float64(dec.decodeFloat32()))
	return nil
}

func (dec *Decoder) decodekFloat64(rv reflect.Value) error {
	rv.SetFloat(dec.decodeFloat64())
	return nil
}

func (dec *Decoder) decodekString(rv reflect.Value) error {
	rv.SetString(dec.decodeString())
	return nil
}

func (dec *Decoder) decodekBool(rv reflect.Value) error {
	rv.SetBool(dec.decodeBool())
	return nil
}

func (dec *Decoder) decodekInterface(rv reflect.Value) error {
	if !rv.IsNil() {
		return dec.decode(rv.Elem())
	}

	// blind decoding
	v, vk, err := dec.blind()
	if err != nil {
		return err
	}
	decodeFurther := false
	if v == nil {
		decodeFurther = true
	}

	// check for nil and undef values
	if vk == reflect.Invalid {
		return nil
	}

	// process the data
	switch vk {
	case reflect.Slice:
		v = new([]interface{})
	case reflect.Map:
		v = new(map[interface{}]interface{})
	}

	if decodeFurther {
		if v != nil {
			dec.decode(reflect.ValueOf(v).Elem())
		}
	}
	if v != nil {
		rv.Set(reflect.ValueOf(v))
	}
	return nil
}

// Decoce into a slice
func (dec *Decoder) decodekSlice(rv reflect.Value) error {
	_, info := dec.parser.parseHeader()
	rvt := rv.Type()
	if info != cborIndefinite {
		length := int(dec.parser.buflen())
		if rv.IsNil() {
			rv.Set(reflect.MakeSlice(rvt, length, length))
		}
		for i := 0; i < length; i++ {
			if _, _, err := dec.parser.parseInformation(); err != nil {
				return err
			}
			if err := dec.decode(rv.Index(i)); err != nil {
				return err
			}
		}
	} else {
		rvti := rvt.Elem() // elements type for the slice
		rv.Set(reflect.MakeSlice(rvt, 0, 0))
		for i := 0; ; i++ {
			if _, _, err := dec.parser.parseInformation(); err != nil {
				return err
			}
			if dec.parser.isBreak() {
				break
			}
			rv.Set(reflect.Append(rv, reflect.Zero(rvti)))
			if err := dec.decode(rv.Index(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (dec *Decoder) decodekArray(rv reflect.Value) error {
	return dec.decodekSlice(rv.Slice(0, rv.Len()))
}

// Decode into a map, if the strict mode is not enforced and
// there is a duplicated key in the map, the behavior is
// undefined, if the values are all of the same type then the
// key should be probably overwritten, if the values are empty
// interfaces, then probably the first value assigned to the
// key will be stuck but there is nothing that guarantee this
//
// For more information about the strict mode take a look at
// the RFC7049 in the secton 3.10. Strict Mode
func (dec *Decoder) decodekMap(rv reflect.Value) error {
	rvt := rv.Type()
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(rvt))
	}
	keytype := rvt.Key()
	valtype := rvt.Elem()

	_, info := dec.parser.parseHeader()
	if info != cborIndefinite {
		lenght := int(dec.parser.buflen())
		for i := 0; i < lenght; i++ {
			if err := dec.generateKeyValue(keytype, valtype, rv); err != nil {
				return err
			}
		}
	} else {
		for {
			if err := dec.generateKeyValue(keytype, valtype, rv); err != nil {
				if err != io.EOF {
					return err
				}
				if dec.parser.isBreak() { // check again to make sure we are all right
					break
				}
				return err
			}
		}
	}
	return nil
}

// Decode into an struct
//
// CBOR arrays and maps can be decoded into structs using a
// simple series of rules and conventions if the strict mode is
// not enforced and there is a duplicated key in the map (or
// array) , the behavior is totally undefined
//
// If the underlying CBOR structure is an array the convention
// is to use odds indexes as keys and even indexes as value as
// it was a map
//
// So the first value read from the CBOR data will be mapped
// into the `Name` field, the second into the Age field using
// their types.
//
// Tags can be used with maps as well in case that the keys
// names doesn't match with out struct fields names,
//		type MyOtherType struct {
//			Name string `cbor:"name"`
//			Age  uint8	`cbor:"how_old"`
//		}
//
// Is the Strict Mode is used, will also fail if it receives a
// key that doesn't match with any field of the struct or if
// there are more indexes than fields in the struct
// Note: that last behavior is not part of the RFC7049
//
// For more information about the strict mode take a look at
// the RFC7049 in the secton 3.10. Strict Mode
func (dec *Decoder) decodekStruct(rv reflect.Value) error {
	rv.Set(reflect.New(rv.Type()).Elem())
	major, _ := dec.parser.parseHeader()
	length := 0
	numFields := rv.NumField()
	array := true
	if major == cborDataMap {
		array = false
	}
	err := dec.checkStructLength(numFields, &length, array)
	if err != nil {
		return err
	}
	return dec.decodeInner(rv, numFields, length, array)
}

func (dec *Decoder) decodeInner(rv reflect.Value, nf, length int, array bool) error {
	shownKeys := map[string]struct{}{}
	for i := 0; ; i++ {
		if length == 0 && !dec.parser.indefinite {
			break
		}
		op, err := dec.checkRtStructLength(i, nf)
		if err != nil {
			return err
		}
		if op == d_BREAK {
			break
		} else if op == d_CONTINUE {
			length--
			continue
		}

		major, _, err := dec.parser.parseInformation()
		if err != nil {
			return err
		}
		if dec.parser.indefinite && dec.parser.isBreak() {
			break
		}

		// key must be a string
		if major < cborByteString || major > cborTextString {
			t := "map"
			if array {
				t = "array"
			}
			return fmt.Errorf("%s keys must be string, %s received", t, major)
		}
		key, err := dec.decodeStructFieldKey(shownKeys)
		if err != nil {
			return err
		}

		// let's decode the value and assign it to the struct field
		if err := dec.decodeStructFieldValue(rv, key, array); err != nil {
			if err == forceContinueError && !dec.strict {
				length--
				continue
			}
			return err
		}
		length--
	}
	return nil
}

// helper function to generate a pair key, value to decode into maps
func (dec *Decoder) generateKeyValue(ktype, vtype reflect.Type, rv reflect.Value) error {
	if _, _, err := dec.parser.parseInformation(); err != nil {
		return err
	}
	if dec.parser.isBreak() {
		return io.EOF
	}
	key := reflect.New(ktype).Elem()
	dec.decode(key)
	// check if the key exists when we are in strict mode
	if dec.strict {
		if rv.MapIndex(key).IsValid() {
			return NewStrictModeError(fmt.Sprintf("duplicated key %s in map", key))
		}
	}
	if _, _, err := dec.parser.parseInformation(); err != nil {
		return err
	}
	val := rv.MapIndex(key)
	if !val.IsValid() {
		val = reflect.New(vtype).Elem()
	}
	dec.decode(val)
	rv.SetMapIndex(key, val)
	return nil
}

// helper function that iterates over the fields
// of a struct looking for a specific tag
func (dec *Decoder) lookupStructTag(st reflect.Value, tag string, array bool) string {
	for i := 0; i < st.NumField(); i++ {
		field := st.Type().Field(i)
		t := field.Tag.Get("cbor")
		if array {
			if strings.Contains(strings.Trim(t, " "), ",index") {
				// just ignore errors
				index, err := strconv.Atoi(strings.Split(t, ",")[0])
				if err == nil {
					if i == index {
						return field.Name
					}
				}
			}
		}
		if t != "" {
			if strings.Contains(t, tag) {
				return field.Name
			}
		}
	}
	return ""
}

// common length checks for struct decoders
func (dec *Decoder) checkStructLength(nf int, length *int, array bool) error {
	if !dec.parser.indefinite {
		l := int(dec.parser.buflen())
		nlen := l
		if array {
			nlen /= 2
		}
		if nlen != nf {
			if dec.strict {
				msg := fmt.Sprintf(
					"destination struct fields num %d doesn't match map length %d",
					nf, nlen,
				)
				return NewStrictModeError(msg)
			}
		}
		*length = nlen
	}
	return nil
}

// common length in runtime check for struct decoders
func (dec *Decoder) checkRtStructLength(i, nf int) (uint, error) {
	if i > nf {
		// if strict mode is on, check for the right number of fields
		msg := fmt.Sprintf(
			"destination struct fields num %d doesn't match map length %d", nf, i)
		if dec.strict {
			return d_NOP, NewStrictModeError(msg)
		}
		log.Printf("warning strict-mode: %s\n", msg)
		if dec.parser.indefinite && dec.parser.isBreak() {
			return d_BREAK, nil
		}
		if _, _, err := dec.parser.parseInformation(); err != nil {
			return d_NOP, err
		}
		return d_CONTINUE, nil
	}
	return d_NOP, nil
}

// decodes a key to be used as a struct field in struct decoders
func (dec *Decoder) decodeStructFieldKey(shownKeys map[string]struct{}) (string, error) {
	key := dec.decodeString()
	if dec.strict {
		if _, ok := shownKeys[key]; ok {
			return "", NewStrictModeError(
				fmt.Sprintf("duplicated key %s in map", key))
		}
		shownKeys[key] = struct{}{}
	}
	return key, nil
}

// decode a value to be used as a struct field value in struct decoders
func (dec *Decoder) decodeStructFieldValue(rv reflect.Value, key string, array bool) error {
	var field reflect.Value
	if field = rv.FieldByName(key); !field.IsValid() {
		if field = rv.FieldByName(dec.lookupStructTag(rv, key, array)); !field.IsValid() {
			msg := fmt.Sprintf("key %s doesn't match with any field", key)
			if dec.strict {
				return NewStrictModeError(msg)
			}
			log.Printf("warning strict-mode: %s skipping...\n", msg)
			if _, _, err := dec.parser.parseInformation(); err != nil {
				return err
			}
			return forceContinueError
		}
	}
	if _, _, err := dec.parser.parseInformation(); err != nil {
		return err
	}
	err := dec.decode(field)
	return err
}
