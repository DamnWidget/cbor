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

import "testing"

// Common message for failed tests
const failed = "%s: expected %v, got %v"

// Interfaces to map results of any type
type CurrentResult interface{}
type ExpectedResult interface{}

// convenience function for simple assertions
func expect(expected ExpectedResult, got CurrentResult, t *testing.T, from ...string) {
	f := "TestTake"
	if len(from) > 0 {
		f = from[0]
	}
	if expected != got {
		t.Errorf(failed, f, expected, got)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
