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
	"testing"
)

func TestWrite(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	c := newComposer(buf)
	n, err := c.write([]byte{absoluteArray})
	check(err)
	expect(n, 1, t, "TestWrite")
	b := []byte{0x10, 0x20, 0x30, 0x40, 0x50}
	expect(buf.Bytes()[0], byte(absoluteArray), t, "TestWrite")
	n, err = c.write(b)
	check(err)
	expect(n, 5, t, "TestWrite")
	for i, elem := range buf.Bytes()[1:] {
		expect(elem, b[i], t, "TestWrite")
	}
}

func TestComposeInt(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	c := newComposer(buf)
	var i int64 = 1936
	n, err := c.composeInt(i)
	check(err)
	expect(n, 2, t, "TestComposeInt")
	expect(buf.Bytes()[0], byte(0x19), t, "TestComposeInt")
	expect(buf.Bytes()[1], byte(0x07), t, "TestComposeInt")
	expect(buf.Bytes()[2], byte(0x90), t, "TestComposeInt")
	buf.Reset()
	i = 56
	n, err = c.composeInt(i)
	check(err)
	expect(n, 1, t, "TestComposeInt")
	expect(buf.Bytes()[0], byte(0x18), t, "TestComposeInt")
	expect(int8(buf.Bytes()[1]), int8(56), t, "TestComposeInt")
	buf.Reset()
	i = -1936
	n, err = c.composeInt(i)
	check(err)
	expect(n, 2, t, "TestComposeInt")
	expect(buf.Bytes()[0], byte(0x39), t, "TestComposeInt")
	expect(buf.Bytes()[1], byte(0x07), t, "TestComposeInt")
	expect(buf.Bytes()[2], byte(0x8f), t, "TestComposeInt")
	buf.Reset()
	i = -56
	n, err = c.composeInt(i)
	check(err)
	expect(n, 1, t, "TestComposeInt")
	expect(buf.Bytes()[0], byte(0x38), t, "TestComposeInt")
	expect(int8(buf.Bytes()[1]), int8(55), t, "TestComposeInt")
}

func TestComposeUint(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	c := newComposer(buf)
	var i uint64 = 67524
	n, err := c.composeUint(i)
	check(err)
	expect(n, 4, t, "TestComposeUint")
	expect(buf.Bytes()[0], byte(0x1a), t, "TestComposeUint")
	expect(buf.Bytes()[1], byte(0x00), t, "TestComposeUint")
	expect(buf.Bytes()[2], byte(0x01), t, "TestComposeUint")
	expect(buf.Bytes()[3], byte(0x07), t, "TestComposeUint")
	expect(buf.Bytes()[4], byte(0xc4), t, "TestComposeUint")
	buf.Reset()
	i = 1
	n, err = c.composeUint(i)
	check(err)
	expect(n, 1, t, "TestComposeUint")
	expect(buf.Bytes()[0], byte(0x01), t, "TestComposeUint")
	expect(uint8(buf.Bytes()[0]), uint8(1), t, "TestComposeUint")
}

func TestComposeBoolean(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	c := newComposer(buf)
	var v bool = false
	check(c.composeBoolean(v))
	expect(buf.Bytes()[0], byte(0xf4), t, "TestComposeBoolean")
	v = true
	check(c.composeBoolean(v))
	expect(buf.Bytes()[1], byte(0xf5), t, "TestComposeBoolean")
}
