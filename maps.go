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

import "reflect"

// type map
type typeMap map[Major]map[byte]reflect.Type

var expectedTypesMap typeMap = typeMap{
	cborUnsignedInt: map[byte]reflect.Type{
		cborUint8:  reflect.TypeOf(uint8(0)),
		cborUint16: reflect.TypeOf(uint16(0)),
		cborUint32: reflect.TypeOf(uint32(0)),
		cborUint64: reflect.TypeOf(uint64(0)),
	},
	cborNegativeInt: map[byte]reflect.Type{
		cborUint8:  reflect.TypeOf(int8(0)),
		cborUint16: reflect.TypeOf(int16(0)),
		cborUint32: reflect.TypeOf(int32(0)),
		cborUint64: reflect.TypeOf(int64(0)),
	},
	cborByteString: map[byte]reflect.Type{
		cborUint8:  reflect.TypeOf([]byte{}),
		cborUint16: reflect.TypeOf([]byte{}),
		cborUint32: reflect.TypeOf([]byte{}),
		cborUint64: reflect.TypeOf([]byte{}),
	},
	cborTextString: map[byte]reflect.Type{
		cborSmallInt: reflect.TypeOf(string("")),
		cborUint8:    reflect.TypeOf(string("")),
		cborUint16:   reflect.TypeOf(string("")),
		cborUint32:   reflect.TypeOf(string("")),
		cborUint64:   reflect.TypeOf(string("")),
	},
	cborNC: map[byte]reflect.Type{
		cborUint8:  reflect.TypeOf(byte(0)),
		cborUint16: reflect.TypeOf(float16(0)),
		cborUint32: reflect.TypeOf(float32(0)),
		cborUint64: reflect.TypeOf(float64(0)),
	},
}
