package parser

import (
	"bytes"
	"errors"
	"io"

	"github.com/DamnWidget/cbor/internal/test"
	"testing"
)

// TestScan tests read operations on a reader value
func TestScan(t *testing.T) {

	tests := []struct {
		name      string
		data      []byte
		expected  [][]byte
		nexpected []int
		n         []int
		err       error
	}{
		{
			name: "Read",
			data: []byte("Read Test"),
			expected: [][]byte{
				[]byte("Read"),
				[]byte(" "),
				[]byte("Test"),
			},
			nexpected: []int{4, 1, 4},
			n:         []int{4, 1, 4},
			err:       nil,
		},
		{
			name: "ScanEncoded",
			data: []byte{0x1b, 0x00, 0x00, 0x00, 0xe8, 0xd4, 0xa5, 0x10, 0x00},
			expected: [][]byte{
				[]byte{0x1b, 0x00, 0x00},
				[]byte{0x00, 0xe8, 0xd4},
				[]byte{0xa5},
				[]byte{0x10, 0x00},
			},
			nexpected: []int{3, 3, 1, 2},
			n:         []int{3, 3, 1, 2},
			err:       nil,
		},
		{
			name: "ReadLeftBytesWhenOutOfBoundaries",
			data: []byte{0xfb, 0x7e, 0x37, 0xe4, 0x3c, 0x88, 0x00, 0x75, 0x9c},
			expected: [][]byte{
				[]byte{0xfb, 0x7e, 0x37, 0xe4, 0x3c, 0x88, 0x00, 0x75, 0x9c, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			},
			nexpected: []int{9},
			n:         []int{15},
			err:       nil,
		},
		{
			name:      "ErrorWhenTryToReadBeyondCapacity",
			data:      []byte{},
			expected:  [][]byte{[]byte{0x0}},
			nexpected: []int{0},
			n:         []int{1},
			err:       io.EOF,
		},
		{
			name:      "ReturnEmptyIsBufLengthIsZero",
			data:      []byte{0x01, 0x02, 0x03},
			expected:  [][]byte{nil},
			nexpected: []int{0},
			n:         []int{0},
			err:       nil,
		},
	}

	t.Log("Given the need to scan n bytes from a CBOR-encoded reader")
	{
		for i, tt := range tests {
			tf := func(t *testing.T) {

				t.Logf("\tTest: %d\tWhen reading %d bytes from io.Reader", i, tt.n)
				{
					r := NewReader(bytes.NewBuffer(tt.data))
					for i := range tt.expected {

						data := make([]byte, tt.n[i])
						n, err := scan(r.r, data)
						if err != nil {
							if !errors.Is(err, tt.err) {
								t.Fatalf("\t%s\tShould not raise an error but %v was raised", test.Failed, err)
							}
						}

						if !test.BytesEqual(data, tt.expected[i]) {
							t.Fatalf("\t%s\tExpected: %v but got %v", test.Failed, tt.expected[i], data)
						}

						if n != tt.nexpected[i] {
							t.Fatalf("\t%s\tExpected to read %d bytes but %d were read", test.Failed, tt.nexpected[i], n)
						}
					}

					t.Logf("\t%s\tShould be scanned successfully", test.Succeed)
				}
			}

			t.Run(tt.name, tf)
		}
	}
}

// TestReader tests CBOR reader operations to guarantee RFC7049 compliance
func TestReader(t *testing.T) {

	tests := []struct {
		name      string
		data      []byte
		expected  []byte
		nexpected int
		offset    int
		n         int
	}{
		{
			name:      "Read",
			data:      []byte("Read Test"),
			expected:  []byte("Read"),
			nexpected: 4,
			offset:    0,
			n:         4,
		},
		{
			name:      "ReadEncodedFromOffset",
			data:      []byte{0x1b, 0x00, 0x00, 0x00, 0xe8, 0xd4, 0xa5, 0x10, 0x00},
			expected:  []byte{0xa5, 0x10},
			nexpected: 2,
			offset:    6,
			n:         2,
		},
		{
			name:      "ReadBytesLeftWhenNIsOutbounds",
			data:      []byte{0xfb, 0x7e, 0x37, 0xe4, 0x3c, 0x88, 0x00, 0x75, 0x9c},
			expected:  []byte{0x75, 0x9c},
			nexpected: 2,
			offset:    7,
			n:         10,
		},
		{
			name:      "ReturnsEmptyIfNothingLeftToRead",
			data:      []byte{0xf9, 0x04, 0x00},
			expected:  nil,
			nexpected: 0,
			offset:    3,
			n:         5,
		},
	}

	t.Log("Given the need to read n bytes from a CBOR-encoded data buffer")
	{
		for i, tt := range tests {
			tf := func(t *testing.T) {

				t.Logf("\tTest: %d\tWhen reading %d bytes from buffer %v", i, tt.n, tt.data)
				{
					r := reader{buf: tt.data, offset: tt.offset}
					n, result := r.take(tt.n)
					if !test.BytesEqual(result, tt.expected) {
						t.Fatalf("\t%s\tExpected: %v but got %v", test.Failed, tt.expected, result)
					}

					if n != tt.nexpected {
						t.Fatalf("\t%s\tExpected to read %d bytes but %d were read", test.Failed, tt.nexpected, n)
					}

					t.Logf("\t%s\tShould be read successfully", test.Succeed)
				}
			}

			t.Run(tt.name, tf)
		}
	}

	t.Log("Given the need to make multiple read of n bytes from a CBOR-encoded data buffer")
	{

		name := "SucesiveReadsDepleteDataAsExpected"
		data := []byte{0xc1, 0xfb, 0x41, 0xd4, 0x52, 0xd9, 0xec, 0x20, 0x00, 0x00}
		expected := [][]byte{
			[]byte{0xc1, 0xfb},
			[]byte{0x41, 0xd4},
			[]byte{0x52, 0xd9},
			[]byte{0xec, 0x20},
			[]byte{0x00, 0x00},
			nil,
		}
		nexpected := 2
		offset := 0
		n := 2

		t.Logf("\tTest: %d\tWhen reading %d bytes multiple times from buffer %v", len(tests)+1, n, data)
		{
			fn := func(t *testing.T) {

				r := reader{buf: data, offset: offset}
				for i := range expected {
					nr, result := r.take(n)
					if !test.BytesEqual(result, expected[i]) {
						t.Fatalf("\t%s\t Expected: %v but got %v", test.Failed, expected[i], result)
					}

					if nr != nexpected && (expected[i] != nil) {
						t.Fatalf("\t%s\tExpected to read %d bytes but %d were read", test.Failed, n, nr)
					}

				}
				t.Logf("\t%s\tShould be read successfully", test.Succeed)
			}

			t.Run(name, fn)
		}
	}
}

//
// Benchmarks
//

var gn int
var ge error

func BenchmarkScan(b *testing.B) {

	var n int
	var e error

	buf := new(bytes.Buffer)
	data := []byte{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}
	for i := 0; i < b.N; i++ {
		buf.Write([]byte{0x1b, 0x00, 0x00, 0x00, 0xe8, 0xd4, 0xa5, 0x10, 0x00})
		r := NewReader(buf)
		n, e = scan(r.r, data)
		if e != nil {
			b.Fatalf("error encountered: %v", e)
		}
		buf.Reset()
	}

	gn = n
	ge = e
}
