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
	"fmt"
	"math/big"
	"reflect"
)

// function used to decode extended tag info
type handleTagDecFn func(*Decoder, interface{}) error

// tag maps, used by user code to register custom extensions
// using the major type 6 Optional Semantic Tagging for more
// information refer to http://tools.ietf.org/html/rfc7049#section-2.4
type extensionTagMap map[uint64]handleTagDecFn

// global extension tags map
var extensionTagDec extensionTagMap = make(extensionTagMap)

// register a new extension information tag in the tags register
func (e *extensionTagMap) register(tagInfo uint64, fn handleTagDecFn) error {
	if _, ok := extensionTagDec[tagInfo]; ok {
		return fmt.Errorf("0x%x tag information is already registered", tagInfo)
	}
	extensionTagDec[tagInfo] = fn
	return nil
}

// Look for a function registered to handle a given tag info
func (e *extensionTagMap) lookup(tagInfo uint64) (handleTagDecFn, error) {
	fn, ok := extensionTagDec[tagInfo]
	if !ok {
		return nil, fmt.Errorf(
			"0x%x not matched as registered tag extension handler", tagInfo)
	}
	return fn, nil
}

// Registers a new funtion to handle decode of tag extensions
func RegisterTagExtensionFn(tagInfo uint64, fn handleTagDecFn) error {
	return extensionTagDec.register(tagInfo, fn)
}

// decodes into v scanning the CBOR data that comes in the encoded data
func (dec *Decoder) blind() (v interface{}, vk reflect.Kind, err error) {
	header := dec.parser.header
	info := header & 0x1f
	switch header {
	case absoluteNil, absoluteUndef:
		vk = reflect.Invalid
	case absoluteFalse:
		vk = reflect.Bool
		v = false
	case absoluteTrue:
		vk = reflect.Bool
		v = true
	case absoluteFloat16, absoluteFloat32:
		vk = reflect.Float32
		if info == cborFloat16 {
			v = dec.decodeFloat16()
		} else {
			v = dec.decodeInt32()
		}
	case absoluteFloat64:
		vk = reflect.Float64
		v = dec.decodeFloat64()
	case absoluteIndefiniteBytes:
		vk = byteString
		v = dec.decodeBytes()
	case absoluteIndefiniteString:
		vk = reflect.String
		v = dec.decodeString()
	case absoluteIndefiniteArray:
		vk = reflect.Slice
	case absoluteIndefiniteMap:
		vk = reflect.Map
	case absolutePositiveBigNum:
		vk = bigNum
		v = dec.decodePositiveBigNum()
	case absoluteNegativeBigNum:
		vk = bigNum
		v = new(big.Int).Neg(dec.decodeNegativeBigNum())
	case absoluteStringDateTime:
		vk = stringDateTime
		v = dec.decodeStringDateTime()
	case absoluteEpochDateTime:
		vk = epochDateTime
		v = dec.decodeEpochDateTime()
	case absoluteDecimalFraction:
		vk = decimalFraction
		v = dec.decodeDecimalFraction()
	case absoluteBigFloat:
		vk = bigFloat
		v = dec.decodeBigFloat()
	case absoluteBase64Url:
		vk = base64Url
		v = dec.decodeBase64Url()
	case absoluteBase64String:
		vk = base64String
		v = dec.decodeBase64()
	case absoluteBase16String:
		vk = base16String
		v = dec.decodeBase16()
	default:
		// unsigned integers
		if header >= absoluteUint && header < absoluteInt {
			switch info {
			case cborSmallInt, cborUint8:
				vk = reflect.Uint8
				v = dec.decodeUint8()
			case cborUint16:
				vk = reflect.Uint16
				v = dec.decodeUint16()
			case cborUint32:
				vk = reflect.Uint32
				v = dec.decodeUint32()
			case cborUint64:
				vk = reflect.Uint64
				v = dec.decodeUint64()
			default:
				if info < cborSmallInt {
					vk = reflect.Uint8
					v = dec.decodeUint8()
				}
			}
		}
		// signed integers
		if header >= absoluteInt && header < absoluteBytes {
			switch info {
			case cborSmallInt, cborUint8:
				vk = reflect.Int8
				v = dec.decodeInt8()
			case cborUint16:
				vk = reflect.Int16
				v = dec.decodeInt16()
			case cborUint32:
				vk = reflect.Int32
				v = dec.decodeInt32()
			case cborUint64:
				vk = reflect.Int64
				v = dec.decodeInt64()
			default:
				if info < cborSmallInt {
					vk = reflect.Int8
					v = dec.decodeInt8()
				}
			}
		}
		// byte strings
		if header >= absoluteBytes && header < absoluteString {
			vk = byteString
			v = dec.decodeBytes()
		}
		// unicode string
		if header >= absoluteString && header < absoluteArray {
			vk = reflect.String
			v = dec.decodeString()
		}
		// slice
		if header >= absoluteArray && header < absoluteMap {
			vk = reflect.Slice
		}
		// map
		if header >= absoluteMap && header < absoluteTag {
			vk = reflect.Map
		}
		// tags
		if header >= absoluteTag && header < absoluteNoContent {
			tagInfo := dec.parser.buflen()
			switch tagInfo {
			case cborURI:
				vk = URI
				v = dec.decodeURI()
			case cborTextBase64Url:
				vk = URI
				v = dec.decodeBase64URI()
			case cborTextBase64:
				vk = reflect.String
				v = dec.decodeBase64String()
			case cborRegexp:
				vk = tagRegexp
				v = dec.decodeRegexp()
			case cborMime:
				vk = MIME
				v = dec.decodeMime()
			default:
				// lookup in the extended user defined tags
				fn, err := extensionTagDec.lookup(tagInfo)
				if err == nil {
					vk = reflect.Invalid
					if err := fn(dec, v); err != nil {
						return nil, 0, err
					}
				} else {
					vk = reflect.Ptr
				}
			}
		}
	}

	if vk == 0 {
		return nil, 0, fmt.Errorf("blind: Unrecognized header 0x%x", header)
	}
	return v, vk, nil
}
