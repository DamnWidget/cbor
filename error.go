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
	"reflect"
)

// An InvalidDecoderError describes an invalid argument passed to Decode
// (The argument to Decode must be a non nil pointer)
type InvalidDecodeError struct {
	Type reflect.Type
}

func (e *InvalidDecodeError) Error() string {
	if e.Type == nil {
		return "cbor: Decocde(nil)"
	}
	if e.Type.Kind() != reflect.Ptr {
		return fmt.Sprintf("cbor: Decode(non-pointer %s\n", e.Type)
	}
	return fmt.Sprintf("cbor: Decode(nil %s)\n", e.Type)
}

// An StrictModeError describes an invalid operation that violates
// the section 3.10. Strict Mode definition of the RFC7049
type StrictModeError struct {
	Msg string
}

func NewStrictModeError(msg string) *StrictModeError {
	return &StrictModeError{Msg: fmt.Sprintf("strict-mode: %s", msg)}
}

func (e *StrictModeError) Error() string {
	return e.Msg
}

// A CanonicalModeError describes an invalid operation that violates
// the section 3.9. Canonical CBOR definition of the RFC7049
type CanonicalModeError struct {
	Msg string
}

func NewCanonicalModeError(msg string) *CanonicalModeError {
	return &CanonicalModeError{Msg: fmt.Sprintf("canonical-mode: %s", msg)}
}

func (e *CanonicalModeError) Error() string {
	return e.Msg
}

type UnsupportedValueError struct {
	Value reflect.Value
	Str   string
}

func (e *UnsupportedValueError) Error() string {
	return "cbor: unsupported value: " + e.Str
}

type BigNumEncodeError struct {
	Value reflect.Value
	Str   string
}

func (e *BigNumEncodeError) Error() string {
	return "cbor: while encoding big num: " + e.Str
}

type BigFloatEncodeError struct {
	Value reflect.Value
	Str   string
}

func (e *BigFloatEncodeError) Error() string {
	return "cbor: while encoding big float: " + e.Str
}

type DateTimeEncodeError struct {
	Value reflect.Value
	Str   string
}

func (e *DateTimeEncodeError) Error() string {
	return "cbor: while encoding time.Time: " + e.Str
}

type StructEncodeError struct {
	Value reflect.Value
	Str   string
}

func (e *StructEncodeError) Error() string {
	return "cbor: while encoding struct type " + e.Value.Type().String() + ": " + e.Str
}
