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

// Header is a custom type to define the initial byte of a DataItem
// It contains all the information we ever need to decodde a chunk
// of data encoded with CBOR
//
// The first (or high order) 3 bits define the major type,
// the last (or low order) 5 bits define the additional information
// For more information about the the CBOR-encoded major types and
// additional information refer to the CBOR Specification Section 2.
// 	https://tools.ietf.org/html/rfc7049#section-2
type header byte

// DataItem is a CBOR-encoded chunk of information
type DataItem struct {
	b1   header
	data []byte
}

//
