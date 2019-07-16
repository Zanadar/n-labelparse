package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		log.Fatal("<USAGE> $ namecompression /path/to/binary OFFSET")
	}
	r, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal("Problem", err)
	}

	offset, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		log.Fatal("Problem", err)
	}
	name, err := NameAt(r, offset)

	fmt.Println("Name:", name)

}

func NameAt(r io.Reader, offset int) (string, error) {
	buf := bufio.NewReader(r)

	lengthHeader, err := buf.ReadByte()
	if err != nil {
		return "", err
	}
	offset, isPointer := IsPointer(int(lengthHeader))
	if isPointer {
		fmt.Println("Pointer!")
	} else {
		fmt.Println("Not a pointer")
	}

	lengthByte := []byte{}
	n, err := r.Read(lengthByte)
	if n != 1 {
		return "", fmt.Errorf("Something went wrong")
	}
	if err != nil {
		return "", err
	}
	return "", nil
}

func IsPointer(b int) (offset int, is bool) {
	pointerMask := 0xC0
	pointerSet := (b & pointerMask) == pointerMask // clear 6 lowest bits, check what's left is the mask
	if !pointerSet {
		return offset, false
	}

	return (b ^ pointerMask), true
}
