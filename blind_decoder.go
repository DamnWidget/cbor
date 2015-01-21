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
		v = dec.decodeBigNum()
	case absoluteNegativeBigNum:
		vk = bigNum
		v = new(big.Int).Neg(dec.decodeBigNum())
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
			vk = reflect.Ptr
		}
	}

	if vk == 0 {
		return nil, 0, fmt.Errorf("blind: Unrecognized header 0x%x", header)
	}
	return v, vk, nil
}
