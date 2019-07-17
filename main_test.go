package main

import (
	"bytes"
	"testing"
)

func TestNameAt(t *testing.T) {
	var testNames = []struct {
		bin      []byte
		offset   int
		expected string
	}{
		{
			bin:      []byte{3, 110, 115, 49, 3, 99, 111, 109, 0},
			offset:   0,
			expected: "ns1.com",
		},
		{
			bin:      []byte{0, 0, 0, 3, 110, 115, 49, 3, 99, 111, 109, 0},
			offset:   3,
			expected: "ns1.com",
		},
		{
			bin:      []byte{0, 0, 0, 3, 110, 115, 49, 0, 3, 102, 111, 111, 192, 3},
			offset:   8,
			expected: "foo.ns1",
		},
	}
	for _, tn := range testNames {
		t.Run(tn.expected, func(t *testing.T) {
			r := bytes.NewBuffer(tn.bin)
			name, err := NameAt(r, tn.offset)
			if err != nil {
				t.Fatal(err)
			}
			if name != tn.expected {
				t.Fatalf("wanted %+v, got %+v\n", tn.expected, name)
			}

		})
	}
}

func TestIsPointer(t *testing.T) {
	aPointer := 0xC0
	notPointer := 0x6

	is := IsPointer(aPointer)
	if !is {
		t.Fatal("Badd")
	}

	is = IsPointer(notPointer)
	if is {
		t.Fatal("Naddasdasd")
	}
}
