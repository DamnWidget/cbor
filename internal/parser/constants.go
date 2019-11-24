// A Golang RFC7049 implementation
// Copyright (C) 2015 - 2020 Oscar Campos

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package parser

// Mejor is a custom type to define possible CBOR-encoded data major types
// https://tools.ietf.org/html/rfc7049#section-2.1
type Major byte

// Major enum of all possible CBOR-encoded data major types
const (
	UnsignedIntMajorType Major = iota  // Non signed integers of any size
	NegativeIntMajorType  // Negative integers of any size
	ByteStringMajorType   // Sequence of bytes that represents a string
	TextStringMajorType   // Sequence of UTF-8 [RFC3629] encoded bytes
	ArrayMajorType  	  // Array of arbitrary data
	MapMajorType  		  // Map of arbitrary data
	TagMajorType  		  // User defined tags
	NoContentMajorType 	  // floating point and simple data type that needs no content
)

// AdditionalInfo is a custom type to define all possible meanings
const (

)