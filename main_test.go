package main

import "testing"

func TestNameCompressions(t *testing.T) {
	t.Skip("Not Implemented")
}

func TestFetchOffset(t *testing.T) {
	aPointer := 0xC6
	anOffset := 6
	notPointer := 0x06

	offset, is := FetchOffset(aPointer)
	if !is {
		t.Fatal("Badd")
	}
	if offset != anOffset {
		t.Fatal("Badd offset")
	}
	offset, is = FetchOffset(notPointer)
	if is {
		t.Fatal("Naddasdasd")
	}
	if offset != anOffset {
		t.Fatal("Incorrect")
	}
}
